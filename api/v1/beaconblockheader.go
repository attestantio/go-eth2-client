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

package v1

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// BeaconBlockHeader is the data providing information about beacon blocks.
type BeaconBlockHeader struct {
	// Root is the root of the beacon block.
	Root phase0.Root
	// Canonical is true if the block is considered canonical.
	Canonical bool
	// Header is the beacon block header.
	Header *phase0.SignedBeaconBlockHeader
}

// beaconBlockHeaderJSON is the spec representation of the struct.
type beaconBlockHeaderJSON struct {
	Root      string                          `json:"root"`
	Canonical bool                            `json:"canonical"`
	Header    *phase0.SignedBeaconBlockHeader `json:"header"`
}

// MarshalJSON implements json.Marshaler.
func (b *BeaconBlockHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(&beaconBlockHeaderJSON{
		Root:      fmt.Sprintf("%#x", b.Root),
		Canonical: b.Canonical,
		Header:    b.Header,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BeaconBlockHeader) UnmarshalJSON(input []byte) error {
	var err error

	var beaconBlockHeaderJSON beaconBlockHeaderJSON
	if err = json.Unmarshal(input, &beaconBlockHeaderJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if beaconBlockHeaderJSON.Root == "" {
		return errors.New("root missing")
	}
	root, err := hex.DecodeString(strings.TrimPrefix(beaconBlockHeaderJSON.Root, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for root")
	}
	if len(root) != rootLength {
		return fmt.Errorf("incorrect length %d for root", len(root))
	}
	copy(b.Root[:], root)

	b.Canonical = beaconBlockHeaderJSON.Canonical
	if beaconBlockHeaderJSON.Header == nil {
		return errors.New("header missing")
	}
	b.Header = beaconBlockHeaderJSON.Header

	return nil
}

// String returns a string version of the structure.
func (b *BeaconBlockHeader) String() string {
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
