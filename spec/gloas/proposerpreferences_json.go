// Copyright © 2025 Attestant Limited.
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

package gloas

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// proposerPreferencesJSON is the spec representation of the struct.
type proposerPreferencesJSON struct {
	ProposalSlot   string `json:"proposal_slot"`
	ValidatorIndex string `json:"validator_index"`
	FeeRecipient   string `json:"fee_recipient"`
	GasLimit       string `json:"gas_limit"`
}

// MarshalJSON implements json.Marshaler.
func (p *ProposerPreferences) MarshalJSON() ([]byte, error) {
	return json.Marshal(&proposerPreferencesJSON{
		ProposalSlot:   fmt.Sprintf("%d", p.ProposalSlot),
		ValidatorIndex: fmt.Sprintf("%d", p.ValidatorIndex),
		FeeRecipient:   fmt.Sprintf("%#x", p.FeeRecipient),
		GasLimit:       fmt.Sprintf("%d", p.GasLimit),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *ProposerPreferences) UnmarshalJSON(input []byte) error {
	var data proposerPreferencesJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	// Proposal slot.
	if data.ProposalSlot == "" {
		return errors.New("proposal slot missing")
	}
	proposalSlot, err := strconv.ParseUint(data.ProposalSlot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid proposal slot")
	}
	p.ProposalSlot = phase0.Slot(proposalSlot)

	// Validator index.
	if data.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(data.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid validator index")
	}
	p.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	// Fee recipient.
	if data.FeeRecipient == "" {
		return errors.New("fee recipient missing")
	}
	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.FeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid fee recipient")
	}
	copy(p.FeeRecipient[:], feeRecipient)

	// Gas limit.
	if data.GasLimit == "" {
		return errors.New("gas limit missing")
	}
	gasLimit, err := strconv.ParseUint(data.GasLimit, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid gas limit")
	}
	p.GasLimit = gasLimit

	return nil
}
