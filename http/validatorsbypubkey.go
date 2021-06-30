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

package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type validatorsByPubKeyJSON struct {
	Data []*api.Validator `json:"data"`
}

// pubKeyChunkSize is the maximum number of public keys to send in each request.
// A request should be no more than 8,000 bytes to work with all currently-supported clients.
// A public key, including 0x header and comma separator, takes up 99 bytes.
// We also need to reserve space for the state ID and the endpoint itself, to be safe we go
// with 500 bytes for this which results in us having space for 75 public keys.
var pubKeyChunkSize = 75

// ValidatorsByPubKey provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validatorPubKeys is a list of validator public keys to restrict the returned values.  If no validators public keys are
// supplied no filter will be applied.
func (s *Service) ValidatorsByPubKey(ctx context.Context, stateID string, validatorPubKeys []spec.BLSPubKey) (map[spec.ValidatorIndex]*api.Validator, error) {
	if stateID == "" {
		return nil, errors.New("no state ID specified")
	}

	if len(validatorPubKeys) > pubKeyChunkSize {
		return s.chunkedValidatorsByPubKey(ctx, stateID, validatorPubKeys)
	}

	url := fmt.Sprintf("/eth/v1/beacon/states/%s/validators", stateID)
	if len(validatorPubKeys) != 0 {
		ids := make([]string, len(validatorPubKeys))
		for i := range validatorPubKeys {
			ids[i] = fmt.Sprintf("%#x", validatorPubKeys[i])
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

	var validatorsByPubKeyJSON validatorsByPubKeyJSON
	if err := json.NewDecoder(respBodyReader).Decode(&validatorsByPubKeyJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse validators")
	}
	if validatorsByPubKeyJSON.Data == nil {
		return nil, errors.New("no validators returned")
	}

	res := make(map[spec.ValidatorIndex]*api.Validator)
	for _, validator := range validatorsByPubKeyJSON.Data {
		res[validator.Index] = validator
	}
	return res, nil
}

// chunkedValidatorsByPubKey obtains the validators a chunk at a time.
func (s *Service) chunkedValidatorsByPubKey(ctx context.Context, stateID string, validatorPubKeys []spec.BLSPubKey) (map[spec.ValidatorIndex]*api.Validator, error) {
	res := make(map[spec.ValidatorIndex]*api.Validator)
	for i := 0; i < len(validatorPubKeys); i += pubKeyChunkSize {
		chunkStart := i
		chunkEnd := i + pubKeyChunkSize
		if len(validatorPubKeys) < chunkEnd {
			chunkEnd = len(validatorPubKeys)
		}
		chunk := validatorPubKeys[chunkStart:chunkEnd]
		chunkRes, err := s.ValidatorsByPubKey(ctx, stateID, chunk)
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain chunk")
		}
		for k, v := range chunkRes {
			res[k] = v
		}
	}
	return res, nil
}
