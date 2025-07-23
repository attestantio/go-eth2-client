// Copyright © 2025 Attestant Limited.
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

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// DataColumnSidecarEvent is the data for the data column sidecar event.
type DataColumnSidecarEvent struct {
	BlockRoot      phase0.Root
	Slot           phase0.Slot
	Index          uint64
	KZGCommitments []deneb.KZGCommitment
}

// dataColumnSidecarEventJSON is the spec representation of the struct.
type dataColumnSidecarEventJSON struct {
	BlockRoot      string   `json:"block_root"`
	Slot           string   `json:"slot"`
	Index          string   `json:"index"`
	KZGCommitments []string `json:"kzg_commitments"`
}

// MarshalJSON implements json.Marshaler.
func (e *DataColumnSidecarEvent) MarshalJSON() ([]byte, error) {
	commitments := make([]string, len(e.KZGCommitments))
	for i, commitment := range e.KZGCommitments {
		commitments[i] = fmt.Sprintf("%#x", commitment)
	}

	return json.Marshal(&dataColumnSidecarEventJSON{
		BlockRoot:      fmt.Sprintf("%#x", e.BlockRoot),
		Slot:           fmt.Sprintf("%d", e.Slot),
		Index:          fmt.Sprintf("%d", e.Index),
		KZGCommitments: commitments,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *DataColumnSidecarEvent) UnmarshalJSON(input []byte) error {
	var err error

	var dataColumnSidecarEventJSON dataColumnSidecarEventJSON
	if err = json.Unmarshal(input, &dataColumnSidecarEventJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if dataColumnSidecarEventJSON.BlockRoot == "" {
		return errors.New("block_root missing")
	}
	err = e.BlockRoot.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, dataColumnSidecarEventJSON.BlockRoot)))
	if err != nil {
		return errors.Wrap(err, "invalid value for block_root")
	}

	if dataColumnSidecarEventJSON.Slot == "" {
		return errors.New("slot missing")
	}
	err = e.Slot.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, dataColumnSidecarEventJSON.Slot)))
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}

	if dataColumnSidecarEventJSON.Index == "" {
		return errors.New("index missing")
	}
	e.Index = 0
	if _, err = fmt.Sscanf(dataColumnSidecarEventJSON.Index, "%d", &e.Index); err != nil {
		return errors.Wrap(err, "invalid value for index")
	}

	if len(dataColumnSidecarEventJSON.KZGCommitments) == 0 {
		return errors.New("kzg_commitments missing")
	}
	e.KZGCommitments = make([]deneb.KZGCommitment, len(dataColumnSidecarEventJSON.KZGCommitments))
	for i, commitment := range dataColumnSidecarEventJSON.KZGCommitments {
		if commitment == "" {
			return fmt.Errorf("kzg_commitments[%d] missing", i)
		}
		err = e.KZGCommitments[i].UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, commitment)))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid value for kzg_commitments[%d]", i))
		}
	}

	return nil
}

// String returns a string version of the structure.
func (e *DataColumnSidecarEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
