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

// BlockRewards are the rewards for proposing a block.
type BlockRewards struct {
	ProposerIndex     phase0.ValidatorIndex
	Total             phase0.Gwei
	Attestations      phase0.Gwei
	SyncAggregate     phase0.Gwei
	ProposerSlashings phase0.Gwei
	AttesterSlashings phase0.Gwei
}

// blockRewardsJSON is the spec representation of the struct.
type blockRewardsJSON struct {
	ProposerIndex     string `json:"proposer_index"`
	Total             string `json:"total"`
	Attestations      string `json:"attestations"`
	SyncAggregate     string `json:"sync_aggregate"`
	ProposerSlashings string `json:"proposer_slashings"`
	AttesterSlashings string `json:"attester_slashings"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlockRewards) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blockRewardsJSON{
		ProposerIndex:     fmt.Sprintf("%d", b.ProposerIndex),
		Total:             fmt.Sprintf("%d", b.Total),
		Attestations:      fmt.Sprintf("%d", b.Attestations),
		SyncAggregate:     fmt.Sprintf("%d", b.SyncAggregate),
		ProposerSlashings: fmt.Sprintf("%d", b.ProposerSlashings),
		AttesterSlashings: fmt.Sprintf("%d", b.AttesterSlashings),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlockRewards) UnmarshalJSON(input []byte) error {
	var err error

	var data blockRewardsJSON
	if err = json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}
	proposerIndex, err := strconv.ParseUint(data.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}
	b.ProposerIndex = phase0.ValidatorIndex(proposerIndex)

	if data.Total == "" {
		return errors.New("total missing")
	}
	total, err := strconv.ParseUint(data.Total, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for total")
	}
	b.Total = phase0.Gwei(total)

	if data.Attestations == "" {
		return errors.New("attestations missing")
	}
	attestations, err := strconv.ParseUint(data.Attestations, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for attestations")
	}
	b.Attestations = phase0.Gwei(attestations)

	if data.SyncAggregate == "" {
		return errors.New("sync aggregate missing")
	}
	syncAggregate, err := strconv.ParseUint(data.SyncAggregate, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for sync aggregate")
	}
	b.SyncAggregate = phase0.Gwei(syncAggregate)

	if data.ProposerSlashings == "" {
		return errors.New("proposer slashings missing")
	}
	proposerSlashings, err := strconv.ParseUint(data.ProposerSlashings, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer slashings")
	}
	b.ProposerSlashings = phase0.Gwei(proposerSlashings)

	if data.AttesterSlashings == "" {
		return errors.New("attester slashings missing")
	}
	attesterSlashings, err := strconv.ParseUint(data.AttesterSlashings, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for attester slashings")
	}
	b.AttesterSlashings = phase0.Gwei(attesterSlashings)

	return nil
}

// String returns a string version of the structure.
func (b *BlockRewards) String() string {
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
