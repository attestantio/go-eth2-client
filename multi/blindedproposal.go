// Copyright © 2021, 2023 Attestant Limited.
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

package multi

import (
	"context"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
)

// BlindedProposal fetches a blinded proposal for signing.
func (s *Service) BlindedProposal(ctx context.Context,
	opts *api.BlindedProposalOpts,
) (
	*api.Response[*api.VersionedBlindedProposal],
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		block, err := client.(consensusclient.BlindedProposalProvider).BlindedProposal(ctx, opts)
		if err != nil {
			return nil, err
		}

		return block, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	return res.(*api.Response[*api.VersionedBlindedProposal]), nil
}
