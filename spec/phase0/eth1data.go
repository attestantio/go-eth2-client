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

// ETH1Data provides information about the state of Ethereum 1 as viewed by the
// Ethereum 2 chain.
type ETH1Data struct {
	DepositRoot  Root `ssz-size:"32"`
	DepositCount uint64
	BlockHash    []byte `ssz-size:"32"`
}

// eth1DataJSON is the spec representation of the struct.
type eth1DataJSON struct {
	DepositRoot  string `json:"deposit_root"`
	DepositCount string `json:"deposit_count"`
	BlockHash    string `json:"block_hash"`
}

// eth1DataYAML is the spec representation of the struct.
type eth1DataYAML struct {
	DepositRoot  string `yaml:"deposit_root"`
	DepositCount uint64 `yaml:"deposit_count"`
	BlockHash    string `yaml:"block_hash"`
}

// MarshalJSON implements json.Marshaler.
func (e *ETH1Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(&eth1DataJSON{
		DepositRoot:  fmt.Sprintf("%#x", e.DepositRoot),
		DepositCount: strconv.FormatUint(e.DepositCount, 10),
		BlockHash:    fmt.Sprintf("%#x", e.BlockHash),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ETH1Data) UnmarshalJSON(input []byte) error {
	var eth1DataJSON eth1DataJSON
	if err := json.Unmarshal(input, &eth1DataJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return e.unpack(&eth1DataJSON)
}

func (e *ETH1Data) unpack(eth1DataJSON *eth1DataJSON) error {
	if eth1DataJSON.DepositRoot == "" {
		return errors.New("deposit root missing")
	}
	depositRoot, err := hex.DecodeString(strings.TrimPrefix(eth1DataJSON.DepositRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for deposit root")
	}
	if len(depositRoot) != RootLength {
		return errors.New("incorrect length for deposit root")
	}
	copy(e.DepositRoot[:], depositRoot)
	if eth1DataJSON.DepositCount == "" {
		return errors.New("deposit count missing")
	}
	if e.DepositCount, err = strconv.ParseUint(eth1DataJSON.DepositCount, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for deposit count")
	}
	if eth1DataJSON.BlockHash == "" {
		return errors.New("block hash missing")
	}
	if e.BlockHash, err = hex.DecodeString(strings.TrimPrefix(eth1DataJSON.BlockHash, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for block hash")
	}
	if len(e.BlockHash) != HashLength {
		return errors.New("incorrect length for block hash")
	}

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (e *ETH1Data) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&eth1DataYAML{
		DepositRoot:  fmt.Sprintf("%#x", e.DepositRoot),
		DepositCount: e.DepositCount,
		BlockHash:    fmt.Sprintf("%#x", e.BlockHash),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (e *ETH1Data) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var eth1DataJSON eth1DataJSON
	if err := yaml.Unmarshal(input, &eth1DataJSON); err != nil {
		return err
	}

	return e.unpack(&eth1DataJSON)
}

// String returns a string version of the structure.
func (e *ETH1Data) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
