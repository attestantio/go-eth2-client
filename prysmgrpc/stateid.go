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
	"fmt"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// EpochFromStateID obtains the epoch given the state ID.
func (s *Service) EpochFromStateID(ctx context.Context, stateID string) (uint64, error) {
	slot, err := s.SlotFromStateID(ctx, stateID)
	if err != nil {
		return 0, err
	}
	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return 0, err
	}
	return slot / slotsPerEpoch, nil
}

// SlotFromStateID obtains the slot given the state ID.
func (s *Service) SlotFromStateID(ctx context.Context, stateID string) (uint64, error) {
	var slot uint64
	var err error
	switch {
	case stateID == "genesis":
		slot = 0
	case stateID == "justified":
		chainHead, err := s.beaconHead(ctx)
		if err != nil {
			return 0, errors.New("failed to obtain chain head")
		}
		return chainHead.JustifiedSlot, nil
	case stateID == "finalized":
		chainHead, err := s.beaconHead(ctx)
		if err != nil {
			return 0, errors.New("failed to obtain chain head")
		}
		return chainHead.FinalizedSlot, nil
	case stateID == "head":
		chainHead, err := s.beaconHead(ctx)
		if err != nil {
			return 0, errors.New("failed to obtain chain head")
		}
		return chainHead.HeadSlot, nil
	case strings.HasPrefix(stateID, "0x"):
		// State root.
		return 0, errors.New("not implemented")
	default:
		// Slot.
		slot, err = strconv.ParseUint(stateID, 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("failed to parse state ID %s as a slot", stateID))
		}
	}
	log.Trace().Str("state", stateID).Uint64("slot", slot).Msg("Calculated from state ID")
	return slot, nil
}

func (s *Service) beaconHead(ctx context.Context) (*ethpb.ChainHead, error) {
	beaconChainClient := ethpb.NewBeaconChainClient(s.conn)
	return beaconChainClient.GetChainHead(ctx, &types.Empty{})
}
