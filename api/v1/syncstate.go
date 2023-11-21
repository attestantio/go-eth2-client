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
	"strconv"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SyncState is the data regarding the node's synchronization state to the chain.
type SyncState struct {
	// HeadSlot is the head slot of the chain as understood by the node.
	HeadSlot phase0.Slot
	// SyncDistance is the distance between the node's highest synced slot and the head slot.
	SyncDistance phase0.Slot
	// IsOptimistic is true if the node is optimistic.
	IsOptimistic bool
	// IsSyncing is true if the node is syncing.
	IsSyncing bool
}

// syncStateJSON is the spec representation of the struct.
type syncStateJSON struct {
	HeadSlot     string `json:"head_slot"`
	SyncDistance string `json:"sync_distance"`
	IsOptimistic bool   `json:"is_optimistic"`
	IsSyncing    bool   `json:"is_syncing"`
}

// MarshalJSON implements json.Marshaler.
func (s *SyncState) MarshalJSON() ([]byte, error) {
	return json.Marshal(&syncStateJSON{
		HeadSlot:     fmt.Sprintf("%d", s.HeadSlot),
		SyncDistance: fmt.Sprintf("%d", s.SyncDistance),
		IsOptimistic: s.IsOptimistic,
		IsSyncing:    s.IsSyncing,
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
	headSlot, err := strconv.ParseUint(syncStateJSON.HeadSlot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for head slot")
	}
	s.HeadSlot = phase0.Slot(headSlot)
	if syncStateJSON.SyncDistance == "" {
		return errors.New("sync distance missing")
	}
	syncDistance, err := strconv.ParseUint(syncStateJSON.SyncDistance, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for sync distance")
	}
	s.SyncDistance = phase0.Slot(syncDistance)
	s.IsOptimistic = syncStateJSON.IsOptimistic
	s.IsSyncing = syncStateJSON.IsSyncing

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
