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

// BeaconBlockHeader represents the header of a beacon block without its content.
type BeaconBlockHeader struct {
	Slot          Slot
	ProposerIndex ValidatorIndex
	ParentRoot    Root `ssz-size:"32"`
	StateRoot     Root `ssz-size:"32"`
	BodyRoot      Root `ssz-size:"32"`
}

// beaconBlockHeaderJSON is a raw representation of the struct.
type beaconBlockHeaderJSON struct {
	Slot          string `json:"slot"`
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
}

// beaconBlockHeaderYAML is a raw representation of the struct.
type beaconBlockHeaderYAML struct {
	Slot          uint64 `yaml:"slot"`
	ProposerIndex uint64 `yaml:"proposer_index"`
	ParentRoot    string `yaml:"parent_root"`
	StateRoot     string `yaml:"state_root"`
	BodyRoot      string `yaml:"body_root"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconBlockHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(&beaconBlockHeaderJSON{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    fmt.Sprintf("%#x", b.ParentRoot),
		StateRoot:     fmt.Sprintf("%#x", b.StateRoot),
		BodyRoot:      fmt.Sprintf("%#x", b.BodyRoot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlockHeader) UnmarshalJSON(input []byte) error {
	var beaconBlockHeaderJSON beaconBlockHeaderJSON
	if err := json.Unmarshal(input, &beaconBlockHeaderJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return b.unpack(&beaconBlockHeaderJSON)
}

func (b *BeaconBlockHeader) unpack(beaconBlockHeaderJSON *beaconBlockHeaderJSON) error {
	if beaconBlockHeaderJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(beaconBlockHeaderJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	b.Slot = Slot(slot)
	if beaconBlockHeaderJSON.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}
	proposerIndex, err := strconv.ParseUint(beaconBlockHeaderJSON.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}
	b.ProposerIndex = ValidatorIndex(proposerIndex)
	if beaconBlockHeaderJSON.ParentRoot == "" {
		return errors.New("parent root missing")
	}
	parentRoot, err := hex.DecodeString(strings.TrimPrefix(beaconBlockHeaderJSON.ParentRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for parent root")
	}
	if len(parentRoot) != RootLength {
		return errors.New("incorrect length for parent root")
	}
	copy(b.ParentRoot[:], parentRoot)
	if beaconBlockHeaderJSON.StateRoot == "" {
		return errors.New("state root missing")
	}
	stateRoot, err := hex.DecodeString(strings.TrimPrefix(beaconBlockHeaderJSON.StateRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for state root")
	}
	if len(stateRoot) != RootLength {
		return errors.New("incorrect length for state root")
	}
	copy(b.StateRoot[:], stateRoot)
	if beaconBlockHeaderJSON.BodyRoot == "" {
		return errors.New("body root missing")
	}
	bodyRoot, err := hex.DecodeString(strings.TrimPrefix(beaconBlockHeaderJSON.BodyRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for body root")
	}
	if len(bodyRoot) != RootLength {
		return errors.New("incorrect length for body root")
	}
	copy(b.BodyRoot[:], bodyRoot)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (b *BeaconBlockHeader) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&beaconBlockHeaderYAML{
		Slot:          uint64(b.Slot),
		ProposerIndex: uint64(b.ProposerIndex),
		ParentRoot:    fmt.Sprintf("%#x", b.ParentRoot),
		StateRoot:     fmt.Sprintf("%#x", b.StateRoot),
		BodyRoot:      fmt.Sprintf("%#x", b.BodyRoot),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BeaconBlockHeader) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var beaconBlockHeaderJSON beaconBlockHeaderJSON
	if err := yaml.Unmarshal(input, &beaconBlockHeaderJSON); err != nil {
		return err
	}

	return b.unpack(&beaconBlockHeaderJSON)
}

// String returns a string representation of the struct.
func (b *BeaconBlockHeader) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
