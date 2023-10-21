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
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// Deposit provides information about a deposit.
type Deposit struct {
	Proof [][]byte `ssz-size:"33,32"`
	Data  *DepositData
}

// depositJSON is the spec representation of the struct.
type depositJSON struct {
	Proof []string     `json:"proof"`
	Data  *DepositData `json:"data"`
}

// depositYAML is the spec representation of the struct.
type depositYAML struct {
	Proof []string     `yaml:"proof"`
	Data  *DepositData `yaml:"data"`
}

// MarshalJSON implements json.Marshaler.
func (d *Deposit) MarshalJSON() ([]byte, error) {
	proof := make([]string, len(d.Proof))
	for i := range d.Proof {
		proof[i] = fmt.Sprintf("%#x", d.Proof[i])
	}

	return json.Marshal(&depositJSON{
		Proof: proof,
		Data:  d.Data,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Deposit) UnmarshalJSON(input []byte) error {
	var depositJSON depositJSON
	if err := json.Unmarshal(input, &depositJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return d.unpack(&depositJSON)
}

func (d *Deposit) unpack(depositJSON *depositJSON) error {
	var err error
	if depositJSON.Proof == nil {
		return errors.New("proof missing")
	}
	if len(depositJSON.Proof) != 33 {
		return errors.New("incorrect length for proof")
	}
	d.Proof = make([][]byte, len(depositJSON.Proof))
	for i := range depositJSON.Proof {
		if depositJSON.Proof[i] == "" {
			return errors.New("proof component missing")
		}
		if d.Proof[i], err = hex.DecodeString(strings.TrimPrefix(depositJSON.Proof[i], "0x")); err != nil {
			return errors.Wrap(err, "invalid value for proof")
		}
		if len(d.Proof[i]) != 32 {
			return fmt.Errorf("incorrect size %d for deposit proof", len(d.Proof[i]))
		}
	}
	if depositJSON.Data == nil {
		return errors.New("data missing")
	}
	d.Data = depositJSON.Data

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (d *Deposit) MarshalYAML() ([]byte, error) {
	proof := make([]string, len(d.Proof))
	for i := range d.Proof {
		proof[i] = fmt.Sprintf("%#x", d.Proof[i])
	}
	yamlBytes, err := yaml.MarshalWithOptions(&depositYAML{
		Proof: proof,
		Data:  d.Data,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (d *Deposit) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var depositJSON depositJSON
	if err := yaml.Unmarshal(input, &depositJSON); err != nil {
		return err
	}

	return d.unpack(&depositJSON)
}

// String returns a string version of the structure.
func (d *Deposit) String() string {
	data, err := yaml.Marshal(d)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
