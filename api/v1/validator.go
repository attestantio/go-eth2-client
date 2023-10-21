// Copyright © 2020, 2021 Attestant Limited.
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
	"strconv"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// Validator contains the spec validator plus additional fields.
type Validator struct {
	Index     phase0.ValidatorIndex
	Balance   phase0.Gwei
	Status    ValidatorState
	Validator *phase0.Validator
}

// validatorJSON is the spec representation of the struct.
type validatorJSON struct {
	Index     string            `json:"index"`
	Balance   string            `json:"balance"`
	Status    ValidatorState    `json:"status"`
	Validator *phase0.Validator `json:"validator"`
}

// MarshalJSON implements json.Marshaler.
func (v *Validator) MarshalJSON() ([]byte, error) {
	return json.Marshal(&validatorJSON{
		Index:     fmt.Sprintf("%d", v.Index),
		Balance:   fmt.Sprintf("%d", v.Balance),
		Status:    v.Status,
		Validator: v.Validator,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *Validator) UnmarshalJSON(input []byte) error {
	var err error

	var validatorJSON validatorJSON
	if err = json.Unmarshal(input, &validatorJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if validatorJSON.Index == "" {
		return errors.New("index missing")
	}
	index, err := strconv.ParseUint(validatorJSON.Index, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for index")
	}
	v.Index = phase0.ValidatorIndex(index)
	if validatorJSON.Balance == "" {
		return errors.New("balance missing")
	}
	balance, err := strconv.ParseUint(validatorJSON.Balance, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for balance")
	}
	v.Balance = phase0.Gwei(balance)
	v.Status = validatorJSON.Status
	if validatorJSON.Validator == nil {
		return errors.New("validator missing")
	}
	v.Validator = validatorJSON.Validator

	return nil
}

// String returns a string version of the structure.
func (v *Validator) String() string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

// PubKey implements ValidatorPubKeyProvider.
func (v *Validator) PubKey(_ context.Context) (phase0.BLSPubKey, error) {
	return v.Validator.PublicKey, nil
}
