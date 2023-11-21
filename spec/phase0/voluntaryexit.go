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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// VoluntaryExit provides information about a voluntary exit.
type VoluntaryExit struct {
	Epoch          Epoch
	ValidatorIndex ValidatorIndex
}

// voluntaryExitJSON is an internal representation of the struct.
type voluntaryExitJSON struct {
	Epoch          string `json:"epoch"`
	ValidatorIndex string `json:"validator_index"`
}

// voluntaryExitYAML is an internal representation of the struct.
type voluntaryExitYAML struct {
	Epoch          uint64 `json:"epoch"`
	ValidatorIndex uint64 `json:"validator_index"`
}

// MarshalJSON implements json.Marshaler.
func (v *VoluntaryExit) MarshalJSON() ([]byte, error) {
	return json.Marshal(&voluntaryExitJSON{
		Epoch:          fmt.Sprintf("%d", v.Epoch),
		ValidatorIndex: fmt.Sprintf("%d", v.ValidatorIndex),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *VoluntaryExit) UnmarshalJSON(input []byte) error {
	var voluntaryExitJSON voluntaryExitJSON
	err := json.Unmarshal(input, &voluntaryExitJSON)
	if err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return v.unpack(&voluntaryExitJSON)
}

func (v *VoluntaryExit) unpack(voluntaryExitJSON *voluntaryExitJSON) error {
	if voluntaryExitJSON.Epoch == "" {
		return errors.New("epoch missing")
	}
	epoch, err := strconv.ParseUint(voluntaryExitJSON.Epoch, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for epoch")
	}
	v.Epoch = Epoch(epoch)
	if voluntaryExitJSON.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(voluntaryExitJSON.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	v.ValidatorIndex = ValidatorIndex(validatorIndex)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (v *VoluntaryExit) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&voluntaryExitYAML{
		Epoch:          uint64(v.Epoch),
		ValidatorIndex: uint64(v.ValidatorIndex),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (v *VoluntaryExit) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var voluntaryExitJSON voluntaryExitJSON
	if err := yaml.Unmarshal(input, &voluntaryExitJSON); err != nil {
		return err
	}

	return v.unpack(&voluntaryExitJSON)
}

// String returns a string version of the structure.
func (v *VoluntaryExit) String() string {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
