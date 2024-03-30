// Copyright © 2023 Attestant Limited.
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

package deneb

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// beaconStateYAML is the spec representation of the struct.
type beaconStateYAML struct {
	GenesisTime                  uint64                       `yaml:"genesis_time"`
	GenesisValidatorsRoot        phase0.Root                  `yaml:"genesis_validators_root"`
	Slot                         phase0.Slot                  `yaml:"slot"`
	Fork                         *phase0.Fork                 `yaml:"fork"`
	LatestBlockHeader            *phase0.BeaconBlockHeader    `yaml:"latest_block_header"`
	BlockRoots                   []phase0.Root                `yaml:"block_roots"`
	StateRoots                   []phase0.Root                `yaml:"state_roots"`
	HistoricalRoots              []phase0.Root                `yaml:"historical_roots"`
	ETH1Data                     *phase0.ETH1Data             `yaml:"eth1_data"`
	ETH1DataVotes                []*phase0.ETH1Data           `yaml:"eth1_data_votes"`
	ETH1DepositIndex             uint64                       `yaml:"eth1_deposit_index"`
	Validators                   []*phase0.Validator          `yaml:"validators"`
	Balances                     []phase0.Gwei                `yaml:"balances"`
	RANDAOMixes                  []phase0.Root                `yaml:"randao_mixes"`
	Slashings                    []phase0.Gwei                `yaml:"slashings"`
	PreviousEpochParticipation   []altair.ParticipationFlags  `yaml:"previous_epoch_participation"`
	CurrentEpochParticipation    []altair.ParticipationFlags  `yaml:"current_epoch_participation"`
	JustificationBits            string                       `yaml:"justification_bits"`
	PreviousJustifiedCheckpoint  *phase0.Checkpoint           `yaml:"previous_justified_checkpoint"`
	CurrentJustifiedCheckpoint   *phase0.Checkpoint           `yaml:"current_justified_checkpoint"`
	FinalizedCheckpoint          *phase0.Checkpoint           `yaml:"finalized_checkpoint"`
	InactivityScores             []uint64                     `yaml:"inactivity_scores"`
	CurrentSyncCommittee         *altair.SyncCommittee        `yaml:"current_sync_committee"`
	NextSyncCommittee            *altair.SyncCommittee        `yaml:"next_sync_committee"`
	LatestExecutionPayloadHeader *ExecutionPayloadHeader      `yaml:"latest_execution_payload_header"`
	NextWithdrawalIndex          capella.WithdrawalIndex      `yaml:"next_withdrawal_index"`
	NextWithdrawalValidatorIndex phase0.ValidatorIndex        `yaml:"next_withdrawal_validator_index"`
	HistoricalSummaries          []*capella.HistoricalSummary `yaml:"historical_summaries"`
}

// MarshalYAML implements yaml.Marshaler.
func (b *BeaconState) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&beaconStateYAML{
		GenesisTime:                  b.GenesisTime,
		GenesisValidatorsRoot:        b.GenesisValidatorsRoot,
		Slot:                         b.Slot,
		Fork:                         b.Fork,
		LatestBlockHeader:            b.LatestBlockHeader,
		BlockRoots:                   b.BlockRoots,
		StateRoots:                   b.StateRoots,
		HistoricalRoots:              b.HistoricalRoots,
		ETH1Data:                     b.ETH1Data,
		ETH1DataVotes:                b.ETH1DataVotes,
		ETH1DepositIndex:             b.ETH1DepositIndex,
		Validators:                   b.Validators,
		Balances:                     b.Balances,
		RANDAOMixes:                  b.RANDAOMixes,
		Slashings:                    b.Slashings,
		PreviousEpochParticipation:   b.PreviousEpochParticipation,
		CurrentEpochParticipation:    b.CurrentEpochParticipation,
		JustificationBits:            fmt.Sprintf("%#x", b.JustificationBits.Bytes()),
		PreviousJustifiedCheckpoint:  b.PreviousJustifiedCheckpoint,
		CurrentJustifiedCheckpoint:   b.CurrentJustifiedCheckpoint,
		FinalizedCheckpoint:          b.FinalizedCheckpoint,
		InactivityScores:             b.InactivityScores,
		CurrentSyncCommittee:         b.CurrentSyncCommittee,
		NextSyncCommittee:            b.NextSyncCommittee,
		LatestExecutionPayloadHeader: b.LatestExecutionPayloadHeader,
		NextWithdrawalIndex:          b.NextWithdrawalIndex,
		NextWithdrawalValidatorIndex: b.NextWithdrawalValidatorIndex,
		HistoricalSummaries:          b.HistoricalSummaries,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (b *BeaconState) UnmarshalYAML(input []byte) error {
	// This is very inefficient, but YAML is only used for spec tests so we do this
	// rather than maintain a custom YAML unmarshaller.
	var data beaconStateJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal YAML")
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	return b.UnmarshalJSON(bytes)
}
