// Copyright © 2020, 2021 Attestant Limited.
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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// AttesterDuty is the data regarding which validators have the duty to attest in a slot.
type AttesterDuty struct {
	// PubKey is the public key of the validator that should attest.
	PubKey phase0.BLSPubKey
	// Slot is the slot in which the validator should attest.
	Slot phase0.Slot
	// ValidatorIndex is the index of the validator that should attest.
	ValidatorIndex phase0.ValidatorIndex
	// CommitteeIndex is the index of the committee in which the attesting validator has been placed.
	CommitteeIndex phase0.CommitteeIndex
	// CommitteeLength is the length of the committee in which the attesting validator has been placed.
	CommitteeLength uint64
	// CommitteesAtSlot is the number of committees in the slot.
	CommitteesAtSlot uint64
	// ValidatorCommitteeIndex is the index of the validator in the list of validators in the committee.
	ValidatorCommitteeIndex uint64
}

// attesterDutyJSON is the spec representation of the struct.
type attesterDutyJSON struct {
	PubKey                  string `json:"pubkey"`
	Slot                    string `json:"slot"`
	ValidatorIndex          string `json:"validator_index"`
	CommitteeIndex          string `json:"committee_index"`
	CommitteeLength         string `json:"committee_length"`
	CommitteesAtSlot        string `json:"committees_at_slot"`
	ValidatorCommitteeIndex string `json:"validator_committee_index"`
}

// MarshalJSON implements json.Marshaler.
func (a *AttesterDuty) MarshalJSON() ([]byte, error) {
	return json.Marshal(&attesterDutyJSON{
		PubKey:                  fmt.Sprintf("%#x", a.PubKey),
		Slot:                    fmt.Sprintf("%d", a.Slot),
		ValidatorIndex:          fmt.Sprintf("%d", a.ValidatorIndex),
		CommitteeIndex:          fmt.Sprintf("%d", a.CommitteeIndex),
		CommitteeLength:         strconv.FormatUint(a.CommitteeLength, 10),
		CommitteesAtSlot:        strconv.FormatUint(a.CommitteesAtSlot, 10),
		ValidatorCommitteeIndex: strconv.FormatUint(a.ValidatorCommitteeIndex, 10),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AttesterDuty) UnmarshalJSON(input []byte) error {
	var err error

	var attesterDutyJSON attesterDutyJSON
	if err = json.Unmarshal(input, &attesterDutyJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if attesterDutyJSON.PubKey == "" {
		return errors.New("public key missing")
	}
	pubKey, err := hex.DecodeString(strings.TrimPrefix(attesterDutyJSON.PubKey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}
	if len(pubKey) != publicKeyLength {
		return errors.New("incorrect length for public key")
	}
	copy(a.PubKey[:], pubKey)
	if attesterDutyJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(attesterDutyJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	a.Slot = phase0.Slot(slot)
	if attesterDutyJSON.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(attesterDutyJSON.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	a.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)
	if attesterDutyJSON.CommitteeIndex == "" {
		return errors.New("committee index missing")
	}
	committeeIndex, err := strconv.ParseUint(attesterDutyJSON.CommitteeIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for committee index")
	}
	a.CommitteeIndex = phase0.CommitteeIndex(committeeIndex)
	if attesterDutyJSON.CommitteeLength == "" {
		return errors.New("committee length missing")
	}
	if a.CommitteeLength, err = strconv.ParseUint(attesterDutyJSON.CommitteeLength, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for committee length")
	}
	if a.CommitteeLength == 0 {
		return errors.New("committee length cannot be 0")
	}
	if attesterDutyJSON.CommitteesAtSlot == "" {
		return errors.New("committees at slot missing")
	}
	if a.CommitteesAtSlot, err = strconv.ParseUint(attesterDutyJSON.CommitteesAtSlot, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for committees at slot")
	}
	if a.CommitteesAtSlot == 0 {
		return errors.New("committees at slot cannot be 0")
	}
	if attesterDutyJSON.ValidatorCommitteeIndex == "" {
		return errors.New("validator committee index missing")
	}
	if a.ValidatorCommitteeIndex, err = strconv.ParseUint(attesterDutyJSON.ValidatorCommitteeIndex, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for validator committee index")
	}

	return nil
}

// String returns a string version of the structure.
func (a *AttesterDuty) String() string {
	data, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
