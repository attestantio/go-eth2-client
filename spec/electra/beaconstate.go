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

package electra

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	bitfield "github.com/prysmaticlabs/go-bitfield"
)

// BeaconState represents a beacon state.
type BeaconState struct {
	GenesisTime                   uint64
	GenesisValidatorsRoot         phase0.Root `ssz-size:"32"`
	Slot                          phase0.Slot
	Fork                          *phase0.Fork
	LatestBlockHeader             *phase0.BeaconBlockHeader
	BlockRoots                    []phase0.Root `dynssz-size:"SLOTS_PER_HISTORICAL_ROOT,32" ssz-size:"8192,32"`
	StateRoots                    []phase0.Root `dynssz-size:"SLOTS_PER_HISTORICAL_ROOT,32" ssz-size:"8192,32"`
	HistoricalRoots               []phase0.Root `ssz-max:"16777216"                         ssz-size:"?,32"`
	ETH1Data                      *phase0.ETH1Data
	ETH1DataVotes                 []*phase0.ETH1Data `ssz-max:"2048"`
	ETH1DepositIndex              uint64
	Validators                    []*phase0.Validator         `ssz-max:"1099511627776"`
	Balances                      []phase0.Gwei               `ssz-max:"1099511627776"`
	RANDAOMixes                   []phase0.Root               `dynssz-size:"EPOCHS_PER_HISTORICAL_VECTOR,32" ssz-size:"65536,32"`
	Slashings                     []phase0.Gwei               `dynssz-size:"EPOCHS_PER_SLASHINGS_VECTOR"     ssz-size:"8192"`
	PreviousEpochParticipation    []altair.ParticipationFlags `ssz-max:"1099511627776"`
	CurrentEpochParticipation     []altair.ParticipationFlags `ssz-max:"1099511627776"`
	JustificationBits             bitfield.Bitvector4         `ssz-size:"1"`
	PreviousJustifiedCheckpoint   *phase0.Checkpoint
	CurrentJustifiedCheckpoint    *phase0.Checkpoint
	FinalizedCheckpoint           *phase0.Checkpoint
	InactivityScores              []uint64 `ssz-max:"1099511627776"`
	CurrentSyncCommittee          *altair.SyncCommittee
	NextSyncCommittee             *altair.SyncCommittee
	LatestExecutionPayloadHeader  *deneb.ExecutionPayloadHeader
	NextWithdrawalIndex           capella.WithdrawalIndex
	NextWithdrawalValidatorIndex  phase0.ValidatorIndex
	HistoricalSummaries           []*capella.HistoricalSummary `ssz-max:"16777216"`
	DepositRequestsStartIndex     uint64
	DepositBalanceToConsume       phase0.Gwei
	ExitBalanceToConsume          phase0.Gwei
	EarliestExitEpoch             phase0.Epoch
	ConsolidationBalanceToConsume phase0.Gwei
	EarliestConsolidationEpoch    phase0.Epoch
	PendingDeposits               []*PendingDeposit           `ssz-max:"134217728"`
	PendingPartialWithdrawals     []*PendingPartialWithdrawal `ssz-max:"134217728"`
	PendingConsolidations         []*PendingConsolidation     `ssz-max:"262144"`
}

// String returns a string version of the structure.
func (b *BeaconState) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
