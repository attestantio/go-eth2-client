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

package v1

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

// BlindedBeaconBlock represents a blinded beacon block.
type BlindedBeaconBlock struct {
	Slot          phase0.Slot
	ProposerIndex phase0.ValidatorIndex
	ParentRoot    phase0.Root `ssz-size:"32"`
	StateRoot     phase0.Root `ssz-size:"32"`
	Body          *BlindedBeaconBlockBody
}

// blindedBeaconBlockJSON is the spec representation of the struct.
type blindedBeaconBlockJSON struct {
	Slot          string                  `json:"slot"`
	ProposerIndex string                  `json:"proposer_index"`
	ParentRoot    string                  `json:"parent_root"`
	StateRoot     string                  `json:"state_root"`
	Body          *BlindedBeaconBlockBody `json:"body"`
}

// blindedBeaconBlockYAML is the spec representation of the struct.
type blindedBeaconBlockYAML struct {
	Slot          uint64                  `yaml:"slot"`
	ProposerIndex uint64                  `yaml:"proposer_index"`
	ParentRoot    string                  `yaml:"parent_root"`
	StateRoot     string                  `yaml:"state_root"`
	Body          *BlindedBeaconBlockBody `yaml:"body"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlindedBeaconBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blindedBeaconBlockJSON{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    b.ParentRoot.String(),
		StateRoot:     b.StateRoot.String(),
		Body:          b.Body,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlindedBeaconBlock) UnmarshalJSON(input []byte) error {
	var data blindedBeaconBlockJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	return b.unpack(&data)
}

func (b *BlindedBeaconBlock) unpack(data *blindedBeaconBlockJSON) error {
	if data.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	b.Slot = phase0.Slot(slot)
	if data.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}
	proposerIndex, err := strconv.ParseUint(data.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}
	b.ProposerIndex = phase0.ValidatorIndex(proposerIndex)
	if data.ParentRoot == "" {
		return errors.New("parent root missing")
	}
	parentRoot, err := hex.DecodeString(strings.TrimPrefix(data.ParentRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for parent root")
	}
	if len(parentRoot) != phase0.RootLength {
		return errors.New("incorrect length for parent root")
	}
	copy(b.ParentRoot[:], parentRoot)
	if data.StateRoot == "" {
		return errors.New("state root missing")
	}
	stateRoot, err := hex.DecodeString(strings.TrimPrefix(data.StateRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for state root")
	}
	if len(stateRoot) != phase0.RootLength {
		return errors.New("incorrect length for state root")
	}
	copy(b.StateRoot[:], stateRoot)
	if data.Body == nil {
		return errors.New("body missing")
	}
	b.Body = data.Body

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (b *BlindedBeaconBlock) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&blindedBeaconBlockYAML{
		Slot:          uint64(b.Slot),
		ProposerIndex: uint64(b.ProposerIndex),
		ParentRoot:    b.ParentRoot.String(),
		StateRoot:     b.StateRoot.String(),
		Body:          b.Body,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}
	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BlindedBeaconBlock) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data blindedBeaconBlockJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}
	return b.unpack(&data)
}

// String returns a string version of the structure.
func (b *BlindedBeaconBlock) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
