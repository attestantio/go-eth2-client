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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// payloadAttestationDataJSON is the spec representation of the struct.
type payloadAttestationDataJSON struct {
	BeaconBlockRoot string `json:"beacon_block_root"`
	Slot            string `json:"slot"`
	PayloadStatus   string `json:"payload_status"`
}

// MarshalJSON implements json.Marshaler.
func (p *PayloadAttestationData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&payloadAttestationDataJSON{
		BeaconBlockRoot: fmt.Sprintf("%#x", p.BeaconBlockRoot),
		Slot:            fmt.Sprintf("%d", p.Slot),
		PayloadStatus:   fmt.Sprintf("%d", p.PayloadStatus),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttestationData) UnmarshalJSON(input []byte) error {
	var data payloadAttestationDataJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.BeaconBlockRoot == "" {
		return errors.New("beacon block root missing")
	}
	root, err := hex.DecodeString(strings.TrimPrefix(data.BeaconBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid beacon block root")
	}
	copy(p.BeaconBlockRoot[:], root)

	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid slot")
	}
	p.Slot = phase0.Slot(slot)

	status, err := strconv.ParseUint(data.PayloadStatus, 10, 8)
	if err != nil {
		return errors.Wrap(err, "invalid payload status")
	}
	p.PayloadStatus = uint8(status)

	return nil
}
