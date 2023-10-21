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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// ForkChoice is the data regarding the node's current fork choice context.
type ForkChoice struct {
	// JustifiedCheckpoint is the current justified checkpoint.
	JustifiedCheckpoint phase0.Checkpoint
	// FInalizedCheckpoint is the current finalized checkpoint.
	FinalizedCheckpoint phase0.Checkpoint
	// ForkChoiceNodes contains the fork choice nodes.
	ForkChoiceNodes []*ForkChoiceNode
}

// MarshalJSON implements json.Marshaler.
func (f *ForkChoice) MarshalJSON() ([]byte, error) {
	return json.Marshal(&forkChoiceJSON{
		JustifiedCheckpoint: &f.JustifiedCheckpoint,
		FinalizedCheckpoint: &f.FinalizedCheckpoint,
		ForkChoiceNodes:     f.ForkChoiceNodes,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *ForkChoice) UnmarshalJSON(input []byte) error {
	var err error

	var forkChoiceJSON forkChoiceJSON
	if err = json.Unmarshal(input, &forkChoiceJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if forkChoiceJSON.JustifiedCheckpoint == nil {
		return errors.New("justified checkpoint missing")
	}
	f.JustifiedCheckpoint = *forkChoiceJSON.JustifiedCheckpoint

	if forkChoiceJSON.FinalizedCheckpoint == nil {
		return errors.New("finalized checkpoint missing")
	}
	f.FinalizedCheckpoint = *forkChoiceJSON.FinalizedCheckpoint

	if forkChoiceJSON.ForkChoiceNodes == nil {
		return errors.New("fork choice nodes missing")
	}
	f.ForkChoiceNodes = forkChoiceJSON.ForkChoiceNodes

	return nil
}

// String returns a string version of the structure.
func (f *ForkChoice) String() string {
	data, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

// forkChoiceJSON is the json representation of the struct.
type forkChoiceJSON struct {
	JustifiedCheckpoint *phase0.Checkpoint `json:"justified_checkpoint"`
	FinalizedCheckpoint *phase0.Checkpoint `json:"finalized_checkpoint"`
	ForkChoiceNodes     []*ForkChoiceNode  `json:"fork_choice_nodes"`
}

// ForkChoiceNodeValidity represents the validity of a fork choice node.
type ForkChoiceNodeValidity uint64

const (
	// ForkChoiceNodeValidityUnknown is an unknown fork choice node.
	ForkChoiceNodeValidityUnknown ForkChoiceNodeValidity = iota
	// ForkChoiceNodeValidityInvalid is an invalid fork choice node.
	ForkChoiceNodeValidityInvalid
	// ForkChoiceNodeValidityValid is a valid fork choice node.
	ForkChoiceNodeValidityValid
	// ForkChoiceNodeValidityOptimistic is an optimistic fork choice node.
	ForkChoiceNodeValidityOptimistic
)

// ForkChoiceNodeValidityStrings are the strings for fork choice validity names.
var ForkChoiceNodeValidityStrings = [...]string{
	"unknown",
	"invalid",
	"valid",
	"optimistic",
}

// ForkChoiceNodeValidityFromString converts a string input to a fork choice.
func ForkChoiceNodeValidityFromString(input string) (ForkChoiceNodeValidity, error) {
	switch strings.ToLower(input) {
	case "invalid":
		return ForkChoiceNodeValidityInvalid, nil
	case "valid":
		return ForkChoiceNodeValidityValid, nil
	case "optimistic":
		return ForkChoiceNodeValidityOptimistic, nil
	default:
		return ForkChoiceNodeValidityUnknown, fmt.Errorf("unrecognised fork choice validity: %s", input)
	}
}

// MarshalJSON implements json.Marshaler.
func (d *ForkChoiceNodeValidity) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", ForkChoiceNodeValidityStrings[*d])), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *ForkChoiceNodeValidity) UnmarshalJSON(input []byte) error {
	var err error

	inputString := strings.Trim(string(input), "\"")
	if *d, err = ForkChoiceNodeValidityFromString(inputString); err != nil {
		return err
	}

	return nil
}

// String returns a string representation of the ForkChoiceNodeValidity.
func (d ForkChoiceNodeValidity) String() string {
	if int(d) >= len(ForkChoiceNodeValidityStrings) {
		return "unknown"
	}

	return ForkChoiceNodeValidityStrings[d]
}

// ForkChoiceNode is a node in the fork choice tree.
type ForkChoiceNode struct {
	// Slot is the slot of the node.
	Slot phase0.Slot
	// BlockRoot is the block root of the node.
	BlockRoot phase0.Root
	// ParentRoot is the parent root of the node.
	ParentRoot phase0.Root
	// JustifiedEpcih is the justified epoch of the node.
	JustifiedEpoch phase0.Epoch
	// FinalizedEpoch is the finalized epoch of the node.
	FinalizedEpoch phase0.Epoch
	// Weight is the weight of the node.
	Weight uint64
	// Validity is the validity of the node.
	Validity ForkChoiceNodeValidity
	// ExecutiionBlockHash is the execution block hash of the node.
	ExecutionBlockHash phase0.Root
	// ExtraData is the extra data of the node.
	ExtraData map[string]interface{}
}

// forkChoiceNodeJSON is the json representation of the struct.
type forkChoiceNodeJSON struct {
	Slot               string                 `json:"slot"`
	BlockRoot          string                 `json:"block_root"`
	ParentRoot         string                 `json:"parent_root"`
	JustifiedEpoch     string                 `json:"justified_epoch"`
	FinalizedEpoch     string                 `json:"finalized_epoch"`
	Weight             string                 `json:"weight"`
	Validity           string                 `json:"validity"`
	ExecutionBlockHash string                 `json:"execution_block_hash"`
	ExtraData          map[string]interface{} `json:"extra_data,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (f *ForkChoiceNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(&forkChoiceNodeJSON{
		Slot:               fmt.Sprintf("%d", f.Slot),
		BlockRoot:          fmt.Sprintf("%#x", f.BlockRoot),
		ParentRoot:         fmt.Sprintf("%#x", f.ParentRoot),
		JustifiedEpoch:     fmt.Sprintf("%d", f.JustifiedEpoch),
		FinalizedEpoch:     fmt.Sprintf("%d", f.FinalizedEpoch),
		Weight:             strconv.FormatUint(f.Weight, 10),
		Validity:           f.Validity.String(),
		ExecutionBlockHash: fmt.Sprintf("%#x", f.ExecutionBlockHash),
		ExtraData:          f.ExtraData,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *ForkChoiceNode) UnmarshalJSON(input []byte) error {
	var err error

	var forkChoiceNodeJSON forkChoiceNodeJSON
	if err = json.Unmarshal(input, &forkChoiceNodeJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	slot, err := strconv.ParseUint(forkChoiceNodeJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for slot: %s", forkChoiceNodeJSON.Slot))
	}
	f.Slot = phase0.Slot(slot)

	blockRoot, err := hex.DecodeString(strings.TrimPrefix(forkChoiceNodeJSON.BlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for block root: %s", forkChoiceNodeJSON.BlockRoot))
	}
	if len(blockRoot) != rootLength {
		return fmt.Errorf("incorrect length %d for block root", len(blockRoot))
	}
	copy(f.BlockRoot[:], blockRoot)

	parentRoot, err := hex.DecodeString(strings.TrimPrefix(forkChoiceNodeJSON.ParentRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for parent root: %s", forkChoiceNodeJSON.ParentRoot))
	}
	copy(f.ParentRoot[:], parentRoot)

	justifiedEpoch, err := strconv.ParseUint(forkChoiceNodeJSON.JustifiedEpoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for justified epoch: %s", forkChoiceNodeJSON.JustifiedEpoch))
	}
	f.JustifiedEpoch = phase0.Epoch(justifiedEpoch)

	finalizedEpoch, err := strconv.ParseUint(forkChoiceNodeJSON.FinalizedEpoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for finalized epoch: %s", forkChoiceNodeJSON.FinalizedEpoch))
	}
	f.FinalizedEpoch = phase0.Epoch(finalizedEpoch)

	weight, err := strconv.ParseUint(forkChoiceNodeJSON.Weight, 10, 64)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for weight: %s", forkChoiceNodeJSON.Weight))
	}
	f.Weight = weight

	validity, err := ForkChoiceNodeValidityFromString(forkChoiceNodeJSON.Validity)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for validity: %s", forkChoiceNodeJSON.Validity))
	}
	f.Validity = validity

	executionBlockHash, err := hex.DecodeString(strings.TrimPrefix(forkChoiceNodeJSON.ExecutionBlockHash, "0x"))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("invalid value for execution block hash: %s", forkChoiceNodeJSON.ExecutionBlockHash))
	}
	if len(executionBlockHash) != rootLength {
		return fmt.Errorf("incorrect length %d for execution block hash", len(executionBlockHash))
	}
	copy(f.ExecutionBlockHash[:], executionBlockHash)

	f.ExtraData = forkChoiceNodeJSON.ExtraData

	return nil
}

// String returns a string version of the structure.
func (f *ForkChoiceNode) String() string {
	data, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
