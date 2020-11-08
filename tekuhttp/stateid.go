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

package tekuhttp

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// SlotFromStateID parses the state ID and returns the relevant slot.
func (s *Service) SlotFromStateID(ctx context.Context, stateID string) (uint64, error) {
	var slot uint64
	var err error
	switch {
	case stateID == "genesis":
		slot = 0
	case stateID == "justified":
		head, err := s.chainHead(ctx)
		if err != nil {
			return 0, errors.Wrap(err, "failed to obtain justified slot from beacon head")
		}
		slot = head.JustifiedSlot
	case stateID == "finalized":
		head, err := s.chainHead(ctx)
		if err != nil {
			return 0, errors.Wrap(err, "failed to obtain finalized slot from beacon head")
		}
		slot = head.FinalizedSlot
	case stateID == "head":
		head, err := s.chainHead(ctx)
		if err != nil {
			return 0, errors.Wrap(err, "failed to obtain head slot from beacon head")
		}
		slot = head.Slot
	case strings.HasPrefix(stateID, "0x"):
		stateRoot, err := hex.DecodeString(strings.TrimPrefix(stateID, "0x"))
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("failed to parse state ID %s as a state root", stateID))
		}
		slot, err = s.stateToSlot(ctx, stateRoot)
		if err != nil {
			return 0, errors.Wrap(err, "failed to obtain slot from state")
		}
	default:
		// State ID should be a slot.
		slot, err = strconv.ParseUint(stateID, 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("failed to parse state ID %s as a slot", stateID))
		}
	}

	log.Trace().Str("state", stateID).Uint64("slot", slot).Msg("Calculated from state ID")
	return slot, nil
}

// EpochFromStateID parses the state ID and returns the relevant epoch.
func (s *Service) EpochFromStateID(ctx context.Context, stateID string) (uint64, error) {
	slot, err := s.SlotFromStateID(ctx, stateID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to obtain slot for state ID")
	}
	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to obtain slot per epoch")
	}
	return slot / slotsPerEpoch, nil
}

// stateRootFromStateID parses the state ID and returns the relevant state root.
// func (s *Service) stateRootFromStateID(ctx context.Context, stateID string) ([]byte, error) {
//  	var stateRoot []byte
//  	var err error
//  	switch {
//  	case stateID == "genesis":
//  		signedBeaconBlock, err := s.SignedBeaconBlockBySlot(ctx, 0)
//  		if err != nil {
//  			return nil, errors.Wrap(err, "failed to obtain genesis beacon block")
//  		}
//  		if signedBeaconBlock == nil {
//  			return nil, errors.New("failed to fetch genesis beacon block")
//  		}
//  		stateRoot = signedBeaconBlock.Message.StateRoot
//  	case stateID == "justified":
//  		head, err := s.chainHead(ctx)
//  		if err != nil {
//  			return nil, errors.Wrap(err, "failed to obtain beacon head for justified state")
//  		}
//  		signedBeaconBlock, err := s.SignedBeaconBlockBySlot(ctx, head.JustifiedSlot)
//  		if err != nil {
//  			return nil, errors.Wrap(err, "failed to obtain justified beacon block")
//  		}
//  		if signedBeaconBlock == nil {
//  			return nil, errors.New("failed to fetch justified beacon block")
//  		}
//  		stateRoot = signedBeaconBlock.Message.StateRoot
//  	case stateID == "finalized":
//  		head, err := s.chainHead(ctx)
//  		if err != nil {
//  			return nil, errors.Wrap(err, "failed to obtain beacon head for finalized state")
//  		}
//  		signedBeaconBlock, err := s.SignedBeaconBlockBySlot(ctx, head.FinalizedSlot)
//  		if err != nil {
//  			return nil, errors.Wrap(err, "failed to obtain finalized beacon block")
//  		}
//  		if signedBeaconBlock == nil {
//  			return nil, errors.New("failed to fetch finalized beacon block")
//  		}
//  		stateRoot = signedBeaconBlock.Message.StateRoot
//  	case stateID == "head":
//  		head, err := s.beaconHead(ctx)
//  		if err != nil {
//  			return nil, errors.Wrap(err, "failed to obtain state from beacon head")
//  		}
//  		stateRoot = head.StateRoot
//  	case strings.HasPrefix(stateID, "0x"):
//  		stateRoot, err = hex.DecodeString(strings.TrimPrefix(stateID, "0x"))
//  		if err != nil {
//  			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse state ID %s as a state root", stateID))
//  		}
//  	default:
//  		// State ID should be a slot.
//  		slot, err := strconv.ParseUint(stateID, 10, 64)
//  		if err != nil {
//  			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse state ID %s as a slot", stateID))
//  		}
//  		signedBeaconBlock, err := s.SignedBeaconBlockBySlot(ctx, slot)
//  		if err != nil {
//  			return nil, errors.Wrap(err, "failed to obtain genesis beacon block")
//  		}
//  		if signedBeaconBlock == nil {
//  			return nil, errors.New("failed to fetch genesis beacon block")
//  		}
//  		stateRoot = signedBeaconBlock.Message.StateRoot
//  	}
//
//  	log.Trace().Str("state", stateID).Str("state_root", fmt.Sprintf("%#x", stateRoot)).Msg("Calculated from state ID")
//  	return stateRoot, nil
//  }

