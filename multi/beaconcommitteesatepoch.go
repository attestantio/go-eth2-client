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

// BeaconCommitteesAtEpoch fetches all beacon committees for the given epoch at the given state.
func (s *Service) BeaconCommitteesAtEpoch(ctx context.Context, stateID string, epoch phase0.Epoch) ([]*api.BeaconCommittee, error) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (interface{}, error) {
		beaconCommittees, err := client.(consensusclient.BeaconCommitteesProvider).BeaconCommitteesAtEpoch(ctx, stateID, epoch)
		if err != nil {
			return nil, err
		}
		return beaconCommittees, nil
	}, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.([]*api.BeaconCommittee), nil
}
