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

// BeaconCommitteeSubscription is the data required for a beacon committee subscription.
type BeaconCommitteeSubscription struct {
	// ValidatorIdex is the index of the validator making the subscription request.
	ValidatorIndex phase0.ValidatorIndex
	// Slot is the slot for which the validator is attesting.
	Slot phase0.Slot
	// CommitteeIndex is the index of the committee of which the validator is a member at the given slot.
	CommitteeIndex phase0.CommitteeIndex
	// CommitteesAtSlot is the number of committees at the given slot.
	CommitteesAtSlot uint64
	// IsAggregator is true if the validator that wishes to subscribe is required to aggregate attestations.
	IsAggregator bool
}

// beaconCommitteeSubscriptionJSON is the spec representation of the struct.
type beaconCommitteeSubscriptionJSON struct {
	ValidatorIndex   string `json:"validator_index"`
	Slot             string `json:"slot"`
	CommitteeIndex   string `json:"committee_index"`
	CommitteesAtSlot string `json:"committees_at_slot"`
	IsAggregator     bool   `json:"is_aggregator"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconCommitteeSubscription) MarshalJSON() ([]byte, error) {
	return json.Marshal(&beaconCommitteeSubscriptionJSON{
		ValidatorIndex:   fmt.Sprintf("%d", b.ValidatorIndex),
		Slot:             fmt.Sprintf("%d", b.Slot),
		CommitteeIndex:   fmt.Sprintf("%d", b.CommitteeIndex),
		CommitteesAtSlot: strconv.FormatUint(b.CommitteesAtSlot, 10),
		IsAggregator:     b.IsAggregator,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconCommitteeSubscription) UnmarshalJSON(input []byte) error {
	var err error

	var beaconCommitteeSubscriptionJSON beaconCommitteeSubscriptionJSON
	if err = json.Unmarshal(input, &beaconCommitteeSubscriptionJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if beaconCommitteeSubscriptionJSON.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(beaconCommitteeSubscriptionJSON.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	b.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)
	if beaconCommitteeSubscriptionJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(beaconCommitteeSubscriptionJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	b.Slot = phase0.Slot(slot)
	if beaconCommitteeSubscriptionJSON.CommitteeIndex == "" {
		return errors.New("committee index missing")
	}
	committeeIndex, err := strconv.ParseUint(beaconCommitteeSubscriptionJSON.CommitteeIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for committee index")
	}
	b.CommitteeIndex = phase0.CommitteeIndex(committeeIndex)
	if beaconCommitteeSubscriptionJSON.CommitteesAtSlot == "" {
		return errors.New("committees at slot missing")
	}
	if b.CommitteesAtSlot, err = strconv.ParseUint(beaconCommitteeSubscriptionJSON.CommitteesAtSlot, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for committees at slot")
	}
	if b.CommitteesAtSlot == 0 {
		return errors.New("committees at slot cannot be 0")
	}
	b.IsAggregator = beaconCommitteeSubscriptionJSON.IsAggregator

	return nil
}

// String returns a string version of the structure.
func (b *BeaconCommitteeSubscription) String() string {
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
