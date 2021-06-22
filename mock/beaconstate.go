// Copyright © 2020 Attestant Limited.
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

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

// BeaconState fetches a beacon state given a state ID.
func (s *Service) BeaconState(ctx context.Context, stateID string) (*spec.BeaconState, error) {
	return &spec.BeaconState{
		LatestBlockHeader:           &spec.BeaconBlockHeader{},
		ETH1Data:                    &spec.ETH1Data{},
		PreviousJustifiedCheckpoint: &spec.Checkpoint{},
		CurrentJustifiedCheckpoint:  &spec.Checkpoint{},
		FinalizedCheckpoint:         &spec.Checkpoint{},
	}, nil
}
