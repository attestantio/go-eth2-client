// Copyright © 2020, 2023 Attestant Limited.
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

package phase0

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	bitfield "github.com/prysmaticlabs/go-bitfield"
)

// BeaconState represents a beacon state.
type BeaconState struct {
	GenesisTime                 uint64
	GenesisValidatorsRoot       Root `ssz-size:"32"`
	Slot                        Slot
	Fork                        *Fork
	LatestBlockHeader           *BeaconBlockHeader
	BlockRoots                  []Root `ssz-size:"8192,32"`
	StateRoots                  []Root `ssz-size:"8192,32"`
	HistoricalRoots             []Root `ssz-size:"?,32" ssz-max:"16777216"`
	ETH1Data                    *ETH1Data
	ETH1DataVotes               []*ETH1Data `ssz-max:"1024"` // Should be 2048 for mainnet?
	ETH1DepositIndex            uint64
	Validators                  []*Validator          `ssz-max:"1099511627776"`
	Balances                    []Gwei                `ssz-max:"1099511627776"`
	RANDAOMixes                 []Root                `ssz-size:"65536,32"`
	Slashings                   []Gwei                `ssz-size:"8192"`
	PreviousEpochAttestations   []*PendingAttestation `ssz-max:"4096"`
	CurrentEpochAttestations    []*PendingAttestation `ssz-max:"4096"`
	JustificationBits           bitfield.Bitvector4   `ssz-size:"1"`
	PreviousJustifiedCheckpoint *Checkpoint
	CurrentJustifiedCheckpoint  *Checkpoint
	FinalizedCheckpoint         *Checkpoint
}

