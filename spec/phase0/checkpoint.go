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
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// Checkpoint provides information about a checkpoint.
type Checkpoint struct {
	Epoch Epoch
	Root  Root `ssz-size:"32"`
}

// checkpointJSON is an internal representation of the struct.
type checkpointJSON struct {
	Epoch string `json:"epoch"`
	Root  string `json:"root"`
}

// checkpointYAML is an internal representation of the struct.
type checkpointYAML struct {
	Epoch uint64 `yaml:"epoch"`
	Root  string `yaml:"root"`
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
	var checkpointJSON checkpointJSON
	err := json.Unmarshal(input, &checkpointJSON)
	if err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return c.unpack(&checkpointJSON)
}

func (c *Checkpoint) unpack(checkpointJSON *checkpointJSON) error {
	if checkpointJSON.Epoch == "" {
		return errors.New("epoch missing")
	}
	epoch, err := strconv.ParseUint(checkpointJSON.Epoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for epoch")
	}
	c.Epoch = Epoch(epoch)
	if checkpointJSON.Root == "" {
		return errors.New("root missing")
	}
	root, err := hex.DecodeString(strings.TrimPrefix(checkpointJSON.Root, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for root")
	}
	if len(root) != RootLength {
		return errors.New("incorrect length for root")
	}
	copy(c.Root[:], root)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (c *Checkpoint) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&checkpointYAML{
		Epoch: uint64(c.Epoch),
		Root:  fmt.Sprintf("%#x", c.Root),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (c *Checkpoint) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var checkpointJSON checkpointJSON
	if err := yaml.Unmarshal(input, &checkpointJSON); err != nil {
		return err
	}

	return c.unpack(&checkpointJSON)
}

// String returns a string version of the structure.
func (c *Checkpoint) String() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
