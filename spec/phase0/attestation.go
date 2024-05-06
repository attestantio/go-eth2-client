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
	bitfield "github.com/prysmaticlabs/go-bitfield"
)

// Attestation is the Ethereum 2 attestation structure.
type Attestation struct {
	AggregationBits bitfield.Bitlist `ssz-max:"2048"`
	Data            *AttestationData
	Signature       BLSSignature `ssz-size:"96"`
}

// attestationJSON is a raw representation of the struct.
type attestationJSON struct {
	AggregationBits string           `json:"aggregation_bits"`
	Data            *AttestationData `json:"data"`
	Signature       string           `json:"signature"`
}

// attestationYAML is a raw representation of the struct.
type attestationYAML struct {
	AggregationBits string           `yaml:"aggregation_bits"`
	Data            *AttestationData `yaml:"data"`
	Signature       string           `yaml:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (a *Attestation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&attestationJSON{
		AggregationBits: fmt.Sprintf("%#x", []byte(a.AggregationBits)),
		Data:            a.Data,
		Signature:       fmt.Sprintf("%#x", a.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Attestation) UnmarshalJSON(input []byte) error {
	var attestationJSON attestationJSON
	err := json.Unmarshal(input, &attestationJSON)
	if err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return a.unpack(&attestationJSON)
}

func (a *Attestation) unpack(attestationJSON *attestationJSON) error {
	var err error
	if attestationJSON.AggregationBits == "" {
		return errors.New("aggregation bits missing")
	}
	if a.AggregationBits, err = hex.DecodeString(strings.TrimPrefix(attestationJSON.AggregationBits, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for beacon block root")
	}
	a.Data = attestationJSON.Data
	if a.Data == nil {
		return errors.New("data missing")
	}
	if attestationJSON.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(attestationJSON.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(signature) != SignatureLength {
		return errors.New("incorrect length for signature")
	}
	copy(a.Signature[:], signature)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (a *Attestation) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&attestationYAML{
		AggregationBits: fmt.Sprintf("%#x", []byte(a.AggregationBits)),
		Data:            a.Data,
		Signature:       fmt.Sprintf("%#x", a.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (a *Attestation) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var attestationJSON attestationJSON
	if err := yaml.Unmarshal(input, &attestationJSON); err != nil {
		return err
	}

	return a.unpack(&attestationJSON)
}

// String returns a string version of the structure.
func (a *Attestation) String() string {
	data, err := yaml.Marshal(a)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
