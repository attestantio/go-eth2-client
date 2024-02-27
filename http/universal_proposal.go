// Copyright © 2020 - 2024 Attestant Limited.
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
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
)

const (
	blindedHeader = "Eth-Execution-Payload-Blinded"
)

// UniversalProposal fetches a potential beacon block for signing.
// The returned block may be blinded or unblinded, depending on the current state of the network
// as decided by the execution and beacon nodes.
func (s *Service) UniversalProposal(ctx context.Context,
	opts *api.UniversalProposalOpts,
) (
	*api.Response[*api.VersionedUniversalProposal],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "UniversalProposal")
	defer span.End()

	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if opts.Slot == 0 {
		return nil, errors.New("no slot specified")
	}

	url := fmt.Sprintf("/eth/v3/validator/blocks/%d?randao_reveal=%#x&graffiti=%#x", opts.Slot, opts.RandaoReveal, opts.Graffiti)

	if opts.SkipRandaoVerification {
		if !opts.RandaoReveal.IsInfinity() {
			return nil, errors.New("randao reveal must be point at infinity if skip randao verification is set")
		}
		url = fmt.Sprintf("%s&skip_randao_verification", url)
	}

	if opts.BuilderBoostFactor != "" {
		url = fmt.Sprintf("%s&builder_boost_factor=%s", url, opts.BuilderBoostFactor)
	}

	httpResponse, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request block proposal v3")
	}

	var response *api.Response[*api.VersionedUniversalProposal]
	switch httpResponse.contentType {
	case ContentTypeSSZ:
		response, err = s.universalBeaconBlockProposalFromSSZ(httpResponse)
	case ContentTypeJSON:
		response, err = s.universalBeaconBlockProposalFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
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
		return nil, errors.New("universal block proposal not for requested slot")
	}

	// Only check the RANDAO reveal and graffiti if we are not connected to DVT middleware,
	// as the returned values will be decided by the middleware.
	if !s.connectedToDVTMiddleware {
		blockRandaoReveal, err := response.Data.RandaoReveal()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockRandaoReveal[:], opts.RandaoReveal[:]) {
			return nil, fmt.Errorf("universal block proposal has RANDAO reveal %#x; expected %#x", blockRandaoReveal[:], opts.RandaoReveal[:])
		}

		blockGraffiti, err := response.Data.Graffiti()
		if err != nil {
			return nil, err
		}
		if !bytes.Equal(blockGraffiti[:], opts.Graffiti[:]) {
			return nil, fmt.Errorf("universal block proposal has graffiti %#x; expected %#x", blockGraffiti[:], opts.Graffiti[:])
		}
	}

	return response, nil
}

func (s *Service) universalBeaconBlockProposalFromSSZ(res *httpResponse) (*api.Response[*api.VersionedUniversalProposal], error) {
	var universalProposal *api.VersionedUniversalProposal
	var metadata map[string]any

	blinded := res.headers[blindedHeader] == "true"
	if blinded {
		blindedProposalResponse, err := s.blindedProposalFromSSZ(res)
		if err != nil {
			return nil, err
		}
		universalProposal = &api.VersionedUniversalProposal{
			Blinded: blindedProposalResponse.Data,
		}
		metadata = blindedProposalResponse.Metadata
	} else {
		proposalResponse, err := s.beaconBlockProposalFromSSZ(res)
		if err != nil {
			return nil, err
		}
		universalProposal = &api.VersionedUniversalProposal{
			Full: proposalResponse.Data,
		}
		metadata = proposalResponse.Metadata
	}

	return &api.Response[*api.VersionedUniversalProposal]{
		Data:     universalProposal,
		Metadata: metadata,
	}, nil
}

func (s *Service) universalBeaconBlockProposalFromJSON(res *httpResponse) (*api.Response[*api.VersionedUniversalProposal], error) {
	var universalProposal *api.VersionedUniversalProposal
	var metadata map[string]any

	blinded := res.headers[blindedHeader] == "true"
	if blinded {
		blindedProposalResponse, err := s.blindedProposalFromJSON(res)
		if err != nil {
			return nil, err
		}
		universalProposal = &api.VersionedUniversalProposal{
			Blinded: blindedProposalResponse.Data,
		}
		metadata = blindedProposalResponse.Metadata
	} else {
		proposalResponse, err := s.beaconBlockProposalFromJSON(res)
		if err != nil {
			return nil, err
		}
		universalProposal = &api.VersionedUniversalProposal{
			Full: proposalResponse.Data,
		}
		metadata = proposalResponse.Metadata
	}

	return &api.Response[*api.VersionedUniversalProposal]{
		Data:     universalProposal,
		Metadata: metadata,
	}, nil
}
