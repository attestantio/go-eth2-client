// Copyright © 2021, 2023 Attestant Limited.
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

package mock

import (
	"context"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// NodeSyncing provides the state of the node's synchronization with the chain.
func (s *Service) NodeSyncing(ctx context.Context,
	opts *api.NodeSyncingOpts,
) (
	*api.Response[*apiv1.SyncState],
	error,
) {
	if s.NodeSyncingFunc != nil {
		return s.NodeSyncingFunc(ctx, opts)
	}

	return &api.Response[*apiv1.SyncState]{
		Data: &apiv1.SyncState{
			HeadSlot:     s.HeadSlot,
			SyncDistance: s.SyncDistance,
			IsSyncing:    s.SyncDistance > 0,
		},
		Metadata: make(map[string]any),
	}, nil
}
