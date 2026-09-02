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
	"unicode/utf8"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	dynssz "github.com/pk910/dynamic-ssz"
	"go.opentelemetry.io/otel"
)

// maxProposalValueDigits bounds the length of a value header before it is
// parsed.  The total supply of ether is under 10^27 wei, so 40 digits leaves
// generous room for any real value while keeping a hostile one cheap to reject.
const maxProposalValueDigits = 40

// maxEPBSResponseSize bounds the bodies of the ePBS fetch endpoints and, via
// post() (see http.go), every POST response body across the library.  Those
// bodies are an execution payload envelope, a payload-attestation datum or
// pool, or a small acknowledgement or error, so 64MiB is orders of magnitude
// above any legitimate response while capping a hostile or runaway one.  Block
// production is the one endpoint whose response can be far larger, and it
// passes its own limit below rather than raising this one for everybody.
const maxEPBSResponseSize = 64 * 1024 * 1024

// maxEPBSProposalResponseSize bounds a full SSZ payload-included block
// production response under the pinned Gloas static bounds: 4096 128KiB blobs
// plus 33,554,432 KZG proofs.
const maxEPBSProposalResponseSize = 2*1024*1024*1024 + 64*1024*1024

// maxEPBSProposalJSONResponseSize accommodates the same bounded values as JSON hex.
const maxEPBSProposalJSONResponseSize = 5 * 1024 * 1024 * 1024

// EPBSProposal fetches a potential ePBS beacon block for signing.
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

	if err := validateBuilderConfig(opts.BuilderConfig); err != nil {
		return nil, err
	}

	body, contentType, err := s.marshalRequestBody(ctx, opts.BuilderConfig)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/eth/v4/validator/blocks/%d", opts.Slot)
	headers := map[string]string{
		"Accept":                contentType.MediaType(),
		"Eth-Consensus-Version": spec.DataVersionGloas.String(),
	}

	responseLimit := maxEPBSProposalResponseSize
	if contentType == ContentTypeJSON {
		responseLimit = maxEPBSProposalJSONResponseSize
	}

	httpResponse, err := s.postWithResponseLimit(
		ctx,
		endpoint,
		query,
		&opts.Common,
		bytes.NewReader(body),
		contentType,
		headers,
		responseLimit,
	)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request epbs beacon block proposal"), err)
	}
	if err := populateConsensusVersionFromHeaders(httpResponse); err != nil {
		return nil, err
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

// assertEPBSProposalMatchesRequest checks the returned proposal matches the request.
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

	// A node must not include a payload that was not requested.
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

	return query, nil
}

func validateBuilderConfig(config *gloas.BuilderConfig) error {
	if config == nil {
		return errors.Join(errors.New("no builder config supplied"), client.ErrInvalidOptions)
	}
	// An empty builders list is a legitimate request -- it solicits no builder
	// bids, leaving only p2p ones -- so only the endpoint's maximum is enforced
	// here.  A nil slice says the same thing as an empty one and encodes
	// identically in both JSON and SSZ, so it is accepted too.
	if len(config.Builders) > 64 {
		return errors.Join(errors.New("too many builders supplied"), client.ErrInvalidOptions)
	}

	for i, builder := range config.Builders {
		if builder == nil {
			return errors.Join(fmt.Errorf("builder %d missing", i), client.ErrInvalidOptions)
		}
		if len(builder.URL) == 0 || len(builder.URL) > 2048 || !utf8.Valid(builder.URL) {
			return errors.Join(fmt.Errorf("builder %d has invalid URL", i), client.ErrInvalidOptions)
		}
		auth := builder.Auth
		if auth == nil || auth.Message == nil || len(auth.Message.Data) == 0 || len(auth.Message.Data) > 4096 {
			return errors.Join(fmt.Errorf("builder %d has invalid authorization", i), client.ErrInvalidOptions)
		}
		if len(builder.BuilderPubkeys) > 64 {
			return errors.Join(fmt.Errorf("builder %d has too many public keys", i), client.ErrInvalidOptions)
		}
	}

	return nil
}

