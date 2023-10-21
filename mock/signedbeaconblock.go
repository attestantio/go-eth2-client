// Copyright © 2020, 2023 Attestant Limited.
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

// SignedBeaconBlock fetches a signed beacon block given a block ID.
func (s *Service) SignedBeaconBlock(_ context.Context, _ *api.SignedBeaconBlockOpts) (*api.Response[*spec.VersionedSignedBeaconBlock], error) {
	return &api.Response[*spec.VersionedSignedBeaconBlock]{
		Data: &spec.VersionedSignedBeaconBlock{
			Version: spec.DataVersionPhase0,
			Phase0: &phase0.SignedBeaconBlock{
				Message: &phase0.BeaconBlock{
					Body: &phase0.BeaconBlockBody{
						ETH1Data: &phase0.ETH1Data{},
					},
				},
			},
		},
		Metadata: make(map[string]any),
	}, nil
}
