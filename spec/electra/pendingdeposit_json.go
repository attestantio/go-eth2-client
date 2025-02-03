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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// pendingDepositJSON is the spec representation of the struct.
type pendingDepositJSON struct {
	Pubkey                string      `json:"pubkey"`
	WithdrawalCredentials string      `json:"withdrawal_credentials"`
	Amount                phase0.Gwei `json:"amount"`
	Signature             string      `json:"signature"`
	Slot                  phase0.Slot `json:"slot"`
}

// MarshalJSON implements json.Marshaler.
func (p *PendingDeposit) MarshalJSON() ([]byte, error) {
	return json.Marshal(&pendingDepositJSON{
		Pubkey:                fmt.Sprintf("%#x", p.Pubkey),
		WithdrawalCredentials: fmt.Sprintf("%#x", p.WithdrawalCredentials),
		Amount:                p.Amount,
		Signature:             fmt.Sprintf("%#x", p.Signature),
		Slot:                  p.Slot,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PendingDeposit) UnmarshalJSON(input []byte) error {
	var pendingDeposit pendingDepositJSON
	if err := json.Unmarshal(input, &pendingDeposit); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&pendingDeposit)
}

func (p *PendingDeposit) unpack(pendingDeposit *pendingDepositJSON) error {
	if pendingDeposit.Pubkey == "" {
		return errors.New("public key missing")
	}
	pubkey, err := hex.DecodeString(strings.TrimPrefix(pendingDeposit.Pubkey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}
	if len(pubkey) != phase0.PublicKeyLength {
		return errors.New("incorrect length for public key")
	}
	copy(p.Pubkey[:], pubkey)

	if pendingDeposit.WithdrawalCredentials == "" {
		return errors.New("withdrawal credentials missing")
	}
	if p.WithdrawalCredentials, err = hex.DecodeString(strings.TrimPrefix(pendingDeposit.WithdrawalCredentials, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for withdrawal credentials")
	}
	if len(p.WithdrawalCredentials) != phase0.HashLength {
		return errors.New("incorrect length for withdrawal credentials")
	}

	p.Amount = pendingDeposit.Amount

	if pendingDeposit.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(pendingDeposit.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(signature) != phase0.SignatureLength {
		return errors.New("incorrect length for signature")
	}
	copy(p.Signature[:], signature)

	p.Slot = pendingDeposit.Slot

	return nil
}
