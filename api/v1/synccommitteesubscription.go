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

package v1

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SyncCommitteeSubscription is the data required for a sync committee subscription.
type SyncCommitteeSubscription struct {
	// ValidatorIdex is the index of the validator making the subscription request.
	ValidatorIndex phase0.ValidatorIndex
	// SyncCommitteeIndices are the indices of the sync committees of which the validator is a member.
	SyncCommitteeIndices []phase0.CommitteeIndex
	// UntilEpoch is the epoch at which the subscription no longer applies.
	UntilEpoch phase0.Epoch
}

// syncCommitteeSubscriptionJSON is the spec representation of the struct.
type syncCommitteeSubscriptionJSON struct {
	ValidatorIndex       string   `json:"validator_index"`
	SyncCommitteeIndices []string `json:"sync_committee_indices"`
	UntilEpoch           string   `json:"until_epoch"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncCommitteeSubscription) MarshalJSON() ([]byte, error) {
	syncCommitteeIndices := make([]string, len(s.SyncCommitteeIndices))
	for i, syncCommitteeIndex := range s.SyncCommitteeIndices {
		syncCommitteeIndices[i] = fmt.Sprintf("%d", syncCommitteeIndex)
	}

	return json.Marshal(&syncCommitteeSubscriptionJSON{
		ValidatorIndex:       fmt.Sprintf("%d", s.ValidatorIndex),
		SyncCommitteeIndices: syncCommitteeIndices,
		UntilEpoch:           fmt.Sprintf("%d", s.UntilEpoch),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncCommitteeSubscription) UnmarshalJSON(input []byte) error {
	var err error

	var syncCommitteeSubscriptionJSON syncCommitteeSubscriptionJSON
	if err = json.Unmarshal(input, &syncCommitteeSubscriptionJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if syncCommitteeSubscriptionJSON.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(syncCommitteeSubscriptionJSON.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	s.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	if len(syncCommitteeSubscriptionJSON.SyncCommitteeIndices) == 0 {
		return errors.New("sync committee indices missing")
	}
	s.SyncCommitteeIndices = make([]phase0.CommitteeIndex, len(syncCommitteeSubscriptionJSON.SyncCommitteeIndices))
	for i, committeeIndex := range syncCommitteeSubscriptionJSON.SyncCommitteeIndices {
		syncCommitteeIndex, err := strconv.ParseUint(committeeIndex, 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid value for sync committee index")
		}
		s.SyncCommitteeIndices[i] = phase0.CommitteeIndex(syncCommitteeIndex)
	}
	if syncCommitteeSubscriptionJSON.UntilEpoch == "" {
		return errors.New("until epoch missing")
	}
	untilEpoch, err := strconv.ParseUint(syncCommitteeSubscriptionJSON.UntilEpoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for until epoch")
	}
	s.UntilEpoch = phase0.Epoch(untilEpoch)

	return nil
}

// String returns a string version of the structure.
func (s *SyncCommitteeSubscription) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
