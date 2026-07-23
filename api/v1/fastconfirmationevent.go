// Copyright © 2026 Attestant Limited.
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

// FastConfirmationEvent is the data for the fast confirmation event. Slot and
// Block identify the most recent confirmed block; CurrentSlot is the wall-clock
// slot at which the fast-confirmation algorithm was executed.
type FastConfirmationEvent struct {
	Slot        phase0.Slot
	Block       phase0.Root
	CurrentSlot phase0.Slot
}

// fastConfirmationEventJSON is the spec representation of the struct.
type fastConfirmationEventJSON struct {
	Slot        string `json:"slot"`
	Block       string `json:"block"`
	CurrentSlot string `json:"current_slot"`
}

// MarshalJSON implements json.Marshaler.
func (e *FastConfirmationEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&fastConfirmationEventJSON{
		Slot:        fmt.Sprintf("%d", e.Slot),
		Block:       fmt.Sprintf("%#x", e.Block),
		CurrentSlot: fmt.Sprintf("%d", e.CurrentSlot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *FastConfirmationEvent) UnmarshalJSON(input []byte) error {
	var data fastConfirmationEventJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.Slot == "" {
		return errors.New("slot missing")
	}

	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}

	e.Slot = phase0.Slot(slot)

	if data.Block == "" {
		return errors.New("block missing")
	}

	block, err := hex.DecodeString(strings.TrimPrefix(data.Block, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block")
	}

	if len(block) != rootLength {
		return fmt.Errorf("incorrect length %d for block", len(block))
	}

	copy(e.Block[:], block)

	// current_slot was added to the spec after slot/block; parse it when
	// present but tolerate clients that do not yet emit it.
	if data.CurrentSlot != "" {
		currentSlot, err := strconv.ParseUint(data.CurrentSlot, 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid value for current slot")
		}

		e.CurrentSlot = phase0.Slot(currentSlot)
	}

	return nil
}

// String returns a string version of the structure.
func (e *FastConfirmationEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
