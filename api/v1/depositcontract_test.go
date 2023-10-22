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

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	require "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestDepositContractJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.depositContractJSON",
		},
		{
			name:  "ChainIDMissing",
			input: []byte(`{"address":"0x07b39f4fde4a38bace212b546dac87c58dfe3fdc"}`),
			err:   "chain ID missing",
		},
		{
			name:  "ChainIDWrongType",
			input: []byte(`{"chain_id":true,"address":"0x07b39f4fde4a38bace212b546dac87c58dfe3fdc"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field depositContractJSON.chain_id of type string",
		},
		{
			name:  "ChainIDInvalid",
			input: []byte(`{"chain_id":"-1","address":"0x07b39f4fde4a38bace212b546dac87c58dfe3fdc"}`),
			err:   "invalid value for chain ID: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "AddressMissing",
			input: []byte(`{"chain_id":"5"}`),
			err:   "address missing",
		},
		{
			name:  "AddressWrongType",
			input: []byte(`{"chain_id":"5","address":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field depositContractJSON.address of type string",
		},
		{
			name:  "AddressInvalid",
			input: []byte(`{"chain_id":"5","address":"invalid"}`),
			err:   "invalid value for address: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "AddressShort",
			input: []byte(`{"chain_id":"5","address":"0xb39f4fde4a38bace212b546dac87c58dfe3fdc"}`),
			err:   "incorrect length 19 for address",
		},
		{
			name:  "AddressLong",
			input: []byte(`{"chain_id":"5","address":"0x0707b39f4fde4a38bace212b546dac87c58dfe3fdc"}`),
			err:   "incorrect length 21 for address",
		},
		{
			name:  "Good",
			input: []byte(`{"chain_id":"5","address":"0x07b39f4fde4a38bace212b546dac87c58dfe3fdc"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.DepositContract
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}
