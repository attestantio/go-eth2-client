// Copyright © 2024 Attestant Limited.
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

package electra

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// pendingPartialWithdrawalJSON is the spec representation of the struct.
type pendingPartialWithdrawalJSON struct {
	ValidatorIndex    phase0.ValidatorIndex `json:"validator_index"`
	Amount            phase0.Gwei           `json:"amount"`
	WithdrawableEpoch phase0.Epoch          `json:"withdrawable_epoch"`
}

// MarshalJSON implements json.Marshaler.
func (p *PendingPartialWithdrawal) MarshalJSON() ([]byte, error) {
	return json.Marshal(&pendingPartialWithdrawalJSON{
		ValidatorIndex:    p.ValidatorIndex,
		Amount:            p.Amount,
		WithdrawableEpoch: p.WithdrawableEpoch,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PendingPartialWithdrawal) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&pendingPartialWithdrawalJSON{}, input)
	if err != nil {
		return err
	}

	if err := p.ValidatorIndex.UnmarshalJSON(raw["validator_index"]); err != nil {
		return errors.Wrap(err, "validator_index")
	}
	if err := p.Amount.UnmarshalJSON(raw["amount"]); err != nil {
		return errors.Wrap(err, "amount")
	}
	if err := p.WithdrawableEpoch.UnmarshalJSON(raw["withdrawable_epoch"]); err != nil {
		return errors.Wrap(err, "withdrawable_epoch")
	}

	return nil
}
