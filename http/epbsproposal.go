// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math/big"
	"strconv"
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"go.opentelemetry.io/otel"
)

// maxProposalValueDigits bounds the length of a value header before it is
// parsed.  The total supply of ether is under 10^27 wei, so 40 digits leaves
// generous room for any real value while keeping a hostile one cheap to reject.
const maxProposalValueDigits = 40

// EPBSProposal fetches a potential ePBS beacon block for signing.
//
// This is the gloas-onwards block production endpoint.  Post-Gloas a proposer
// commits to an execution payload bid rather than to a payload, so there is no
// blinded proposal; the axis that replaced it is opts.IncludePayload, which
// decides whether the execution payload envelope travels with the block or
// stays cached on the producing node.
func (s *Service) EPBSProposal(ctx context.Context,
	opts *api.EPBSProposalOpts,
) (
	*api.Response[*api.VersionedEPBSProposal],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "EPBSProposal")
	defer span.End()

	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}

	query, err := epbsProposalQuery(opts)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/eth/v4/validator/blocks/%d", opts.Slot)

	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common, true)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request epbs beacon block proposal"), err)
	}

	response, err := s.epbsProposalFromResponse(ctx, httpResponse)
	if err != nil {
		return nil, err
	}

	if err := s.assertEPBSProposalMatchesRequest(response.Data, opts); err != nil {
		return nil, err
	}

	return response, nil
}

// assertEPBSProposalMatchesRequest checks that the proposal a node returned is
// the one that was asked for.  Everything checked here is something the caller is
// about to sign or to rely on when publishing, so a node that answers with
// something else must be refused rather than obeyed.
func (s *Service) assertEPBSProposalMatchesRequest(proposal *api.VersionedEPBSProposal,
	opts *api.EPBSProposalOpts,
) error {
	blockSlot, err := proposal.Slot()
	if err != nil {
		return err
	}

	if blockSlot != opts.Slot {
		return errors.Join(
			fmt.Errorf("epbs beacon block proposal for slot %d; expected %d", blockSlot, opts.Slot),
			client.ErrInconsistentResult,
		)
	}

	// Only the node volunteering a payload is a fault.  Including one that was
	// not asked for changes where the caller is able to publish, which is not
	// the node's to decide.  The other direction is what the spec requires:
	// include_payload only governs self-building, and a block built on an
	// external builder's bid comes back alone whatever was asked for, since the
	// node does not hold the builder's payload.  So an excluded payload can no
	// longer be told apart from a node quietly defaulting the parameter, and the
	// builder path is worth more than catching that.
	if proposal.ExecutionPayloadIncluded && !*opts.IncludePayload {
		return errors.Join(
			fmt.Errorf("epbs beacon block proposal has execution payload included %t; expected %t",
				proposal.ExecutionPayloadIncluded, *opts.IncludePayload),
			client.ErrInconsistentResult,
		)
	}

	// Only check the RANDAO reveal if we are not connected to DVT middleware,
	// as the returned values will be decided by the middleware.
	if s.connectedToDVTMiddleware {
		return nil
	}

	blockRandaoReveal, err := proposal.RandaoReveal()
	if err != nil {
		return err
	}

	if !bytes.Equal(blockRandaoReveal[:], opts.RandaoReveal[:]) {
		return errors.Join(
			fmt.Errorf("epbs beacon block proposal has RANDAO reveal %#x; expected %#x",
				blockRandaoReveal[:], opts.RandaoReveal[:]),
			client.ErrInconsistentResult,
		)
	}

	return nil
}

