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

// AttestationData is the Ethereum 2 specification structure.
type AttestationData struct {
	Slot            Slot
	Index           CommitteeIndex
	BeaconBlockRoot Root `ssz-size:"32"`
	Source          *Checkpoint
	Target          *Checkpoint
}

// attestationDataJSON is an internal representation of the struct.
type attestationDataJSON struct {
	Slot            string      `json:"slot"`
	Index           string      `json:"index"`
	BeaconBlockRoot string      `json:"beacon_block_root"`
	Source          *Checkpoint `json:"source"`
	Target          *Checkpoint `json:"target"`
}

// attestationDataYAML is an internal representation of the struct.
type attestationDataYAML struct {
	Slot            uint64      `yaml:"slot"`
	Index           uint64      `yaml:"index"`
	BeaconBlockRoot string      `yaml:"beacon_block_root"`
	Source          *Checkpoint `json:"source"`
	Target          *Checkpoint `json:"target"`
}

// MarshalJSON implements json.Marshaler.
func (a *AttestationData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&attestationDataJSON{
		Slot:            fmt.Sprintf("%d", a.Slot),
		Index:           fmt.Sprintf("%d", a.Index),
		BeaconBlockRoot: fmt.Sprintf("%#x", a.BeaconBlockRoot),
		Source:          a.Source,
		Target:          a.Target,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AttestationData) UnmarshalJSON(input []byte) error {
	var attestationDataJSON attestationDataJSON
	if err := json.Unmarshal(input, &attestationDataJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return a.unpack(&attestationDataJSON)
}

func (a *AttestationData) unpack(attestationDataJSON *attestationDataJSON) error {
	if attestationDataJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(attestationDataJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	a.Slot = Slot(slot)
	if attestationDataJSON.Index == "" {
		return errors.New("index missing")
	}
	index, err := strconv.ParseUint(attestationDataJSON.Index, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for index")
	}
	a.Index = CommitteeIndex(index)
	if attestationDataJSON.BeaconBlockRoot == "" {
		return errors.New("beacon block root missing")
	}
	beaconBlockRoot, err := hex.DecodeString(strings.TrimPrefix(attestationDataJSON.BeaconBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for beacon block root")
	}
	if len(beaconBlockRoot) != RootLength {
		return errors.New("incorrect length for beacon block root")
	}
	copy(a.BeaconBlockRoot[:], beaconBlockRoot)
	if attestationDataJSON.Source == nil {
		return errors.New("source missing")
	}
	a.Source = attestationDataJSON.Source
	if attestationDataJSON.Target == nil {
		return errors.New("target missing")
	}
	a.Target = attestationDataJSON.Target

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (a *AttestationData) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&attestationDataYAML{
		Slot:            uint64(a.Slot),
		Index:           uint64(a.Index),
		BeaconBlockRoot: fmt.Sprintf("%#x", a.BeaconBlockRoot),
		Source:          a.Source,
		Target:          a.Target,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (a *AttestationData) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var attestationDataJSON attestationDataJSON
	if err := yaml.Unmarshal(input, &attestationDataJSON); err != nil {
		return err
	}

	return a.unpack(&attestationDataJSON)
}

// String provids a string representation of the struct.
func (a *AttestationData) String() string {
	data, err := yaml.Marshal(a)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
