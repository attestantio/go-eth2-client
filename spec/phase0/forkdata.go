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
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// ForkData provides data about a fork.
type ForkData struct {
	// Current version is the current fork version.
	CurrentVersion Version `ssz-size:"4"`
	// GenesisValidatorsRoot is the hash tree root of the validators at genesis.
	GenesisValidatorsRoot Root `ssz-size:"32"`
}

// forkDataJSON is the spec representation of the struct.
type forkDataJSON struct {
	CurrentVersion        string `json:"current_version"`
	GenesisValidatorsRoot string `json:"genesis_validators_root"`
}

// forkDataYAML is the spec representation of the struct.
type forkDataYAML struct {
	CurrentVersion        string `yaml:"current_version"`
	GenesisValidatorsRoot string `yaml:"genesis_validators_root"`
}

// MarshalJSON implements json.Marshaler.
func (f *ForkData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&forkDataJSON{
		CurrentVersion:        fmt.Sprintf("%#x", f.CurrentVersion),
		GenesisValidatorsRoot: fmt.Sprintf("%#x", f.GenesisValidatorsRoot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *ForkData) UnmarshalJSON(input []byte) error {
	var forkDataJSON forkDataJSON
	if err := json.Unmarshal(input, &forkDataJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return f.unpack(&forkDataJSON)
}

func (f *ForkData) unpack(forkDataJSON *forkDataJSON) error {
	if forkDataJSON.CurrentVersion == "" {
		return errors.New("current version missing")
	}
	currentVersion, err := hex.DecodeString(strings.TrimPrefix(forkDataJSON.CurrentVersion, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for current version")
	}
	if len(currentVersion) != ForkVersionLength {
		return errors.New("incorrect length for current version")
	}
	copy(f.CurrentVersion[:], currentVersion)
	if forkDataJSON.GenesisValidatorsRoot == "" {
		return errors.New("genesis validators root missing")
	}
	genesisValidatorsRoot, err := hex.DecodeString(strings.TrimPrefix(forkDataJSON.GenesisValidatorsRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for genesis validators root")
	}
	if len(genesisValidatorsRoot) != RootLength {
		return errors.New("incorrect length for genesis validators root")
	}
	copy(f.GenesisValidatorsRoot[:], genesisValidatorsRoot)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (f *ForkData) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&forkDataYAML{
		CurrentVersion:        fmt.Sprintf("%#x", f.CurrentVersion),
		GenesisValidatorsRoot: fmt.Sprintf("%#x", f.GenesisValidatorsRoot),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (f *ForkData) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var forkDataJSON forkDataJSON
	if err := yaml.Unmarshal(input, &forkDataJSON); err != nil {
		return err
	}

	return f.unpack(&forkDataJSON)
}

// String returns a string version of the structure.
func (f *ForkData) String() string {
	data, err := yaml.Marshal(f)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
