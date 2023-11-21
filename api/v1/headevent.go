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
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// HeadEvent is the data for the head event.
type HeadEvent struct {
	Slot                      phase0.Slot
	Block                     phase0.Root
	State                     phase0.Root
	EpochTransition           bool
	CurrentDutyDependentRoot  phase0.Root
	PreviousDutyDependentRoot phase0.Root
}

// headEventJSON is the spec representation of the struct.
type headEventJSON struct {
	Slot                      string `json:"slot"`
	Block                     string `json:"block"`
	State                     string `json:"state"`
	EpochTransition           bool   `json:"epoch_transition"`
	CurrentDutyDependentRoot  string `json:"current_duty_dependent_root,omitempty"`
	PreviousDutyDependentRoot string `json:"previous_duty_dependent_root,omitempty"`
}

// MarshalJSON implements json.Marshaler.
func (e *HeadEvent) MarshalJSON() ([]byte, error) {
	data := &headEventJSON{
		Slot:            fmt.Sprintf("%d", e.Slot),
		Block:           fmt.Sprintf("%#x", e.Block),
		State:           fmt.Sprintf("%#x", e.State),
		EpochTransition: e.EpochTransition,
	}
	// Optional fields (for now).
	var zeroRoot phase0.Root
	if !bytes.Equal(zeroRoot[:], e.CurrentDutyDependentRoot[:]) {
		data.CurrentDutyDependentRoot = fmt.Sprintf("%#x", e.CurrentDutyDependentRoot)
	}
	if !bytes.Equal(zeroRoot[:], e.PreviousDutyDependentRoot[:]) {
		data.PreviousDutyDependentRoot = fmt.Sprintf("%#x", e.PreviousDutyDependentRoot)
	}

	return json.Marshal(data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *HeadEvent) UnmarshalJSON(input []byte) error {
	var err error

	var headEventJSON headEventJSON
	if err = json.Unmarshal(input, &headEventJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if headEventJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(headEventJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	e.Slot = phase0.Slot(slot)
	if headEventJSON.Block == "" {
		return errors.New("block missing")
	}
	block, err := hex.DecodeString(strings.TrimPrefix(headEventJSON.Block, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block")
	}
	if len(block) != rootLength {
		return fmt.Errorf("incorrect length %d for block", len(block))
	}
	copy(e.Block[:], block)
	if headEventJSON.State == "" {
		return errors.New("state missing")
	}
	state, err := hex.DecodeString(strings.TrimPrefix(headEventJSON.State, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for state")
	}
	if len(state) != rootLength {
		return fmt.Errorf("incorrect length %d for state", len(state))
	}
	copy(e.State[:], state)
	e.EpochTransition = headEventJSON.EpochTransition
	// CurrentDutyDependentRoot only has partial coverage so do not complain if not present.
	if headEventJSON.CurrentDutyDependentRoot != "" {
		currentDutyDependentRoot, err := hex.DecodeString(strings.TrimPrefix(headEventJSON.CurrentDutyDependentRoot, "0x"))
		if err != nil {
			return errors.Wrap(err, "invalid value for current duty dependent root")
		}
		if len(currentDutyDependentRoot) != rootLength {
			return fmt.Errorf("incorrect length %d for current duty dependent root", len(currentDutyDependentRoot))
		}
		copy(e.CurrentDutyDependentRoot[:], currentDutyDependentRoot)
	}
	// PreviousDutyDependentRoot only has partial coverage so do not complain if not present.
	if headEventJSON.PreviousDutyDependentRoot != "" {
		previousDutyDependentRoot, err := hex.DecodeString(strings.TrimPrefix(headEventJSON.PreviousDutyDependentRoot, "0x"))
		if err != nil {
			return errors.Wrap(err, "invalid value for previous duty dependent root")
		}
		if len(previousDutyDependentRoot) != rootLength {
			return fmt.Errorf("incorrect length %d for previous duty dependent root", len(previousDutyDependentRoot))
		}
		copy(e.PreviousDutyDependentRoot[:], previousDutyDependentRoot)
	}

	return nil
}

// String returns a string version of the structure.
func (e *HeadEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
