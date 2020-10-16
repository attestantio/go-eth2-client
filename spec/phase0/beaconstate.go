// Copyright © 2020 Attestant Limited.
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
	"encoding/json"
	"fmt"

	bitfield "github.com/prysmaticlabs/go-bitfield"
)

// BeaconState represents a beacon state.
type BeaconState struct {
	GenesisTime                 uint64
	GenesisValidatorsRoot       []byte `ssz-size:"32"`
	Slot                        uint64
	Fork                        *Fork
	LatestBlockHeader           *BeaconBlockHeader
	BlockRoots                  [][]byte `ssz-size:"8192,32"`
	StateRoots                  [][]byte `ssz-size:"8192,32"`
	HistoricalRoots             [][]byte `ssz-size:"?,32" ssz-max:"16777216"`
	ETH1Data                    *ETH1Data
	ETH1DataVotes               []*ETH1Data           `ssz-max:"1024"` // TODO should be 2048 for mainnet?
	Validators                  []*Validator          `ssz-max:"1099511627776"`
	Balances                    []uint64              `ssz-max:"1099511627776"`
	RANDAOMixes                 [][]byte              `ssz-size:"65536,32"`
	Slashings                   []uint64              `ssz-size:"8192"`
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
func (b *BeaconState) MarshalJSON() ([]byte, error) {
	blockRoots := make([]string, len(b.BlockRoots))
	for i := range b.BlockRoots {
		blockRoots[i] = fmt.Sprintf("%#x", b.BlockRoots[i])
	}
	stateRoots := make([]string, len(b.StateRoots))
	for i := range b.StateRoots {
		stateRoots[i] = fmt.Sprintf("%#x", b.StateRoots[i])
	}
	historicalRoots := make([]string, len(b.HistoricalRoots))
	for i := range b.HistoricalRoots {
		historicalRoots[i] = fmt.Sprintf("%#x", b.HistoricalRoots[i])
	}
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
	return json.Marshal(&beaconStateJSON{
		GenesisTime:                 fmt.Sprintf("%d", b.GenesisTime),
		GenesisValidatorsRoot:       fmt.Sprintf("%#x", b.GenesisValidatorsRoot),
		Slot:                        fmt.Sprintf("%d", b.Slot),
		Fork:                        b.Fork,
		LatestBlockHeader:           b.LatestBlockHeader,
		BlockRoots:                  blockRoots,
		StateRoots:                  stateRoots,
		HistoricalRoots:             historicalRoots,
		ETH1Data:                    b.ETH1Data,
		ETH1DataVotes:               b.ETH1DataVotes,
		Validators:                  b.Validators,
		Balances:                    balances,
		RANDAOMixes:                 randaoMixes,
		Slashings:                   slashings,
		PreviousEpochAttestations:   b.PreviousEpochAttestations,
		CurrentEpochAttestations:    b.CurrentEpochAttestations,
		JustificationBits:           fmt.Sprintf("%#x", b.JustificationBits.Bytes()),
		PreviousJustifiedCheckpoint: b.PreviousJustifiedCheckpoint,
		CurrentJustifiedCheckpoint:  b.CurrentJustifiedCheckpoint,
		FinalizedCheckpoint:         b.FinalizedCheckpoint,
	})
}

//// UnmarshalJSON implements json.Unmarshaler.
//func (b *BeaconBlock) UnmarshalJSON(input []byte) error {
//	var err error
//
//	var beaconBlockJSON beaconBlockJSON
//	if err = json.Unmarshal(input, &beaconBlockJSON); err != nil {
//		return errors.Wrap(err, "invalid JSON")
//	}
//	if beaconBlockJSON.Slot == "" {
//		return errors.New("slot missing")
//	}
//	if b.Slot, err = strconv.ParseUint(beaconBlockJSON.Slot, 10, 64); err != nil {
//		return errors.Wrap(err, "invalid value for slot")
//	}
//	if beaconBlockJSON.ProposerIndex == "" {
//		return errors.New("proposer index missing")
//	}
//	if b.ProposerIndex, err = strconv.ParseUint(beaconBlockJSON.ProposerIndex, 10, 64); err != nil {
//		return errors.Wrap(err, "invalid value for proposer index")
//	}
//	if beaconBlockJSON.ParentRoot == "" {
//		return errors.New("parent root missing")
//	}
//	if b.ParentRoot, err = hex.DecodeString(strings.TrimPrefix(beaconBlockJSON.ParentRoot, "0x")); err != nil {
//		return errors.Wrap(err, "invalid value for parent root")
//	}
//	if len(b.ParentRoot) != rootLength {
//		return errors.New("incorrect length for parent root")
//	}
//	if beaconBlockJSON.StateRoot == "" {
//		return errors.New("state root missing")
//	}
//	if b.StateRoot, err = hex.DecodeString(strings.TrimPrefix(beaconBlockJSON.StateRoot, "0x")); err != nil {
//		return errors.Wrap(err, "invalid value for state root")
//	}
//	if len(b.StateRoot) != rootLength {
//		return errors.New("incorrect length for state root")
//	}
//	if beaconBlockJSON.Body == nil {
//		return errors.New("body missing")
//	}
//	b.Body = beaconBlockJSON.Body
//
//	return nil
//}
//
//// String returns a string version of the structure.
//func (b *BeaconBlock) String() string {
//	data, err := json.Marshal(b)
//	if err != nil {
//		return fmt.Sprintf("ERR: %v", err)
//	}
//	return string(data)
//}
