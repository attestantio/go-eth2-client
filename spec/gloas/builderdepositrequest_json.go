// Copyright © 2026 Attestant Limited.
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

// builderDepositRequestJSON is the spec representation of the struct.
type builderDepositRequestJSON struct {
	Pubkey                string `json:"pubkey"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	Amount                string `json:"amount"`
	Signature             string `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (b *BuilderDepositRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&builderDepositRequestJSON{
		Pubkey:                fmt.Sprintf("%#x", b.Pubkey),
		WithdrawalCredentials: fmt.Sprintf("%#x", b.WithdrawalCredentials),
		Amount:                fmt.Sprintf("%d", b.Amount),
		Signature:             fmt.Sprintf("%#x", b.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BuilderDepositRequest) UnmarshalJSON(input []byte) error {
	var data builderDepositRequestJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return b.unpack(&data)
}

func (b *BuilderDepositRequest) unpack(data *builderDepositRequestJSON) error {
	if data.Pubkey == "" {
		return errors.New("public key missing")
	}

	pubkey, err := hex.DecodeString(strings.TrimPrefix(data.Pubkey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}

	if len(pubkey) != phase0.PublicKeyLength {
		return errors.New("incorrect length for public key")
	}

	copy(b.Pubkey[:], pubkey)

	if data.WithdrawalCredentials == "" {
		return errors.New("withdrawal credentials missing")
	}

	if b.WithdrawalCredentials, err = hex.DecodeString(strings.TrimPrefix(data.WithdrawalCredentials, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for withdrawal credentials")
	}

	if len(b.WithdrawalCredentials) != phase0.HashLength {
		return errors.New("incorrect length for withdrawal credentials")
	}

	if data.Amount == "" {
		return errors.New("amount missing")
	}

	amount, err := strconv.ParseUint(data.Amount, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for amount")
	}

	b.Amount = phase0.Gwei(amount)

	if data.Signature == "" {
		return errors.New("signature missing")
	}

	signature, err := hex.DecodeString(strings.TrimPrefix(data.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}

	if len(signature) != phase0.SignatureLength {
		return errors.New("incorrect length for signature")
	}

	copy(b.Signature[:], signature)

	return nil
}
