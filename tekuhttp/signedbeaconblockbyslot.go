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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// Teku wraps the beacon block.
type tekuSignedBeaconBlockJSON struct {
	SignedBlock *spec.SignedBeaconBlock `json:"beacon_block"`
}

// Teku has 'header_1' and 'header_2' for proposer slashings; spec is 'signed_header_1' and 'signed_header_2'.
var signedBeaconBlockRe1 = regexp.MustCompile(`"header_([12])"`)

// SignedBeaconBlockBySlot fetches a signed beacon block given its slot.
func (s *Service) SignedBeaconBlockBySlot(ctx context.Context, slot spec.Slot) (*spec.SignedBeaconBlock, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/beacon/block?slot=%d", slot))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request signed beacon block")
	}

	// Need to munge the data, so read it in.
	data, err := ioutil.ReadAll(respBodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read signed beacon block")
	}
	data = signedBeaconBlockRe1.ReplaceAll(data, []byte(`"signed_header_$1"`))

	response := &tekuSignedBeaconBlockJSON{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to parse signed beacon block")
	}
	block := response.SignedBlock

	// Ensure the data returned to us is as expected given our input.
	if block.Message.Slot != slot {
		if block.Message.Slot < slot {
			// If teku does not have a block in a slot it will return an earlier one; treat this as not found.
			log.Trace().Uint64("requested_slot", uint64(slot)).Uint64("returned_slot", uint64(block.Message.Slot)).Msg("Block returned for earlier slot; ignoring")
			return nil, nil
		}

		return nil, fmt.Errorf("failed to obtain correct block (requested %d, returned %d)", slot, block.Message.Slot)
	}

	return block, nil
}
