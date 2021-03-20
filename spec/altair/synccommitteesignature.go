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

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// SyncCommitteeSignature is the Ethereum 2 sync committee signature structure.
type SyncCommitteeSignature struct {
	Slot            Slot
	BeaconBlockRoot Root `ssz-size:"32"`
	ValidatorIndex  ValidatorIndex
	Signature       BLSSignature `ssz-size:"96"`
}

// syncCommitteeSignatureJSON is the spec representation of the struct.
type syncCommitteeSignatureJSON struct {
	Slot            string `json:"slot"`
	BeaconBlockRoot string `json:"beacon_block_root"`
	ValidatorIndex  string `json:"validator_index"`
	Signature       string `json:"signature"`
}

// syncCommitteeSignatureYAML is the spec representation of the struct.
type syncCommitteeSignatureYAML struct {
	Slot            uint64 `yaml:"slot"`
	BeaconBlockRoot string `yaml:"beacon_block_root"`
	ValidatorIndex  uint64 `yaml:"validator_index"`
	Signature       string `yaml:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncCommitteeSignature) MarshalJSON() ([]byte, error) {
	return json.Marshal(&syncCommitteeSignatureJSON{
		Slot:            fmt.Sprintf("%d", s.Slot),
		BeaconBlockRoot: fmt.Sprintf("%#x", s.BeaconBlockRoot),
		ValidatorIndex:  fmt.Sprintf("%d", s.ValidatorIndex),
		Signature:       fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncCommitteeSignature) UnmarshalJSON(input []byte) error {
	var syncCommitteeSignatureJSON syncCommitteeSignatureJSON
	if err := json.Unmarshal(input, &syncCommitteeSignatureJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	return s.unpack(&syncCommitteeSignatureJSON)
}

func (s *SyncCommitteeSignature) unpack(syncCommitteeSignatureJSON *syncCommitteeSignatureJSON) error {
	if syncCommitteeSignatureJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(syncCommitteeSignatureJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	s.Slot = Slot(slot)
	if syncCommitteeSignatureJSON.BeaconBlockRoot == "" {
		return errors.New("beacon block root missing")
	}
	beaconBlockRoot, err := hex.DecodeString(strings.TrimPrefix(syncCommitteeSignatureJSON.BeaconBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for beacon block root")
	}
	if len(beaconBlockRoot) != RootLength {
		return errors.New("incorrect length for beacon block root")
	}
	copy(s.BeaconBlockRoot[:], beaconBlockRoot)
	if syncCommitteeSignatureJSON.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(syncCommitteeSignatureJSON.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	s.ValidatorIndex = ValidatorIndex(validatorIndex)
	if syncCommitteeSignatureJSON.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(syncCommitteeSignatureJSON.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(signature) != SignatureLength {
		return errors.New("incorrect length for signature")
	}
	copy(s.Signature[:], signature)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (s *SyncCommitteeSignature) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&syncCommitteeSignatureYAML{
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
func (s *SyncCommitteeSignature) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var syncCommitteeSignatureJSON syncCommitteeSignatureJSON
	if err := yaml.Unmarshal(input, &syncCommitteeSignatureJSON); err != nil {
		return err
	}
	return s.unpack(&syncCommitteeSignatureJSON)
}

// String returns a string version of the structure.
func (s *SyncCommitteeSignature) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
