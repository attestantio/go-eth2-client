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
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// epochFromStateID obtains the epoch given the state ID.
func (s *Service) epochFromStateID(ctx context.Context, stateID string) (uint64, error) {
	slot, err := s.slotFromStateID(ctx, stateID)
	if err != nil {
		return 0, err
	}
	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return 0, err
	}
	return slot / slotsPerEpoch, nil
}

// slotFromStateID obtains the slot given the state ID.
func (s *Service) slotFromStateID(ctx context.Context, stateID string) (uint64, error) {
	var slot uint64
	var err error
	switch {
	case stateID == "genesis":
		slot = 0
	case stateID == "justified":
		// TODO
		// Fetch epoch from /beacon/{stateid}/finality_checkpoints
		// Epoch -> slot
	case stateID == "finalized":
		// TODO
		// Fetch epoch from /beacon/{stateid}/finality_checkpoints
		// Epoch -> slot
	case stateID == "head":
		// Canonical head in node's view.
		// fetch from /node/syncing
		// TODO
	case strings.HasPrefix(stateID, "0x"):
		// State root.
		// TODO
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

// parseStateID parses the state ID and returns the relevant state root, block root and slot.
func (s *Service) parseStateID(ctx context.Context, stateID string) ([]byte, []byte, uint64, error) {
	var stateRoot []byte
	var blockRoot []byte
	var slot uint64
	var err error
	switch {
	case stateID == "genesis":
		signedBeaconBlock, err := s.SignedBeaconBlockBySlot(ctx, 0)
		if err != nil {
			return nil, nil, 0, errors.Wrap(err, "failed to obtain genesis beacon block")
		}
		if signedBeaconBlock == nil {
			return nil, nil, 0, errors.New("failed to fetch genesis beacon block")
		}
		stateRoot = signedBeaconBlock.Message.StateRoot
		root, err := signedBeaconBlock.Message.Body.HashTreeRoot()
		if err != nil {
			return nil, nil, 0, errors.Wrap(err, "failed to calculate hash tree root of beacon block")
		}
		blockRoot = root[:]
	case stateID == "justified":
		// TODO fetch stateRoot, bodyRoot, slot
	case stateID == "finalized":
		// TODO fetch stateRoot, bodyRoot, slot
	case stateID == "head":
		// TODO fetch stateRoot, bodyRoot, slot
	case strings.HasPrefix(stateID, "0x"):
		stateRoot, err = hex.DecodeString(strings.TrimPrefix(stateID, "0x"))
		if err != nil {
			return nil, nil, 0, errors.Wrap(err, fmt.Sprintf("failed to parse state ID %s as a state root", stateID))
		}
		// TODO bodyRoot, slot
	default:
		// State ID should be a slot.
		slot, err = strconv.ParseUint(stateID, 10, 64)
		if err != nil {
			return nil, nil, 0, errors.Wrap(err, fmt.Sprintf("failed to parse state ID %s as a slot", stateID))
		}
		signedBeaconBlock, err := s.SignedBeaconBlockBySlot(ctx, slot)
		if err != nil {
			return nil, nil, 0, errors.Wrap(err, "failed to obtain beacon block")
		}
		if signedBeaconBlock == nil {
			// No block for this slot, but can return slot.
			return nil, nil, slot, nil
		}
		stateRoot = signedBeaconBlock.Message.StateRoot
		root, err := signedBeaconBlock.Message.Body.HashTreeRoot()
		if err != nil {
			return nil, nil, 0, errors.Wrap(err, "failed to calculate hash tree root of beacon block")
		}
		blockRoot = root[:]
	}
	log.Trace().Str("state", stateID).Str("state_root", fmt.Sprintf("%#x", stateRoot)).Str("block_root", fmt.Sprintf("%#x", blockRoot)).Uint64("slot", slot).Msg("Calculated from state ID")
	return stateRoot, blockRoot, slot, nil
}
