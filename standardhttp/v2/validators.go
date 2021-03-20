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

package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	api "github.com/attestantio/go-eth2-client/api/v2"
	spec "github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/pkg/errors"
)

type validatorsJSON struct {
	Data []*api.Validator `json:"data"`
}

// indexChunkSize is the maximum number of validator indices to send in each request.
// A request should be no more than 8,000 bytes to work with all currently-supported clients.
// An index has variable size, but assuming 7 characters, including the comma separator, is safe.
// We also need to reserve space for the state ID and the endpoint itself, to be safe we go
// with 500 bytes for this which results in us having comfortable space for 1,000 public keys.
var indexChunkSize = 1000

// Validators provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIndices is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
func (s *Service) Validators(ctx context.Context, stateID string, validatorIndices []spec.ValidatorIndex) (map[spec.ValidatorIndex]*api.Validator, error) {
	if stateID == "" {
		return nil, errors.New("no state ID specified")
	}

	if len(validatorIndices) > indexChunkSize {
		return s.chunkedValidators(ctx, stateID, validatorIndices)
	}

	url := fmt.Sprintf("/eth/v1/beacon/states/%s/validators", stateID)
	if len(validatorIndices) != 0 {
		ids := make([]string, len(validatorIndices))
		for i := range validatorIndices {
			ids[i] = fmt.Sprintf("%d", validatorIndices[i])
		}
		url = fmt.Sprintf("%s?id=%s", url, strings.Join(ids, ","))
	}

	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request validators")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain validators")
	}

	var validatorsJSON validatorsJSON
	if err := json.NewDecoder(respBodyReader).Decode(&validatorsJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse validators")
	}
	if validatorsJSON.Data == nil {
		return nil, errors.New("no validators returned")
	}

	res := make(map[spec.ValidatorIndex]*api.Validator)
	for _, validator := range validatorsJSON.Data {
		res[validator.Index] = validator
	}
	return res, nil
}

// chunkedValidators obtains the validators a chunk at a time.
func (s *Service) chunkedValidators(ctx context.Context, stateID string, validatorIndices []spec.ValidatorIndex) (map[spec.ValidatorIndex]*api.Validator, error) {
	res := make(map[spec.ValidatorIndex]*api.Validator)
	for i := 0; i < len(validatorIndices); i += indexChunkSize {
		chunkStart := i
		chunkEnd := i + indexChunkSize
		if len(validatorIndices) < chunkEnd {
			chunkEnd = len(validatorIndices)
		}
		chunk := validatorIndices[chunkStart:chunkEnd]
		chunkRes, err := s.Validators(ctx, stateID, chunk)
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain chunk")
		}
		for k, v := range chunkRes {
			res[k] = v
		}
	}
	return res, nil
}
