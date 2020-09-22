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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// validatorsJSON handles the JSON returned from lighthouse.
type validatorsJSON struct {
	PubKey    string          `json:"pubkey"`
	Index     string          `json:"validator_index"`
	Balance   string          `json:"balance"`
	Validator *spec.Validator `json:"validator"`
}

// Validators provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validators is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
func (s *Service) Validators(ctx context.Context, stateID string, validatorIDs []client.ValidatorIDProvider) (map[uint64]*api.Validator, error) {
	stateRoot, err := s.StateRootFromStateID(ctx, stateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain state root")
	}

	// Lighthouse has different calls depending on if the caller wants all or a filtered set of validators.
	var respBodyReader io.ReadCloser
	if len(validatorIDs) == 0 {
		respBodyReader, err = s.get(ctx, fmt.Sprintf("/beacon/validators/all?state_root=%#x", stateRoot))
	} else {
		reqBodyReader := new(bytes.Buffer)
		if _, err := reqBodyReader.WriteString(`{"pubkeys":[`); err != nil {
			return nil, errors.Wrap(err, "failed to write pubkeys")
		}
		for i := range validatorIDs {
			pubKey, err := validatorIDs[i].PubKey(ctx)
			if err != nil {
				// Warn but continue.
				log.Warn().Err(err).Msg("Failed to obtain public key for validator; skipping")
				continue
			}
			if _, err := reqBodyReader.WriteString(fmt.Sprintf(`"%#x"`, pubKey)); err != nil {
				return nil, errors.Wrap(err, "failed to write pubkey")
			}
			if i != len(validatorIDs)-1 {
				if _, err := reqBodyReader.WriteString(`,`); err != nil {
					return nil, errors.Wrap(err, "failed to write separator")
				}
			}
		}
		if _, err := reqBodyReader.WriteString(`]}`); err != nil {
			return nil, errors.Wrap(err, "failed to write end of pubkeys")
		}
		respBodyReader, err = s.post(ctx, fmt.Sprintf("/beacon/validators?state_root=%#x", stateRoot), reqBodyReader)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to request validators")
	}

	specReader, err := lhToSpec(ctx, respBodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert lighthouse response to spec response")
	}

	var validators []*validatorsJSON
	if err := json.NewDecoder(specReader).Decode(&validators); err != nil {
		return nil, errors.Wrap(err, "failed to parse validators")
	}

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
	res := make(map[uint64]*api.Validator)
	for _, validator := range validators {
		if validator.Index == "" {
			// Validator does not have an index yet; ignore.
			continue
		}
		index, err := strconv.ParseUint(validator.Index, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse validator index")
		}
		res[index] = &api.Validator{
			Index:     index,
			State:     api.ValidatorToState(validator.Validator, epoch, farFutureEpoch),
			Validator: validator.Validator,
		}
		if validator.Balance != "" {
			balance, err := strconv.ParseUint(validator.Balance, 10, 64)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse validator balance")
			}
			res[index].Balance = balance
		}
	}

	return res, nil
}
