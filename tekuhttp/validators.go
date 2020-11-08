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

package tekuhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// validatorsJSON handles the JSON returned from lighthouse.
type validatorsJSON struct {
	Validators []*validatorsValidatorJSON `json:"validators"`
}

type validatorsValidatorJSON struct {
	PubKey    string          `json:"pubkey"`
	Index     uint64          `json:"validator_index"`
	Balance   string          `json:"balance"`
	Validator *spec.Validator `json:"validator"`
}

// Validators provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIndices is a list of validator indices to restrict the returned values.  If no validators IDs are supplied no filter
// will be applied.
func (s *Service) Validators(ctx context.Context, stateID string, validatorIndices []spec.ValidatorIndex) (map[spec.ValidatorIndex]*api.Validator, error) {
	slot, err := s.SlotFromStateID(ctx, stateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain slot")
	}

	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain slots per epoch")
	}
	epoch := slot / slotsPerEpoch
	farFutureEpoch, err := s.FarFutureEpoch(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain far future epoch")
	}
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/beacon/validators?pageSize=9999999&epoch=%d", epoch))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request validators")
	}

	validatorsResponse := &validatorsJSON{}
	if err := json.NewDecoder(respBodyReader).Decode(&validatorsResponse); err != nil {
		return nil, errors.Wrap(err, "failed to parse validators")
	}

	res := make(map[spec.ValidatorIndex]*api.Validator, len(validatorsResponse.Validators))
	for _, validatorResp := range validatorsResponse.Validators {
		balance := uint64(0)
		if validatorResp.Balance != "" {
			balance, err = strconv.ParseUint(validatorResp.Balance, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse validator balance")
			}
		}
		res[spec.ValidatorIndex(validatorResp.Index)] = &api.Validator{
			Index:     spec.ValidatorIndex(validatorResp.Index),
			Status:    api.ValidatorToState(validatorResp.Validator, spec.Epoch(epoch), spec.Epoch(farFutureEpoch)),
			Validator: validatorResp.Validator,
			Balance:   spec.Gwei(balance),
		}
	}
	return res, nil
}
