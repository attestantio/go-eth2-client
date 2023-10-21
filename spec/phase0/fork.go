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

// Fork provides information about a fork.
type Fork struct {
	// Previous version is the previous fork version.
	PreviousVersion Version `ssz-size:"4"`
	// Current version is the current fork version.
	CurrentVersion Version `ssz-size:"4"`
	// Epoch is the epoch at which the current fork version took effect.
	Epoch Epoch
}

// forkJSON is the spec representation of the struct.
type forkJSON struct {
	PreviousVersion string `json:"previous_version"`
	CurrentVersion  string `json:"current_version"`
	Epoch           string `json:"epoch"`
}

// forkYAML is the spec representation of the struct.
type forkYAML struct {
	PreviousVersion string `yaml:"previous_version"`
	CurrentVersion  string `yaml:"current_version"`
	Epoch           uint64 `yaml:"epoch"`
}

// MarshalJSON implements json.Marshaler.
func (f *Fork) MarshalJSON() ([]byte, error) {
	return json.Marshal(&forkJSON{
		PreviousVersion: fmt.Sprintf("%#x", f.PreviousVersion),
		CurrentVersion:  fmt.Sprintf("%#x", f.CurrentVersion),
		Epoch:           fmt.Sprintf("%d", f.Epoch),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Fork) UnmarshalJSON(input []byte) error {
	var forkJSON forkJSON
	if err := json.Unmarshal(input, &forkJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return f.unpack(&forkJSON)
}

func (f *Fork) unpack(forkJSON *forkJSON) error {
	if forkJSON.PreviousVersion == "" {
		return errors.New("previous version missing")
	}
	previousVersion, err := hex.DecodeString(strings.TrimPrefix(forkJSON.PreviousVersion, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for previous version")
	}
	if len(previousVersion) != ForkVersionLength {
		return errors.New("incorrect length for previous version")
	}
	copy(f.PreviousVersion[:], previousVersion)
	if forkJSON.CurrentVersion == "" {
		return errors.New("current version missing")
	}
	currentVersion, err := hex.DecodeString(strings.TrimPrefix(forkJSON.CurrentVersion, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for current version")
	}
	if len(currentVersion) != ForkVersionLength {
		return errors.New("incorrect length for current version")
	}
	copy(f.CurrentVersion[:], currentVersion)
	if forkJSON.Epoch == "" {
		return errors.New("epoch missing")
	}
	epoch, err := strconv.ParseUint(forkJSON.Epoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for epoch")
	}
	f.Epoch = Epoch(epoch)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (f *Fork) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&forkYAML{
		PreviousVersion: fmt.Sprintf("%#x", f.PreviousVersion),
		CurrentVersion:  fmt.Sprintf("%#x", f.CurrentVersion),
		Epoch:           uint64(f.Epoch),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (f *Fork) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var forkJSON forkJSON
	if err := yaml.Unmarshal(input, &forkJSON); err != nil {
		return err
	}

	return f.unpack(&forkJSON)
}

// String returns a string version of the structure.
func (f *Fork) String() string {
	data, err := yaml.Marshal(f)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
