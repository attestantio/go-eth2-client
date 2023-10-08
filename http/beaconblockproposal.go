// Copyright © 2020 - 2023 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

// BeaconBlockProposal fetches a potential beacon block for signing.
//
//nolint:gocyclo
func (s *Service) BeaconBlockProposal(ctx context.Context,
	opts *api.BeaconBlockProposalOpts,
) (
	*api.Response[*spec.VersionedBeaconBlock],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "BeaconBlockProposal")
	defer span.End()

	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if opts.Slot == 0 {
		return nil, errors.New("no slot specified")
	}

	url := fmt.Sprintf("/eth/v2/validator/blocks/%d?randao_reveal=%#x&graffiti=%#x", opts.Slot, opts.RandaoReveal, opts.Graffiti)

	if opts.SkipRandaoVerification {
		if !opts.RandaoReveal.IsInfinity() {
			return nil, errors.New("randao reveal must be point at infinity if skip randao verification is set")
		}
		url = fmt.Sprintf("%s&skip_randao_verification", url)
	}

	res, err := s.get2(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon block proposal")
	}

	var response *api.Response[*spec.VersionedBeaconBlock]
	switch res.contentType {
	case ContentTypeSSZ:
		response, err = s.beaconBlockProposalFromSSZ(res)
	case ContentTypeJSON:
		response, err = s.beaconBlockProposalFromJSON(res)
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
		return nil, errors.New("beacon block proposal not for requested slot")
	}

	// Only check the RANDAO reveal and graffiti if we are not connected to DVT middleware,
	// as the returned values will be decided by the middleware.
	if !s.connectedToDVTMiddleware {
		blockRandaoReveal, err := response.Data.RandaoReveal()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockRandaoReveal[:], opts.RandaoReveal[:]) {
			return nil, fmt.Errorf("beacon block proposal has RANDAO reveal %#x; expected %#x", blockRandaoReveal[:], opts.RandaoReveal[:])
		}

		blockGraffiti, err := response.Data.Graffiti()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockGraffiti[:], opts.Graffiti[:]) {
			return nil, fmt.Errorf("beacon block proposal has graffiti %#x; expected %#x", blockGraffiti[:], opts.Graffiti[:])
		}
	}

	return response, nil
}

func (s *Service) beaconBlockProposalFromSSZ(res *httpResponse) (*api.Response[*spec.VersionedBeaconBlock], error) {
	response := &api.Response[*spec.VersionedBeaconBlock]{
		Data: &spec.VersionedBeaconBlock{
			Version: res.consensusVersion,
		},
		Metadata: metadataFromHeaders(res.headers),
	}

	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		response.Data.Phase0 = &phase0.BeaconBlock{}
		if err := response.Data.Phase0.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode phase0 beacon block proposal")
		}
	case spec.DataVersionAltair:
		response.Data.Altair = &altair.BeaconBlock{}
		if err := response.Data.Altair.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode altair beacon block proposal")
		}
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix = &bellatrix.BeaconBlock{}
		if err := response.Data.Bellatrix.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode bellatrix beacon block proposal")
		}
	case spec.DataVersionCapella:
		response.Data.Capella = &capella.BeaconBlock{}
		if err := response.Data.Capella.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode capella beacon block proposal")
		}
	case spec.DataVersionDeneb:
		response.Data.Deneb = &deneb.BeaconBlock{}
		if err := response.Data.Deneb.UnmarshalSSZ(res.body); err != nil {
			return nil, errors.Wrap(err, "failed to decode deneb beacon block proposal")
		}
	default:
		return nil, fmt.Errorf("unhandled block proposal version %s", res.consensusVersion)
	}

	return response, nil
}

func (s *Service) beaconBlockProposalFromJSON(res *httpResponse) (*api.Response[*spec.VersionedBeaconBlock], error) {
	response := &api.Response[*spec.VersionedBeaconBlock]{
		Data: &spec.VersionedBeaconBlock{
			Version: res.consensusVersion,
		},
	}

	var err error
	switch res.consensusVersion {
	case spec.DataVersionPhase0:
		response.Data.Phase0, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &phase0.BeaconBlock{})
	case spec.DataVersionAltair:
		response.Data.Altair, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &altair.BeaconBlock{})
	case spec.DataVersionBellatrix:
		response.Data.Bellatrix, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &bellatrix.BeaconBlock{})
	case spec.DataVersionCapella:
		response.Data.Capella, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &capella.BeaconBlock{})
	case spec.DataVersionDeneb:
		response.Data.Deneb, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body), &deneb.BeaconBlock{})
	default:
		err = fmt.Errorf("unsupported version %s", res.consensusVersion)
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}
