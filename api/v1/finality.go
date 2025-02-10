// Copyright © 2020, 2021 Attestant Limited.
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

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// Finality is the data regarding finality checkpoints at a given state.
type Finality struct {
	// Finalized is the finalized checkpoint.
	Finalized *phase0.Checkpoint
	// Justified is the justified checkpoint.
	Justified *phase0.Checkpoint
	// PreviousJustified is the previous justified checkpoint.
	PreviousJustified *phase0.Checkpoint
}

// finalityJSON is the spec representation of the struct.
type finalityJSON struct {
	Finalized         *phase0.Checkpoint `json:"finalized"`
	Justified         *phase0.Checkpoint `json:"current_justified"`
	PreviousJustified *phase0.Checkpoint `json:"previous_justified"`
}

// MarshalJSON implements json.Marshaler.
func (f *Finality) MarshalJSON() ([]byte, error) {
	return json.Marshal(&finalityJSON{
		Finalized:         f.Finalized,
		Justified:         f.Justified,
		PreviousJustified: f.PreviousJustified,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Finality) UnmarshalJSON(input []byte) error {
	var err error

	var finalityJSON finalityJSON
	if err = json.Unmarshal(input, &finalityJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if finalityJSON.Finalized == nil {
		return errors.New("finalized checkpoint missing")
	}
	f.Finalized = finalityJSON.Finalized
	if finalityJSON.Justified == nil {
		return errors.New("justified checkpoint missing")
	}
	f.Justified = finalityJSON.Justified
	if finalityJSON.PreviousJustified == nil {
		return errors.New("previous justified checkpoint missing")
	}
	f.PreviousJustified = finalityJSON.PreviousJustified

	return nil
}

// String returns a string version of the structure.
func (f *Finality) String() string {
	data, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
