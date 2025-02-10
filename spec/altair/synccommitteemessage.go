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

// SyncCommitteeMessage is the Ethereum 2 sync committee message structure.
type SyncCommitteeMessage struct {
	Slot            phase0.Slot
	BeaconBlockRoot phase0.Root `ssz-size:"32"`
	ValidatorIndex  phase0.ValidatorIndex
	Signature       phase0.BLSSignature `ssz-size:"96"`
}

// syncCommitteeMessageJSON is the spec representation of the struct.
type syncCommitteeMessageJSON struct {
	Slot            string `json:"slot"`
	BeaconBlockRoot string `json:"beacon_block_root"`
	ValidatorIndex  string `json:"validator_index"`
	Signature       string `json:"signature"`
}

// syncCommitteeMessageYAML is the spec representation of the struct.
type syncCommitteeMessageYAML struct {
	Slot            uint64 `yaml:"slot"`
	BeaconBlockRoot string `yaml:"beacon_block_root"`
	ValidatorIndex  uint64 `yaml:"validator_index"`
	Signature       string `yaml:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncCommitteeMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&syncCommitteeMessageJSON{
		Slot:            fmt.Sprintf("%d", s.Slot),
		BeaconBlockRoot: fmt.Sprintf("%#x", s.BeaconBlockRoot),
		ValidatorIndex:  fmt.Sprintf("%d", s.ValidatorIndex),
		Signature:       fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncCommitteeMessage) UnmarshalJSON(input []byte) error {
	var syncCommitteeMessageJSON syncCommitteeMessageJSON
	if err := json.Unmarshal(input, &syncCommitteeMessageJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return s.unpack(&syncCommitteeMessageJSON)
}

func (s *SyncCommitteeMessage) unpack(syncCommitteeMessageJSON *syncCommitteeMessageJSON) error {
	if syncCommitteeMessageJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(syncCommitteeMessageJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	s.Slot = phase0.Slot(slot)
	if syncCommitteeMessageJSON.BeaconBlockRoot == "" {
		return errors.New("beacon block root missing")
	}
	beaconBlockRoot, err := hex.DecodeString(strings.TrimPrefix(syncCommitteeMessageJSON.BeaconBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for beacon block root")
	}
	if len(beaconBlockRoot) != phase0.RootLength {
		return errors.New("incorrect length for beacon block root")
	}
	copy(s.BeaconBlockRoot[:], beaconBlockRoot)
	if syncCommitteeMessageJSON.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(syncCommitteeMessageJSON.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	s.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)
	if syncCommitteeMessageJSON.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(syncCommitteeMessageJSON.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(signature) != phase0.SignatureLength {
		return errors.New("incorrect length for signature")
	}
	copy(s.Signature[:], signature)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (s *SyncCommitteeMessage) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&syncCommitteeMessageYAML{
		Slot:            uint64(s.Slot),
		BeaconBlockRoot: fmt.Sprintf("%#x", s.BeaconBlockRoot),
		ValidatorIndex:  uint64(s.ValidatorIndex),
		Signature:       fmt.Sprintf("%#x", s.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SyncCommitteeMessage) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var syncCommitteeMessageJSON syncCommitteeMessageJSON
	if err := yaml.Unmarshal(input, &syncCommitteeMessageJSON); err != nil {
		return err
	}

	return s.unpack(&syncCommitteeMessageJSON)
}

// String returns a string version of the structure.
func (s *SyncCommitteeMessage) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
