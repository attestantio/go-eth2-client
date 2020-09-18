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
	"strings"

	"github.com/pkg/errors"
)

// ForkData provides data about a fork.
type ForkData struct {
	// Current version is the current fork version.
	CurrentVersion []byte `ssz-size:"4"`
	// GenesisValidatorsRoot is the hash tree root of the validators at genesis.
	GenesisValidatorsRoot []byte `ssz-size:"32"`
}

// forkDataJSON is the spec representation of the struct.
type forkDataJSON struct {
	CurrentVersion        string `json:"current_version"`
	GenesisValidatorsRoot string `json:"genesis_validators_root"`
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
	var err error

	var forkDataJSON forkDataJSON
	if err = json.Unmarshal(input, &forkDataJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if forkDataJSON.CurrentVersion == "" {
		return errors.New("current version missing")
	}
	if f.CurrentVersion, err = hex.DecodeString(strings.TrimPrefix(forkDataJSON.CurrentVersion, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for current version")
	}
	if len(f.CurrentVersion) != forkVersionLength {
		return errors.New("incorrect length for current version")
	}
	if forkDataJSON.GenesisValidatorsRoot == "" {
		return errors.New("genesis validators root missing")
	}
	if f.GenesisValidatorsRoot, err = hex.DecodeString(strings.TrimPrefix(forkDataJSON.GenesisValidatorsRoot, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for genesis validators root")
	}
	if len(f.GenesisValidatorsRoot) != rootLength {
		return errors.New("incorrect length for genesis validators root")
	}

	return nil
}

// String returns a string version of the structure.
func (f *ForkData) String() string {
	data, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
