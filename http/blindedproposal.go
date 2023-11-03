// Copyright © 2022, 2023 Attestant Limited.
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
	"fmt"

	"github.com/attestantio/go-eth2-client/api"
	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

// BlindedProposal fetches a proposal for signing.
func (s *Service) BlindedProposal(ctx context.Context,
	opts *api.BlindedProposalOpts,
) (
	*api.Response[*api.VersionedBlindedProposal],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "BlindedProposal")
	defer span.End()

	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if opts.Slot == 0 {
		return nil, errors.New("no slot specified")
	}

	url := fmt.Sprintf("/eth/v1/validator/blinded_blocks/%d?randao_reveal=%#x&graffiti=%#x", opts.Slot, opts.RandaoReveal, opts.Graffiti)

	if opts.SkipRandaoVerification {
		if !opts.RandaoReveal.IsInfinity() {
			return nil, errors.New("randao reveal must be point at infinity if skip randao verification is set")
		}
		url = fmt.Sprintf("%s&skip_randao_verification", url)
	}

	res, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request blinded beacon block proposal")
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
		return nil, errors.New("blinded beacon block proposal not for requested slot")
	}

	// Only check the RANDAO reveal and graffiti if we are not connected to DVT middleware,
	// as the returned values will be decided by the middleware.
	if !s.connectedToDVTMiddleware {
		blockRandaoReveal, err := response.Data.RandaoReveal()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockRandaoReveal[:], opts.RandaoReveal[:]) {
			return nil, fmt.Errorf("blinded beacon block proposal has RANDAO reveal %#x; expected %#x", blockRandaoReveal[:], opts.RandaoReveal[:])
		}

		blockGraffiti, err := response.Data.Graffiti()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockGraffiti[:], opts.Graffiti[:]) {
			return nil, fmt.Errorf("blinded beacon block proposal has graffiti %#x; expected %#x", blockGraffiti[:], opts.Graffiti[:])
		}
	}

	return response, nil
}

func (s *Service) blindedProposalFromSSZ(res *httpResponse) (*api.Response[*api.VersionedBlindedProposal], error) {
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
			return nil, errors.Wrap(err, "failed to decode bellatrix blinded beacon block proposal")
		}
	case spec.DataVersionCapella:
		response.Data.Capella = &apiv1capella.BlindedBeaconBlock{}
		if err := response.Data.Capella.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode capella blinded beacon block proposal")
		}
	case spec.DataVersionDeneb:
		response.Data.Deneb = &apiv1deneb.BlindedBeaconBlock{}
		if err := response.Data.Deneb.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode deneb blinded beacon block proposal")
		}
	default:
		return nil, fmt.Errorf("unhandled block proposal version %s", res.consensusVersion)
	}

	return response, nil
}

func (s *Service) blindedProposalFromJSON(res *httpResponse) (*api.Response[*api.VersionedBlindedProposal], error) {
	response := &api.Response[*api.VersionedBlindedProposal]{
		Data: &api.VersionedBlindedProposal{
			Version: res.consensusVersion,
		},
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &apiv1bellatrix.BlindedBeaconBlock{})
	case spec.DataVersionCapella:
		response.Data.Capella, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &apiv1capella.BlindedBeaconBlock{})
	case spec.DataVersionDeneb:
		response.Data.Deneb, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &apiv1deneb.BlindedBeaconBlock{})
	default:
		return nil, fmt.Errorf("unsupported version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}
