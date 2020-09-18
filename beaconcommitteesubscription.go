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

package client

// BeaconCommitteeSubscription is the data required for a beacon committee subscription.
type BeaconCommitteeSubscription struct {
	// Slot is the slot for which the validator is attesting.
	Slot uint64
	// CommitteeIndex is the index of the committee of which the validator is a member at the given slot.
	CommitteeIndex uint64
	// CommitteeSize is the number of validators in the committee at the given slot.
	CommitteeSize uint64
	// ValidatorIndex is the index of the valdiator that wishes to subscribe.
	ValidatorIndex uint64
	// ValidatorPubKey is the public key of the valdiator that wishes to subscribe.
	ValidatorPubKey []byte
	// Aggregate is true if the validator that wishes to subscribe also needs to aggregate attestations.
	Aggregate bool
	// SlotSelectionSignature is the result of the validator signing the slot with the "selection proof" domain.
	SlotSelectionSignature []byte
}
