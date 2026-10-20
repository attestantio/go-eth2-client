// Copyright © 2026 Attestant Limited.
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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/pkg/errors"
)

// Builder contains the spec builder plus its index and status.
type Builder struct {
	Index   gloas.BuilderIndex
	Status  BuilderStatus
	Builder *gloas.Builder
}

// builderJSON is the spec representation of the struct.
type builderJSON struct {
	Index   string         `json:"index"`
	Status  BuilderStatus  `json:"status"`
	Builder *gloas.Builder `json:"builder"`
}

// MarshalJSON implements json.Marshaler.
func (v *Builder) MarshalJSON() ([]byte, error) {
	return json.Marshal(&builderJSON{
		Index:   fmt.Sprintf("%d", v.Index),
		Status:  v.Status,
		Builder: v.Builder,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *Builder) UnmarshalJSON(input []byte) error {
	var data builderJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.Index == "" {
		return errors.New("index missing")
	}

	index, err := strconv.ParseUint(data.Index, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for index")
	}

	if data.Status == BuilderStatusUnknown {
		return errors.New("status missing")
	}

	if data.Builder == nil {
		return errors.New("builder missing")
	}

	v.Index = gloas.BuilderIndex(index)
	v.Status = data.Status
	v.Builder = data.Builder

	return nil
}

// String returns a string version of the structure.
func (v *Builder) String() string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
