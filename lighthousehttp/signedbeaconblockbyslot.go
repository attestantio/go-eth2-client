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
	"fmt"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// Lighthouse returns the block in a sub-object.
type blockResponseJSON struct {
	SignedBlock *spec.SignedBeaconBlock `json:"beacon_block"`
}

// SignedBeaconBlockBySlot fetches a signed beacon block given its slot.
func (s *Service) SignedBeaconBlockBySlot(ctx context.Context, slot uint64) (*spec.SignedBeaconBlock, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/beacon/block?slot=%d", slot))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request signed beacon block")
	}

	specReader, err := lhToSpec(ctx, respBodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert lighthouse response to spec response")
	}

	blockResponse := blockResponseJSON{}
	if err := json.NewDecoder(specReader).Decode(&blockResponse); err != nil {
		return nil, errors.Wrap(err, "failed to parse signed beacon block")
	}
	if blockResponse.SignedBlock == nil {
		// No block.
		return nil, nil
	}
	// Ensure the data returned to us is as expected given our input.
	if blockResponse.SignedBlock.Message.Slot != slot {
		if blockResponse.SignedBlock.Message.Slot < slot {
			// If lighthouse does not have a block in a slot it will return an earlier one; treat this as not found.
			log.Trace().Uint64("requested_slot", slot).Uint64("returned_slot", blockResponse.SignedBlock.Message.Slot).Msg("Block returned for earlier slot; ignoring")
			return nil, nil
		}

		return nil, fmt.Errorf("failed to obtain correct block (requested %d, returned %d)", slot, blockResponse.SignedBlock.Message.Slot)
	}

	return blockResponse.SignedBlock, nil
}