// beaconStateJSON is the spec representation of the struct.
type beaconStateJSON struct {
	GenesisTime                 string                `json:"genesis_time"`
	GenesisValidatorsRoot       string                `json:"genesis_validators_root"`
	Slot                        string                `json:"slot"`
	Fork                        *Fork                 `json:"fork"`
	LatestBlockHeader           *BeaconBlockHeader    `json:"latest_block_header"`
	BlockRoots                  []string              `json:"block_roots"`
	StateRoots                  []string              `json:"state_roots"`
	HistoricalRoots             []string              `json:"historical_roots"`
	ETH1Data                    *ETH1Data             `json:"eth1_data"`
	ETH1DataVotes               []*ETH1Data           `json:"eth1_data_votes"`
	ETH1DepositIndex            string                `json:"eth1_deposit_index"`
	Validators                  []*Validator          `json:"validators"`
	Balances                    []string              `json:"balances"`
	RANDAOMixes                 []string              `json:"randao_mixes"`
	Slashings                   []string              `json:"slashings"`
	PreviousEpochAttestations   []*PendingAttestation `json:"previous_epoch_attestations"`
	CurrentEpochAttestations    []*PendingAttestation `json:"current_epoch_attestations"`
	JustificationBits           string                `json:"justification_bits"`
	PreviousJustifiedCheckpoint *Checkpoint           `json:"previous_justified_checkpoint"`
	CurrentJustifiedCheckpoint  *Checkpoint           `json:"current_justified_checkpoint"`
	FinalizedCheckpoint         *Checkpoint           `json:"finalized_checkpoint"`
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
	return json.Marshal(&beaconStateJSON{
		GenesisTime:                 fmt.Sprintf("%d", s.GenesisTime),
		GenesisValidatorsRoot:       fmt.Sprintf("%#x", s.GenesisValidatorsRoot),
		Slot:                        fmt.Sprintf("%d", s.Slot),
		Fork:                        s.Fork,
		LatestBlockHeader:           s.LatestBlockHeader,
		BlockRoots:                  blockRoots,
		StateRoots:                  stateRoots,
		HistoricalRoots:             historicalRoots,
		ETH1Data:                    s.ETH1Data,
		ETH1DataVotes:               s.ETH1DataVotes,
		Validators:                  s.Validators,
		Balances:                    balances,
		RANDAOMixes:                 randaoMixes,
		Slashings:                   slashings,
		PreviousEpochAttestations:   s.PreviousEpochAttestations,
		CurrentEpochAttestations:    s.CurrentEpochAttestations,
		JustificationBits:           fmt.Sprintf("%#x", s.JustificationBits.Bytes()),
		PreviousJustifiedCheckpoint: s.PreviousJustifiedCheckpoint,
		CurrentJustifiedCheckpoint:  s.CurrentJustifiedCheckpoint,
		FinalizedCheckpoint:         s.FinalizedCheckpoint,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
// nolint:gocyclo
func (s *BeaconState) UnmarshalJSON(input []byte) error {
	var err error

	var beaconStateJSON beaconStateJSON
	if err = json.Unmarshal(input, &beaconStateJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if beaconStateJSON.GenesisTime == "" {
		return errors.New("genesis time missing")
	}
	if s.GenesisTime, err = strconv.ParseUint(beaconStateJSON.GenesisTime, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for genesis time")
	}
	if beaconStateJSON.GenesisValidatorsRoot == "" {
		return errors.New("genesis validators root missing")
	}
	genesisValidatorsRoot, err := hex.DecodeString(strings.TrimPrefix(beaconStateJSON.GenesisValidatorsRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for genesis validators root")
	}
	if len(genesisValidatorsRoot) != RootLength {
		return fmt.Errorf("incorrect length %d for genesis validators root", len(genesisValidatorsRoot))
	}
	copy(s.GenesisValidatorsRoot[:], genesisValidatorsRoot)
	if beaconStateJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(beaconStateJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	s.Slot = Slot(slot)
	if beaconStateJSON.Fork == nil {
		return errors.New("fork missing")
	}
	s.Fork = beaconStateJSON.Fork
	if beaconStateJSON.LatestBlockHeader == nil {
		return errors.New("latest block header missing")
	}
	s.LatestBlockHeader = beaconStateJSON.LatestBlockHeader
	if len(beaconStateJSON.BlockRoots) == 0 {
		return errors.New("block roots missing")
	}
	s.BlockRoots = make([]Root, len(beaconStateJSON.BlockRoots))
	for i := range beaconStateJSON.BlockRoots {
		if beaconStateJSON.BlockRoots[i] == "" {
			return fmt.Errorf("block root %d missing", i)
		}
		blockRoot, err := hex.DecodeString(strings.TrimPrefix(beaconStateJSON.BlockRoots[i], "0x"))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for block root %d", i))
		}
		if len(blockRoot) != RootLength {
			return fmt.Errorf("incorrect length %d for block root %d", len(blockRoot), i)
		}
		copy(s.BlockRoots[i][:], blockRoot)
	}
	s.StateRoots = make([]Root, len(beaconStateJSON.StateRoots))
	for i := range beaconStateJSON.StateRoots {
		if beaconStateJSON.StateRoots[i] == "" {
			return fmt.Errorf("state root %d missing", i)
		}
		stateRoot, err := hex.DecodeString(strings.TrimPrefix(beaconStateJSON.StateRoots[i], "0x"))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for state root %d", i))
		}
		if len(stateRoot) != RootLength {
			return fmt.Errorf("incorrect length %d for state root %d", len(stateRoot), i)
		}
		copy(s.StateRoots[i][:], stateRoot)
	}
	s.HistoricalRoots = make([]Root, len(beaconStateJSON.HistoricalRoots))
	for i := range beaconStateJSON.HistoricalRoots {
		if beaconStateJSON.HistoricalRoots[i] == "" {
			return fmt.Errorf("historical root %d missing", i)
		}
		historicalRoot, err := hex.DecodeString(strings.TrimPrefix(beaconStateJSON.HistoricalRoots[i], "0x"))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for historical root %d", i))
		}
		if len(historicalRoot) != RootLength {
			return fmt.Errorf("incorrect length %d for historical root %d", len(historicalRoot), i)
		}
		copy(s.HistoricalRoots[i][:], historicalRoot)
	}
	if beaconStateJSON.ETH1Data == nil {
		return errors.New("eth1 data missing")
	}
	s.ETH1Data = beaconStateJSON.ETH1Data
	// ETH1DataVotes can be empty.
	s.ETH1DataVotes = beaconStateJSON.ETH1DataVotes
	if beaconStateJSON.Validators == nil {
		return errors.New("validators missing")
	}
	if beaconStateJSON.ETH1DepositIndex == "" {
		return errors.New("eth1 deposit index missing")
	}
	if s.ETH1DepositIndex, err = strconv.ParseUint(beaconStateJSON.ETH1DepositIndex, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for eth1 deposit index")
	}
	s.Validators = beaconStateJSON.Validators
	s.Balances = make([]Gwei, len(beaconStateJSON.Balances))
	for i := range beaconStateJSON.Balances {
		if beaconStateJSON.Balances[i] == "" {
			return fmt.Errorf("balance %d missing", i)
		}
		balance, err := strconv.ParseUint(beaconStateJSON.Balances[i], 10, 64)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for balance %d", i))
		}
		s.Balances[i] = Gwei(balance)
	}
	s.RANDAOMixes = make([]Root, len(beaconStateJSON.RANDAOMixes))
	for i := range beaconStateJSON.RANDAOMixes {
		if beaconStateJSON.RANDAOMixes[i] == "" {
			return fmt.Errorf("RANDAO mix %d missing", i)
		}
		randaoMix, err := hex.DecodeString(strings.TrimPrefix(beaconStateJSON.RANDAOMixes[i], "0x"))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for RANDAO mix %d", i))
		}
		if len(randaoMix) != RootLength {
			return fmt.Errorf("incorrect length %d for RANDAO mix %d", len(randaoMix), i)
		}
		copy(s.RANDAOMixes[i][:], randaoMix)
	}
	s.Slashings = make([]Gwei, len(beaconStateJSON.Slashings))
	for i := range beaconStateJSON.Slashings {
		if beaconStateJSON.Slashings[i] == "" {
			return fmt.Errorf("slashing %d missing", i)
		}
		slashings, err := strconv.ParseUint(beaconStateJSON.Slashings[i], 10, 64)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for slashing %d", i))
		}
		s.Slashings[i] = Gwei(slashings)
	}
	s.PreviousEpochAttestations = beaconStateJSON.PreviousEpochAttestations
	s.CurrentEpochAttestations = beaconStateJSON.CurrentEpochAttestations
	if beaconStateJSON.JustificationBits == "" {
		return errors.New("justification bits missing")
	}
	if s.JustificationBits, err = hex.DecodeString(strings.TrimPrefix(beaconStateJSON.JustificationBits, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for justification bits")
	}
	if beaconStateJSON.PreviousJustifiedCheckpoint == nil {
		return errors.New("previous justified checkpoint missing")
	}
	s.PreviousJustifiedCheckpoint = beaconStateJSON.PreviousJustifiedCheckpoint
	if beaconStateJSON.CurrentJustifiedCheckpoint == nil {
		return errors.New("current justified checkpoint missing")
	}
	s.CurrentJustifiedCheckpoint = beaconStateJSON.CurrentJustifiedCheckpoint
	if beaconStateJSON.FinalizedCheckpoint == nil {
		return errors.New("finalized checkpoint missing")
	}
	s.FinalizedCheckpoint = beaconStateJSON.FinalizedCheckpoint

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
