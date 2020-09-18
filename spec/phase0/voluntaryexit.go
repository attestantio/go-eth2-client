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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// VoluntaryExit provides information about a voluntary exit.
type VoluntaryExit struct {
	Epoch          uint64
	ValidatorIndex uint64
}

// voluntaryExitJSON is the spec representation of the struct.
type voluntaryExitJSON struct {
	Epoch          string `json:"epoch"`
	ValidatorIndex string `json:"validator_index"`
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
	var err error

	var voluntaryExitJSON voluntaryExitJSON
	if err = json.Unmarshal(input, &voluntaryExitJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if voluntaryExitJSON.Epoch == "" {
		return errors.New("epoch missing")
	}
	if v.Epoch, err = strconv.ParseUint(voluntaryExitJSON.Epoch, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for epoch")
	}
	if voluntaryExitJSON.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	if v.ValidatorIndex, err = strconv.ParseUint(voluntaryExitJSON.ValidatorIndex, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}

	return nil
}

// String returns a string version of the structure.
func (v *VoluntaryExit) String() string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
