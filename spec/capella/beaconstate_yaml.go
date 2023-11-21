// Copyright © 2022 Attestant Limited.
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

package capella

import (
	"bytes"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// beaconStateYAML is the spec representation of the struct.
type beaconStateYAML struct {
	GenesisTime                  uint64                    `json:"genesis_time"`
	GenesisValidatorsRoot        string                    `json:"genesis_validators_root"`
	Slot                         uint64                    `json:"slot"`
	Fork                         *phase0.Fork              `json:"fork"`
	LatestBlockHeader            *phase0.BeaconBlockHeader `json:"latest_block_header"`
	BlockRoots                   []string                  `json:"block_roots"`
	StateRoots                   []string                  `json:"state_roots"`
	HistoricalRoots              []string                  `json:"historical_roots"`
	ETH1Data                     *phase0.ETH1Data          `json:"eth1_data"`
	ETH1DataVotes                []*phase0.ETH1Data        `json:"eth1_data_votes"`
	ETH1DepositIndex             uint64                    `json:"eth1_deposit_index"`
	Validators                   []*phase0.Validator       `json:"validators"`
	Balances                     []uint64                  `json:"balances"`
	RANDAOMixes                  []string                  `json:"randao_mixes"`
	Slashings                    []uint64                  `json:"slashings"`
	PreviousEpochParticipation   []uint8                   `json:"previous_epoch_participation"`
	CurrentEpochParticipation    []uint8                   `json:"current_epoch_participation"`
	JustificationBits            string                    `json:"justification_bits"`
	PreviousJustifiedCheckpoint  *phase0.Checkpoint        `json:"previous_justified_checkpoint"`
	CurrentJustifiedCheckpoint   *phase0.Checkpoint        `json:"current_justified_checkpoint"`
	FinalizedCheckpoint          *phase0.Checkpoint        `json:"finalized_checkpoint"`
	InactivityScores             []uint64                  `json:"inactivity_scores"`
	CurrentSyncCommittee         *altair.SyncCommittee     `json:"current_sync_committee"`
	NextSyncCommittee            *altair.SyncCommittee     `json:"next_sync_committee"`
	LatestExecutionPayloadHeader *ExecutionPayloadHeader   `json:"latest_execution_payload_header"`
	NextWithdrawalIndex          uint64                    `json:"next_withdrawal_index"`
	NextWithdrawalValidatorIndex uint64                    `json:"next_withdrawal_validator_index"`
	HistoricalSummaries          []*HistoricalSummary      `json:"historical_summaries"`
}

// MarshalYAML implements yaml.Marshaler.
func (s *BeaconState) MarshalYAML() ([]byte, error) {
	blockRoots := make([]string, len(s.BlockRoots))
	for i := range s.BlockRoots {
		blockRoots[i] = fmt.Sprintf("%#x", s.BlockRoots[i])
	}
	stateRoots := make([]string, len(s.StateRoots))
	for i := range s.StateRoots {
		stateRoots[i] = fmt.Sprintf("%#x", s.StateRoots[i])
	}
	historicalRoots := make([]string, len(s.HistoricalRoots))
	for i := range s.HistoricalRoots {
		historicalRoots[i] = fmt.Sprintf("%#x", s.HistoricalRoots[i])
	}
	balances := make([]uint64, len(s.Balances))
	for i := range s.Balances {
		balances[i] = uint64(s.Balances[i])
	}
	randaoMixes := make([]string, len(s.RANDAOMixes))
	for i := range s.RANDAOMixes {
		randaoMixes[i] = fmt.Sprintf("%#x", s.RANDAOMixes[i])
	}
	slashings := make([]uint64, len(s.Slashings))
	for i := range s.Slashings {
		slashings[i] = uint64(s.Slashings[i])
	}
	PreviousEpochParticipation := make([]uint8, len(s.PreviousEpochParticipation))
	for i := range s.PreviousEpochParticipation {
		PreviousEpochParticipation[i] = uint8(s.PreviousEpochParticipation[i])
	}
	CurrentEpochParticipation := make([]uint8, len(s.CurrentEpochParticipation))
	for i := range s.CurrentEpochParticipation {
		CurrentEpochParticipation[i] = uint8(s.CurrentEpochParticipation[i])
	}
	yamlBytes, err := yaml.MarshalWithOptions(&beaconStateYAML{
		GenesisTime:                  s.GenesisTime,
		GenesisValidatorsRoot:        fmt.Sprintf("%#x", s.GenesisValidatorsRoot),
		Slot:                         uint64(s.Slot),
		Fork:                         s.Fork,
		LatestBlockHeader:            s.LatestBlockHeader,
		BlockRoots:                   blockRoots,
		StateRoots:                   stateRoots,
		HistoricalRoots:              historicalRoots,
		ETH1Data:                     s.ETH1Data,
		ETH1DataVotes:                s.ETH1DataVotes,
		ETH1DepositIndex:             s.ETH1DepositIndex,
		Validators:                   s.Validators,
		Balances:                     balances,
		RANDAOMixes:                  randaoMixes,
		Slashings:                    slashings,
		PreviousEpochParticipation:   PreviousEpochParticipation,
		CurrentEpochParticipation:    CurrentEpochParticipation,
		JustificationBits:            fmt.Sprintf("%#x", s.JustificationBits.Bytes()),
		PreviousJustifiedCheckpoint:  s.PreviousJustifiedCheckpoint,
		CurrentJustifiedCheckpoint:   s.CurrentJustifiedCheckpoint,
		FinalizedCheckpoint:          s.FinalizedCheckpoint,
		InactivityScores:             s.InactivityScores,
		CurrentSyncCommittee:         s.CurrentSyncCommittee,
		NextSyncCommittee:            s.NextSyncCommittee,
		LatestExecutionPayloadHeader: s.LatestExecutionPayloadHeader,
		NextWithdrawalIndex:          uint64(s.NextWithdrawalIndex),
		NextWithdrawalValidatorIndex: uint64(s.NextWithdrawalValidatorIndex),
		HistoricalSummaries:          s.HistoricalSummaries,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *BeaconState) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data beaconStateJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return s.unpack(&data)
}
