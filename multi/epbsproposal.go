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

package multi

import (
	"context"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
)

// EPBSProposal fetches an ePBS proposal for signing.
func (s *Service) EPBSProposal(ctx context.Context,
	opts *api.EPBSProposalOpts,
) (
	*api.Response[*api.VersionedEPBSProposal],
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		proposal, err := client.(consensusclient.EPBSProposalProvider).EPBSProposal(ctx, opts)
		if err != nil {
			return nil, err
		}

		return proposal, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	response, isResponse := res.(*api.Response[*api.VersionedEPBSProposal])
	if !isResponse {
		return nil, ErrIncorrectType
	}

	return response, nil
}
