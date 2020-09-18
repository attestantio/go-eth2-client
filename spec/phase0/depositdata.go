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

// DepositData provides information about a deposit made on Ethereum 1.
type DepositData struct {
	PublicKey             []byte `ssz-size:"48"`
	WithdrawalCredentials []byte `ssz-size:"32"`
	Amount                uint64
	Signature             []byte `ssz-size:"96"`
}

// depositDataJSON is the spec representation of the struct.
type depositDataJSON struct {
	PublicKey             string `json:"pubkey"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	Amount                string `json:"amount"`
	Signature             string `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (d *DepositData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&depositDataJSON{
		PublicKey:             fmt.Sprintf("%#x", d.PublicKey),
		WithdrawalCredentials: fmt.Sprintf("%#x", d.WithdrawalCredentials),
		Amount:                fmt.Sprintf("%d", d.Amount),
		Signature:             fmt.Sprintf("%#x", d.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DepositData) UnmarshalJSON(input []byte) error {
	var err error

	var depositDataJSON depositDataJSON
	if err = json.Unmarshal(input, &depositDataJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if depositDataJSON.PublicKey == "" {
		return errors.New("public key missing")
	}
	if d.PublicKey, err = hex.DecodeString(strings.TrimPrefix(depositDataJSON.PublicKey, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}
	if len(d.PublicKey) != publicKeyLength {
		return errors.New("incorrect length for public key")
	}
	if depositDataJSON.WithdrawalCredentials == "" {
		return errors.New("withdrawal credentials missing")
	}
	if d.WithdrawalCredentials, err = hex.DecodeString(strings.TrimPrefix(depositDataJSON.WithdrawalCredentials, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for withdrawal credentials")
	}
	if len(d.WithdrawalCredentials) != hashLength {
		return errors.New("incorrect length for withdrawal credentials")
	}
	if depositDataJSON.Amount == "" {
		return errors.New("amount missing")
	}
	if d.Amount, err = strconv.ParseUint(depositDataJSON.Amount, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for amount")
	}
	if depositDataJSON.Signature == "" {
		return errors.New("signature missing")
	}
	if d.Signature, err = hex.DecodeString(strings.TrimPrefix(depositDataJSON.Signature, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(d.Signature) != signatureLength {
		return errors.New("incorrect length for signature")
	}

	return nil
}

// String returns a string version of the structure.
func (d *DepositData) String() string {
	data, err := json.Marshal(d)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
