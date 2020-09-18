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

// AttestationData is the Ethereum 2 specification structure.
type AttestationData struct {
	Slot uint64
	// Index is the committee index.
	Index           uint64
	BeaconBlockRoot []byte `ssz-size:"32"`
	Source          *Checkpoint
	Target          *Checkpoint
}

// attestationDataJSON is the spec representation of the struct.
type attestationDataJSON struct {
	Slot            string      `json:"slot"`
	Index           string      `json:"index"`
	BeaconBlockRoot string      `json:"beacon_block_root"`
	Source          *Checkpoint `json:"source"`
	Target          *Checkpoint `json:"target"`
}

// MarshalJSON implements json.Marshaler.
func (a *AttestationData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&attestationDataJSON{
		Slot:            fmt.Sprintf("%d", a.Slot),
		Index:           fmt.Sprintf("%d", a.Index),
		BeaconBlockRoot: fmt.Sprintf("%#x", a.BeaconBlockRoot),
		Source:          a.Source,
		Target:          a.Target,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AttestationData) UnmarshalJSON(input []byte) error {
	var err error

	var attestationDataJSON attestationDataJSON
	if err = json.Unmarshal(input, &attestationDataJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if a.Slot, err = strconv.ParseUint(attestationDataJSON.Slot, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	if a.Index, err = strconv.ParseUint(attestationDataJSON.Index, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for index")
	}
	if a.BeaconBlockRoot, err = hex.DecodeString(strings.TrimPrefix(attestationDataJSON.BeaconBlockRoot, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for beacon block root")
	}
	if len(a.BeaconBlockRoot) != rootLength {
		return errors.New("incorrect length for beacon block root")
	}
	a.Source = attestationDataJSON.Source
	if a.Source == nil {
		return errors.New("source missing")
	}
	a.Target = attestationDataJSON.Target
	if a.Target == nil {
		return errors.New("target missing")
	}

	return nil
}

func (a *AttestationData) String() string {
	data, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