func populateConsensusVersionFromHeaders(res *httpResponse) error {
	for key, value := range res.headers {
		if !strings.EqualFold(key, "Eth-Consensus-Version") {
			continue
		}

		if err := res.consensusVersion.UnmarshalJSON(fmt.Appendf(nil, "%q", value)); err != nil {
			return errors.Join(errors.New("failed to parse consensus version"), err)
		}

		return nil
	}

	return errors.New("no Eth-Consensus-Version header in epbs proposal response")
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
		Version: res.consensusVersion,
	}
	metadata := metadataFromHeaders(res.headers)

	if err := populateEPBSProposalValuesFromHeaders(proposal, res.headers); err != nil {
		return nil, err
	}

	// Custom presets require a request-scoped codec.  The same codec that
	// decodes the proposal has to compute its body root below: the generated
	// HashTreeRoot methods inline mainnet preset sizes and are wrong on any
	// other preset, and that applies just as much to a JSON-decoded proposal
	// as to an SSZ-decoded one.
	ds, err := s.dynSSZForRequest(ctx)
	if err != nil {
		return nil, err
	}

	switch res.contentType {
	case ContentTypeSSZ:
		included, err := epbsPayloadIncludedFromHeaders(res.headers)
		if err != nil {
			return nil, err
		}

		proposal.ExecutionPayloadIncluded = included

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
		jsonMetadata, err := decodeEPBSProposalJSON(res.body, proposal)
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

	block := proposal.Gloas
	if proposal.ExecutionPayloadIncluded && proposal.GloasContents != nil {
		block = proposal.GloasContents.Block
	}

	if block == nil || block.Body == nil {
		return nil, fmt.Errorf("no %s beacon block body in response", res.consensusVersion)
	}
	if block.Body.SignedExecutionPayloadBid == nil || block.Body.SignedExecutionPayloadBid.Message == nil {
		return nil, errors.Join(errors.New("no execution payload bid in epbs proposal response"), client.ErrInconsistentResult)
	}

	builderIndex := block.Body.SignedExecutionPayloadBid.Message.BuilderIndex
	proposal.BuilderIndex = &builderIndex

	// This is the one place BeaconBlockBodyRoot is set: ds is the codec that
	// decoded the block above, so this is the only spec-aware root available.
	// Every downstream consumer -- BodyRoot(), Root() and the guard below --
	// relies on this being correct rather than falling back to the generated,
	// mainnet-baked Body.HashTreeRoot().
	bodyRootRaw, err := ds.HashTreeRoot(block.Body)
	if err != nil {
		return nil, errors.Join(errors.New("failed to compute epbs proposal body root"), err)
	}

	bodyRoot := phase0.Root(bodyRootRaw)
	proposal.BeaconBlockBodyRoot = &bodyRoot

	if proposal.ExecutionPayloadIncluded {
		if err := assertIncludedEPBSProposalEnvelopeMatchesBlock(proposal, ds); err != nil {
			return nil, err
		}
	}

	return &api.Response[*api.VersionedEPBSProposal]{
		Data:     proposal,
		Metadata: metadata,
	}, nil
}

// populateEPBSProposalValuesFromHeaders reads the value components from response headers.
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

