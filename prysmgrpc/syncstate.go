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

package prysmgrpc

import (
	"context"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// SyncState provides the state of the node's synchronization with the chain.
func (s *Service) SyncState(ctx context.Context) (*api.SyncState, error) {
	conn := ethpb.NewBeaconChainClient(s.conn)

	// Work out expected head slot.
	slot, err := s.CurrentSlot(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain current slot")
	}

	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	head, err := conn.GetChainHead(opCtx, &types.Empty{})
	cancel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain current head")
	}

	syncDistance := uint64(0)
	if head.HeadSlot < slot {
		syncDistance = slot - head.HeadSlot
	}

	return &api.SyncState{
		HeadSlot:     slot,
		SyncDistance: syncDistance,
	}, nil
}
