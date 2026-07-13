// Copyright © 2025 Attestant Limited.
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

// ExecutionPayloadAvailableEvent is the data for the execution payload available event.
type ExecutionPayloadAvailableEvent struct {
	BlockRoot phase0.Root
	Slot      phase0.Slot
}

// executionPayloadAvailableEventJSON is the spec representation of the struct.
type executionPayloadAvailableEventJSON struct {
	BlockRoot string `json:"block_root"`
	Slot      string `json:"slot"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayloadAvailableEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&executionPayloadAvailableEventJSON{
		BlockRoot: fmt.Sprintf("%#x", e.BlockRoot),
		Slot:      fmt.Sprintf("%d", e.Slot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionPayloadAvailableEvent) UnmarshalJSON(input []byte) error {
	var err error

	var executionPayloadAvailableEventJSON executionPayloadAvailableEventJSON
	if err = json.Unmarshal(input, &executionPayloadAvailableEventJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if executionPayloadAvailableEventJSON.BlockRoot == "" {
		return errors.New("block root missing")
	}
	block, err := hex.DecodeString(strings.TrimPrefix(executionPayloadAvailableEventJSON.BlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block root")
	}
	if len(block) != rootLength {
		return fmt.Errorf("incorrect length %d for block root", len(block))
	}
	copy(e.BlockRoot[:], block)
	if executionPayloadAvailableEventJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(executionPayloadAvailableEventJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	e.Slot = phase0.Slot(slot)

	return nil
}

// String returns a string version of the structure.
func (e *ExecutionPayloadAvailableEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
