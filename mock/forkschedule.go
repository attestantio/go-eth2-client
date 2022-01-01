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

// ForkSchedule provides details of past and future changes in the chain's fork version.
func (s *Service) ForkSchedule(ctx context.Context) ([]*spec.Fork, error) {
	return []*spec.Fork{
		{
			PreviousVersion: spec.Version{0x01, 0x02, 0x03, 0x04},
			CurrentVersion:  spec.Version{0x01, 0x02, 0x03, 0x04},
			Epoch:           0,
		},
		{
			PreviousVersion: spec.Version{0x01, 0x02, 0x03, 0x04},
			CurrentVersion:  spec.Version{0x11, 0x12, 0x13, 0x14},
			Epoch:           1024,
		},
	}, nil
}
