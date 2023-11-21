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

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// BeaconState fetches a beacon state given a state ID.
func (s *Service) BeaconState(_ context.Context, _ *api.BeaconStateOpts) (*api.Response[*spec.VersionedBeaconState], error) {
	data := &spec.VersionedBeaconState{
		Version: spec.DataVersionPhase0,
		Phase0: &phase0.BeaconState{
			LatestBlockHeader:           &phase0.BeaconBlockHeader{},
			ETH1Data:                    &phase0.ETH1Data{},
			PreviousJustifiedCheckpoint: &phase0.Checkpoint{},
			CurrentJustifiedCheckpoint:  &phase0.Checkpoint{},
			FinalizedCheckpoint:         &phase0.Checkpoint{},
		},
	}

	return &api.Response[*spec.VersionedBeaconState]{
		Data:     data,
		Metadata: make(map[string]any),
	}, nil
}
