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

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// EpochFromStateID obtains the epoch given the state ID.
func (s *Service) EpochFromStateID(ctx context.Context, stateID string) (spec.Epoch, error) {
	if stateID == "head" {
		epoch, err := s.CurrentEpoch(ctx)
		return spec.Epoch(epoch), err
	}

	slot, err := s.SlotFromStateID(ctx, stateID)
	if err != nil {
		return 0, err
	}
	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return 0, err
	}
	return spec.Epoch(uint64(slot) / slotsPerEpoch), nil
}

// SlotFromStateID obtains the slot given the state ID.
func (s *Service) SlotFromStateID(ctx context.Context, stateID string) (spec.Slot, error) {
	var slot spec.Slot
	switch {
	case stateID == "genesis":
		slot = 0
	case stateID == "justified":
		chainHead, err := s.beaconHead(ctx)
		if err != nil {
			return 0, errors.New("failed to obtain chain head")
		}
		return spec.Slot(chainHead.JustifiedSlot), nil
	case stateID == "finalized":
		chainHead, err := s.beaconHead(ctx)
		if err != nil {
			return 0, errors.New("failed to obtain chain head")
		}
		return spec.Slot(chainHead.FinalizedSlot), nil
	case stateID == "head":
		chainHead, err := s.beaconHead(ctx)
		if err != nil {
			return 0, errors.New("failed to obtain chain head")
		}
		return spec.Slot(chainHead.HeadSlot), nil
	case strings.HasPrefix(stateID, "0x"):
		// State root.
		return 0, errors.New("not implemented")
	default:
		// Slot.
		tmp, err := strconv.ParseUint(stateID, 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("failed to parse state ID %s as a slot", stateID))
		}
		slot = spec.Slot(tmp)
	}
	log.Trace().Str("state", stateID).Uint64("slot", uint64(slot)).Msg("Calculated from state ID")
	return slot, nil
}

func (s *Service) beaconHead(ctx context.Context) (*ethpb.ChainHead, error) {
	conn := ethpb.NewBeaconChainClient(s.conn)
	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	head, err := conn.GetChainHead(opCtx, &types.Empty{})
	cancel()

	return head, err
}
