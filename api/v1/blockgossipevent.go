// Copyright © 2024 Attestant Limited.
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

// BlockGossipEvent is the data for the block gossip event.
type BlockGossipEvent struct {
	Slot  phase0.Slot
	Block phase0.Root
}

// blockGossipEventJSON is the spec representation of the struct.
type blockGossipEventJSON struct {
	Slot  string `json:"slot"`
	Block string `json:"block"`
}

// MarshalJSON implements json.Marshaler.
func (e *BlockGossipEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blockGossipEventJSON{
		Slot:  fmt.Sprintf("%d", e.Slot),
		Block: fmt.Sprintf("%#x", e.Block),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *BlockGossipEvent) UnmarshalJSON(input []byte) error {
	var err error

	var data blockGossipEventJSON
	if err = json.Unmarshal(input, &data); err != nil {
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

	return nil
}

// String returns a string version of the structure.
func (e *BlockGossipEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
