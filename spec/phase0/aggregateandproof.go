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

// AggregateAndProof is the Ethereum 2 attestation structure.
type AggregateAndProof struct {
	AggregatorIndex ValidatorIndex
	Aggregate       *Attestation
	SelectionProof  BLSSignature `ssz-size:"96"`
}

// aggregateAndProofJSON is the spec representation of the struct.
type aggregateAndProofJSON struct {
	AggregatorIndex string       `json:"aggregator_index"`
	Aggregate       *Attestation `json:"aggregate"`
	SelectionProof  string       `json:"selection_proof"`
}

// aggregateAndProofYAML is the spec representation of the struct.
type aggregateAndProofYAML struct {
	AggregatorIndex uint64       `yaml:"aggregator_index"`
	Aggregate       *Attestation `yaml:"aggregate"`
	SelectionProof  string       `yaml:"selection_proof"`
}

// MarshalJSON implements json.Marshaler.
func (a *AggregateAndProof) MarshalJSON() ([]byte, error) {
	return json.Marshal(&aggregateAndProofJSON{
		AggregatorIndex: fmt.Sprintf("%d", a.AggregatorIndex),
		Aggregate:       a.Aggregate,
		SelectionProof:  fmt.Sprintf("%#x", a.SelectionProof),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AggregateAndProof) UnmarshalJSON(input []byte) error {
	var aggregateAndProofJSON aggregateAndProofJSON
	if err := json.Unmarshal(input, &aggregateAndProofJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return a.unpack(&aggregateAndProofJSON)
}

func (a *AggregateAndProof) unpack(aggregateAndProofJSON *aggregateAndProofJSON) error {
	if aggregateAndProofJSON.AggregatorIndex == "" {
		return errors.New("aggregator index missing")
	}
	aggregatorIndex, err := strconv.ParseUint(aggregateAndProofJSON.AggregatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for aggregator index")
	}
	a.AggregatorIndex = ValidatorIndex(aggregatorIndex)
	if aggregateAndProofJSON.Aggregate == nil {
		return errors.New("aggregate missing")
	}
	a.Aggregate = aggregateAndProofJSON.Aggregate
	if aggregateAndProofJSON.SelectionProof == "" {
		return errors.New("selection proof missing")
	}
	selectionProof, err := hex.DecodeString(strings.TrimPrefix(aggregateAndProofJSON.SelectionProof, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for selection proof")
	}
	if len(selectionProof) != SignatureLength {
		return errors.New("incorrect length for selection proof")
	}
	copy(a.SelectionProof[:], selectionProof)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (a *AggregateAndProof) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&aggregateAndProofYAML{
		AggregatorIndex: uint64(a.AggregatorIndex),
		Aggregate:       a.Aggregate,
		SelectionProof:  fmt.Sprintf("%#x", a.SelectionProof),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (a *AggregateAndProof) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var aggregateAndProofJSON aggregateAndProofJSON
	if err := yaml.Unmarshal(input, &aggregateAndProofJSON); err != nil {
		return err
	}

	return a.unpack(&aggregateAndProofJSON)
}

// String returns a string version of the structure.
func (a *AggregateAndProof) String() string {
	data, err := yaml.Marshal(a)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
