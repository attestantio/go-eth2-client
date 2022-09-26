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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	bitfield "github.com/prysmaticlabs/go-bitfield"
)

// BeaconState represents a beacon state.
type BeaconState struct {
	GenesisTime                         uint64
	GenesisValidatorsRoot               []byte `ssz-size:"32"`
	Slot                                uint64
	Fork                                *phase0.Fork
	LatestBlockHeader                   *phase0.BeaconBlockHeader
	BlockRoots                          [][]byte `ssz-size:"8192,32"`
	StateRoots                          [][]byte `ssz-size:"8192,32"`
	HistoricalRoots                     [][]byte `ssz-size:"?,32" ssz-max:"16777216"`
	ETH1Data                            *phase0.ETH1Data
	ETH1DataVotes                       []*phase0.ETH1Data `ssz-max:"2048"`
	ETH1DepositIndex                    uint64
	Validators                          []*phase0.Validator         `ssz-max:"1099511627776"`
	Balances                            []uint64                    `ssz-max:"1099511627776"`
	RANDAOMixes                         [][]byte                    `ssz-size:"65536,32"`
	Slashings                           []uint64                    `ssz-size:"8192"`
	PreviousEpochParticipation          []altair.ParticipationFlags `ssz-size:"1099511627776"`
	CurrentEpochParticipation           []altair.ParticipationFlags `ssz-size:"1099511627776"`
	JustificationBits                   bitfield.Bitvector4         `ssz-size:"1"`
	PreviousJustifiedCheckpoint         *phase0.Checkpoint
	CurrentJustifiedCheckpoint          *phase0.Checkpoint
	FinalizedCheckpoint                 *phase0.Checkpoint
	InactivityScores                    []uint64 `ssz-size:"1099511627776"`
	CurrentSyncCommittee                *altair.SyncCommittee
	NextSyncCommittee                   *altair.SyncCommittee
	LatestExecutionPayloadHeader        *ExecutionPayloadHeader
	WithdrawalQueue                     []*Withdrawal `ssz-size:"1099511627776"`
	NextWithdrawalIndex                 WithdrawalIndex
	NextPartialWithdrawalValidatorIndex phase0.ValidatorIndex
}

// beaconStateJSON is the spec representation of the struct.
type beaconStateJSON struct {
	GenesisTime                         string                    `json:"genesis_time"`
	GenesisValidatorsRoot               string                    `json:"genesis_validators_root"`
	Slot                                string                    `json:"slot"`
	Fork                                *phase0.Fork              `json:"fork"`
	LatestBlockHeader                   *phase0.BeaconBlockHeader `json:"latest_block_header"`
	BlockRoots                          []string                  `json:"block_roots"`
	StateRoots                          []string                  `json:"state_roots"`
	HistoricalRoots                     []string                  `json:"historical_roots"`
	ETH1Data                            *phase0.ETH1Data          `json:"eth1_data"`
	ETH1DataVotes                       []*phase0.ETH1Data        `json:"eth1_data_votes"`
	ETH1DepositIndex                    string                    `json:"eth1_deposit_index"`
	Validators                          []*phase0.Validator       `json:"validators"`
	Balances                            []string                  `json:"balances"`
	RANDAOMixes                         []string                  `json:"randao_mixes"`
	Slashings                           []string                  `json:"slashings"`
	PreviousEpochParticipation          []string                  `json:"previous_epoch_participation"`
	CurrentEpochParticipation           []string                  `json:"current_epoch_participation"`
	JustificationBits                   string                    `json:"justification_bits"`
	PreviousJustifiedCheckpoint         *phase0.Checkpoint        `json:"previous_justified_checkpoint"`
	CurrentJustifiedCheckpoint          *phase0.Checkpoint        `json:"current_justified_checkpoint"`
	FinalizedCheckpoint                 *phase0.Checkpoint        `json:"finalized_checkpoint"`
	InactivityScores                    []string                  `json:"inactivity_scores"`
	CurrentSyncCommittee                *altair.SyncCommittee     `json:"current_sync_committee"`
	NextSyncCommittee                   *altair.SyncCommittee     `json:"next_sync_committee"`
	LatestExecutionPayloadHeader        *ExecutionPayloadHeader   `json:"latest_execution_payload_header"`
	WithdrawalQueue                     []*Withdrawal             `json:"withdrawal_queue"`
	NextWithdrawalIndex                 string                    `json:"next_withdrawal_index"`
	NextPartialWithdrawalValidatorIndex string                    `json:"next_partial_withdrawal_validator_index"`
}

