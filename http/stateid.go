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

package http

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SlotFromStateID parses the state ID and returns the relevant slot.
func (s *Service) SlotFromStateID(_ context.Context, stateID string) (phase0.Slot, error) {
	var slot phase0.Slot
	switch {
	case stateID == "genesis":
		slot = 0
	case stateID == "justified":
		return 0, errors.New("state from justified not implemented")
	case stateID == "finalized":
		return 0, errors.New("state from finalized not implemented")
	case stateID == "head":
		return 0, errors.New("state from head not implemented")
	case strings.HasPrefix(stateID, "0x"):
		return 0, errors.New("state from state root not implemented")
	default:
		// State ID should be a slot.
		tmp, err := strconv.ParseUint(stateID, 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("failed to parse state %s as a slot", stateID))
		}
		slot = phase0.Slot(tmp)
	}

	return slot, nil
}

// EpochFromStateID parses the state ID and returns the relevant epoch.
func (s *Service) EpochFromStateID(ctx context.Context, stateID string) (phase0.Epoch, error) {
	var epoch phase0.Epoch
	switch {
	case stateID == "genesis":
		epoch = 0
	case stateID == "justified":
		response, err := s.Finality(ctx, &api.FinalityOpts{State: stateID})
		if err != nil {
			return 0, errors.Wrap(err, "failed to obtain finality for justified epoch")
		}
		epoch = response.Data.Justified.Epoch
	case stateID == "finalized":
		response, err := s.Finality(ctx, &api.FinalityOpts{State: stateID})
		if err != nil {
			return 0, errors.Wrap(err, "failed to obtain finality for finalized epoch")
		}
		epoch = response.Data.Justified.Epoch
	case stateID == "head":
		return 0, errors.New("epoch from head not implemented")
	case strings.HasPrefix(stateID, "0x"):
		return 0, errors.New("epoch from state root not implemented")
	default:
		// State ID should be a slot.
		tmp, err := strconv.ParseUint(stateID, 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("failed to parse state %s as a slot", stateID))
		}
		slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
		if err != nil {
			return 0, errors.Wrap(err, "failed to obtain slots per epoch")
		}
		epoch = phase0.Epoch(tmp / slotsPerEpoch)
	}

	return epoch, nil
}
