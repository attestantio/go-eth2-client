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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// SyncAggregatorSelectionDataSignature is the Ethereum 2 sync aggregator selection data structure.
type SyncAggregatorSelectionDataSignature struct {
	Slot              phase0.Slot
	SubcommitteeIndex uint64
}

// syncAggregatorSelectionDataJSON is the spec representation of the struct.
type syncAggregatorSelectionDataJSON struct {
	Slot              string `json:"slot"`
	SubcommitteeIndex string `json:"subcommittee_index"`
}

// syncAggregatorSelectionDataYAML is the spec representation of the struct.
type syncAggregatorSelectionDataYAML struct {
	Slot              uint64 `yaml:"slot"`
	SubcommitteeIndex uint64 `yaml:"subcommittee_index"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncAggregatorSelectionDataSignature) MarshalJSON() ([]byte, error) {
	return json.Marshal(&syncAggregatorSelectionDataJSON{
		Slot:              fmt.Sprintf("%d", s.Slot),
		SubcommitteeIndex: fmt.Sprintf("%d", s.SubcommitteeIndex),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncAggregatorSelectionDataSignature) UnmarshalJSON(input []byte) error {
	var syncAggregatorSelectionDataJSON syncAggregatorSelectionDataJSON
	if err := json.Unmarshal(input, &syncAggregatorSelectionDataJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	return s.unpack(&syncAggregatorSelectionDataJSON)
}

func (s *SyncAggregatorSelectionDataSignature) unpack(syncAggregatorSelectionDataJSON *syncAggregatorSelectionDataJSON) error {
	if syncAggregatorSelectionDataJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(syncAggregatorSelectionDataJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	s.Slot = phase0.Slot(slot)
	if syncAggregatorSelectionDataJSON.SubcommitteeIndex == "" {
		return errors.New("subcommittee index missing")
	}
	subcommitteeIndex, err := strconv.ParseUint(syncAggregatorSelectionDataJSON.SubcommitteeIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for subcommittee index")
	}
	s.SubcommitteeIndex = subcommitteeIndex

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (s *SyncAggregatorSelectionDataSignature) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&syncAggregatorSelectionDataYAML{
		Slot:              uint64(s.Slot),
		SubcommitteeIndex: s.SubcommitteeIndex,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}
	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SyncAggregatorSelectionDataSignature) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var syncAggregatorSelectionDataJSON syncAggregatorSelectionDataJSON
	if err := yaml.Unmarshal(input, &syncAggregatorSelectionDataJSON); err != nil {
		return err
	}
	return s.unpack(&syncAggregatorSelectionDataJSON)
}

// String returns a string version of the structure.
func (s *SyncAggregatorSelectionDataSignature) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
