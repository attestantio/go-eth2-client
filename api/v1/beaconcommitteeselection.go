// Copyright © 2025 Attestant Limited.
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

// BeaconCommitteeSelection is the data required for a beacon committee selection.
type BeaconCommitteeSelection struct {
	// ValidatorIndex is the index of the validator making the selection request.
	ValidatorIndex phase0.ValidatorIndex
	// Slot is the slot for which the validator is attesting.
	Slot phase0.Slot
	// SelectionProof is the proof of the validator being selected for beacon committee aggregation.
	SelectionProof phase0.BLSSignature
}

// beaconCommitteeSelectionJSON is the spec representation of the struct.
type beaconCommitteeSelectionJSON struct {
	ValidatorIndex string `json:"validator_index"`
	Slot           string `json:"slot"`
	SelectionProof string `json:"selection_proof"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconCommitteeSelection) MarshalJSON() ([]byte, error) {
	return json.Marshal(&beaconCommitteeSelectionJSON{
		ValidatorIndex: fmt.Sprintf("%d", b.ValidatorIndex),
		Slot:           fmt.Sprintf("%d", b.Slot),
		SelectionProof: fmt.Sprintf("%#x", b.SelectionProof),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconCommitteeSelection) UnmarshalJSON(input []byte) error {
	var err error

	var beaconCommitteeSelectionJSON beaconCommitteeSelectionJSON
	if err = json.Unmarshal(input, &beaconCommitteeSelectionJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if beaconCommitteeSelectionJSON.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(beaconCommitteeSelectionJSON.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	b.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)
	if beaconCommitteeSelectionJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(beaconCommitteeSelectionJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	b.Slot = phase0.Slot(slot)
	selectionProof, err := hex.DecodeString(strings.TrimPrefix(beaconCommitteeSelectionJSON.SelectionProof, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for selection proof")
	}
	if len(selectionProof) != phase0.SignatureLength {
		return errors.New("incorrect length for selection proof")
	}
	copy(b.SelectionProof[:], selectionProof)

	return nil
}

// String returns a string version of the structure.
func (b *BeaconCommitteeSelection) String() string {
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
