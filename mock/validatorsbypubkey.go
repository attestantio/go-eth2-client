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

// ValidatorsByPubKey provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorPubKeys is a list of validator public keys to restrict the returned values.  If no validators public keys are
// supplied no filter will be applied.
func (s *Service) ValidatorsByPubKey(ctx context.Context, stateID string, validatorPubKeys []spec.BLSPubKey) (map[spec.ValidatorIndex]*api.Validator, error) {
	return map[spec.ValidatorIndex]*api.Validator{}, nil
}
