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

package v1

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

// SyncState is the data regarding the node's synchronization state to the chain.
type SyncState struct {
	// HeadSlot is the head slot of the chain as understood by the node.
	HeadSlot uint64
	// SyncDistance is the distance between the node's highest synced slot and the head slot.
	SyncDistance uint64
}

// syncStateJSON is the spec representation of the struct.
type syncStateJSON struct {
	HeadSlot     string `json:"head_slot"`
	SyncDistance string `json:"sync_distance"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncState) MarshalJSON() ([]byte, error) {
	return json.Marshal(&syncStateJSON{
		HeadSlot:     fmt.Sprintf("%d", s.HeadSlot),
		SyncDistance: fmt.Sprintf("%d", s.SyncDistance),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SyncState) UnmarshalJSON(input []byte) error {
	var err error

	var syncStateJSON syncStateJSON
	if err = json.Unmarshal(input, &syncStateJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if syncStateJSON.HeadSlot == "" {
		return errors.New("head slot missing")
	}
	if s.HeadSlot, err = strconv.ParseUint(syncStateJSON.HeadSlot, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for head slot")
	}
	if syncStateJSON.SyncDistance == "" {
		return errors.New("sync distance missing")
	}
	if s.SyncDistance, err = strconv.ParseUint(syncStateJSON.SyncDistance, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for sync distance")
	}

	return nil
}

// String returns a string version of the structure.
func (s *SyncState) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
