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

package phase0

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Checkpoint provides information about a checkpoint.
type Checkpoint struct {
	Epoch uint64
	Root  []byte `ssz-size:"32"`
}

// checkpointJSON is the spec representation of the struct.
type checkpointJSON struct {
	Epoch string `json:"epoch"`
	Root  string `json:"root"`
}

// MarshalJSON implements json.Marshaler.
func (c *Checkpoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(&checkpointJSON{
		Epoch: fmt.Sprintf("%d", c.Epoch),
		Root:  fmt.Sprintf("%#x", c.Root),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Checkpoint) UnmarshalJSON(input []byte) error {
	var err error

	var checkpointJSON checkpointJSON
	if err = json.Unmarshal(input, &checkpointJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if checkpointJSON.Epoch == "" {
		return errors.New("epoch missing")
	}
	if c.Epoch, err = strconv.ParseUint(checkpointJSON.Epoch, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for epoch")
	}
	if checkpointJSON.Root == "" {
		return errors.New("root missing")
	}
	if c.Root, err = hex.DecodeString(strings.TrimPrefix(checkpointJSON.Root, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for root")
	}
	if len(c.Root) != rootLength {
		return errors.New("incorrect length for root")
	}

	return nil
}

// String returns a string version of the structure.
func (c *Checkpoint) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
