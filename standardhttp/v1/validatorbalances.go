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

package v1

import (
	"context"
	"encoding/json"
	"fmt"

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
)

type validatorBalancesJSON struct {
	Data []*api.ValidatorBalance `json:"data"`
}

// ValidatorBalances provides the validator balances for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validators is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
func (s *Service) ValidatorBalances(ctx context.Context, stateID string, validatorIDs []client.ValidatorIDProvider) (map[uint64]uint64, error) {
	if stateID == "" {
		return nil, errors.New("no state ID specified")
	}

	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/states/%s/validator_balances", stateID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request state validator balances")
	}

	var validatorBalancesJSON validatorBalancesJSON
	if err := json.NewDecoder(respBodyReader).Decode(&validatorBalancesJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse validator balances")
	}
	if validatorBalancesJSON.Data == nil {
		return nil, errors.New("no state validator balances returned")
	}

	res := make(map[uint64]uint64)
	for _, validatorBalance := range validatorBalancesJSON.Data {
		res[validatorBalance.Index] = validatorBalance.Balance
	}
	return res, nil
}
