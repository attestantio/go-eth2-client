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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// beaconStateJSON is the spec representation of the struct.
type beaconStateJSON struct {
	GenesisTime           string                    `json:"genesis_time"`
	GenesisValidatorsRoot phase0.Root               `json:"genesis_validators_root"`
	Slot                  phase0.Slot               `json:"slot"`
	Fork                  *phase0.Fork              `json:"fork"`
	LatestBlockHeader     *phase0.BeaconBlockHeader `json:"latest_block_header"`
	BlockRoots            []phase0.Root             `json:"block_roots"`
	StateRoots            []phase0.Root             `json:"state_roots"`
	HistoricalRoots       []phase0.Root             `json:"historical_roots"`
	ETH1Data              *phase0.ETH1Data          `json:"eth1_data"`
	//nolint:staticcheck
	ETH1DataVotes                []*phase0.ETH1Data           `json:"eth1_data_votes,allowempty"`
	ETH1DepositIndex             string                       `json:"eth1_deposit_index"`
	Validators                   []*phase0.Validator          `json:"validators"`
	Balances                     []string                     `json:"balances"`
	RANDAOMixes                  []string                     `json:"randao_mixes"`
	Slashings                    []string                     `json:"slashings"`
	PreviousEpochParticipation   []string                     `json:"previous_epoch_participation"`
	CurrentEpochParticipation    []string                     `json:"current_epoch_participation"`
	JustificationBits            string                       `json:"justification_bits"`
	PreviousJustifiedCheckpoint  *phase0.Checkpoint           `json:"previous_justified_checkpoint"`
	CurrentJustifiedCheckpoint   *phase0.Checkpoint           `json:"current_justified_checkpoint"`
	FinalizedCheckpoint          *phase0.Checkpoint           `json:"finalized_checkpoint"`
	InactivityScores             []string                     `json:"inactivity_scores"`
	CurrentSyncCommittee         *altair.SyncCommittee        `json:"current_sync_committee"`
	NextSyncCommittee            *altair.SyncCommittee        `json:"next_sync_committee"`
	LatestExecutionPayloadHeader *ExecutionPayloadHeader      `json:"latest_execution_payload_header"`
	NextWithdrawalIndex          string                       `json:"next_withdrawal_index"`
	NextWithdrawalValidatorIndex string                       `json:"next_withdrawal_validator_index"`
	HistoricalSummaries          []*capella.HistoricalSummary `json:"historical_summaries"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconState) MarshalJSON() ([]byte, error) {
	balances := make([]string, len(b.Balances))
	for i := range b.Balances {
		balances[i] = fmt.Sprintf("%d", b.Balances[i])
	}
	randaoMixes := make([]string, len(b.RANDAOMixes))
	for i := range b.RANDAOMixes {
		randaoMixes[i] = fmt.Sprintf("%#x", b.RANDAOMixes[i])
	}
	slashings := make([]string, len(b.Slashings))
	for i := range b.Slashings {
		slashings[i] = fmt.Sprintf("%d", b.Slashings[i])
	}
	PreviousEpochParticipation := make([]string, len(b.PreviousEpochParticipation))
	for i := range b.PreviousEpochParticipation {
		PreviousEpochParticipation[i] = fmt.Sprintf("%d", b.PreviousEpochParticipation[i])
	}
	CurrentEpochParticipation := make([]string, len(b.CurrentEpochParticipation))
	for i := range b.CurrentEpochParticipation {
		CurrentEpochParticipation[i] = fmt.Sprintf("%d", b.CurrentEpochParticipation[i])
	}
	inactivityScores := make([]string, len(b.InactivityScores))
	for i := range b.InactivityScores {
		inactivityScores[i] = strconv.FormatUint(b.InactivityScores[i], 10)
	}

	return json.Marshal(&beaconStateJSON{
		GenesisTime:                  strconv.FormatUint(b.GenesisTime, 10),
		GenesisValidatorsRoot:        b.GenesisValidatorsRoot,
		Slot:                         b.Slot,
		Fork:                         b.Fork,
		LatestBlockHeader:            b.LatestBlockHeader,
		BlockRoots:                   b.BlockRoots,
		StateRoots:                   b.StateRoots,
		HistoricalRoots:              b.HistoricalRoots,
		ETH1Data:                     b.ETH1Data,
		ETH1DataVotes:                b.ETH1DataVotes,
		ETH1DepositIndex:             strconv.FormatUint(b.ETH1DepositIndex, 10),
		Validators:                   b.Validators,
		Balances:                     balances,
		RANDAOMixes:                  randaoMixes,
		Slashings:                    slashings,
		PreviousEpochParticipation:   PreviousEpochParticipation,
		CurrentEpochParticipation:    CurrentEpochParticipation,
		JustificationBits:            fmt.Sprintf("%#x", b.JustificationBits.Bytes()),
		PreviousJustifiedCheckpoint:  b.PreviousJustifiedCheckpoint,
		CurrentJustifiedCheckpoint:   b.CurrentJustifiedCheckpoint,
		FinalizedCheckpoint:          b.FinalizedCheckpoint,
		InactivityScores:             inactivityScores,
		CurrentSyncCommittee:         b.CurrentSyncCommittee,
		NextSyncCommittee:            b.NextSyncCommittee,
		LatestExecutionPayloadHeader: b.LatestExecutionPayloadHeader,
		NextWithdrawalIndex:          fmt.Sprintf("%d", b.NextWithdrawalIndex),
		NextWithdrawalValidatorIndex: fmt.Sprintf("%d", b.NextWithdrawalValidatorIndex),
		HistoricalSummaries:          b.HistoricalSummaries,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
//
//nolint:gocyclo
func (b *BeaconState) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&beaconStateJSON{}, input)
	if err != nil {
		return err
	}

	genesisTime := string(bytes.Trim(raw["genesis_time"], `"`))
	if b.GenesisTime, err = strconv.ParseUint(genesisTime, 10, 64); err != nil {
		return errors.Wrap(err, "genesis_time")
	}

	if err := b.GenesisValidatorsRoot.UnmarshalJSON(raw["genesis_validators_root"]); err != nil {
		return errors.Wrap(err, "genesis_validators_root")
	}

	if err := b.Slot.UnmarshalJSON(raw["slot"]); err != nil {
		return errors.Wrap(err, "slot")
	}

	b.Fork = &phase0.Fork{}
	if err := b.Fork.UnmarshalJSON(raw["fork"]); err != nil {
		return errors.Wrap(err, "fork")
	}

	b.LatestBlockHeader = &phase0.BeaconBlockHeader{}
	if err := b.LatestBlockHeader.UnmarshalJSON(raw["latest_block_header"]); err != nil {
		return errors.Wrap(err, "latest_block_header")
	}

	if err := json.Unmarshal(raw["block_roots"], &b.BlockRoots); err != nil {
		return errors.Wrap(err, "block_roots")
	}

	if err := json.Unmarshal(raw["state_roots"], &b.StateRoots); err != nil {
		return errors.Wrap(err, "state_roots")
	}

	if err := json.Unmarshal(raw["historical_roots"], &b.HistoricalRoots); err != nil {
		return errors.Wrap(err, "historical_roots")
	}

	b.ETH1Data = &phase0.ETH1Data{}
	if err := b.ETH1Data.UnmarshalJSON(raw["eth1_data"]); err != nil {
		return errors.Wrap(err, "eth1_data")
	}

	if err := json.Unmarshal(raw["eth1_data_votes"], &b.ETH1DataVotes); err != nil {
		return errors.Wrap(err, "eth1_data_votes")
	}

	eth1DepositIndex := string(bytes.Trim(raw["eth1_deposit_index"], `"`))
	if b.ETH1DepositIndex, err = strconv.ParseUint(eth1DepositIndex, 10, 64); err != nil {
		return errors.Wrap(err, "eth1_deposit_index")
	}

	if err := json.Unmarshal(raw["validators"], &b.Validators); err != nil {
		return errors.Wrap(err, "validators")
	}

	if err := json.Unmarshal(raw["balances"], &b.Balances); err != nil {
		return errors.Wrap(err, "balances")
	}

	if err := json.Unmarshal(raw["randao_mixes"], &b.RANDAOMixes); err != nil {
		return errors.Wrap(err, "randao_mixes")
	}

	if err := json.Unmarshal(raw["slashings"], &b.Slashings); err != nil {
		return errors.Wrap(err, "slashings")
	}

	if err := json.Unmarshal(raw["previous_epoch_participation"], &b.PreviousEpochParticipation); err != nil {
		return errors.Wrap(err, "previous_epoch_participation")
	}

	if err := json.Unmarshal(raw["current_epoch_participation"], &b.CurrentEpochParticipation); err != nil {
		return errors.Wrap(err, "current_epoch_participation")
	}

	justificationBits := string(bytes.TrimPrefix(bytes.Trim(raw["justification_bits"], `"`), []byte{'0', 'x'}))
	if b.JustificationBits, err = hex.DecodeString(justificationBits); err != nil {
		return errors.Wrap(err, "justification_bits")
	}

	b.PreviousJustifiedCheckpoint = &phase0.Checkpoint{}
	if err := b.PreviousJustifiedCheckpoint.UnmarshalJSON(raw["previous_justified_checkpoint"]); err != nil {
		return errors.Wrap(err, "previous_justified_checkpoint")
	}

	b.CurrentJustifiedCheckpoint = &phase0.Checkpoint{}
	if err := b.CurrentJustifiedCheckpoint.UnmarshalJSON(raw["current_justified_checkpoint"]); err != nil {
		return errors.Wrap(err, "current_justified_checkpoint")
	}

	b.FinalizedCheckpoint = &phase0.Checkpoint{}
	if err := b.FinalizedCheckpoint.UnmarshalJSON(raw["finalized_checkpoint"]); err != nil {
		return errors.Wrap(err, "finalized_checkpoint")
	}

	inactivityScores := make([]string, 0)
	if err := json.Unmarshal(raw["inactivity_scores"], &inactivityScores); err != nil {
		return errors.Wrap(err, "inactivity_scores")
	}
	b.InactivityScores = make([]uint64, len(inactivityScores))
	for i := range inactivityScores {
		if inactivityScores[i] == "" {
			return fmt.Errorf("inactivity score %d missing", i)
		}
		if b.InactivityScores[i], err = strconv.ParseUint(inactivityScores[i], 10, 64); err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for inactivity score %d", i))
		}
	}

	b.CurrentSyncCommittee = &altair.SyncCommittee{}
	if err := b.CurrentSyncCommittee.UnmarshalJSON(raw["current_sync_committee"]); err != nil {
		return errors.Wrap(err, "current_sync_committee")
	}

	b.NextSyncCommittee = &altair.SyncCommittee{}
	if err := b.NextSyncCommittee.UnmarshalJSON(raw["next_sync_committee"]); err != nil {
		return errors.Wrap(err, "next_sync_committee")
	}

	b.LatestExecutionPayloadHeader = &ExecutionPayloadHeader{}
	if err := b.LatestExecutionPayloadHeader.UnmarshalJSON(raw["latest_execution_payload_header"]); err != nil {
		return errors.Wrap(err, "latest_execution_payload_header")
	}

	if err := b.NextWithdrawalIndex.UnmarshalJSON(raw["next_withdrawal_index"]); err != nil {
		return errors.Wrap(err, "next_withdrawal_index")
	}

	if err := b.NextWithdrawalValidatorIndex.UnmarshalJSON(raw["next_withdrawal_validator_index"]); err != nil {
		return errors.Wrap(err, "next_withdrawal_validator_index")
	}

	if err := json.Unmarshal(raw["historical_summaries"], &b.HistoricalSummaries); err != nil {
		return errors.Wrap(err, "historical_summaries")
	}

	return nil
}
