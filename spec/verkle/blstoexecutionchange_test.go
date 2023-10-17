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

package capella_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBLSToExecutionChangeJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type capella.blsToExecutionChangeJSON",
		},
		{
			name:  "ValidatorIndexMissing",
			input: []byte(`{"from_bls_pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","to_execution_address":"0x0102030405060708090a0B0c0d0e0f1011121314"}`),
			err:   "validator index missing",
		},
		{
			name:  "ValidatorIndexWrongType",
			input: []byte(`{"validator_index":true,"from_bls_pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","to_execution_address":"0x0102030405060708090a0B0c0d0e0f1011121314"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blsToExecutionChangeJSON.validator_index of type string",
		},
		{
			name:  "ValidatorIndexInvalid",
			input: []byte(`{"validator_index":"true","from_bls_pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","to_execution_address":"0x0102030405060708090a0B0c0d0e0f1011121314"}`),
			err:   "invalid value for validator index: strconv.ParseUint: parsing \"true\": invalid syntax",
		},
		{
			name:  "FromBLSPubkeyMissing",
			input: []byte(`{"validator_index":"2","to_execution_address":"0x0102030405060708090a0B0c0d0e0f1011121314"}`),
			err:   "from BLS public key missing",
		},
		{
			name:  "FromBLSPubkeyWrongType",
			input: []byte(`{"validator_index":"2","from_bls_pubkey":true,"to_execution_address":"0x0102030405060708090a0B0c0d0e0f1011121314"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blsToExecutionChangeJSON.from_bls_pubkey of type string",
		},
		{
			name:  "FromBLSPubkeyInvalid",
			input: []byte(`{"validator_index":"2","from_bls_pubkey":"invalid","to_execution_address":"0x0102030405060708090a0B0c0d0e0f1011121314"}`),
			err:   "invalid value for from BLS public key: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "FromBLSPubkeyWrongLength",
			input: []byte(`{"validator_index":"2","from_bls_pubkey":"0x9bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","to_execution_address":"0x0102030405060708090a0B0c0d0e0f1011121314"}`),
			err:   "incorrect length for from BLS public key",
		},
		{
			name:  "ToExecutionAddressMissing",
			input: []byte(`{"validator_index":"2","from_bls_pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b"}`),
			err:   "to execution address missing",
		},
		{
			name:  "ToExecutionAddressWrongType",
			input: []byte(`{"validator_index":"2","from_bls_pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","to_execution_address":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blsToExecutionChangeJSON.to_execution_address of type string",
		},
		{
			name:  "ToExecutionAddressInvalid",
			input: []byte(`{"validator_index":"2","from_bls_pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","to_execution_address":"invalid"}`),
			err:   "invalid value for to execution address: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "ToExecutionAddressWrongLength",
			input: []byte(`{"validator_index":"2","from_bls_pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","to_execution_address":"0x0102030405060708090a0b0c0d0e0f10111213"}`),
			err:   "incorrect length for to execution address",
		},
		{
			name:  "Good",
			input: []byte(`{"validator_index":"2","from_bls_pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","to_execution_address":"0x0102030405060708090a0B0c0d0e0f1011121314"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.BLSToExecutionChange
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

func TestBLSToExecutionChangeYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{validator_index: 2, from_bls_pubkey: '0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b', to_execution_address: '0x0102030405060708090a0B0c0d0e0f1011121314'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.BLSToExecutionChange
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
