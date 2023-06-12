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

package deneb

import (
	"encoding/json"
	"fmt"

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/pkg/errors"
)

// beaconBlockJSON is the spec representation of the struct.
type beaconBlockJSON struct {
	Slot          string           `json:"slot"`
	ProposerIndex string           `json:"proposer_index"`
	ParentRoot    string           `json:"parent_root"`
	StateRoot     string           `json:"state_root"`
	Body          *BeaconBlockBody `json:"body"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&beaconBlockJSON{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    b.ParentRoot.String(),
		StateRoot:     b.StateRoot.String(),
		Body:          b.Body,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlock) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&beaconBlockJSON{}, input)
	if err != nil {
		return err
	}

	if err := b.Slot.UnmarshalJSON(raw["slot"]); err != nil {
		return errors.Wrap(err, "slot")
	}

	if err := b.ProposerIndex.UnmarshalJSON(raw["proposer_index"]); err != nil {
		return errors.Wrap(err, "proposer_index")
	}

	if err := b.ParentRoot.UnmarshalJSON(raw["parent_root"]); err != nil {
		return errors.Wrap(err, "parent_root")
	}

	if err := b.StateRoot.UnmarshalJSON(raw["state_root"]); err != nil {
		return errors.Wrap(err, "state_root")
	}

	b.Body = &BeaconBlockBody{}
	if err := b.Body.UnmarshalJSON(raw["body"]); err != nil {
		return errors.Wrap(err, "body")
	}

	return nil
}
