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
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
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
	HistoricalRoots             []Root `ssz-max:"16777216" ssz-size:"?,32"`
	ETH1Data                    *ETH1Data
	ETH1DataVotes               []*ETH1Data `ssz-max:"2048"`
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

// beaconStateYAML is the spec representation of the struct.
type beaconStateYAML struct {
	GenesisTime                 uint64                `json:"genesis_time"`
	GenesisValidatorsRoot       Root                  `json:"genesis_validators_root"`
	Slot                        uint64                `json:"slot"`
	Fork                        *Fork                 `json:"fork"`
	LatestBlockHeader           *BeaconBlockHeader    `json:"latest_block_header"`
	BlockRoots                  []Root                `json:"block_roots"`
	StateRoots                  []Root                `json:"state_roots"`
	HistoricalRoots             []Root                `json:"historical_roots"`
	ETH1Data                    *ETH1Data             `json:"eth1_data"`
	ETH1DataVotes               []*ETH1Data           `json:"eth1_data_votes"`
	ETH1DepositIndex            uint64                `json:"eth1_deposit_index"`
	Validators                  []*Validator          `json:"validators"`
	Balances                    []Gwei                `json:"balances"`
	RANDAOMixes                 []Root                `json:"randao_mixes"`
	Slashings                   []Gwei                `json:"slashings"`
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
		GenesisTime:                 strconv.FormatUint(s.GenesisTime, 10),
		GenesisValidatorsRoot:       fmt.Sprintf("%#x", s.GenesisValidatorsRoot),
		Slot:                        fmt.Sprintf("%d", s.Slot),
		Fork:                        s.Fork,
		LatestBlockHeader:           s.LatestBlockHeader,
		BlockRoots:                  blockRoots,
		StateRoots:                  stateRoots,
		HistoricalRoots:             historicalRoots,
		ETH1Data:                    s.ETH1Data,
		ETH1DataVotes:               s.ETH1DataVotes,
		ETH1DepositIndex:            strconv.FormatUint(s.ETH1DepositIndex, 10),
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
func (s *BeaconState) UnmarshalJSON(input []byte) error {
	var err error

	var data beaconStateJSON
	if err = json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return s.unpack(&data)
}

//nolint:gocyclo
func (s *BeaconState) unpack(data *beaconStateJSON) error {
	var err error

	if data.GenesisTime == "" {
		return errors.New("genesis time missing")
	}
	if s.GenesisTime, err = strconv.ParseUint(data.GenesisTime, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for genesis time")
	}
	if data.GenesisValidatorsRoot == "" {
		return errors.New("genesis validators root missing")
	}
	genesisValidatorsRoot, err := hex.DecodeString(strings.TrimPrefix(data.GenesisValidatorsRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for genesis validators root")
	}
	if len(genesisValidatorsRoot) != RootLength {
		return fmt.Errorf("incorrect length %d for genesis validators root", len(genesisValidatorsRoot))
	}
	copy(s.GenesisValidatorsRoot[:], genesisValidatorsRoot)
	if data.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	s.Slot = Slot(slot)
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
	s.BlockRoots = make([]Root, len(data.BlockRoots))
	for i := range data.BlockRoots {
		if data.BlockRoots[i] == "" {
			return fmt.Errorf("block root %d missing", i)
		}
		blockRoot, err := hex.DecodeString(strings.TrimPrefix(data.BlockRoots[i], "0x"))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for block root %d", i))
		}
		if len(blockRoot) != RootLength {
			return fmt.Errorf("incorrect length %d for block root %d", len(blockRoot), i)
		}
		copy(s.BlockRoots[i][:], blockRoot)
	}
	s.StateRoots = make([]Root, len(data.StateRoots))
	for i := range data.StateRoots {
		if data.StateRoots[i] == "" {
			return fmt.Errorf("state root %d missing", i)
		}
		stateRoot, err := hex.DecodeString(strings.TrimPrefix(data.StateRoots[i], "0x"))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for state root %d", i))
		}
		if len(stateRoot) != RootLength {
			return fmt.Errorf("incorrect length %d for state root %d", len(stateRoot), i)
		}
		copy(s.StateRoots[i][:], stateRoot)
	}
	s.HistoricalRoots = make([]Root, len(data.HistoricalRoots))
	for i := range data.HistoricalRoots {
		if data.HistoricalRoots[i] == "" {
			return fmt.Errorf("historical root %d missing", i)
		}
		historicalRoot, err := hex.DecodeString(strings.TrimPrefix(data.HistoricalRoots[i], "0x"))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for historical root %d", i))
		}
		if len(historicalRoot) != RootLength {
			return fmt.Errorf("incorrect length %d for historical root %d", len(historicalRoot), i)
		}
		copy(s.HistoricalRoots[i][:], historicalRoot)
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
	s.Balances = make([]Gwei, len(data.Balances))
	for i := range data.Balances {
		if data.Balances[i] == "" {
			return fmt.Errorf("balance %d missing", i)
		}
		balance, err := strconv.ParseUint(data.Balances[i], 10, 64)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for balance %d", i))
		}
		s.Balances[i] = Gwei(balance)
	}
	s.RANDAOMixes = make([]Root, len(data.RANDAOMixes))
	for i := range data.RANDAOMixes {
		if data.RANDAOMixes[i] == "" {
			return fmt.Errorf("RANDAO mix %d missing", i)
		}
		randaoMix, err := hex.DecodeString(strings.TrimPrefix(data.RANDAOMixes[i], "0x"))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for RANDAO mix %d", i))
		}
		if len(randaoMix) != RootLength {
			return fmt.Errorf("incorrect length %d for RANDAO mix %d", len(randaoMix), i)
		}
		copy(s.RANDAOMixes[i][:], randaoMix)
	}
	s.Slashings = make([]Gwei, len(data.Slashings))
	for i := range data.Slashings {
		if data.Slashings[i] == "" {
			return fmt.Errorf("slashing %d missing", i)
		}
		slashings, err := strconv.ParseUint(data.Slashings[i], 10, 64)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for slashing %d", i))
		}
		s.Slashings[i] = Gwei(slashings)
	}
	s.PreviousEpochAttestations = data.PreviousEpochAttestations
	s.CurrentEpochAttestations = data.CurrentEpochAttestations
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

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (s *BeaconState) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&beaconStateYAML{
		GenesisTime:                 s.GenesisTime,
		GenesisValidatorsRoot:       s.GenesisValidatorsRoot,
		Slot:                        uint64(s.Slot),
		Fork:                        s.Fork,
		LatestBlockHeader:           s.LatestBlockHeader,
		BlockRoots:                  s.BlockRoots,
		StateRoots:                  s.StateRoots,
		HistoricalRoots:             s.HistoricalRoots,
		ETH1Data:                    s.ETH1Data,
		ETH1DataVotes:               s.ETH1DataVotes,
		ETH1DepositIndex:            s.ETH1DepositIndex,
		Validators:                  s.Validators,
		Balances:                    s.Balances,
		RANDAOMixes:                 s.RANDAOMixes,
		Slashings:                   s.Slashings,
		PreviousEpochAttestations:   s.PreviousEpochAttestations,
		CurrentEpochAttestations:    s.CurrentEpochAttestations,
		JustificationBits:           fmt.Sprintf("%#x", s.JustificationBits.Bytes()),
		PreviousJustifiedCheckpoint: s.PreviousJustifiedCheckpoint,
		CurrentJustifiedCheckpoint:  s.CurrentJustifiedCheckpoint,
		FinalizedCheckpoint:         s.FinalizedCheckpoint,
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

// String returns a string version of the structure.
func (s *BeaconState) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
