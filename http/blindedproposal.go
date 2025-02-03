// Copyright © 2022 - 2024 Attestant Limited.
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
	"errors"
	"fmt"

	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"go.opentelemetry.io/otel"
)

// BlindedProposal fetches a proposal for signing.
//
// Deprecated: use `Proposal` instead.
func (s *Service) BlindedProposal(ctx context.Context,
	opts *api.BlindedProposalOpts,
) (
	*api.Response[*api.VersionedBlindedProposal],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "BlindedProposal")
	defer span.End()

	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}
	if opts.Slot == 0 {
		return nil, errors.Join(errors.New("no slot specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/eth/v1/validator/blinded_blocks/%d", opts.Slot)
	query := fmt.Sprintf("randao_reveal=%#x&graffiti=%#x", opts.RandaoReveal, opts.Graffiti)

	if opts.SkipRandaoVerification {
		if !opts.RandaoReveal.IsInfinity() {
			return nil, errors.Join(
				errors.New("randao reveal must be point at infinity if skip randao verification is set"),
				client.ErrInvalidOptions,
			)
		}
		query = fmt.Sprintf("%s&skip_randao_verification", query)
	}

	res, err := s.get(ctx, endpoint, query, &opts.Common, true)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request blinded beacon block proposal"), err)
	}

	var response *api.Response[*api.VersionedBlindedProposal]
	switch res.contentType {
	case ContentTypeSSZ:
		response, err = s.blindedProposalFromSSZ(res)
	case ContentTypeJSON:
		response, err = s.blindedProposalFromJSON(res)
	default:
		return nil, fmt.Errorf("unhandled content type %v", res.contentType)
	}
	if err != nil {
		return nil, err
	}

	// Ensure the data returned to us is as expected given our input.
	blockSlot, err := response.Data.Slot()
	if err != nil {
		return nil, err
	}
	if blockSlot != opts.Slot {
		return nil, errors.Join(
			fmt.Errorf("blinded beacon block proposal for slot %d; expected %d", blockSlot, opts.Slot),
			client.ErrInconsistentResult,
		)
	}

	// Only check the RANDAO reveal if we are not connected to DVT middleware,
	// as the returned values will be decided by the middleware.
	if !s.connectedToDVTMiddleware {
		blockRandaoReveal, err := response.Data.RandaoReveal()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockRandaoReveal[:], opts.RandaoReveal[:]) {
			return nil, errors.Join(
				fmt.Errorf(
					"blinded beacon block proposal has RANDAO reveal %#x; expected %#x",
					blockRandaoReveal[:],
					opts.RandaoReveal[:],
				),
				client.ErrInconsistentResult,
			)
		}
	}

	return response, nil
}

func (*Service) blindedProposalFromSSZ(res *httpResponse) (*api.Response[*api.VersionedBlindedProposal], error) {
	response := &api.Response[*api.VersionedBlindedProposal]{
		Data: &api.VersionedBlindedProposal{
			Version: res.consensusVersion,
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	switch res.consensusVersion {
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix = &apiv1bellatrix.BlindedBeaconBlock{}
		if err := response.Data.Bellatrix.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Join(errors.New("failed to decode bellatrix blinded beacon block proposal"), err)
		}
	case spec.DataVersionCapella:
		response.Data.Capella = &apiv1capella.BlindedBeaconBlock{}
		if err := response.Data.Capella.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Join(errors.New("failed to decode capella blinded beacon block proposal"), err)
		}
	case spec.DataVersionDeneb:
		response.Data.Deneb = &apiv1deneb.BlindedBeaconBlock{}
		if err := response.Data.Deneb.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Join(errors.New("failed to decode deneb blinded beacon block proposal"), err)
		}
	case spec.DataVersionElectra:
		response.Data.Electra = &apiv1electra.BlindedBeaconBlock{}
		if err := response.Data.Electra.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Join(errors.New("failed to decode electra blinded beacon block proposal"), err)
		}
	default:
		return nil, fmt.Errorf("unhandled block proposal version %s", res.consensusVersion)
	}

	return response, nil
}

func (*Service) blindedProposalFromJSON(res *httpResponse) (*api.Response[*api.VersionedBlindedProposal], error) {
	response := &api.Response[*api.VersionedBlindedProposal]{
		Data: &api.VersionedBlindedProposal{
			Version: res.consensusVersion,
		},
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix, response.Metadata, err = decodeJSONResponse(
			bytes.NewReader(res.body),
			&apiv1bellatrix.BlindedBeaconBlock{},
		)
	case spec.DataVersionCapella:
		response.Data.Capella, response.Metadata, err = decodeJSONResponse(
			bytes.NewReader(res.body),
			&apiv1capella.BlindedBeaconBlock{},
		)
	case spec.DataVersionDeneb:
		response.Data.Deneb, response.Metadata, err = decodeJSONResponse(
			bytes.NewReader(res.body),
			&apiv1deneb.BlindedBeaconBlock{},
		)
	case spec.DataVersionElectra:
		response.Data.Electra, response.Metadata, err = decodeJSONResponse(
			bytes.NewReader(res.body),
			&apiv1electra.BlindedBeaconBlock{},
		)
	default:
		return nil, fmt.Errorf("unsupported version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}
