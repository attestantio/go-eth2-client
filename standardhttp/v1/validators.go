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

type validatorsJSON struct {
	Data []*api.Validator `json:"data"`
}

// Validators provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validators is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
func (s *Service) Validators(ctx context.Context, stateID string, validators []client.ValidatorIDProvider) (map[uint64]*api.Validator, error) {
	if stateID == "" {
		return nil, errors.New("no state ID specified")
	}

	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/states/%s/validators", stateID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request state validators")
	}

	var validatorsJSON validatorsJSON
	if err := json.NewDecoder(respBodyReader).Decode(&validatorsJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse validators")
	}
	if validatorsJSON.Data == nil {
		return nil, errors.New("no state validators returned")
	}

	res := make(map[uint64]*api.Validator)
	for _, validator := range validatorsJSON.Data {
		res[validator.Index] = validator
	}
	return res, nil
}
