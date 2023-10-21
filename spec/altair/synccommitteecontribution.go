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
	bitfield "github.com/prysmaticlabs/go-bitfield"
)

// SyncCommitteeContribution is the Ethereum 2 sync committee contribution structure.
type SyncCommitteeContribution struct {
	Slot              phase0.Slot
	BeaconBlockRoot   phase0.Root `ssz-size:"32"`
	SubcommitteeIndex uint64
	// AggregationBits size is SYNC_COMMITTEE_SIZE // SYNC_COMMITTEE_SUBNET_COUNT
	AggregationBits bitfield.Bitvector128 `ssz-size:"16"`
	Signature       phase0.BLSSignature   `ssz-size:"96"`
}

// syncCommitteeContributionJSON is the spec representation of the struct.
type syncCommitteeContributionJSON struct {
	Slot              string `json:"slot"`
	BeaconBlockRoot   string `json:"beacon_block_root"`
	SubcommitteeIndex string `json:"subcommittee_index"`
	AggregationBits   string `json:"aggregation_bits"`
	Signature         string `json:"signature"`
}

// syncCommitteeContributionYAML is the spec representation of the struct.
type syncCommitteeContributionYAML struct {
	Slot              uint64 `yaml:"slot"`
	BeaconBlockRoot   string `yaml:"beacon_block_root"`
	SubcommitteeIndex uint64 `yaml:"subcommittee_index"`
	AggregationBits   string `json:"aggregation_bits"`
	Signature         string `yaml:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncCommitteeContribution) MarshalJSON() ([]byte, error) {
	return json.Marshal(&syncCommitteeContributionJSON{
		Slot:              fmt.Sprintf("%d", s.Slot),
		BeaconBlockRoot:   fmt.Sprintf("%#x", s.BeaconBlockRoot),
		SubcommitteeIndex: strconv.FormatUint(s.SubcommitteeIndex, 10),
		AggregationBits:   fmt.Sprintf("%#x", []byte(s.AggregationBits)),
		Signature:         fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncCommitteeContribution) UnmarshalJSON(input []byte) error {
	var syncCommitteeContributionJSON syncCommitteeContributionJSON
	if err := json.Unmarshal(input, &syncCommitteeContributionJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return s.unpack(&syncCommitteeContributionJSON)
}

func (s *SyncCommitteeContribution) unpack(syncCommitteeContributionJSON *syncCommitteeContributionJSON) error {
	if syncCommitteeContributionJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(syncCommitteeContributionJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	s.Slot = phase0.Slot(slot)
	if syncCommitteeContributionJSON.BeaconBlockRoot == "" {
		return errors.New("beacon block root missing")
	}
	beaconBlockRoot, err := hex.DecodeString(strings.TrimPrefix(syncCommitteeContributionJSON.BeaconBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for beacon block root")
	}
	if len(beaconBlockRoot) != phase0.RootLength {
		return errors.New("incorrect length for beacon block root")
	}
	copy(s.BeaconBlockRoot[:], beaconBlockRoot)
	if syncCommitteeContributionJSON.SubcommitteeIndex == "" {
		return errors.New("subcommittee index missing")
	}
	subCommitteeIndex, err := strconv.ParseUint(syncCommitteeContributionJSON.SubcommitteeIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for subcommittee index")
	}
	s.SubcommitteeIndex = subCommitteeIndex
	if syncCommitteeContributionJSON.AggregationBits == "" {
		return errors.New("aggregation bits missing")
	}
	if s.AggregationBits, err = hex.DecodeString(strings.TrimPrefix(syncCommitteeContributionJSON.AggregationBits, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for aggregation bits")
	}
	if syncCommitteeContributionJSON.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(syncCommitteeContributionJSON.Signature, "0x"))
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
func (s *SyncCommitteeContribution) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&syncCommitteeContributionYAML{
		Slot:              uint64(s.Slot),
		BeaconBlockRoot:   fmt.Sprintf("%#x", s.BeaconBlockRoot),
		SubcommitteeIndex: s.SubcommitteeIndex,
		AggregationBits:   fmt.Sprintf("%#x", []byte(s.AggregationBits)),
		Signature:         fmt.Sprintf("%#x", s.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SyncCommitteeContribution) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var syncCommitteeContributionJSON syncCommitteeContributionJSON
	if err := yaml.Unmarshal(input, &syncCommitteeContributionJSON); err != nil {
		return err
	}

	return s.unpack(&syncCommitteeContributionJSON)
}

// String returns a string version of the structure.
func (s *SyncCommitteeContribution) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
