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

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// HeadEventV2 is the data for the head_v2 event.
type HeadEventV2 struct {
	Version                   spec.DataVersion
	Slot                      phase0.Slot
	Block                     phase0.Root
	State                     phase0.Root
	PayloadStatus             string
	EpochTransition           bool
	CurrentEpochDependentRoot phase0.Root
	NextEpochDependentRoot    phase0.Root
	ExecutionOptimistic       bool
}

// headEventV2JSON is the spec representation of the versioned event.
type headEventV2JSON struct {
	Version spec.DataVersion     `json:"version"`
	Data    *headEventV2DataJSON `json:"data"`
}

// headEventV2DataJSON is the spec representation of the event data.
type headEventV2DataJSON struct {
	Slot                      string `json:"slot"`
	Block                     string `json:"block"`
	State                     string `json:"state"`
	PayloadStatus             string `json:"payload_status"`
	EpochTransition           bool   `json:"epoch_transition"`
	CurrentEpochDependentRoot string `json:"current_epoch_dependent_root"`
	NextEpochDependentRoot    string `json:"next_epoch_dependent_root"`
	ExecutionOptimistic       bool   `json:"execution_optimistic"`
}

// MarshalJSON implements json.Marshaler.
func (e *HeadEventV2) MarshalJSON() ([]byte, error) {
	return json.Marshal(&headEventV2JSON{
		Version: e.Version,
		Data: &headEventV2DataJSON{
			Slot:                      fmt.Sprintf("%d", e.Slot),
			Block:                     fmt.Sprintf("%#x", e.Block),
			State:                     fmt.Sprintf("%#x", e.State),
			PayloadStatus:             e.PayloadStatus,
			EpochTransition:           e.EpochTransition,
			CurrentEpochDependentRoot: fmt.Sprintf("%#x", e.CurrentEpochDependentRoot),
			NextEpochDependentRoot:    fmt.Sprintf("%#x", e.NextEpochDependentRoot),
			ExecutionOptimistic:       e.ExecutionOptimistic,
		},
	})
}

// String returns a string version of the structure.
func (e *HeadEventV2) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *HeadEventV2) UnmarshalJSON(input []byte) error {
	var event headEventV2JSON
	if err := json.Unmarshal(input, &event); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if event.Version == spec.DataVersionUnknown {
		return errors.New("version missing")
	}
	if event.Data == nil {
		return errors.New("data missing")
	}
	if event.Data.PayloadStatus == "" {
		return errors.New("payload status missing")
	}
	if event.Data.Slot == "" {
		return errors.New("slot missing")
	}

	slot, err := strconv.ParseUint(event.Data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	if event.Data.Block == "" {
		return errors.New("block missing")
	}
	if err := decodeFixedBytes(e.Block[:], event.Data.Block, "block"); err != nil {
		return err
	}
	if event.Data.State == "" {
		return errors.New("state missing")
	}
	if err := decodeFixedBytes(e.State[:], event.Data.State, "state"); err != nil {
		return err
	}
	if event.Data.CurrentEpochDependentRoot == "" {
		return errors.New("current epoch dependent root missing")
	}
	if err := decodeFixedBytes(
		e.CurrentEpochDependentRoot[:],
		event.Data.CurrentEpochDependentRoot,
		"current epoch dependent root",
	); err != nil {
		return err
	}
	if event.Data.NextEpochDependentRoot == "" {
		return errors.New("next epoch dependent root missing")
	}
	if err := decodeFixedBytes(
		e.NextEpochDependentRoot[:],
		event.Data.NextEpochDependentRoot,
		"next epoch dependent root",
	); err != nil {
		return err
	}

	e.Version = event.Version
	e.Slot = phase0.Slot(slot)
	e.PayloadStatus = event.Data.PayloadStatus
	e.EpochTransition = event.Data.EpochTransition
	e.ExecutionOptimistic = event.Data.ExecutionOptimistic

	return nil
}
