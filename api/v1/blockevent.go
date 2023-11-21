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

// BlockEvent is the data for the block event.
type BlockEvent struct {
	Slot                phase0.Slot
	Block               phase0.Root
	ExecutionOptimistic bool
}

// blockEventJSON is the spec representation of the struct.
type blockEventJSON struct {
	Slot                string `json:"slot"`
	Block               string `json:"block"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
}

// MarshalJSON implements json.Marshaler.
func (e *BlockEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blockEventJSON{
		Slot:                fmt.Sprintf("%d", e.Slot),
		Block:               fmt.Sprintf("%#x", e.Block),
		ExecutionOptimistic: e.ExecutionOptimistic,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *BlockEvent) UnmarshalJSON(input []byte) error {
	var err error

	var blockEventJSON blockEventJSON
	if err = json.Unmarshal(input, &blockEventJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if blockEventJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(blockEventJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	e.Slot = phase0.Slot(slot)
	if blockEventJSON.Block == "" {
		return errors.New("block missing")
	}
	block, err := hex.DecodeString(strings.TrimPrefix(blockEventJSON.Block, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block")
	}
	if len(block) != rootLength {
		return fmt.Errorf("incorrect length %d for block", len(block))
	}
	copy(e.Block[:], block)
	e.ExecutionOptimistic = blockEventJSON.ExecutionOptimistic

	return nil
}

// String returns a string version of the structure.
func (e *BlockEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
