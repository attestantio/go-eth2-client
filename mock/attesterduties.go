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

	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

// AttesterDuties obtains attester duties.
// If validatorIndicess is nil it will return all duties for the given epoch.
func (s *Service) AttesterDuties(ctx context.Context, epoch spec.Epoch, validatorIndices []spec.ValidatorIndex) ([]*api.AttesterDuty, error) {
	res := make([]*api.AttesterDuty, len(validatorIndices))
	for i := range validatorIndices {
		res[i] = &api.AttesterDuty{
			ValidatorIndex: validatorIndices[i],
		}
	}

	return res, nil
}
