// Copyright © 2021 Attestant Limited.
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

package altair

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// ContributionAndProof is the Ethereum 2 contribution and proof structure.
type ContributionAndProof struct {
	AggregatorIndex phase0.ValidatorIndex
	Contribution    *SyncCommitteeContribution
	SelectionProof  phase0.BLSSignature `ssz-size:"96"`
}

// contributionAndProofJSON is the spec representation of the struct.
type contributionAndProofJSON struct {
	AggregatorIndex string                     `json:"aggregator_index"`
	Contribution    *SyncCommitteeContribution `json:"contribution"`
	SelectionProof  string                     `json:"selection_proof"`
}

// contributionAndProofYAML is the spec representation of the struct.
type contributionAndProofYAML struct {
	AggregatorIndex uint64                     `yaml:"aggregator_index"`
	Contribution    *SyncCommitteeContribution `yaml:"contribution"`
	SelectionProof  string                     `yaml:"selection_proof"`
}

// MarshalJSON implements json.Marshaler.
func (a *ContributionAndProof) MarshalJSON() ([]byte, error) {
	return json.Marshal(&contributionAndProofJSON{
		AggregatorIndex: fmt.Sprintf("%d", a.AggregatorIndex),
		Contribution:    a.Contribution,
		SelectionProof:  fmt.Sprintf("%#x", a.SelectionProof),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *ContributionAndProof) UnmarshalJSON(input []byte) error {
	var contributionAndProofJSON contributionAndProofJSON
	if err := json.Unmarshal(input, &contributionAndProofJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return a.unpack(&contributionAndProofJSON)
}

func (a *ContributionAndProof) unpack(contributionAndProofJSON *contributionAndProofJSON) error {
	if contributionAndProofJSON.AggregatorIndex == "" {
		return errors.New("aggregator index missing")
	}
	aggregatorIndex, err := strconv.ParseUint(contributionAndProofJSON.AggregatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for aggregator index")
	}
	a.AggregatorIndex = phase0.ValidatorIndex(aggregatorIndex)
	if contributionAndProofJSON.Contribution == nil {
		return errors.New("contribution missing")
	}
	a.Contribution = contributionAndProofJSON.Contribution
	if contributionAndProofJSON.SelectionProof == "" {
		return errors.New("selection proof missing")
	}
	selectionProof, err := hex.DecodeString(strings.TrimPrefix(contributionAndProofJSON.SelectionProof, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for selection proof")
	}
	if len(selectionProof) != phase0.SignatureLength {
		return errors.New("incorrect length for selection proof")
	}
	copy(a.SelectionProof[:], selectionProof)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (a *ContributionAndProof) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&contributionAndProofYAML{
		AggregatorIndex: uint64(a.AggregatorIndex),
		Contribution:    a.Contribution,
		SelectionProof:  fmt.Sprintf("%#x", a.SelectionProof),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (a *ContributionAndProof) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var contributionAndProofJSON contributionAndProofJSON
	if err := yaml.Unmarshal(input, &contributionAndProofJSON); err != nil {
		return err
	}

	return a.unpack(&contributionAndProofJSON)
}

// String returns a string version of the structure.
func (a *ContributionAndProof) String() string {
	data, err := yaml.Marshal(a)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
