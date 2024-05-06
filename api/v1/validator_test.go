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

func TestValidatorJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.validatorJSON",
		},
		{
			name:  "IndexMissing",
			input: []byte(`{"balance":"32000000000","status":"active_ongoing","validator":{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}}`),
			err:   "index missing",
		},
		{
			name:  "IndexWrongType",
			input: []byte(`{"index":true,"balance":"32000000000","status":"active_ongoing","validator":{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.index of type string",
		},
		{
			name:  "IndexInvalid",
			input: []byte(`{"index":"-1","balance":"32000000000","status":"active_ongoing","validator":{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}}`),
			err:   "invalid value for index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BalanceMissing",
			input: []byte(`{"index":"1","status":"active_ongoing","validator":{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}}`),
			err:   "balance missing",
		},
		{
			name:  "BalanceWrongType",
			input: []byte(`{"index":"1","balance":true,"status":"active_ongoing","validator":{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorJSON.balance of type string",
		},
		{
			name:  "BalanceInvalid",
			input: []byte(`{"index":"1","balance":"-1","status":"active_ongoing","validator":{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}}`),
			err:   "invalid value for balance: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "StateWrongType",
			input: []byte(`{"index":"1","balance":"32000000000","status":true,"validator":{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}}`),
			err:   "invalid JSON: unrecognised validator state true",
		},
		{
			name:  "StateInvalid",
			input: []byte(`{"index":"1","balance":"32000000000","status":"invalid","validator":{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}}`),
			err:   "invalid JSON: unrecognised validator state \"invalid\"",
		},
		{
			name:  "ValidatorMissing",
			input: []byte(`{"index":"1","balance":"32000000000","status":"active_ongoing"}`),
			err:   "validator missing",
		},
		{
			name:  "ValidatorWrongType",
			input: []byte(`{"index":"1","balance":"32000000000","status":"active_ongoing","validator":true}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.validatorJSON",
		},
		{
			name:  "ValidatorInvalid",
			input: []byte(`{"index":"1","balance":"32000000000","status":"active_ongoing","validator":{}}`),
			err:   "invalid JSON: public key missing",
		},
		{
			name:  "Good",
			input: []byte(`{"index":"1","balance":"32000000000","status":"active_ongoing","validator":{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","withdrawal_credentials":"0x00ec7ef7780c9d151597924036262dd28dc60e1228f4da6fecf9d402cb3f3594","effective_balance":"32000000000","slashed":false,"activation_eligibility_epoch":"0","activation_epoch":"0","exit_epoch":"18446744073709551615","withdrawable_epoch":"18446744073709551615"}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.Validator
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
