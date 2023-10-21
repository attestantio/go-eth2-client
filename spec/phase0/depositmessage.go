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

// DepositMessage provides information about a deposit made on Ethereum 1.
type DepositMessage struct {
	PublicKey             BLSPubKey `ssz-size:"48"`
	WithdrawalCredentials []byte    `ssz-size:"32"`
	Amount                Gwei
}

// depositMessageJSON is the spec representation of the struct.
type depositMessageJSON struct {
	PublicKey             string `json:"pubkey"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	Amount                string `json:"amount"`
}

// depositMessageYAML is the spec representation of the struct.
type depositMessageYAML struct {
	PublicKey             string `yaml:"pubkey"`
	WithdrawalCredentials string `yaml:"withdrawal_credentials"`
	Amount                uint64 `yaml:"amount"`
}

// MarshalJSON implements json.Marshaler.
func (d *DepositMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&depositMessageJSON{
		PublicKey:             fmt.Sprintf("%#x", d.PublicKey),
		WithdrawalCredentials: fmt.Sprintf("%#x", d.WithdrawalCredentials),
		Amount:                fmt.Sprintf("%d", d.Amount),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DepositMessage) UnmarshalJSON(input []byte) error {
	var depositMessageJSON depositMessageJSON
	if err := json.Unmarshal(input, &depositMessageJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return d.unpack(&depositMessageJSON)
}

func (d *DepositMessage) unpack(depositMessageJSON *depositMessageJSON) error {
	if depositMessageJSON.PublicKey == "" {
		return errors.New("public key missing")
	}
	publicKey, err := hex.DecodeString(strings.TrimPrefix(depositMessageJSON.PublicKey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}
	if len(publicKey) != PublicKeyLength {
		return errors.New("incorrect length for public key")
	}
	copy(d.PublicKey[:], publicKey)
	if depositMessageJSON.WithdrawalCredentials == "" {
		return errors.New("withdrawal credentials missing")
	}
	if d.WithdrawalCredentials, err = hex.DecodeString(strings.TrimPrefix(depositMessageJSON.WithdrawalCredentials, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for withdrawal credentials")
	}
	if len(d.WithdrawalCredentials) != HashLength {
		return errors.New("incorrect length for withdrawal credentials")
	}
	if depositMessageJSON.Amount == "" {
		return errors.New("amount missing")
	}
	amount, err := strconv.ParseUint(depositMessageJSON.Amount, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for amount")
	}
	d.Amount = Gwei(amount)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (d *DepositMessage) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&depositMessageYAML{
		PublicKey:             fmt.Sprintf("%#x", d.PublicKey),
		WithdrawalCredentials: fmt.Sprintf("%#x", d.WithdrawalCredentials),
		Amount:                uint64(d.Amount),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (d *DepositMessage) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var depositMessageJSON depositMessageJSON
	if err := yaml.Unmarshal(input, &depositMessageJSON); err != nil {
		return err
	}

	return d.unpack(&depositMessageJSON)
}

// String returns a string version of the structure.
func (d *DepositMessage) String() string {
	data, err := yaml.Marshal(d)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
