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

package v1

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// ValidatorRegistration represents a ValidatorRegistrationV1.
type ValidatorRegistration struct {
	FeeRecipient bellatrix.ExecutionAddress `ssz-size:"20"`
	GasLimit     uint64
	Timestamp    time.Time
	Pubkey       phase0.BLSPubKey `ssz-size:"48"`
}

// validatorRegistrationJSON is the spec representation of the struct.
type validatorRegistrationJSON struct {
	FeeRecipient string `json:"fee_recipient"`
	GasLimit     string `json:"gas_limit"`
	Timestamp    string `json:"timestamp"`
	Pubkey       string `json:"pubkey"`
}

// validatorRegistrationYAML is the spec representation of the struct.
type validatorRegistrationYAML struct {
	FeeRecipient string `yaml:"fee_recipient"`
	GasLimit     uint64 `yaml:"gas_limit"`
	Timestamp    uint64 `yaml:"timestamp"`
	Pubkey       string `yaml:"pubkey"`
}

// MarshalJSON implements json.Marshaler.
func (v *ValidatorRegistration) MarshalJSON() ([]byte, error) {
	return json.Marshal(&validatorRegistrationJSON{
		FeeRecipient: v.FeeRecipient.String(),
		GasLimit:     strconv.FormatUint(v.GasLimit, 10),
		Timestamp:    strconv.FormatInt(v.Timestamp.Unix(), 10),
		Pubkey:       fmt.Sprintf("%#x", v.Pubkey),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *ValidatorRegistration) UnmarshalJSON(input []byte) error {
	var data validatorRegistrationJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return v.unpack(&data)
}

func (v *ValidatorRegistration) unpack(data *validatorRegistrationJSON) error {
	if data.FeeRecipient == "" {
		return errors.New("fee recipient missing")
	}
	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.FeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for fee recipient")
	}
	copy(v.FeeRecipient[:], feeRecipient)

	if data.GasLimit == "" {
		return errors.New("gas limit missing")
	}
	gasLimit, err := strconv.ParseUint(data.GasLimit, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for gas limit")
	}
	v.GasLimit = gasLimit

	if data.Timestamp == "" {
		return errors.New("timestamp missing")
	}
	timestamp, err := strconv.ParseInt(data.Timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for timestamp")
	}
	v.Timestamp = time.Unix(timestamp, 0)

	if data.Pubkey == "" {
		return errors.New("public key missing")
	}
	pubKey, err := hex.DecodeString(strings.TrimPrefix(data.Pubkey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}
	if len(pubKey) != publicKeyLength {
		return errors.New("incorrect length for public key")
	}
	copy(v.Pubkey[:], pubKey)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (v *ValidatorRegistration) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&validatorRegistrationYAML{
		FeeRecipient: v.FeeRecipient.String(),
		GasLimit:     v.GasLimit,
		Timestamp:    uint64(v.Timestamp.Unix()),
		Pubkey:       fmt.Sprintf("%#x", v.Pubkey),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (v *ValidatorRegistration) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data validatorRegistrationJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return v.unpack(&data)
}

// String returns a string version of the structure.
func (v *ValidatorRegistration) String() string {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
