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
	"encoding/json"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// ExecutionPayloadGossipEvent is a payload that passed gossip validation, not fork-choice import.
type ExecutionPayloadGossipEvent struct {
	Slot         phase0.Slot
	BuilderIndex gloas.BuilderIndex
	BlockHash    phase0.Hash32
	BlockRoot    phase0.Root
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayloadGossipEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Slot         string `json:"slot"`
		BuilderIndex string `json:"builder_index"`
		BlockHash    string `json:"block_hash"`
		BlockRoot    string `json:"block_root"`
	}{
		Slot:         fmt.Sprintf("%d", e.Slot),
		BuilderIndex: fmt.Sprintf("%d", e.BuilderIndex),
		BlockHash:    fmt.Sprintf("%#x", e.BlockHash),
		BlockRoot:    fmt.Sprintf("%#x", e.BlockRoot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionPayloadGossipEvent) UnmarshalJSON(input []byte) error {
	var imported ExecutionPayloadEvent
	if err := json.Unmarshal(input, &imported); err != nil {
		return err
	}
	e.Slot = imported.Slot
	e.BuilderIndex = imported.BuilderIndex
	e.BlockHash = imported.BlockHash
	e.BlockRoot = imported.BlockRoot

	return nil
}

// String returns a string version of the structure.
func (e *ExecutionPayloadGossipEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
