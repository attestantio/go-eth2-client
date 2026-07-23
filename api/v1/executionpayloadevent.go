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

// ExecutionPayloadEvent is the data for the `execution_payload` and
// `execution_payload_gossip` EIP-7732 events. Both carry a flat summary of a
// revealed execution payload (not the full SignedExecutionPayloadEnvelope): the
// `execution_payload` event fires when the envelope is imported into the
// fork-choice store, and `execution_payload_gossip` when it passes gossip
// validation. ExecutionOptimistic is only present on the `execution_payload`
// event.
type ExecutionPayloadEvent struct {
	Slot                phase0.Slot
	BuilderIndex        uint64
	BlockHash           phase0.Hash32
	BlockRoot           phase0.Root
	ExecutionOptimistic bool
}

// executionPayloadEventJSON is the spec representation of the struct.
type executionPayloadEventJSON struct {
	Slot                string `json:"slot"`
	BuilderIndex        string `json:"builder_index"`
	BlockHash           string `json:"block_hash"`
	BlockRoot           string `json:"block_root"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayloadEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&executionPayloadEventJSON{
		Slot:                fmt.Sprintf("%d", e.Slot),
		BuilderIndex:        fmt.Sprintf("%d", e.BuilderIndex),
		BlockHash:           fmt.Sprintf("%#x", e.BlockHash),
		BlockRoot:           fmt.Sprintf("%#x", e.BlockRoot),
		ExecutionOptimistic: e.ExecutionOptimistic,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
//
// Only the identifying slot and block root are required. builder_index and
// block_hash are parsed when present and validated for length, tolerating
// per-client field divergence (e.g. some clients emit an extra, non-spec
// state_root, which is ignored).
func (e *ExecutionPayloadEvent) UnmarshalJSON(input []byte) error {
	var data executionPayloadEventJSON
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

	if data.BlockRoot == "" {
		return errors.New("block root missing")
	}
	if err := decodeFixedBytes(e.BlockRoot[:], data.BlockRoot, rootLength, "block root"); err != nil {
		return err
	}

	if data.BuilderIndex != "" {
		builderIndex, err := strconv.ParseUint(data.BuilderIndex, 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid value for builder index")
		}
		e.BuilderIndex = builderIndex
	}

	if data.BlockHash != "" {
		if err := decodeFixedBytes(e.BlockHash[:], data.BlockHash, phase0.Hash32Length, "block hash"); err != nil {
			return err
		}
	}

	e.ExecutionOptimistic = data.ExecutionOptimistic

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

// decodeFixedBytes hex-decodes a 0x-prefixed value into dst, requiring exactly
// wantLen bytes.
func decodeFixedBytes(dst []byte, value string, wantLen int, name string) error {
	decoded, err := hex.DecodeString(strings.TrimPrefix(value, "0x"))
	if err != nil {
		return errors.Wrapf(err, "invalid value for %s", name)
	}
	if len(decoded) != wantLen {
		return fmt.Errorf("incorrect length %d for %s", len(decoded), name)
	}
	copy(dst, decoded)

	return nil
}
