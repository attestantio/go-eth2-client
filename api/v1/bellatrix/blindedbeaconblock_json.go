// Copyright © 2022 Attestant Limited.
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

package bellatrix

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// blindedBeaconBlockJSON is the spec representation of the struct.
type blindedBeaconBlockJSON struct {
	Slot          string                  `json:"slot"`
	ProposerIndex string                  `json:"proposer_index"`
	ParentRoot    string                  `json:"parent_root"`
	StateRoot     string                  `json:"state_root"`
	Body          *BlindedBeaconBlockBody `json:"body"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlindedBeaconBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blindedBeaconBlockJSON{
		Slot:          fmt.Sprintf("%d", b.Slot),
		ProposerIndex: fmt.Sprintf("%d", b.ProposerIndex),
		ParentRoot:    fmt.Sprintf("%#x", b.ParentRoot),
		StateRoot:     fmt.Sprintf("%#x", b.StateRoot),
		Body:          b.Body,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlindedBeaconBlock) UnmarshalJSON(input []byte) error {
	var data blindedBeaconBlockJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return b.unpack(&data)
}

func (b *BlindedBeaconBlock) unpack(data *blindedBeaconBlockJSON) error {
	if data.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	b.Slot = phase0.Slot(slot)
	if data.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}
	proposerIndex, err := strconv.ParseUint(data.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}
	b.ProposerIndex = phase0.ValidatorIndex(proposerIndex)
	if data.ParentRoot == "" {
		return errors.New("parent root missing")
	}
	parentRoot, err := hex.DecodeString(strings.TrimPrefix(data.ParentRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for parent root")
	}
	if len(parentRoot) != phase0.RootLength {
		return errors.New("incorrect length for parent root")
	}
	copy(b.ParentRoot[:], parentRoot)
	if data.StateRoot == "" {
		return errors.New("state root missing")
	}
	stateRoot, err := hex.DecodeString(strings.TrimPrefix(data.StateRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for state root")
	}
	if len(stateRoot) != phase0.RootLength {
		return errors.New("incorrect length for state root")
	}
	copy(b.StateRoot[:], stateRoot)
	if data.Body == nil {
		return errors.New("body missing")
	}
	b.Body = data.Body

	return nil
}
