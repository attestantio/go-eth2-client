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

package phase0_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepositMessageJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type phase0.depositMessageJSON",
		},
		{
			name:  "PublicKeyMissing",
			input: []byte(`{"withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "public key missing",
		},
		{
			name:  "PublicKeyWrongType",
			input: []byte(`{"pubkey":true,"withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field depositMessageJSON.pubkey of type string",
		},
		{
			name:  "PublicKeyInvalid",
			input: []byte(`{"pubkey":"invalid","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "invalid value for public key: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "PublicKeyShort",
			input: []byte(`{"pubkey":"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "incorrect length for public key",
		},
		{
			name:  "PublicKeyLong",
			input: []byte(`{"pubkey":"0x00000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "incorrect length for public key",
		},
		{
			name:  "WithdrawalCredentialsMissing",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","amount":"32000000000"}`),
			err:   "withdrawal credentials missing",
		},
		{
			name:  "WithdrawalCredentialsWrongType",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":true,"amount":"32000000000"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field depositMessageJSON.withdrawal_credentials of type string",
		},
		{
			name:  "WithdrawalCredentialsInvalid",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"invalid","amount":"32000000000"}`),
			err:   "invalid value for withdrawal credentials: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "WithdrawalCredentialsShort",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x2122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "incorrect length for withdrawal credentials",
		},
		{
			name:  "WithdrawalCredentialsLong",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x20202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
			err:   "incorrect length for withdrawal credentials",
		},
		{
			name:  "AmountMissing",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "amount missing",
		},
		{
			name:  "AmountWrongType",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field depositMessageJSON.amount of type string",
		},
		{
			name:  "AmountInvalid",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"-1"}`),
			err:   "invalid value for amount: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.DepositMessage
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestDepositMessageYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{pubkey: '0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f', withdrawal_credentials: '0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f', amount: 32000000000}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.DepositMessage
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(rt), res.String())
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}
