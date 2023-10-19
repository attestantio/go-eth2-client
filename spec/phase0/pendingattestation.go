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
	bitfield "github.com/prysmaticlabs/go-bitfield"
)

// PendingAttestation is the Ethereum 2 pending attestation structure.
type PendingAttestation struct {
	AggregationBits bitfield.Bitlist `ssz-max:"2048"`
	Data            *AttestationData
	InclusionDelay  Slot
	ProposerIndex   ValidatorIndex
}

// pendingAttestationJSON is the spec representation of the struct.
type pendingAttestationJSON struct {
	AggregationBits string           `json:"aggregation_bits"`
	Data            *AttestationData `json:"data"`
	InclusionDelay  string           `json:"inclusion_delay"`
	ProposerIndex   string           `json:"proposer_index"`
}

// pendingAttestationYAML is the spec representation of the struct.
type pendingAttestationYAML struct {
	AggregationBits string           `yaml:"aggregation_bits"`
	Data            *AttestationData `yaml:"data"`
	InclusionDelay  uint64           `yaml:"inclusion_delay"`
	ProposerIndex   uint64           `yaml:"proposer_index"`
}

// MarshalJSON implements json.Marshaler.
func (p *PendingAttestation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&pendingAttestationJSON{
		AggregationBits: fmt.Sprintf("%#x", []byte(p.AggregationBits)),
		Data:            p.Data,
		InclusionDelay:  fmt.Sprintf("%d", p.InclusionDelay),
		ProposerIndex:   fmt.Sprintf("%d", p.ProposerIndex),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PendingAttestation) UnmarshalJSON(input []byte) error {
	var pendingAttestationJSON pendingAttestationJSON
	if err := json.Unmarshal(input, &pendingAttestationJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&pendingAttestationJSON)
}

func (p *PendingAttestation) unpack(pendingAttestationJSON *pendingAttestationJSON) error {
	var err error
	if pendingAttestationJSON.AggregationBits == "" {
		return errors.New("aggregation bits missing")
	}
	if p.AggregationBits, err = hex.DecodeString(strings.TrimPrefix(pendingAttestationJSON.AggregationBits, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for aggregation bits")
	}
	p.Data = pendingAttestationJSON.Data
	if p.Data == nil {
		return errors.New("data missing")
	}
	if pendingAttestationJSON.InclusionDelay == "" {
		return errors.New("inclusion delay missing")
	}
	inclusionDelay, err := strconv.ParseUint(pendingAttestationJSON.InclusionDelay, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for inclusion delay")
	}
	p.InclusionDelay = Slot(inclusionDelay)
	if pendingAttestationJSON.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}
	proposerIndex, err := strconv.ParseUint(pendingAttestationJSON.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}
	p.ProposerIndex = ValidatorIndex(proposerIndex)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (p *PendingAttestation) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&pendingAttestationYAML{
		AggregationBits: fmt.Sprintf("%#x", []byte(p.AggregationBits)),
		Data:            p.Data,
		InclusionDelay:  uint64(p.InclusionDelay),
		ProposerIndex:   uint64(p.ProposerIndex),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (p *PendingAttestation) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var pendingAttestationJSON pendingAttestationJSON
	if err := yaml.Unmarshal(input, &pendingAttestationJSON); err != nil {
		return err
	}

	return p.unpack(&pendingAttestationJSON)
}

// String returns a string version of the structure.
func (p *PendingAttestation) String() string {
	data, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
