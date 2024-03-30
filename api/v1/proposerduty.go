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
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// ProposerDuty represents a duty of a validator to propose a slot.
type ProposerDuty struct {
	PubKey         phase0.BLSPubKey
	Slot           phase0.Slot
	ValidatorIndex phase0.ValidatorIndex
}

// proposerDutyJSON is the standard API representation of the struct.
type proposerDutyJSON struct {
	PubKey         string `json:"pubkey"`
	Slot           string `json:"slot"`
	ValidatorIndex string `json:"validator_index"`
}

// MarshalJSON implements json.Marshaler.
func (p *ProposerDuty) MarshalJSON() ([]byte, error) {
	return json.Marshal(&proposerDutyJSON{
		PubKey:         fmt.Sprintf("%#x", p.PubKey),
		Slot:           fmt.Sprintf("%d", p.Slot),
		ValidatorIndex: fmt.Sprintf("%d", p.ValidatorIndex),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *ProposerDuty) UnmarshalJSON(input []byte) error {
	var err error

	var proposerDutyJSON proposerDutyJSON
	if err = json.Unmarshal(input, &proposerDutyJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if proposerDutyJSON.PubKey == "" {
		return errors.New("public key missing")
	}
	pubKey, err := hex.DecodeString(strings.TrimPrefix(proposerDutyJSON.PubKey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}
	if len(pubKey) != publicKeyLength {
		return fmt.Errorf("incorrect length %d for public key", len(pubKey))
	}
	copy(p.PubKey[:], pubKey)
	if proposerDutyJSON.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(proposerDutyJSON.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	p.Slot = phase0.Slot(slot)
	if proposerDutyJSON.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(proposerDutyJSON.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	p.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	return nil
}

// String returns a string version of the structure.
func (p *ProposerDuty) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
