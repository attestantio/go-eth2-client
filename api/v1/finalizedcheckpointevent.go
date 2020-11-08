// Copyright © 2020 Attestant Limited.
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

	"github.com/pkg/errors"
)

// FinalizedCheckpointEvent is the data for the finalized checkpoint event.
type FinalizedCheckpointEvent struct {
	Block []byte
	State []byte
	Epoch uint64
}

// finalizedCheckpointEventJSON is the spec representation of the struct.
type finalizedCheckpointEventJSON struct {
	Block string `json:"block"`
	State string `json:"state"`
	Epoch string `json:"epoch"`
}

// MarshalJSON implements json.Marshaler.
func (e *FinalizedCheckpointEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(&finalizedCheckpointEventJSON{
		Block: fmt.Sprintf("%#x", e.Block),
		State: fmt.Sprintf("%#x", e.State),
		Epoch: fmt.Sprintf("%d", e.Epoch),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *FinalizedCheckpointEvent) UnmarshalJSON(input []byte) error {
	var err error

	var finalizedCheckpointEventJSON finalizedCheckpointEventJSON
	if err = json.Unmarshal(input, &finalizedCheckpointEventJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if finalizedCheckpointEventJSON.Block == "" {
		return errors.New("block missing")
	}
	if e.Block, err = hex.DecodeString(strings.TrimPrefix(finalizedCheckpointEventJSON.Block, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for block")
	}
	if len(e.Block) != rootLength {
		return fmt.Errorf("incorrect length %d for block", len(e.Block))
	}
	if finalizedCheckpointEventJSON.State == "" {
		return errors.New("state missing")
	}
	if e.State, err = hex.DecodeString(strings.TrimPrefix(finalizedCheckpointEventJSON.State, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for state")
	}
	if len(e.State) != rootLength {
		return fmt.Errorf("incorrect length %d for state", len(e.State))
	}
	if finalizedCheckpointEventJSON.Epoch == "" {
		return errors.New("epoch missing")
	}
	if e.Epoch, err = strconv.ParseUint(finalizedCheckpointEventJSON.Epoch, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for epoch")
	}

	return nil
}

// String returns a string version of the structure.
func (e *FinalizedCheckpointEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
