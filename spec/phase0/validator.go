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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Validator is the Ethereum 2 validator structure.
type Validator struct {
	PublicKey                  []byte `ssz-size:"48"`
	WithdrawalCredentials      []byte `ssz-size:"32"`
	EffectiveBalance           uint64
	Slashed                    bool
	ActivationEligibilityEpoch uint64
	ActivationEpoch            uint64
	ExitEpoch                  uint64
	WithdrawableEpoch          uint64
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
	var err error

	var validatorJSON validatorJSON
	if err = json.Unmarshal(input, &validatorJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if validatorJSON.PublicKey == "" {
		return errors.New("public key missing")
	}
	if v.PublicKey, err = hex.DecodeString(strings.TrimPrefix(validatorJSON.PublicKey, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}
	if len(v.PublicKey) != publicKeyLength {
		return fmt.Errorf("incorrect length %d for public key", len(v.PublicKey))
	}

	if validatorJSON.WithdrawalCredentials == "" {
		return errors.New("withdrawal credentials missing")
	}
	if v.WithdrawalCredentials, err = hex.DecodeString(strings.TrimPrefix(validatorJSON.WithdrawalCredentials, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for withdrawal credentials")
	}
	if len(v.WithdrawalCredentials) != hashLength {
		return fmt.Errorf("incorrect length %d for withdrawal credentials", len(v.WithdrawalCredentials))
	}
	if validatorJSON.EffectiveBalance == "" {
		return errors.New("effective balance missing")
	}
	if v.EffectiveBalance, err = strconv.ParseUint(validatorJSON.EffectiveBalance, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for effective balance")
	}
	v.Slashed = validatorJSON.Slashed
	if validatorJSON.ActivationEligibilityEpoch == "" {
		return errors.New("activation eligibility epoch missing")
	}
	if v.ActivationEligibilityEpoch, err = strconv.ParseUint(validatorJSON.ActivationEligibilityEpoch, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for activation eligibility epoch")
	}
	if validatorJSON.ActivationEpoch == "" {
		return errors.New("activation epoch missing")
	}
	if v.ActivationEpoch, err = strconv.ParseUint(validatorJSON.ActivationEpoch, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for activation epoch")
	}
	if validatorJSON.ExitEpoch == "" {
		return errors.New("exit epoch missing")
	}
	if v.ExitEpoch, err = strconv.ParseUint(validatorJSON.ExitEpoch, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for exit epoch")
	}
	if validatorJSON.WithdrawableEpoch == "" {
		return errors.New("withdrawable epoch missing")
	}
	if v.WithdrawableEpoch, err = strconv.ParseUint(validatorJSON.WithdrawableEpoch, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for withdrawable epoch")
	}

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
