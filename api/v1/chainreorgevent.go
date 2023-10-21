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

// ChainReorgEvent is the data for the head event.
type ChainReorgEvent struct {
	Slot         phase0.Slot
	Depth        uint64
	OldHeadBlock phase0.Root
	NewHeadBlock phase0.Root
	OldHeadState phase0.Root
	NewHeadState phase0.Root
	Epoch        phase0.Epoch
}

// chainReorgEventJSON is the spec representation of the struct.
type chainReorgEventJSON struct {
	Slot         string `json:"slot"`
	Depth        string `json:"depth"`
	OldHeadBlock string `json:"old_head_block"`
	NewHeadBlock string `json:"new_head_block"`
	OldHeadState string `json:"old_head_state"`
	NewHeadState string `json:"new_head_state"`
	Epoch        string `json:"epoch"`
}

// MarshalJSON implements json.Marshaler.
func (e *ChainReorgEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&chainReorgEventJSON{
		Slot:         fmt.Sprintf("%d", e.Slot),
		Depth:        strconv.FormatUint(e.Depth, 10),
		OldHeadBlock: fmt.Sprintf("%#x", e.OldHeadBlock),
		NewHeadBlock: fmt.Sprintf("%#x", e.NewHeadBlock),
		OldHeadState: fmt.Sprintf("%#x", e.OldHeadState),
		NewHeadState: fmt.Sprintf("%#x", e.NewHeadState),
		Epoch:        fmt.Sprintf("%d", e.Epoch),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ChainReorgEvent) UnmarshalJSON(input []byte) error {
	var err error

	var chainReorgEventJSON chainReorgEventJSON
	if err = json.Unmarshal(input, &chainReorgEventJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if chainReorgEventJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(chainReorgEventJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	e.Slot = phase0.Slot(slot)
	if chainReorgEventJSON.Depth == "" {
		return errors.New("depth missing")
	}
	if e.Depth, err = strconv.ParseUint(chainReorgEventJSON.Depth, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for depth")
	}
	if chainReorgEventJSON.OldHeadBlock == "" {
		return errors.New("old head block missing")
	}
	oldHeadBlock, err := hex.DecodeString(strings.TrimPrefix(chainReorgEventJSON.OldHeadBlock, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for old head block")
	}
	if len(oldHeadBlock) != rootLength {
		return fmt.Errorf("incorrect length %d for old head block", len(oldHeadBlock))
	}
	copy(e.OldHeadBlock[:], oldHeadBlock)
	if chainReorgEventJSON.NewHeadBlock == "" {
		return errors.New("new head block missing")
	}
	newHeadBlock, err := hex.DecodeString(strings.TrimPrefix(chainReorgEventJSON.NewHeadBlock, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for new head block")
	}
	if len(newHeadBlock) != rootLength {
		return fmt.Errorf("incorrect length %d for new head block", len(newHeadBlock))
	}
	copy(e.NewHeadBlock[:], newHeadBlock)
	if chainReorgEventJSON.OldHeadState == "" {
		return errors.New("old head state missing")
	}
	oldHeadState, err := hex.DecodeString(strings.TrimPrefix(chainReorgEventJSON.OldHeadState, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for old head state")
	}
	if len(oldHeadState) != rootLength {
		return fmt.Errorf("incorrect length %d for old head state", len(oldHeadState))
	}
	copy(e.OldHeadState[:], oldHeadState)
	if chainReorgEventJSON.NewHeadState == "" {
		return errors.New("new head state missing")
	}
	newHeadState, err := hex.DecodeString(strings.TrimPrefix(chainReorgEventJSON.NewHeadState, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for new head state")
	}
	if len(newHeadState) != rootLength {
		return fmt.Errorf("incorrect length %d for new head state", len(newHeadState))
	}
	copy(e.NewHeadState[:], newHeadState)
	if chainReorgEventJSON.Epoch == "" {
		return errors.New("epoch missing")
	}
	epoch, err := strconv.ParseUint(chainReorgEventJSON.Epoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for epoch")
	}
	e.Epoch = phase0.Epoch(epoch)

	return nil
}

// String returns a string version of the structure.
func (e *ChainReorgEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
