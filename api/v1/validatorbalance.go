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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// ValidatorBalance contains the balance of a validator.
type ValidatorBalance struct {
	Index   phase0.ValidatorIndex
	Balance phase0.Gwei
}

// validatorBalanceJSON is the spec representation of the struct.
type validatorBalanceJSON struct {
	Index   string `json:"index"`
	Balance string `json:"balance"`
}

// MarshalJSON implements json.Marshaler.
func (v *ValidatorBalance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&validatorBalanceJSON{
		Index:   fmt.Sprintf("%d", v.Index),
		Balance: fmt.Sprintf("%d", v.Balance),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *ValidatorBalance) UnmarshalJSON(input []byte) error {
	var err error

	var validatorBalanceJSON validatorBalanceJSON
	if err = json.Unmarshal(input, &validatorBalanceJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if validatorBalanceJSON.Index == "" {
		return errors.New("index missing")
	}
	index, err := strconv.ParseUint(validatorBalanceJSON.Index, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for index")
	}
	v.Index = phase0.ValidatorIndex(index)
	if validatorBalanceJSON.Balance == "" {
		return errors.New("balance missing")
	}
	balance, err := strconv.ParseUint(validatorBalanceJSON.Balance, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for balance")
	}
	v.Balance = phase0.Gwei(balance)

	return nil
}

// String returns a string version of the structure.
func (v *ValidatorBalance) String() string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
