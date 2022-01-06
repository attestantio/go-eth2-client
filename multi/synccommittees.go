// Copyright © 2021 Attestant Limited.
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
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// SyncCommittee fetches the sync committee for the given state.
func (s *Service) SyncCommittee(ctx context.Context, stateID string) (*api.SyncCommittee, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		block, err := client.(consensusclient.SyncCommitteesProvider).SyncCommittee(ctx, stateID)
		if err != nil {
			return nil, err
		}
		return block, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(*api.SyncCommittee), nil
}

// SyncCommitteeAtEpoch fetches the sync committee for the given epoch at the given state.
func (s *Service) SyncCommitteeAtEpoch(ctx context.Context, stateID string, epoch phase0.Epoch) (*api.SyncCommittee, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		block, err := client.(consensusclient.SyncCommitteesProvider).SyncCommitteeAtEpoch(ctx, stateID, epoch)
		if err != nil {
			return nil, err
		}
		return block, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(*api.SyncCommittee), nil
}
