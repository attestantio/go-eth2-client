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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// BeaconCommittee is the data providing information validator membership of committees.
type BeaconCommittee struct {
	// Slot is the slot in which the committee attests.
	Slot phase0.Slot
	// Index is the index of the committee.
	Index phase0.CommitteeIndex
	// Validators is the list of validator indices in the committee.
	Validators []phase0.ValidatorIndex
}

// beaconCommitteeJSON is the spec representation of the struct.
type beaconCommitteeJSON struct {
	Slot       string   `json:"slot"`
	Index      string   `json:"index"`
	Validators []string `json:"validators"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconCommittee) MarshalJSON() ([]byte, error) {
	validators := make([]string, len(b.Validators))
	for i := range b.Validators {
		validators[i] = fmt.Sprintf("%d", b.Validators[i])
	}

	return json.Marshal(&beaconCommitteeJSON{
		Slot:       fmt.Sprintf("%d", b.Slot),
		Index:      fmt.Sprintf("%d", b.Index),
		Validators: validators,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconCommittee) UnmarshalJSON(input []byte) error {
	var err error

	var beaconCommitteeJSON beaconCommitteeJSON
	if err = json.Unmarshal(input, &beaconCommitteeJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if beaconCommitteeJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(beaconCommitteeJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	b.Slot = phase0.Slot(slot)
	if beaconCommitteeJSON.Index == "" {
		return errors.New("index missing")
	}
	index, err := strconv.ParseUint(beaconCommitteeJSON.Index, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for index")
	}
	b.Index = phase0.CommitteeIndex(index)
	if beaconCommitteeJSON.Validators == nil {
		return errors.New("validators missing")
	}
	if len(beaconCommitteeJSON.Validators) == 0 {
		return errors.New("validators length cannot be 0")
	}
	b.Validators = make([]phase0.ValidatorIndex, len(beaconCommitteeJSON.Validators))
	for i := range beaconCommitteeJSON.Validators {
		validator, err := strconv.ParseUint(beaconCommitteeJSON.Validators[i], 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid value for validator")
		}
		b.Validators[i] = phase0.ValidatorIndex(validator)
	}

	return nil
}

// String returns a string version of the structure.
func (b *BeaconCommittee) String() string {
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
