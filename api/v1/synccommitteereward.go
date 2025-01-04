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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SyncCommitteeReward is the rewards for a validator in a sync committee.
type SyncCommitteeReward struct {
	ValidatorIndex phase0.ValidatorIndex
	// Reward can be negative, so it is an int64 (but still a Gwei value).
	Reward int64
}

// syncCommitteeRewardJSON is the spec representation of the struct.
type syncCommitteeRewardJSON struct {
	ValidatorIndex string `json:"validator_index"`
	Reward         string `json:"reward"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncCommitteeReward) MarshalJSON() ([]byte, error) {
	return json.Marshal(&syncCommitteeRewardJSON{
		ValidatorIndex: fmt.Sprintf("%d", s.ValidatorIndex),
		Reward:         fmt.Sprintf("%d", s.Reward),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncCommitteeReward) UnmarshalJSON(input []byte) error {
	var err error

	var data syncCommitteeRewardJSON
	if err = json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(data.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	s.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	if data.Reward == "" {
		return errors.New("reward missing")
	}
	s.Reward, err = strconv.ParseInt(data.Reward, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for reward")
	}

	return nil
}

// String returns a string version of the structure.
func (s *SyncCommitteeReward) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