// MarshalJSON implements json.Marshaler.
func (s *BeaconState) MarshalJSON() ([]byte, error) {
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
	balances := make([]string, len(s.Balances))
	for i := range s.Balances {
		balances[i] = fmt.Sprintf("%d", s.Balances[i])
	}
	randaoMixes := make([]string, len(s.RANDAOMixes))
	for i := range s.RANDAOMixes {
		randaoMixes[i] = fmt.Sprintf("%#x", s.RANDAOMixes[i])
	}
	slashings := make([]string, len(s.Slashings))
	for i := range s.Slashings {
		slashings[i] = fmt.Sprintf("%d", s.Slashings[i])
	}
	PreviousEpochParticipation := make([]string, len(s.PreviousEpochParticipation))
	for i := range s.PreviousEpochParticipation {
		PreviousEpochParticipation[i] = fmt.Sprintf("%d", s.PreviousEpochParticipation[i])
	}
	CurrentEpochParticipation := make([]string, len(s.CurrentEpochParticipation))
	for i := range s.CurrentEpochParticipation {
		CurrentEpochParticipation[i] = fmt.Sprintf("%d", s.CurrentEpochParticipation[i])
	}
	inactivityScores := make([]string, len(s.InactivityScores))
	for i := range s.InactivityScores {
		inactivityScores[i] = fmt.Sprintf("%d", s.InactivityScores[i])
	}
	return json.Marshal(&beaconStateJSON{
		GenesisTime:                         fmt.Sprintf("%d", s.GenesisTime),
		GenesisValidatorsRoot:               fmt.Sprintf("%#x", s.GenesisValidatorsRoot),
		Slot:                                fmt.Sprintf("%d", s.Slot),
		Fork:                                s.Fork,
		LatestBlockHeader:                   s.LatestBlockHeader,
		BlockRoots:                          blockRoots,
		StateRoots:                          stateRoots,
		HistoricalRoots:                     historicalRoots,
		ETH1Data:                            s.ETH1Data,
		ETH1DataVotes:                       s.ETH1DataVotes,
		ETH1DepositIndex:                    fmt.Sprintf("%d", s.ETH1DepositIndex),
		Validators:                          s.Validators,
		Balances:                            balances,
		RANDAOMixes:                         randaoMixes,
		Slashings:                           slashings,
		PreviousEpochParticipation:          PreviousEpochParticipation,
		CurrentEpochParticipation:           CurrentEpochParticipation,
		JustificationBits:                   fmt.Sprintf("%#x", s.JustificationBits.Bytes()),
		PreviousJustifiedCheckpoint:         s.PreviousJustifiedCheckpoint,
		CurrentJustifiedCheckpoint:          s.CurrentJustifiedCheckpoint,
		FinalizedCheckpoint:                 s.FinalizedCheckpoint,
		InactivityScores:                    inactivityScores,
		CurrentSyncCommittee:                s.CurrentSyncCommittee,
		NextSyncCommittee:                   s.NextSyncCommittee,
		LatestExecutionPayloadHeader:        s.LatestExecutionPayloadHeader,
		WithdrawalQueue:                     s.WithdrawalQueue,
		NextWithdrawalIndex:                 fmt.Sprintf("%d", s.NextWithdrawalIndex),
		NextPartialWithdrawalValidatorIndex: fmt.Sprintf("%d", s.NextPartialWithdrawalValidatorIndex),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
// nolint:gocyclo
func (s *BeaconState) UnmarshalJSON(input []byte) error {
	var err error

	var data beaconStateJSON
	if err = json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if data.GenesisTime == "" {
		return errors.New("genesis time missing")
	}
	if s.GenesisTime, err = strconv.ParseUint(data.GenesisTime, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for genesis time")
	}
	if data.GenesisValidatorsRoot == "" {
		return errors.New("genesis validators root missing")
	}
	if s.GenesisValidatorsRoot, err = hex.DecodeString(strings.TrimPrefix(data.GenesisValidatorsRoot, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for genesis validators root")
	}
	if len(s.GenesisValidatorsRoot) != phase0.RootLength {
		return fmt.Errorf("incorrect length %d for genesis validators root", len(s.GenesisValidatorsRoot))
	}
	if data.Slot == "" {
		return errors.New("slot missing")
	}
	if s.Slot, err = strconv.ParseUint(data.Slot, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	if data.Fork == nil {
		return errors.New("fork missing")
	}
	s.Fork = data.Fork
	if data.LatestBlockHeader == nil {
		return errors.New("latest block header missing")
	}
	s.LatestBlockHeader = data.LatestBlockHeader
	if len(data.BlockRoots) == 0 {
		return errors.New("block roots missing")
	}
	s.BlockRoots = make([][]byte, len(data.BlockRoots))
	for i := range data.BlockRoots {
		if data.BlockRoots[i] == "" {
			return fmt.Errorf("block root %d missing", i)
		}
		if s.BlockRoots[i], err = hex.DecodeString(strings.TrimPrefix(data.BlockRoots[i], "0x")); err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for block root %d", i))
		}
		if len(s.BlockRoots[i]) != phase0.RootLength {
			return fmt.Errorf("incorrect length %d for block root %d", len(s.BlockRoots[i]), i)
		}
	}
	s.StateRoots = make([][]byte, len(data.StateRoots))
	for i := range data.StateRoots {
		if data.StateRoots[i] == "" {
			return fmt.Errorf("state root %d missing", i)
		}
		if s.StateRoots[i], err = hex.DecodeString(strings.TrimPrefix(data.StateRoots[i], "0x")); err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for state root %d", i))
		}
		if len(s.StateRoots[i]) != phase0.RootLength {
			return fmt.Errorf("incorrect length %d for state root %d", len(s.StateRoots[i]), i)
		}
	}
	s.HistoricalRoots = make([][]byte, len(data.HistoricalRoots))
	for i := range data.HistoricalRoots {
		if data.HistoricalRoots[i] == "" {
			return fmt.Errorf("historical root %d missing", i)
		}
		if s.HistoricalRoots[i], err = hex.DecodeString(strings.TrimPrefix(data.HistoricalRoots[i], "0x")); err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for historical root %d", i))
		}
		if len(s.HistoricalRoots[i]) != phase0.RootLength {
			return fmt.Errorf("incorrect length %d for historical root %d", len(s.HistoricalRoots[i]), i)
		}
	}
	if data.ETH1Data == nil {
		return errors.New("eth1 data missing")
	}
	s.ETH1Data = data.ETH1Data
	// ETH1DataVotes can be empty.
	s.ETH1DataVotes = data.ETH1DataVotes
	if data.Validators == nil {
		return errors.New("validators missing")
	}
	if data.ETH1DepositIndex == "" {
		return errors.New("eth1 deposit index missing")
	}
	if s.ETH1DepositIndex, err = strconv.ParseUint(data.ETH1DepositIndex, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for eth1 deposit index")
	}
	s.Validators = data.Validators
	s.Balances = make([]uint64, len(data.Balances))
	for i := range data.Balances {
		if data.Balances[i] == "" {
			return fmt.Errorf("balance %d missing", i)
		}
		if s.Balances[i], err = strconv.ParseUint(data.Balances[i], 10, 64); err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for balance %d", i))
		}
	}
	s.RANDAOMixes = make([][]byte, len(data.RANDAOMixes))
	for i := range data.RANDAOMixes {
		if data.RANDAOMixes[i] == "" {
			return fmt.Errorf("RANDAO mix %d missing", i)
		}
		if s.RANDAOMixes[i], err = hex.DecodeString(strings.TrimPrefix(data.RANDAOMixes[i], "0x")); err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for RANDAO mix %d", i))
		}
		if len(s.RANDAOMixes[i]) != phase0.RootLength {
			return fmt.Errorf("incorrect length %d for RANDAO mix %d", len(s.RANDAOMixes[i]), i)
		}
	}
	s.Slashings = make([]uint64, len(data.Slashings))
	for i := range data.Slashings {
		if data.Slashings[i] == "" {
			return fmt.Errorf("slashing %d missing", i)
		}
		if s.Slashings[i], err = strconv.ParseUint(data.Slashings[i], 10, 64); err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for slashing %d", i))
		}
	}
	s.PreviousEpochParticipation = make([]altair.ParticipationFlags, len(data.PreviousEpochParticipation))
	for i := range data.PreviousEpochParticipation {
		if data.PreviousEpochParticipation[i] == "" {
			return fmt.Errorf("previous epoch attestation %d missing", i)
		}
		previousEpochAttestation, err := strconv.ParseUint(data.PreviousEpochParticipation[i], 10, 8)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for previous epoch attestation %d", i))
		}
		s.PreviousEpochParticipation[i] = altair.ParticipationFlags(previousEpochAttestation)
	}
	s.CurrentEpochParticipation = make([]altair.ParticipationFlags, len(data.CurrentEpochParticipation))
	for i := range data.CurrentEpochParticipation {
		if data.CurrentEpochParticipation[i] == "" {
			return fmt.Errorf("current epoch attestation %d missing", i)
		}
		currentEpochAttestation, err := strconv.ParseUint(data.CurrentEpochParticipation[i], 10, 8)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for current epoch attestation %d", i))
		}
		s.CurrentEpochParticipation[i] = altair.ParticipationFlags(currentEpochAttestation)
	}
	if data.JustificationBits == "" {
		return errors.New("justification bits missing")
	}
	if s.JustificationBits, err = hex.DecodeString(strings.TrimPrefix(data.JustificationBits, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for justification bits")
	}
	if data.PreviousJustifiedCheckpoint == nil {
		return errors.New("previous justified checkpoint missing")
	}
	s.PreviousJustifiedCheckpoint = data.PreviousJustifiedCheckpoint
	if data.CurrentJustifiedCheckpoint == nil {
		return errors.New("current justified checkpoint missing")
	}
	s.CurrentJustifiedCheckpoint = data.CurrentJustifiedCheckpoint
	if data.FinalizedCheckpoint == nil {
		return errors.New("finalized checkpoint missing")
	}
	s.FinalizedCheckpoint = data.FinalizedCheckpoint
	s.InactivityScores = make([]uint64, len(data.InactivityScores))
	for i := range data.InactivityScores {
		if data.InactivityScores[i] == "" {
			return fmt.Errorf("inactivity score %d missing", i)
		}
		if s.InactivityScores[i], err = strconv.ParseUint(data.InactivityScores[i], 10, 64); err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for inactivity score %d", i))
		}
	}
	if data.CurrentSyncCommittee == nil {
		return errors.New("current sync committee missing")
	}
	s.CurrentSyncCommittee = data.CurrentSyncCommittee
	if data.NextSyncCommittee == nil {
		return errors.New("next sync committee missing")
	}
	s.NextSyncCommittee = data.NextSyncCommittee
	s.LatestExecutionPayloadHeader = data.LatestExecutionPayloadHeader
	s.WithdrawalQueue = data.WithdrawalQueue
	if data.NextWithdrawalIndex == "" {
		return errors.New("next withdrawal index missing")
	}
	nextWithdrawalIndex, err := strconv.ParseUint(data.NextWithdrawalIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for next withdrawal index")
	}
	s.NextWithdrawalIndex = WithdrawalIndex(nextWithdrawalIndex)
	if data.NextPartialWithdrawalValidatorIndex == "" {
		return errors.New("next partial validator validator index missing")
	}
	nextPartialWithdrawalValidatorIndex, err := strconv.ParseUint(data.NextPartialWithdrawalValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for next partial withdrawal validator index")
	}
	s.NextPartialWithdrawalValidatorIndex = phase0.ValidatorIndex(nextPartialWithdrawalValidatorIndex)

	return nil
}

// String returns a string version of the structure.
func (s *BeaconState) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