// epbsProposalQuery validates the options and builds the endpoint's query
// string.  Both required parameters are checked here, before anything reaches
// the network.
func epbsProposalQuery(opts *api.EPBSProposalOpts) (string, error) {
	if opts == nil {
		return "", client.ErrNoOptions
	}

	if opts.Slot == 0 {
		return "", errors.Join(errors.New("no slot specified"), client.ErrInvalidOptions)
	}

	// The spec marks include_payload required with no default, and the two
	// modes are not interchangeable: false constrains the block to being
	// published through the node that produced it.  There is no safe value to
	// assume, so an unset one is refused.
	if opts.IncludePayload == nil {
		return "", errors.Join(errors.New("no payload inclusion specified"), client.ErrInvalidOptions)
	}

	query := fmt.Sprintf("randao_reveal=%#x&graffiti=%#x&include_payload=%t",
		opts.RandaoReveal, opts.Graffiti, *opts.IncludePayload)

	if opts.SkipRandaoVerification {
		if !opts.RandaoReveal.IsInfinity() {
			return "", errors.Join(
				errors.New("randao reveal must be point at infinity if skip randao verification is set"),
				client.ErrInvalidOptions,
			)
		}

		query += "&skip_randao_verification"
	}

	// Unlike the v3 endpoint, builder_boost_factor is optional here with a
	// server-side default, so an unset one is left off the query rather than
	// materialised client-side.
	if opts.BuilderBoostFactor != nil {
		query = fmt.Sprintf("%s&builder_boost_factor=%d", query, *opts.BuilderBoostFactor)
	}

	return query, nil
}

// epbsProposalFromResponse decodes a fetched ePBS block-production response.
// It is separate from the fetch so that the paths a live node will not produce
// on demand remain testable without one.
func (s *Service) epbsProposalFromResponse(ctx context.Context,
	res *httpResponse,
) (
	*api.Response[*api.VersionedEPBSProposal],
	error,
) {
	if res.consensusVersion != spec.DataVersionGloas {
		return nil, fmt.Errorf("epbs proposal not available for version %s", res.consensusVersion)
	}

	proposal := &api.VersionedEPBSProposal{
		Version:        res.consensusVersion,
		ConsensusValue: big.NewInt(0),
		ExecutionValue: big.NewInt(0),
	}
	metadata := metadataFromHeaders(res.headers)

	if err := populateEPBSProposalValuesFromHeaders(proposal, res.headers); err != nil {
		return nil, err
	}

	switch res.contentType {
	case ContentTypeSSZ:
		included, err := epbsPayloadIncludedFromHeaders(res.headers)
		if err != nil {
			return nil, err
		}

		proposal.ExecutionPayloadIncluded = included

		// A gloas block carries preset-derived fixed-size fields, so the
		// generated codec only reads bodies encoded at the compiled-in mainnet
		// preset.  The request-scoped codec is built from the spec the node
		// reports, which is what makes a non-mainnet preset decodable.
		ds, err := s.dynSSZForRequest(ctx)
		if err != nil {
			return nil, err
		}

		if included {
			proposal.GloasContents = &apiv1gloas.BlockContents{}
			err = ds.UnmarshalSSZ(proposal.GloasContents, res.body)
		} else {
			proposal.Gloas = &gloas.BeaconBlock{}
			err = ds.UnmarshalSSZ(proposal.Gloas, res.body)
		}

		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to decode %s SSZ epbs proposal", res.consensusVersion),
				err,
			)
		}
	case ContentTypeJSON:
		included, err := epbsPayloadIncludedFromBody(res.body)
		if err != nil {
			return nil, err
		}

		proposal.ExecutionPayloadIncluded = included

		// The decoder only writes the keys it finds, so seeding it with a typed
		// nil pointer is what makes a body with no data key in it
		// distinguishable from one carrying a zero-valued datum.
		var jsonMetadata map[string]any

		if included {
			proposal.GloasContents, jsonMetadata, err = decodeJSONResponse(
				bytes.NewReader(res.body),
				(*apiv1gloas.BlockContents)(nil),
			)
		} else {
			proposal.Gloas, jsonMetadata, err = decodeJSONResponse(
				bytes.NewReader(res.body),
				(*gloas.BeaconBlock)(nil),
			)
		}

		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to decode %s JSON epbs proposal", res.consensusVersion),
				err,
			)
		}

		maps.Copy(metadata, jsonMetadata)
	default:
		return nil, fmt.Errorf("unhandled content type %v", res.contentType)
	}

	// Neither an absent data key nor an explicit null leaves a proposal behind.
	// Returning success here would hand back a value whose every accessor
	// errors, discovered only by whichever the caller reaches first.
	if proposal.IsEmpty() {
		return nil, fmt.Errorf("no %s epbs proposal in response", res.consensusVersion)
	}

	return &api.Response[*api.VersionedEPBSProposal]{
		Data:     proposal,
		Metadata: metadata,
	}, nil
}

