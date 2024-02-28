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

package mock

import (
	"context"

	"github.com/attestantio/go-eth2-client/api"
)

// UniversalProposal fetches a universal proposal for signing.
func (s *Service) UniversalProposal(ctx context.Context,
	opts *api.UniversalProposalOpts,
) (
	*api.Response[*api.VersionedUniversalProposal], error,
) {
	if opts.BuilderBoostFactor != "" {
		blindedBlock, err := s.BlindedProposal(ctx, &api.BlindedProposalOpts{
			Slot:         opts.Slot,
			RandaoReveal: opts.RandaoReveal,
			Graffiti:     opts.Graffiti,
		})

		if err != nil {
			return nil, err
		}

		return &api.Response[*api.VersionedUniversalProposal]{
			Data: &api.VersionedUniversalProposal{
				BlindedProposal: blindedBlock.Data,
			},
			Metadata: blindedBlock.Metadata,
		}, nil

	} else {
		beaconBlock, err := s.Proposal(ctx, &api.ProposalOpts{
			Slot:         opts.Slot,
			RandaoReveal: opts.RandaoReveal,
			Graffiti:     opts.Graffiti,
		})

		if err != nil {
			return nil, err
		}

		return &api.Response[*api.VersionedUniversalProposal]{
			Data: &api.VersionedUniversalProposal{
				Proposal: beaconBlock.Data,
			},
			Metadata: beaconBlock.Metadata,
		}, nil
	}
}
