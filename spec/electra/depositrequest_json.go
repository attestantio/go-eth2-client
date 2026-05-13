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
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// depositRequestJSON is the spec representation of the struct.
// Amount and Index are typed as json.RawMessage so we accept both the
// spec form ("42") and the bare-number form (42) that Erigon's Caplin
// emits in execution_requests.
type depositRequestJSON struct {
	Pubkey                string          `json:"pubkey"`
	WithdrawalCredentials string          `json:"withdrawal_credentials"`
	Amount                json.RawMessage `json:"amount"`
	Signature             string          `json:"signature"`
	Index                 json.RawMessage `json:"index"`
}

// MarshalJSON implements json.Marshaler.
func (d *DepositRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&depositRequestJSON{
		Pubkey:                fmt.Sprintf("%#x", d.Pubkey),
		WithdrawalCredentials: fmt.Sprintf("%#x", d.WithdrawalCredentials),
		Amount:                json.RawMessage(fmt.Sprintf(`"%d"`, d.Amount)),
		Signature:             fmt.Sprintf("%#x", d.Signature),
		Index:                 json.RawMessage(fmt.Sprintf(`"%d"`, d.Index)),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DepositRequest) UnmarshalJSON(input []byte) error {
	var depositReceipt depositRequestJSON
	if err := json.Unmarshal(input, &depositReceipt); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return d.unpack(&depositReceipt)
}

func (d *DepositRequest) unpack(depositReceipt *depositRequestJSON) error {
	if depositReceipt.Pubkey == "" {
		return errors.New("public key missing")
	}

	pubkey, err := hex.DecodeString(strings.TrimPrefix(depositReceipt.Pubkey, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for public key")
	}

	if len(pubkey) != phase0.PublicKeyLength {
		return errors.New("incorrect length for public key")
	}

	copy(d.Pubkey[:], pubkey)

	if depositReceipt.WithdrawalCredentials == "" {
		return errors.New("withdrawal credentials missing")
	}

	if d.WithdrawalCredentials, err = hex.DecodeString(strings.TrimPrefix(depositReceipt.WithdrawalCredentials, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for withdrawal credentials")
	}

	if len(d.WithdrawalCredentials) != phase0.HashLength {
		return errors.New("incorrect length for withdrawal credentials")
	}

	amount, err := parseLenientUint64(depositReceipt.Amount, "amount")
	if err != nil {
		return err
	}

	d.Amount = phase0.Gwei(amount)

	if depositReceipt.Signature == "" {
		return errors.New("signature missing")
	}

	signature, err := hex.DecodeString(strings.TrimPrefix(depositReceipt.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}

	if len(signature) != phase0.SignatureLength {
		return errors.New("incorrect length for signature")
	}

	copy(d.Signature[:], signature)

	index, err := parseLenientUint64(depositReceipt.Index, "index")
	if err != nil {
		return err
	}

	d.Index = index

	return nil
}

// parseLenientUint64 parses a JSON-encoded uint64 from raw bytes,
// accepting either a quoted decimal string ("42") or a bare number (42).
// Returns "<field> missing" if the raw bytes are empty or an empty string.
func parseLenientUint64(raw json.RawMessage, field string) (uint64, error) {
	if len(raw) == 0 {
		return 0, errors.Errorf("%s missing", field)
	}

	s := string(raw)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	if s == "" {
		return 0, errors.Errorf("%s missing", field)
	}

	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid value for %s", field)
	}

	return v, nil
}
