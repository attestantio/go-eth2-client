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

package phase0

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// BeaconBlock represents a beacon block.
type BeaconBlock struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    []byte `ssz-size:"32"`
	StateRoot     []byte `ssz-size:"32"`
	Body          *BeaconBlockBody
}

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
		ParentRoot:    fmt.Sprintf("%#x", b.ParentRoot),
		StateRoot:     fmt.Sprintf("%#x", b.StateRoot),
		Body:          b.Body,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlock) UnmarshalJSON(input []byte) error {
	var err error

	var beaconBlockJSON beaconBlockJSON
	if err = json.Unmarshal(input, &beaconBlockJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if beaconBlockJSON.Slot == "" {
		return errors.New("slot missing")
	}
	if b.Slot, err = strconv.ParseUint(beaconBlockJSON.Slot, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	if beaconBlockJSON.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}
	if b.ProposerIndex, err = strconv.ParseUint(beaconBlockJSON.ProposerIndex, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}
	if beaconBlockJSON.ParentRoot == "" {
		return errors.New("parent root missing")
	}
	if b.ParentRoot, err = hex.DecodeString(strings.TrimPrefix(beaconBlockJSON.ParentRoot, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for parent root")
	}
	if len(b.ParentRoot) != rootLength {
		return errors.New("incorrect length for parent root")
	}
	if beaconBlockJSON.StateRoot == "" {
		return errors.New("state root missing")
	}
	if b.StateRoot, err = hex.DecodeString(strings.TrimPrefix(beaconBlockJSON.StateRoot, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for state root")
	}
	if len(b.StateRoot) != rootLength {
		return errors.New("incorrect length for state root")
	}
	if beaconBlockJSON.Body == nil {
		return errors.New("body missing")
	}
	b.Body = beaconBlockJSON.Body

	return nil
}

// String returns a string version of the structure.
func (b *BeaconBlock) String() string {
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
