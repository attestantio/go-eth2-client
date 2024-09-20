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

	"github.com/pkg/errors"
)

// BeaconBlockReward Rewards info for a single block
type BeaconBlockReward struct {
	// Proposer of the block, the proposer index who receives these rewards
	ProposerIndex uint64
	// Total block reward in gwei, equal to attestations + sync_aggregate + proposer_slashings + attester_slashings
	Total uint64
	// Block reward component due to included attestations in gwei
	Attestations uint64
	// Block reward component due to included sync_aggregate in gwei
	SyncAggregate uint64
	// Block reward component due to included proposer_slashings in gwei
	ProposerSlashings uint64
	// Block reward component due to included attester_slashings in gwei
	AttesterSlashings uint64
}

// beaconBlockRewardJSON is the spec representation of the struct.
type beaconBlockRewardJSON struct {
	ProposerIndex     string `json:"proposer_index"`
	Total             string `json:"total"`
	Attestations      string `json:"attestations"`
	SyncAggregate     string `json:"sync_aggregate"`
	ProposerSlashings string `json:"proposer_slashings"`
	AttesterSlashings string `json:"attester_slashings"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconBlockReward) MarshalJSON() ([]byte, error) {
	return json.Marshal(&beaconBlockRewardJSON{
		ProposerIndex:     fmt.Sprintf("%d", b.ProposerIndex),
		Total:             fmt.Sprintf("%d", b.Total),
		Attestations:      fmt.Sprintf("%d", b.Attestations),
		SyncAggregate:     fmt.Sprintf("%d", b.SyncAggregate),
		ProposerSlashings: fmt.Sprintf("%d", b.ProposerSlashings),
		AttesterSlashings: fmt.Sprintf("%d", b.AttesterSlashings),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlockReward) UnmarshalJSON(input []byte) error {
	var err error

	var beaconBlockRewardJSON beaconBlockRewardJSON
	if err = json.Unmarshal(input, &beaconBlockRewardJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return b.unpack(&beaconBlockRewardJSON)
}

func (b *BeaconBlockReward) unpack(beaconBlockRewardJSON *beaconBlockRewardJSON) error {
	if beaconBlockRewardJSON.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}

	proposerIndex, err := strconv.ParseUint(beaconBlockRewardJSON.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}
	b.ProposerIndex = proposerIndex

	if beaconBlockRewardJSON.Total == "" {
		return errors.New("total missing")
	}

	total, err := strconv.ParseUint(beaconBlockRewardJSON.Total, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for total")
	}
	b.Total = total

	if beaconBlockRewardJSON.Attestations == "" {
		return errors.New("total missing")
	}

	attestations, err := strconv.ParseUint(beaconBlockRewardJSON.Attestations, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for attestations")
	}
	b.Attestations = attestations

	if beaconBlockRewardJSON.SyncAggregate == "" {
		return errors.New("sync aggregate missing")
	}

	syncAggregate, err := strconv.ParseUint(beaconBlockRewardJSON.SyncAggregate, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for sync aggregate")
	}
	b.SyncAggregate = syncAggregate

	if beaconBlockRewardJSON.ProposerSlashings == "" {
		return errors.New("proposer slashing missing")
	}

	proposerSlashings, err := strconv.ParseUint(beaconBlockRewardJSON.ProposerSlashings, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer slashings")
	}
	b.ProposerSlashings = proposerSlashings

	if beaconBlockRewardJSON.AttesterSlashings == "" {
		return errors.New("proposer slashing missing")
	}

	attesterSlashings, err := strconv.ParseUint(beaconBlockRewardJSON.AttesterSlashings, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for attester slashings")
	}
	b.AttesterSlashings = attesterSlashings

	return nil
}

// String returns a string version of the structure.
func (b *BeaconBlockReward) String() string {
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
