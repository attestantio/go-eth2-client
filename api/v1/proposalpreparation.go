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

package v1

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// ProposalPreparation is the data required for proposal preparation.
type ProposalPreparation struct {
	// ValidatorIdex is the index of the validator making the proposal request.
	ValidatorIndex phase0.ValidatorIndex
	// FeeRecipient is the execution address to be used with preparing blocks.
	FeeRecipient bellatrix.ExecutionAddress `ssz-size:"20"`
}

// proposalPreparationJSON is the spec representation of the struct.
type proposalPreparationJSON struct {
	ValidatorIndex string `json:"validator_index"`
	FeeRecipient   string `json:"fee_recipient"`
}

// MarshalJSON implements json.Marshaler.
func (p *ProposalPreparation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&proposalPreparationJSON{
		ValidatorIndex: fmt.Sprintf("%d", p.ValidatorIndex),
		FeeRecipient:   p.FeeRecipient.String(),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *ProposalPreparation) UnmarshalJSON(input []byte) error {
	var err error

	var data proposalPreparationJSON
	if err = json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.ValidatorIndex == "" {
		return errors.New("validator index missing")
	}
	validatorIndex, err := strconv.ParseUint(data.ValidatorIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for validator index")
	}
	p.ValidatorIndex = phase0.ValidatorIndex(validatorIndex)

	if data.FeeRecipient == "" {
		return errors.New("fee recipient is missing")
	}
	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.FeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for fee recipient")
	}
	copy(p.FeeRecipient[:], feeRecipient)

	return nil
}

// String returns a string version of the structure.
func (p *ProposalPreparation) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
