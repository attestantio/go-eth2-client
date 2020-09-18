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

// Fork provides information about a fork.
type Fork struct {
	// Previous version is the previous fork version.
	PreviousVersion []byte `ssz-size:"4"`
	// Current version is the current fork version.
	CurrentVersion []byte `ssz-size:"4"`
	// Epoch is the epoch at which the current fork version took effect.
	Epoch uint64
}

// forkJSON is the spec representation of the struct.
type forkJSON struct {
	PreviousVersion string `json:"previous_version"`
	CurrentVersion  string `json:"current_version"`
	Epoch           string `json:"epoch"`
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
	var err error

	var forkJSON forkJSON
	if err = json.Unmarshal(input, &forkJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if forkJSON.PreviousVersion == "" {
		return errors.New("previous version missing")
	}
	if f.PreviousVersion, err = hex.DecodeString(strings.TrimPrefix(forkJSON.PreviousVersion, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for previous version")
	}
	if len(f.PreviousVersion) != forkVersionLength {
		return errors.New("incorrect length for previous version")
	}
	if forkJSON.CurrentVersion == "" {
		return errors.New("current version missing")
	}
	if f.CurrentVersion, err = hex.DecodeString(strings.TrimPrefix(forkJSON.CurrentVersion, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for current version")
	}
	if len(f.CurrentVersion) != forkVersionLength {
		return errors.New("incorrect length for current version")
	}
	if forkJSON.Epoch == "" {
		return errors.New("epoch missing")
	}
	if f.Epoch, err = strconv.ParseUint(forkJSON.Epoch, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for epoch")
	}

	return nil
}

// String returns a string version of the structure.
func (f *Fork) String() string {
	data, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
