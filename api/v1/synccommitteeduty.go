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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SyncCommitteeDuty is the data regarding which validators have the duty to contribute to sync committees in a slot.
type SyncCommitteeDuty struct {
	// PubKey is the public key of the validator that should contribute.
	PubKey phase0.BLSPubKey
	// ValidatorIndex is the index of the validator that should contribute.
	ValidatorIndex phase0.ValidatorIndex
	// ValidatorSyncCommitteeIndices is the index of the validator in the list of validators in the committee.
	ValidatorSyncCommitteeIndices []phase0.CommitteeIndex
}

// syncCommitteeDutyJSON is the spec representation of the struct.
type syncCommitteeDutyJSON struct {
	PubKey                        string   `json:"pubkey"`
	ValidatorIndex                string   `json:"validator_index"`
	ValidatorSyncCommitteeIndices []string `json:"validator_sync_committee_indices"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncCommitteeDuty) MarshalJSON() ([]byte, error) {
	validatorSyncCommitteeIndices := make([]string, len(s.ValidatorSyncCommitteeIndices))
	for i := range s.ValidatorSyncCommitteeIndices {
		validatorSyncCommitteeIndices[i] = fmt.Sprintf("%d", s.ValidatorSyncCommitteeIndices[i])
	}

	return json.Marshal(&syncCommitteeDutyJSON{
		PubKey:                        fmt.Sprintf("%#x", s.PubKey),
		ValidatorIndex:                fmt.Sprintf("%d", s.ValidatorIndex),
		ValidatorSyncCommitteeIndices: validatorSyncCommitteeIndices,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncCommitteeDuty) UnmarshalJSON(input []byte) error {
	var err error

	var syncCommitteeDutyJSON syncCommitteeDutyJSON
	if err = json.Unmarshal(input, &syncCommitteeDutyJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if syncCommitteeDutyJSON.PubKey == "" {
		return errors.New("public key missing")
	}
	pubKey, err := hex.DecodeString(strings.TrimPrefix(syncCommitteeDutyJSON.PubKey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}
	if len(pubKey) != publicKeyLength {
		return errors.New("incorrect length for public key")
	}
	copy(s.PubKey[:], pubKey)
	if syncCommitteeDutyJSON.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(syncCommitteeDutyJSON.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	s.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	if len(syncCommitteeDutyJSON.ValidatorSyncCommitteeIndices) == 0 {
		return errors.New("validator sync committee indices missing")
	}
	s.ValidatorSyncCommitteeIndices = make([]phase0.CommitteeIndex, len(syncCommitteeDutyJSON.ValidatorSyncCommitteeIndices))
	for i := range syncCommitteeDutyJSON.ValidatorSyncCommitteeIndices {
		committeeIndex, err := strconv.ParseUint(syncCommitteeDutyJSON.ValidatorSyncCommitteeIndices[i], 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid value for sync committee index")
		}
		s.ValidatorSyncCommitteeIndices[i] = phase0.CommitteeIndex(committeeIndex)
	}

	return nil
}

// String returns a string version of the structure.
func (s *SyncCommitteeDuty) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
