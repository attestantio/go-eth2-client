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
//
// It emits the bare data object, not the {"version": ..., "data": ...} wrapper.
// The topic is registered as versioned, so eventtopic.Topic and Event.MarshalJSON
// add the wrapper; emitting one here as well would nest it twice.
func (e *HeadEventV2) MarshalJSON() ([]byte, error) {
	return json.Marshal(&headEventV2DataJSON{
		Slot:                      fmt.Sprintf("%d", e.Slot),
		Block:                     fmt.Sprintf("%#x", e.Block),
		State:                     fmt.Sprintf("%#x", e.State),
		PayloadStatus:             e.PayloadStatus,
		EpochTransition:           e.EpochTransition,
		CurrentEpochDependentRoot: fmt.Sprintf("%#x", e.CurrentEpochDependentRoot),
		NextEpochDependentRoot:    fmt.Sprintf("%#x", e.NextEpochDependentRoot),
		ExecutionOptimistic:       e.ExecutionOptimistic,
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
//
// Both the bare data object and the {"version": ..., "data": ...} wrapper are
// accepted.  eventtopic.Topic.Decode strips the wrapper and checks the fork
// before this is reached on the events stream, so the wrapped form is here for
// callers decoding a HeadEventV2 directly.
func (e *HeadEventV2) UnmarshalJSON(input []byte) error {
	version := spec.DataVersionGloas
	data := input

	var wrapper headEventV2JSON
	if err := json.Unmarshal(input, &wrapper); err == nil &&
		wrapper.Version != spec.DataVersionUnknown && wrapper.Data != nil {
		version = wrapper.Version
		if data, err = json.Marshal(wrapper.Data); err != nil {
			return errors.Wrap(err, "invalid JSON")
		}
	}

	var event headEventV2DataJSON
	if err := json.Unmarshal(data, &event); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if event.PayloadStatus == "" {
		return errors.New("payload status missing")
	}
	if event.Slot == "" {
		return errors.New("slot missing")
	}

	slot, err := strconv.ParseUint(event.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	if event.Block == "" {
		return errors.New("block missing")
	}
	if err := decodeFixedBytes(e.Block[:], event.Block, "block"); err != nil {
		return err
	}
	if event.State == "" {
		return errors.New("state missing")
	}
	if err := decodeFixedBytes(e.State[:], event.State, "state"); err != nil {
		return err
	}

	// The dependent roots are optional, as they are on the head event: node
	// implementations of head_v2 are still landing, and a missing root would
	// otherwise fail Topic.Decode, which makes http.handleEvent log "Failed to
	// parse event" and drop it -- so a consumer subscribed only to head_v2
	// would silently receive no head events at all.  They are cleared when
	// absent so that a value reused for an event without them does not keep an
	// earlier one.
	e.CurrentEpochDependentRoot = phase0.Root{}
	if event.CurrentEpochDependentRoot != "" {
		if err := decodeFixedBytes(
			e.CurrentEpochDependentRoot[:],
			event.CurrentEpochDependentRoot,
			"current epoch dependent root",
		); err != nil {
			return err
		}
	}

	e.NextEpochDependentRoot = phase0.Root{}
	if event.NextEpochDependentRoot != "" {
		if err := decodeFixedBytes(
			e.NextEpochDependentRoot[:],
			event.NextEpochDependentRoot,
			"next epoch dependent root",
		); err != nil {
			return err
		}
	}

	e.Version = version
	e.Slot = phase0.Slot(slot)
	e.PayloadStatus = event.PayloadStatus
	e.EpochTransition = event.EpochTransition
	e.ExecutionOptimistic = event.ExecutionOptimistic

	return nil
}
