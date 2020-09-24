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

package lighthousehttp

import (
	"context"
	"encoding/json"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
)

type syncingJSON struct {
	Status *syncingStatusJSON `json:"sync_status"`
}
type syncingStatusJSON struct {
	CurrentSlot uint64 `json:"current_slot"`
}

// SyncState provides the state of the node's synchronization with the chain.
func (s *Service) SyncState(ctx context.Context) (*api.SyncState, error) {
	respBodyReader, cancel, err := s.get(ctx, "/node/syncing")
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain sync status")
	}
	defer cancel()

	// Work out expected head slot.
	slot, err := s.CurrentSlot(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain current slot")
	}

	// Fetch the sync slot.
	syncingResponse := syncingJSON{}
	if err := json.NewDecoder(respBodyReader).Decode(&syncingResponse); err != nil {
		return nil, errors.Wrap(err, "failed to parse syncing response")
	}

	syncDistance := uint64(0)
	if syncingResponse.Status.CurrentSlot < slot {
		syncDistance = slot - syncingResponse.Status.CurrentSlot
	}

	return &api.SyncState{
		HeadSlot:     slot,
		SyncDistance: syncDistance,
	}, nil
}
