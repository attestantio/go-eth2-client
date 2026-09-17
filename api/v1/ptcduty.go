// Copyright © 2026 Attestant Limited.
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

// PTCDuty is the data regarding which validators have the duty to attest to
// payload timeliness in a slot.
//
// The payload timeliness committee is a flat, per-slot committee, so unlike
// AttesterDuty this carries no committee index, committee length or
// validator-committee index.
type PTCDuty struct {
	// PubKey is the public key of the validator that should attest.
	PubKey phase0.BLSPubKey
	// Slot is the slot in which the validator should attest.
	Slot phase0.Slot
	// ValidatorIndex is the index of the validator that should attest.
	ValidatorIndex phase0.ValidatorIndex
}

// ptcDutyJSON is the spec representation of the struct.
type ptcDutyJSON struct {
	PubKey         string `json:"pubkey"`
	ValidatorIndex string `json:"validator_index"`
	Slot           string `json:"slot"`
}

// MarshalJSON implements json.Marshaler.
func (p *PTCDuty) MarshalJSON() ([]byte, error) {
	return json.Marshal(&ptcDutyJSON{
		PubKey:         fmt.Sprintf("%#x", p.PubKey),
		ValidatorIndex: fmt.Sprintf("%d", p.ValidatorIndex),
		Slot:           fmt.Sprintf("%d", p.Slot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PTCDuty) UnmarshalJSON(input []byte) error {
	var err error

	var ptcDutyJSON ptcDutyJSON
	if err = json.Unmarshal(input, &ptcDutyJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if ptcDutyJSON.PubKey == "" {
		return errors.New("public key missing")
	}

	pubKey, err := hex.DecodeString(strings.TrimPrefix(ptcDutyJSON.PubKey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}

	if len(pubKey) != publicKeyLength {
		return errors.New("incorrect length for public key")
	}

	copy(p.PubKey[:], pubKey)

	if ptcDutyJSON.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}

	validatorIndex, err := strconv.ParseUint(ptcDutyJSON.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}

	p.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	if ptcDutyJSON.Slot == "" {
		return errors.New("slot missing")
	}

	slot, err := strconv.ParseUint(ptcDutyJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}

	p.Slot = phase0.Slot(slot)

	return nil
}

// String returns a string version of the structure.
func (p *PTCDuty) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
