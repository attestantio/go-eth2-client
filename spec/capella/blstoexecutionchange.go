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

// BLSToExecutionChange provides information about a change of withdrawal credentials.
type BLSToExecutionChange struct {
	ValidatorIndex     phase0.ValidatorIndex
	FromBLSPubkey      phase0.BLSPubKey           `ssz-size:"48"`
	ToExecutionAddress bellatrix.ExecutionAddress `ssz-size:"20"`
}

// blsToExecutionChangeJSON is an internal representation of the struct.
type blsToExecutionChangeJSON struct {
	ValidatorIndex     string `json:"validator_index"`
	FromBLSPubkey      string `json:"from_bls_pubkey"`
	ToExecutionAddress string `json:"to_execution_address"`
}

// blsToExecutionChangeYAML is an internal representation of the struct.
type blsToExecutionChangeYAML struct {
	ValidatorIndex     uint64 `yaml:"validator_index"`
	FromBLSPubkey      string `yaml:"from_bls_pubkey"`
	ToExecutionAddress string `yaml:"to_execution_address"`
}

// MarshalJSON implements json.Marshaler.
func (b *BLSToExecutionChange) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blsToExecutionChangeJSON{
		ValidatorIndex:     fmt.Sprintf("%d", b.ValidatorIndex),
		FromBLSPubkey:      fmt.Sprintf("%#x", b.FromBLSPubkey),
		ToExecutionAddress: b.ToExecutionAddress.String(),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BLSToExecutionChange) UnmarshalJSON(input []byte) error {
	var data blsToExecutionChangeJSON
	err := json.Unmarshal(input, &data)
	if err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return b.unpack(&data)
}

func (b *BLSToExecutionChange) unpack(data *blsToExecutionChangeJSON) error {
	if data.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(data.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	b.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	if data.FromBLSPubkey == "" {
		return errors.New("from BLS public key missing")
	}
	fromBLSPubkey, err := hex.DecodeString(strings.TrimPrefix(data.FromBLSPubkey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for from BLS public key")
	}
	if len(fromBLSPubkey) != phase0.PublicKeyLength {
		return errors.New("incorrect length for from BLS public key")
	}
	copy(b.FromBLSPubkey[:], fromBLSPubkey)

	if data.ToExecutionAddress == "" {
		return errors.New("to execution address missing")
	}
	toExecutionAddress, err := hex.DecodeString(strings.TrimPrefix(data.ToExecutionAddress, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for to execution address")
	}
	if len(toExecutionAddress) != bellatrix.ExecutionAddressLength {
		return errors.New("incorrect length for to execution address")
	}
	copy(b.ToExecutionAddress[:], toExecutionAddress)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (b *BLSToExecutionChange) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&blsToExecutionChangeYAML{
		ValidatorIndex:     uint64(b.ValidatorIndex),
		FromBLSPubkey:      fmt.Sprintf("%#x", b.FromBLSPubkey),
		ToExecutionAddress: b.ToExecutionAddress.String(),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BLSToExecutionChange) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data blsToExecutionChangeJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return b.unpack(&data)
}

// String returns a string version of the structure.
func (b *BLSToExecutionChange) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
