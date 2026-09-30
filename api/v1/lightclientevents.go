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
	"errors"

	"github.com/attestantio/go-eth2-client/spec"
)

// LightClientFinalityUpdateEvent is a versioned light client finality update event.
type LightClientFinalityUpdateEvent struct {
	Version spec.DataVersion `json:"version"`
	Data    json.RawMessage  `json:"data"`
}

// LightClientOptimisticUpdateEvent is a versioned light client optimistic update event.
type LightClientOptimisticUpdateEvent struct {
	Version spec.DataVersion `json:"version"`
	Data    json.RawMessage  `json:"data"`
}

func (e *LightClientFinalityUpdateEvent) UnmarshalJSON(input []byte) error {
	version, data, err := decodeLightClientEvent(input)
	if err != nil {
		return err
	}
	e.Version, e.Data = version, data

	return nil
}

func (e *LightClientOptimisticUpdateEvent) UnmarshalJSON(input []byte) error {
	version, data, err := decodeLightClientEvent(input)
	if err != nil {
		return err
	}
	e.Version, e.Data = version, data

	return nil
}

func decodeLightClientEvent(input []byte) (spec.DataVersion, json.RawMessage, error) {
	var event struct {
		Version spec.DataVersion `json:"version"`
		Data    json.RawMessage  `json:"data"`
	}
	if err := json.Unmarshal(input, &event); err != nil {
		return 0, nil, err
	}
	if event.Version == spec.DataVersionUnknown {
		return 0, nil, errors.New("version missing")
	}
	if len(event.Data) == 0 || string(event.Data) == "null" {
		return 0, nil, errors.New("data missing")
	}

	return event.Version, event.Data, nil
}
