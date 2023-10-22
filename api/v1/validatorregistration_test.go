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

package v1_test

import (
	"bytes"
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/goccy/go-yaml"
	require "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestValidatorRegistrationJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.validatorRegistrationJSON",
		},

		{
			name:  "FeeRecipientMissing",
			input: []byte(`{"gas_limit":"100","timestamp":"100","pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
			err:   "fee recipient missing",
		},
		{
			name:  "FeeRecipientWrongType",
			input: []byte(`{"fee_recipient":true,"gas_limit":"100","timestamp":"100","pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorRegistrationJSON.fee_recipient of type string",
		},
		{
			name:  "FeeRecipientInvalid",
			input: []byte(`{"fee_recipient":"invalid","gas_limit":"100","timestamp":"100","pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
			err:   "invalid value for fee recipient: encoding/hex: invalid byte: U+0069 'i'",
		},

		{
			name:  "GasLimitMissing",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","timestamp":"100","pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
			err:   "gas limit missing",
		},
		{
			name:  "GasLimitWrongType",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","gas_limit":true,"timestamp":"100","pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorRegistrationJSON.gas_limit of type string",
		},
		{
			name:  "GasLimitInvalid",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","gas_limit":"-1","timestamp":"100","pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
			err:   "invalid value for gas limit: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},

		{
			name:  "TimestampMissing",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","gas_limit":"100","pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
			err:   "timestamp missing",
		},
		{
			name:  "TimestampWrongType",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","gas_limit":"100","timestamp":true,"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorRegistrationJSON.timestamp of type string",
		},
		{
			name:  "TimestampInvalid",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","gas_limit":"100","timestamp":"invalid","pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
			err:   "invalid value for timestamp: strconv.ParseInt: parsing \"invalid\": invalid syntax",
		},

		{
			name:  "PublicKeyMissing",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","gas_limit":"100","timestamp":"100"}`),
			err:   "public key missing",
		},
		{
			name:  "PublicKeyWrongType",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","gas_limit":"100","timestamp":"100","pubkey":[true]}`),
			err:   "invalid JSON: json: cannot unmarshal array into Go struct field validatorRegistrationJSON.pubkey of type string",
		},
		{
			name:  "PublicKeyInvalid",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","gas_limit":"100","timestamp":"100","pubkey":"invalid"}`),
			err:   "invalid value for public key: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "PublicKeyShort",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","gas_limit":"100","timestamp":"100","pubkey":"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
			err:   "incorrect length for public key",
		},
		{
			name:  "PublicKeyLong",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","gas_limit":"100","timestamp":"100","pubkey":"0x00000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
			err:   "incorrect length for public key",
		},

		{
			name:  "Good",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213","gas_limit":"100","timestamp":"100","pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.ValidatorRegistration
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

func TestValidatorRegistrationYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{fee_recipient: '0x000102030405060708090a0b0c0d0e0f10111213', gas_limit: 100, timestamp: 100, pubkey: '0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.ValidatorRegistration
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, res.String(), string(rt))
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}
