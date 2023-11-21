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

package phase0

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// Validator is the Ethereum 2 validator structure.
type Validator struct {
	PublicKey                  BLSPubKey `ssz-size:"48"`
	WithdrawalCredentials      []byte    `ssz-size:"32"`
	EffectiveBalance           Gwei
	Slashed                    bool
	ActivationEligibilityEpoch Epoch
	ActivationEpoch            Epoch
	ExitEpoch                  Epoch
	WithdrawableEpoch          Epoch
}

// validatorJSON is the spec representation of the struct.
type validatorJSON struct {
	PublicKey                  string `json:"pubkey"`
	WithdrawalCredentials      string `json:"withdrawal_credentials"`
	EffectiveBalance           string `json:"effective_balance"`
	Slashed                    bool   `json:"slashed"`
	ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
	ActivationEpoch            string `json:"activation_epoch"`
	ExitEpoch                  string `json:"exit_epoch"`
	WithdrawableEpoch          string `json:"withdrawable_epoch"`
}

// validatorYAML is the spec representation of the struct.
type validatorYAML struct {
	PublicKey                  string `yaml:"pubkey"`
	WithdrawalCredentials      string `yaml:"withdrawal_credentials"`
	EffectiveBalance           uint64 `yaml:"effective_balance"`
	Slashed                    bool   `yaml:"slashed"`
	ActivationEligibilityEpoch uint64 `yaml:"activation_eligibility_epoch"`
	ActivationEpoch            uint64 `yaml:"activation_epoch"`
	ExitEpoch                  uint64 `yaml:"exit_epoch"`
	WithdrawableEpoch          uint64 `yaml:"withdrawable_epoch"`
}

// MarshalJSON implements json.Marshaler.
func (v *Validator) MarshalJSON() ([]byte, error) {
	return json.Marshal(&validatorJSON{
		PublicKey:                  fmt.Sprintf("%#x", v.PublicKey),
		WithdrawalCredentials:      fmt.Sprintf("%#x", v.WithdrawalCredentials),
		EffectiveBalance:           fmt.Sprintf("%d", v.EffectiveBalance),
		Slashed:                    v.Slashed,
		ActivationEligibilityEpoch: fmt.Sprintf("%d", v.ActivationEligibilityEpoch),
		ActivationEpoch:            fmt.Sprintf("%d", v.ActivationEpoch),
		ExitEpoch:                  fmt.Sprintf("%d", v.ExitEpoch),
		WithdrawableEpoch:          fmt.Sprintf("%d", v.WithdrawableEpoch),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *Validator) UnmarshalJSON(input []byte) error {
	var validatorJSON validatorJSON
	if err := json.Unmarshal(input, &validatorJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return v.unpack(&validatorJSON)
}

func (v *Validator) unpack(validatorJSON *validatorJSON) error {
	if validatorJSON.PublicKey == "" {
		return errors.New("public key missing")
	}
	publicKey, err := hex.DecodeString(strings.TrimPrefix(validatorJSON.PublicKey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}
	if len(publicKey) != PublicKeyLength {
		return fmt.Errorf("incorrect length %d for public key", len(publicKey))
	}
	copy(v.PublicKey[:], publicKey)
	if validatorJSON.WithdrawalCredentials == "" {
		return errors.New("withdrawal credentials missing")
	}
	v.WithdrawalCredentials, err = hex.DecodeString(strings.TrimPrefix(validatorJSON.WithdrawalCredentials, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for withdrawal credentials")
	}
	if len(v.WithdrawalCredentials) != HashLength {
		return fmt.Errorf("incorrect length %d for withdrawal credentials", len(v.WithdrawalCredentials))
	}
	if validatorJSON.EffectiveBalance == "" {
		return errors.New("effective balance missing")
	}
	effectiveBalance, err := strconv.ParseUint(validatorJSON.EffectiveBalance, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for effective balance")
	}
	v.EffectiveBalance = Gwei(effectiveBalance)
	v.Slashed = validatorJSON.Slashed
	if validatorJSON.ActivationEligibilityEpoch == "" {
		return errors.New("activation eligibility epoch missing")
	}
	activationEligibilityEpoch, err := strconv.ParseUint(validatorJSON.ActivationEligibilityEpoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for activation eligibility epoch")
	}
	v.ActivationEligibilityEpoch = Epoch(activationEligibilityEpoch)
	if validatorJSON.ActivationEpoch == "" {
		return errors.New("activation epoch missing")
	}
	activationEpoch, err := strconv.ParseUint(validatorJSON.ActivationEpoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for activation epoch")
	}
	v.ActivationEpoch = Epoch(activationEpoch)
	if validatorJSON.ExitEpoch == "" {
		return errors.New("exit epoch missing")
	}
	exitEpoch, err := strconv.ParseUint(validatorJSON.ExitEpoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for exit epoch")
	}
	v.ExitEpoch = Epoch(exitEpoch)
	if validatorJSON.WithdrawableEpoch == "" {
		return errors.New("withdrawable epoch missing")
	}
	withdrawableEpoch, err := strconv.ParseUint(validatorJSON.WithdrawableEpoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for withdrawable epoch")
	}
	v.WithdrawableEpoch = Epoch(withdrawableEpoch)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (v *Validator) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&validatorYAML{
		PublicKey:                  fmt.Sprintf("%#x", v.PublicKey),
		WithdrawalCredentials:      fmt.Sprintf("%#x", v.WithdrawalCredentials),
		EffectiveBalance:           uint64(v.EffectiveBalance),
		Slashed:                    v.Slashed,
		ActivationEligibilityEpoch: uint64(v.ActivationEligibilityEpoch),
		ActivationEpoch:            uint64(v.ActivationEpoch),
		ExitEpoch:                  uint64(v.ExitEpoch),
		WithdrawableEpoch:          uint64(v.WithdrawableEpoch),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (v *Validator) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var validatorJSON validatorJSON
	if err := yaml.Unmarshal(input, &validatorJSON); err != nil {
		return err
	}

	return v.unpack(&validatorJSON)
}

// String returns a string version of the structure.
func (v *Validator) String() string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
