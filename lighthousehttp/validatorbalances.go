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

package lighthousehttp

import (
	"context"

	client "github.com/attestantio/go-eth2-client"
	"github.com/pkg/errors"
)

// ValidatorBalances provides the validator balances for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validators is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
func (s *Service) ValidatorBalances(ctx context.Context, stateID string, validatorIDs []client.ValidatorIDProvider) (map[uint64]uint64, error) {
	// Lighthouse does not have a separate balances endpoint, so fetch validators and pull out the balance info.
	validators, err := s.Validators(ctx, stateID, validatorIDs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain validators for balances")
	}

	res := make(map[uint64]uint64)
	for index, validator := range validators {
		res[index] = validator.Balance
	}

	return res, nil
}
