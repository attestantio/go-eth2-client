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

package gloas

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

// Builder is the Ethereum 2 builder structure.
type Builder struct {
	PublicKey         phase0.BLSPubKey `ssz-size:"48"`
	Version           uint8
	ExecutionAddress  bellatrix.ExecutionAddress `ssz-size:"20"`
	Balance           phase0.Gwei
	DepositEpoch      phase0.Epoch
	WithdrawableEpoch phase0.Epoch
}

// builderJSON is the spec representation of the struct.
type builderJSON struct {
	PublicKey         string `json:"pubkey"`
	Version           uint8  `json:"version"`
	ExecutionAddress  string `json:"execution_address"`
	Balance           string `json:"balance"`
	DepositEpoch      string `json:"deposit_epoch"`
	WithdrawableEpoch string `json:"withdrawable_epoch"`
}

// builderYAML is the spec representation of the struct.
type builderYAML struct {
	PublicKey         string `yaml:"pubkey"`
	Version           uint8  `yaml:"version"`
	ExecutionAddress  string `yaml:"execution_address"`
	Balance           string `yaml:"balance"`
	DepositEpoch      string `yaml:"deposit_epoch"`
	WithdrawableEpoch string `yaml:"withdrawable_epoch"`
}

// MarshalJSON implements json.Marshaler.
func (v *Builder) MarshalJSON() ([]byte, error) {
	return json.Marshal(&builderJSON{
		PublicKey:         fmt.Sprintf("%#x", v.PublicKey),
		Version:           v.Version,
		ExecutionAddress:  v.ExecutionAddress.String(),
		Balance:           fmt.Sprintf("%d", v.Balance),
		DepositEpoch:      fmt.Sprintf("%d", v.DepositEpoch),
		WithdrawableEpoch: fmt.Sprintf("%d", v.WithdrawableEpoch),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *Builder) UnmarshalJSON(input []byte) error {
	var builderJSON builderJSON
	if err := json.Unmarshal(input, &builderJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return v.unpack(&builderJSON)
}

func (v *Builder) unpack(builderJSON *builderJSON) error {
	if builderJSON.PublicKey == "" {
		return errors.New("public key missing")
	}

	publicKey, err := hex.DecodeString(strings.TrimPrefix(builderJSON.PublicKey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}

	if len(publicKey) != phase0.PublicKeyLength {
		return fmt.Errorf("incorrect length %d for public key", len(publicKey))
	}

	copy(v.PublicKey[:], publicKey)

	if builderJSON.ExecutionAddress == "" {
		return errors.New("execution address missing")
	}

	executionAddress, err := hex.DecodeString(strings.TrimPrefix(builderJSON.ExecutionAddress, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for execution address")
	}

	if len(executionAddress) != bellatrix.ExecutionAddressLength {
		return errors.New("incorrect length for fee recipient")
	}

	copy(v.ExecutionAddress[:], executionAddress)

	if len(executionAddress) != bellatrix.ExecutionAddressLength {
		return fmt.Errorf("incorrect length %d for execution address", len(executionAddress))
	}

	if builderJSON.Balance == "" {
		return errors.New("balance missing")
	}

	balance, err := strconv.ParseUint(builderJSON.Balance, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for effective balance")
	}

	v.Balance = phase0.Gwei(balance)

	if builderJSON.DepositEpoch == "" {
		return errors.New("deposit epoch missing")
	}

	depositEpoch, err := strconv.ParseUint(builderJSON.DepositEpoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for deposit epoch")
	}

	v.DepositEpoch = phase0.Epoch(depositEpoch)

	if builderJSON.WithdrawableEpoch == "" {
		return errors.New("withdrawable epoch missing")
	}

	withdrawableEpoch, err := strconv.ParseUint(builderJSON.WithdrawableEpoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for withdrawable epoch")
	}

	v.WithdrawableEpoch = phase0.Epoch(withdrawableEpoch)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (v *Builder) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&builderYAML{
		PublicKey:         fmt.Sprintf("%#x", v.PublicKey),
		ExecutionAddress:  v.ExecutionAddress.String(),
		Balance:           fmt.Sprintf("%d", v.Balance),
		DepositEpoch:      fmt.Sprintf("%d", v.DepositEpoch),
		WithdrawableEpoch: fmt.Sprintf("%d", v.WithdrawableEpoch),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (v *Builder) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var builderJSON builderJSON
	if err := yaml.Unmarshal(input, &builderJSON); err != nil {
		return err
	}

	return v.unpack(&builderJSON)
}

// String returns a string version of the structure.
func (v *Builder) String() string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
