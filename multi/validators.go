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

	eth2client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

// Validators provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIndices is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
func (s *Service) Validators(ctx context.Context,
	stateID string,
	validatorIndices []spec.ValidatorIndex,
) (
	map[spec.ValidatorIndex]*api.Validator,
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client eth2client.Service) (interface{}, error) {
		block, err := client.(eth2client.ValidatorsProvider).Validators(ctx, stateID, validatorIndices)
		if err != nil {
			return nil, err
		}
		return block, nil
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return res.(map[spec.ValidatorIndex]*api.Validator), nil
}
