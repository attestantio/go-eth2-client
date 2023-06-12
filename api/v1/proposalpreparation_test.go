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
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	require "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestProposalPreparationJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.proposalPreparationJSON",
		},
		{
			name:  "ValidatorIndexMissing",
			input: []byte(`{"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213"}`),
			err:   "validator index missing",
		},
		{
			name:  "ValidatorIndexWrongType",
			input: []byte(`{"validator_index":true,"fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field proposalPreparationJSON.validator_index of type string",
		},
		{
			name:  "ValidatorIndexInvalid",
			input: []byte(`{"validator_index":"-1","fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213"}`),
			err:   "invalid value for validator index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "FeeRecipientMissing",
			input: []byte(`{"validator_index":"1"}`),
			err:   "fee recipient is missing",
		},
		{
			name:  "FeeRecipientWrongType",
			input: []byte(`{"validator_index":"1","fee_recipient":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field proposalPreparationJSON.fee_recipient of type string",
		},
		{
			name:  "FeeRecipientInvalid",
			input: []byte(`{"validator_index":"1","fee_recipient":"invalid"}`),
			err:   "invalid value for fee recipient: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "Good",
			input: []byte(`{"validator_index":"1","fee_recipient":"0x000102030405060708090a0b0c0d0e0f10111213"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.ProposalPreparation
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
