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

// SyncCommittee is the data providing validator membership of sync committees.
type SyncCommittee struct {
	// Validators is the list of validator indices in the committee.
	Validators []phase0.ValidatorIndex
	// ValidatorAggregates are the lists of validators in each aggregate.
	ValidatorAggregates [][]phase0.ValidatorIndex
}

// syncCommitteeJSON is the spec representation of the struct.
type syncCommitteeJSON struct {
	Validators          []string   `json:"validators"`
	ValidatorAggregates [][]string `json:"validator_aggregates"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncCommittee) MarshalJSON() ([]byte, error) {
	validators := make([]string, len(s.Validators))
	for i := range s.Validators {
		validators[i] = fmt.Sprintf("%d", s.Validators[i])
	}
	validatorAggregates := make([][]string, len(s.ValidatorAggregates))
	for i := range s.ValidatorAggregates {
		validatorAggregates[i] = make([]string, len(s.ValidatorAggregates[i]))
		for j := range s.ValidatorAggregates[i] {
			validatorAggregates[i][j] = fmt.Sprintf("%d", s.ValidatorAggregates[i][j])
		}
	}

	return json.Marshal(&syncCommitteeJSON{
		Validators:          validators,
		ValidatorAggregates: validatorAggregates,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncCommittee) UnmarshalJSON(input []byte) error {
	var err error

	var syncCommitteeJSON syncCommitteeJSON
	if err = json.Unmarshal(input, &syncCommitteeJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if syncCommitteeJSON.Validators == nil {
		return errors.New("validators missing")
	}
	if len(syncCommitteeJSON.Validators) == 0 {
		return errors.New("validators length cannot be 0")
	}
	s.Validators = make([]phase0.ValidatorIndex, len(syncCommitteeJSON.Validators))
	for i := range syncCommitteeJSON.Validators {
		validator, err := strconv.ParseUint(syncCommitteeJSON.Validators[i], 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid value for validator")
		}
		s.Validators[i] = phase0.ValidatorIndex(validator)
	}
	if syncCommitteeJSON.ValidatorAggregates == nil {
		return errors.New("validator aggregates missing")
	}
	if len(syncCommitteeJSON.ValidatorAggregates) == 0 {
		return errors.New("validator aggregates length cannot be 0")
	}
	s.ValidatorAggregates = make([][]phase0.ValidatorIndex, len(syncCommitteeJSON.ValidatorAggregates))
	for i := range syncCommitteeJSON.ValidatorAggregates {
		if len(syncCommitteeJSON.ValidatorAggregates[i]) == 0 {
			return errors.New("validator aggregate length cannot be 0")
		}
		s.ValidatorAggregates[i] = make([]phase0.ValidatorIndex, len(syncCommitteeJSON.ValidatorAggregates[i]))
		for j := range syncCommitteeJSON.ValidatorAggregates[i] {
			validator, err := strconv.ParseUint(syncCommitteeJSON.ValidatorAggregates[i][j], 10, 64)
			if err != nil {
				return errors.Wrap(err, "invalid value for validator aggregate")
			}
			s.ValidatorAggregates[i][j] = phase0.ValidatorIndex(validator)
		}
	}

	return nil
}

// String returns a string version of the structure.
func (s *SyncCommittee) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
