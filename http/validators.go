// Copyright © 2020 - 2023 Attestant Limited.
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

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	api "github.com/attestantio/go-eth2-client/api/v1"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type validatorsJSON struct {
	Data []*api.Validator `json:"data"`
}

// indexChunkSizes defines the per-beacon-node size of an index chunk.
// A request should be no more than 8,000 bytes to work with all currently-supported clients.
// An index has variable size, but assuming 7 characters, including the comma separator, is safe.
// We also need to reserve space for the state ID and the endpoint itself, to be safe we go
// with 500 bytes for this which results in us having comfortable space for 1,000 indices.
// That said, some nodes have their own built-in limits so use them where appropriate.
var indexChunkSizes = map[string]int{
	"default":    1000,
	"lighthouse": 1000,
	"nimbus":     1000,
	"prysm":      1000,
	"teku":       1000,
}

// indexChunkSize is the maximum number of validator indices to send in each request.
func (s *Service) indexChunkSize(ctx context.Context) int {
	if s.userIndexChunkSize > 0 {
		return s.userIndexChunkSize
	}

	// If this errors it will use the default so not a concern.
	nodeVersion, _ := s.NodeVersion(ctx)

	nodeVersion = strings.ToLower(nodeVersion)

	switch {
	case strings.Contains(nodeVersion, "lighthouse"):
		return indexChunkSizes["lighthouse"]
	case strings.Contains(nodeVersion, "nimbus"):
		return indexChunkSizes["nimbus"]
	case strings.Contains(nodeVersion, "prysm"):
		return indexChunkSizes["prysm"]
	case strings.Contains(nodeVersion, "teku"):
		return indexChunkSizes["teku"]
	default:
		return indexChunkSizes["default"]
	}
}

// Validators provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorIndices is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
// validatorStates is a list of validator states to restrict the returned values.  If no states are supplied no filter will be applied.
func (s *Service) Validators(ctx context.Context, stateID string, validatorIndices []phase0.ValidatorIndex, validatorStates []v1.ValidatorState) (map[phase0.ValidatorIndex]*api.Validator, error) {
	if stateID == "" {
		return nil, errors.New("no state ID specified")
	}

	if len(validatorIndices) == 0 {
		return s.validatorsFromState(ctx, stateID)
	}

	if len(validatorIndices) > s.indexChunkSize(ctx) {
		return s.chunkedValidators(ctx, stateID, validatorIndices, validatorStates)
	}

	url := fmt.Sprintf("/eth/v1/beacon/states/%s/validators", stateID)
	if len(validatorIndices) != 0 {
		ids := make([]string, len(validatorIndices))
		for i := range validatorIndices {
			ids[i] = fmt.Sprintf("%d", validatorIndices[i])
		}
		url = fmt.Sprintf("%s?id=%s", url, strings.Join(ids, ","))
	}
	if len(validatorStates) != 0 {
		states := make([]string, len(validatorStates))
		for i := range validatorStates {
			states[i] = fmt.Sprintf("%s", validatorStates[i])
		}
		if len(validatorIndices) != 0 {
			url = fmt.Sprintf("%s&status=%s", url, strings.Join(states, ","))
		} else {
			url = fmt.Sprintf("%s?status=%s", url, strings.Join(states, ","))
		}
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

	res := make(map[phase0.ValidatorIndex]*api.Validator)
	for _, validator := range validatorsJSON.Data {
		res[validator.Index] = validator
	}
	return res, nil
}

// validatorsFromState fetches all validators from state.
// This is more efficient than fetching the validators endpoint, as validators uses JSON only,
// whereas state can be provided using SSZ.
func (s *Service) validatorsFromState(ctx context.Context, stateID string) (map[phase0.ValidatorIndex]*api.Validator, error) {
	state, err := s.BeaconState(ctx, stateID)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return nil, errors.New("no beacon state")
	}

	validators, err := state.Validators()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain validators from state")
	}

	balances, err := state.ValidatorBalances()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain validator balances from state")
	}

	slot, err := state.Slot()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain slot from state")
	}
	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return nil, err
	}
	epoch := phase0.Epoch(uint64(slot) / slotsPerEpoch)

	farFutureEpoch, err := s.FarFutureEpoch(ctx)
	if err != nil {
		return nil, err
	}

	res := make(map[phase0.ValidatorIndex]*api.Validator, len(validators))
	for i, validator := range validators {
		index := phase0.ValidatorIndex(i)
		state := api.ValidatorToState(validator, &balances[i], epoch, farFutureEpoch)
		res[index] = &api.Validator{
			Index:     index,
			Balance:   balances[i],
			Status:    state,
			Validator: validator,
		}
	}

	return res, nil
}

// chunkedValidators obtains the validators a chunk at a time.
func (s *Service) chunkedValidators(ctx context.Context, stateID string, validatorIndices []phase0.ValidatorIndex, validatorStates []v1.ValidatorState) (map[phase0.ValidatorIndex]*api.Validator, error) {
	res := make(map[phase0.ValidatorIndex]*api.Validator)
	indexChunkSize := s.indexChunkSize(ctx)
	for i := 0; i < len(validatorIndices); i += indexChunkSize {
		chunkStart := i
		chunkEnd := i + indexChunkSize
		if len(validatorIndices) < chunkEnd {
			chunkEnd = len(validatorIndices)
		}
		chunk := validatorIndices[chunkStart:chunkEnd]
		chunkRes, err := s.Validators(ctx, stateID, chunk, validatorStates)
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain chunk")
		}
		for k, v := range chunkRes {
			res[k] = v
		}
	}
	return res, nil
}