type chainHead struct {
	Slot                       uint64
	BlockRoot                  []byte
	FinalizedSlot              uint64
	FinalizedBlockRoot         []byte
	JustifiedSlot              uint64
	JustifiedBlockRoot         []byte
	PreviousJustifiedSlot      uint64
	PreviousJustifiedBlockRoot []byte
}

type chainHeadJSON struct {
	HeadSlot                   string `json:"head_slot"`
	HeadBlockRoot              string `json:"head_block_root"`
	FinalizedSlot              string `json:"finalized_slot"`
	FinalizedBlockRoot         string `json:"finalized_block"`
	JustifiedSlot              string `json:"justified_slot"`
	JustifiedBlockRoot         string `json:"justified_block"`
	PreviousJustifiedSlot      string `json:"previous_justified_slot"`
	PreviousJustifiedBlockRoot string `json:"previous_justified_block"`
}

func (s *Service) chainHead(ctx context.Context) (*chainHead, error) {
	respBodyReader, err := s.get(ctx, "/beacon/chainhead")
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon head")
	}

	chainHeadResponse := chainHeadJSON{}
	if err := json.NewDecoder(respBodyReader).Decode(&chainHeadResponse); err != nil {
		return nil, errors.Wrap(err, "failed to parse signed beacon head")
	}

	slot, err := strconv.ParseUint(chainHeadResponse.HeadSlot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode slot")
	}
	blockRoot, err := hex.DecodeString(strings.TrimPrefix(chainHeadResponse.HeadBlockRoot, "0x"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode block root")
	}
	finalizedSlot, err := strconv.ParseUint(chainHeadResponse.FinalizedSlot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode finalized slot")
	}
	finalizedBlockRoot, err := hex.DecodeString(strings.TrimPrefix(chainHeadResponse.FinalizedBlockRoot, "0x"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode finalized block root")
	}
	justifiedSlot, err := strconv.ParseUint(chainHeadResponse.JustifiedSlot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode justified slot")
	}
	justifiedBlockRoot, err := hex.DecodeString(strings.TrimPrefix(chainHeadResponse.JustifiedBlockRoot, "0x"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode justified block root")
	}
	previousJustifiedSlot, err := strconv.ParseUint(chainHeadResponse.PreviousJustifiedSlot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode previous justified slot")
	}
	previousJustifiedBlockRoot, err := hex.DecodeString(strings.TrimPrefix(chainHeadResponse.PreviousJustifiedBlockRoot, "0x"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode previous justified block root")
	}

	return &chainHead{
		Slot:                       slot,
		BlockRoot:                  blockRoot,
		FinalizedSlot:              finalizedSlot,
		FinalizedBlockRoot:         finalizedBlockRoot,
		JustifiedSlot:              justifiedSlot,
		JustifiedBlockRoot:         justifiedBlockRoot,
		PreviousJustifiedSlot:      previousJustifiedSlot,
		PreviousJustifiedBlockRoot: previousJustifiedBlockRoot,
	}, nil
}

type beaconHead struct {
	Slot      uint64
	BlockRoot []byte
	StateRoot []byte
}

type beaconHeadJSON struct {
	Slot      string `json:"slot"`
	BlockRoot string `json:"block_root"`
	StateRoot string `json:"state_root"`
}

func (s *Service) beaconHead(ctx context.Context) (*beaconHead, error) {
	respBodyReader, err := s.get(ctx, "/beacon/head")
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon head")
	}

	beaconHeadResponse := beaconHeadJSON{}
	if err := json.NewDecoder(respBodyReader).Decode(&beaconHeadResponse); err != nil {
		return nil, errors.Wrap(err, "failed to parse signed beacon head")
	}

	slot, err := strconv.ParseUint(beaconHeadResponse.Slot, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode slot")
	}
	blockRoot, err := hex.DecodeString(strings.TrimPrefix(beaconHeadResponse.BlockRoot, "0x"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode block root")
	}
	stateRoot, err := hex.DecodeString(strings.TrimPrefix(beaconHeadResponse.StateRoot, "0x"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode state root")
	}

	return &beaconHead{
		Slot:      slot,
		BlockRoot: blockRoot,
		StateRoot: stateRoot,
	}, nil
}

type beaconStateJSON struct {
	Slot uint64 `json:"slot"`
}

func (s *Service) stateToSlot(ctx context.Context, stateRoot []byte) (uint64, error) {
	respBodyReader, err := s.get(ctx, "/beacon/state")
	if err != nil {
		return 0, errors.Wrap(err, "failed to request beacon state")
	}

	beaconStateResponse := beaconStateJSON{}
	if err := json.NewDecoder(respBodyReader).Decode(&beaconStateResponse); err != nil {
		return 0, errors.Wrap(err, "failed to parse signed beacon state")
	}

	return beaconStateResponse.Slot, nil
}
