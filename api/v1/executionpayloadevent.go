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

// ExecutionPayloadEvent is the data for the execution payload event.
type ExecutionPayloadEvent struct {
	BlockRoot           phase0.Root
	Slot                phase0.Slot
	ExecutionBlockHash  phase0.Hash32
	ExecutionOptimistic bool
}

// executionPayloadEventJSON is the spec representation of the struct.
type executionPayloadEventJSON struct {
	BlockRoot           string `json:"block_root"`
	Slot                string `json:"slot"`
	ExecutionBlockHash  string `json:"execution_block_hash"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayloadEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&executionPayloadEventJSON{
		BlockRoot:           fmt.Sprintf("%#x", e.BlockRoot),
		Slot:                fmt.Sprintf("%d", e.Slot),
		ExecutionBlockHash:  fmt.Sprintf("%#x", e.ExecutionBlockHash),
		ExecutionOptimistic: e.ExecutionOptimistic,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionPayloadEvent) UnmarshalJSON(input []byte) error {
	var err error

	var executionPayloadEventJSON executionPayloadEventJSON
	if err = json.Unmarshal(input, &executionPayloadEventJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if executionPayloadEventJSON.BlockRoot == "" {
		return errors.New("block root missing")
	}
	block, err := hex.DecodeString(strings.TrimPrefix(executionPayloadEventJSON.BlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block root")
	}
	if len(block) != rootLength {
		return fmt.Errorf("incorrect length %d for block root", len(block))
	}
	copy(e.BlockRoot[:], block)
	if executionPayloadEventJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(executionPayloadEventJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	e.Slot = phase0.Slot(slot)
	if executionPayloadEventJSON.ExecutionBlockHash == "" {
		return errors.New("payload hash missing")
	}
	payloadHash, err := hex.DecodeString(strings.TrimPrefix(executionPayloadEventJSON.ExecutionBlockHash, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload hash")
	}
	if len(payloadHash) != rootLength {
		return fmt.Errorf("incorrect length %d for payload hash", len(payloadHash))
	}
	copy(e.ExecutionBlockHash[:], payloadHash)
	e.ExecutionOptimistic = executionPayloadEventJSON.ExecutionOptimistic

	return nil
}

// String returns a string version of the structure.
func (e *ExecutionPayloadEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
