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

func TestValidatorBalanceJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.validatorBalanceJSON",
		},
		{
			name:  "IndexMissing",
			input: []byte(`{"balance":"32000000000"}`),
			err:   "index missing",
		},
		{
			name:  "IndexWrongType",
			input: []byte(`{"index":true,"balance":"32000000000"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorBalanceJSON.index of type string",
		},
		{
			name:  "IndexInvalid",
			input: []byte(`{"index":"-1","balance":"32000000000"}`),
			err:   "invalid value for index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BalanceMissing",
			input: []byte(`{"index":"1"}`),
			err:   "balance missing",
		},
		{
			name:  "BalanceWrongType",
			input: []byte(`{"index":"1","balance":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field validatorBalanceJSON.balance of type string",
		},
		{
			name:  "BalanceInvalid",
			input: []byte(`{"index":"1","balance":"-1"}`),
			err:   "invalid value for balance: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"index":"1","balance":"32000000000"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.ValidatorBalance
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
