// Copyright © 2023 Attestant Limited.
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

package eip7732

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// beaconBlockJSON is the spec representation of the struct.
type beaconBlockJSON struct {
	Slot          string           `json:"slot"`
	ProposerIndex string           `json:"proposer_index"`
	ParentRoot    phase0.Root      `json:"parent_root"`
	StateRoot     phase0.Root      `json:"state_root"`
	Body          *BeaconBlockBody `json:"body"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&beaconBlockJSON{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    b.ParentRoot,
		StateRoot:     b.StateRoot,
		Body:          b.Body,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlock) UnmarshalJSON(input []byte) error {
	var data beaconBlockJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid slot")
	}
	b.Slot = phase0.Slot(slot)

	proposerIndex, err := strconv.ParseUint(data.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid proposer index")
	}
	b.ProposerIndex = phase0.ValidatorIndex(proposerIndex)

	b.ParentRoot = data.ParentRoot
	b.StateRoot = data.StateRoot
	b.Body = data.Body

	return nil
}