// populateEPBSProposalValuesFromHeaders reads the two value components from the
// response headers.  Both are left at the zero they were seeded with when their
// header is absent, which is what a node that predates the addition of
// execution_payload_value to the endpoint sends; a malformed value is an error,
// since a proposal whose worth silently reads as zero would lose an auction it
// should have won.
func populateEPBSProposalValuesFromHeaders(proposal *api.VersionedEPBSProposal,
	headers map[string]string,
) error {
	for k, v := range headers {
		var target **big.Int

		switch {
		case strings.EqualFold(k, "Eth-Consensus-Block-Value"):
			target = &proposal.ConsensusValue
		case strings.EqualFold(k, "Eth-Execution-Payload-Value"):
			target = &proposal.ExecutionValue
		default:
			continue
		}

		// Bound the input before parsing it.  big.Int.SetString is quadratic in
		// the length of its input, and this is a response header, so its length
		// is chosen by the node: a multi-megabyte value takes minutes to parse,
		// and it is parsed after the request's own deadline has been released,
		// so no caller timeout interrupts it.  A block's worth of wei needs 30
		// digits at the very most, well inside this.
		if len(v) > maxProposalValueDigits {
			return fmt.Errorf("proposal header %s has %d digits; more than the %d a value can need",
				k, len(v), maxProposalValueDigits)
		}

		value, isValid := new(big.Int).SetString(v, 10)
		if !isValid {
			return fmt.Errorf("proposal header %s %s not a valid integer", k, v)
		}

		if value.Sign() < 0 {
			return fmt.Errorf("proposal header %s %s is negative", k, v)
		}

		*target = value
	}

	return nil
}

// epbsPayloadIncludedFromHeaders reads the payload-inclusion flag from the
// Eth-Execution-Payload-Included response header, which is the only place it
// appears when the body is SSZ.  A missing or malformed value is an error
// rather than a default: it selects which container the body holds, so guessing
// would decode block contents as a bare block and discard the envelope, blobs
// and proofs the caller needs to publish.
func epbsPayloadIncludedFromHeaders(headers map[string]string) (bool, error) {
	for k, v := range headers {
		if !strings.EqualFold(k, "Eth-Execution-Payload-Included") {
			continue
		}

		included, err := strconv.ParseBool(v)
		if err != nil {
			// Report the header's name but not its value.  Its length is the
			// node's to choose, so echoing it puts however many megabytes the
			// node felt like sending into an error string and then into a log.
			return false, errors.New("proposal header Eth-Execution-Payload-Included is not a valid boolean")
		}

		return included, nil
	}

	return false, errors.New("no Eth-Execution-Payload-Included header in epbs proposal response")
}

// epbsPayloadIncludedFromBody reads the payload-inclusion flag from a JSON
// response body.  The flag selects which of the two containers the data field
// carries, so it has to be read before the data itself can be decoded.
func epbsPayloadIncludedFromBody(body []byte) (bool, error) {
	var wrapper struct {
		ExecutionPayloadIncluded *bool `json:"execution_payload_included"`
	}

	if err := json.Unmarshal(body, &wrapper); err != nil {
		return false, errors.Join(errors.New("failed to parse epbs proposal response"), err)
	}

	if wrapper.ExecutionPayloadIncluded == nil {
		return false, errors.New("no execution_payload_included in epbs proposal response")
	}

	return *wrapper.ExecutionPayloadIncluded, nil
}