// decodeEPBSProposalJSON reads the payload-inclusion flag from a JSON response
// body and decodes the data into the container it selects.  The flag selects
// which of the two containers the data field carries, so it has to be read
// before the data itself can be decoded.
func decodeEPBSProposalJSON(body []byte, proposal *api.VersionedEPBSProposal) (map[string]any, error) {
	response := make(map[string]json.RawMessage)
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.Join(errors.New("failed to parse epbs proposal response"), err)
	}

	raw, exists := response["execution_payload_included"]
	if !exists {
		return nil, errors.New("no execution_payload_included in epbs proposal response")
	}
	var included *bool
	if err := json.Unmarshal(raw, &included); err != nil {
		return nil, errors.Join(errors.New("failed to unmarshal execution_payload_included"), err)
	}
	if included == nil {
		return nil, errors.New("execution_payload_included cannot be null")
	}
	proposal.ExecutionPayloadIncluded = *included

	data, exists := response["data"]
	if !exists || string(data) == "null" {
		return nil, fmt.Errorf("no %s epbs proposal in response", proposal.Version)
	}
	if proposal.ExecutionPayloadIncluded {
		proposal.GloasContents = new(apiv1gloas.BlockContents)
		if err := json.Unmarshal(data, proposal.GloasContents); err != nil {
			return nil, errors.Join(errors.New("failed to unmarshal data"), err)
		}
	} else {
		proposal.Gloas = new(gloas.BeaconBlock)
		if err := json.Unmarshal(data, proposal.Gloas); err != nil {
			return nil, errors.Join(errors.New("failed to unmarshal data"), err)
		}
	}

	metadata := make(map[string]any, len(response))
	for key, raw := range response {
		if key == "data" {
			continue
		}
		var value any
		if err := json.Unmarshal(raw, &value); err != nil {
			return nil, errors.Join(fmt.Errorf("failed to unmarshal metadata %s", key), err)
		}
		metadata[key] = value
	}

	return metadata, nil
}

// assertIncludedEPBSProposalEnvelopeMatchesBlock checks that the execution
// payload envelope that travelled with a proposal is for the block it
// travelled with.  It uses proposal.Root() rather than
// proposal.GloasContents.Block.HashTreeRoot(): the latter is the generated,
// mainnet-baked hasher this package's body-root fix exists to avoid, so
// calling it here would let the exact bug back in through the guard meant to
// catch it.  Root() shares its computation with every other consumer of the
// block root, so the guard and a proposer's signature cannot disagree.
func assertIncludedEPBSProposalEnvelopeMatchesBlock(proposal *api.VersionedEPBSProposal, ds *dynssz.DynSsz) error {
	contents := proposal.GloasContents
	if contents == nil || contents.Block == nil || contents.ExecutionPayloadEnvelope == nil {
		return errors.New("no block contents in epbs proposal response")
	}

	blockRoot, err := proposal.Root()
	if err != nil {
		return errors.Join(errors.New("failed to hash epbs proposal block"), err)
	}
	if contents.ExecutionPayloadEnvelope.BeaconBlockRoot != blockRoot {
		return errors.Join(errors.New("execution payload envelope is for a different block"), client.ErrInconsistentResult)
	}
	envelope := contents.ExecutionPayloadEnvelope
	bid := contents.Block.Body.SignedExecutionPayloadBid.Message
	if envelope.BuilderIndex != bid.BuilderIndex {
		return errors.Join(errors.New("execution payload envelope builder index does not match bid"), client.ErrInconsistentResult)
	}
	if envelope.Payload == nil || envelope.Payload.BlockHash != bid.BlockHash {
		return errors.Join(errors.New("execution payload block hash does not match bid"), client.ErrInconsistentResult)
	}
	if envelope.ParentBeaconBlockRoot != bid.ParentBlockRoot {
		return errors.Join(errors.New("execution payload envelope parent root does not match bid"), client.ErrInconsistentResult)
	}
	if contents.ExecutionPayloadEnvelope.ExecutionRequests == nil {
		return errors.Join(errors.New("execution payload envelope has no execution requests"), client.ErrInconsistentResult)
	}
	executionRequestsRoot, err := ds.HashTreeRoot(contents.ExecutionPayloadEnvelope.ExecutionRequests)
	if err != nil {
		return errors.Join(errors.New("failed to hash execution payload envelope requests"), err)
	}
	if phase0.Root(executionRequestsRoot) != contents.Block.Body.SignedExecutionPayloadBid.Message.ExecutionRequestsRoot {
		return errors.Join(errors.New("execution payload envelope requests do not match bid"), client.ErrInconsistentResult)
	}

	return nil
}
