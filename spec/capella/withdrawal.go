// Copyright © 2022 Attestant Limited.
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

package capella

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// Withdrawal provides information about a withdrawal.
type Withdrawal struct {
	Index          WithdrawalIndex
	ValidatorIndex phase0.ValidatorIndex
	Address        bellatrix.ExecutionAddress `ssz-size:"20"`
	Amount         phase0.Gwei
}

// withdrawalJSON is an internal representation of the struct.
type withdrawalJSON struct {
	Index          string `json:"index"`
	ValidatorIndex string `json:"validator_index"`
	Address        string `json:"address"`
	Amount         string `json:"amount"`
}

// withdrawalYAML is an internal representation of the struct.
type withdrawalYAML struct {
	Index          uint64 `yaml:"index"`
	ValidatorIndex uint64 `yaml:"validator_index"`
	Address        string `yaml:"address"`
	Amount         uint64 `yaml:"amount"`
}

// MarshalJSON implements json.Marshaler.
func (w *Withdrawal) MarshalJSON() ([]byte, error) {
	return json.Marshal(&withdrawalJSON{
		Index:          fmt.Sprintf("%d", w.Index),
		ValidatorIndex: fmt.Sprintf("%d", w.ValidatorIndex),
		Address:        fmt.Sprintf("%#x", w.Address),
		Amount:         fmt.Sprintf("%d", w.Amount),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (w *Withdrawal) UnmarshalJSON(input []byte) error {
	var data withdrawalJSON
	err := json.Unmarshal(input, &data)
	if err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return w.unpack(&data)
}

func (w *Withdrawal) unpack(data *withdrawalJSON) error {
	if data.Index == "" {
		return errors.New("index missing")
	}
	index, err := strconv.ParseUint(data.Index, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for index")
	}
	w.Index = WithdrawalIndex(index)

	if data.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(data.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	w.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	if data.Address == "" {
		return errors.New("address missing")
	}
	address, err := hex.DecodeString(strings.TrimPrefix(data.Address, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for address")
	}
	if len(address) != bellatrix.ExecutionAddressLength {
		return errors.New("incorrect length for address")
	}
	copy(w.Address[:], address)

	if data.Amount == "" {
		return errors.New("amount missing")
	}
	amount, err := strconv.ParseUint(data.Amount, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for amount")
	}
	w.Amount = phase0.Gwei(amount)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (w *Withdrawal) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&withdrawalYAML{
		Index:          uint64(w.Index),
		ValidatorIndex: uint64(w.ValidatorIndex),
		Address:        fmt.Sprintf("%#x", w.Address),
		Amount:         uint64(w.Amount),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (w *Withdrawal) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data withdrawalJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return w.unpack(&data)
}

// String returns a string version of the structure.
func (w *Withdrawal) String() string {
	data, err := yaml.Marshal(w)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
