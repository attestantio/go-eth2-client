// Copyright © 2021 Attestant Limited.
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

func TestSyncCommitteeJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.syncCommitteeJSON",
		},
		{
			name:  "ValidatorsMissing",
			input: []byte(`{"validator_aggregates":[["2","128"],["4","61"]]}`),
			err:   "validators missing",
		},
		{
			name:  "ValidatorsEmpty",
			input: []byte(`{"validators":[],"validator_aggregates":[["2","128"],["4","61"]]}`),
			err:   "validators length cannot be 0",
		},
		{
			name:  "ValidatorsWrongType",
			input: []byte(`{"validators":true,"validator_aggregates":[["2","128"],["4","61"]]}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeJSON.validators of type []string",
		},
		{
			name:  "ValidatorsIndexWrongType",
			input: []byte(`{"validators":["2","128","4",true],"validator_aggregates":[["2","128"],["4","61"]]}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeJSON.validators of type string",
		},
		{
			name:  "ValidatorsIndexInvalid",
			input: []byte(`{"validators":["2","128","4","invalid"],"validator_aggregates":[["2","128"],["4","61"]]}`),
			err:   "invalid value for validator: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "ValidatorAggregatesMissing",
			input: []byte(`{"validators":["2","128","4","61"]}`),
			err:   "validator aggregates missing",
		},
		{
			name:  "ValidatorAggregatesEmpty",
			input: []byte(`{"validators":["2","128","4","61"],"validator_aggregates":[]}`),
			err:   "validator aggregates length cannot be 0",
		},
		{
			name:  "ValidatorAggregatesWrongType",
			input: []byte(`{"validators":["2","128","4","61"],"validator_aggregates":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeJSON.validator_aggregates of type [][]string",
		},
		{
			name:  "ValidatorAggregateWrongType",
			input: []byte(`{"validators":["2","128","4","61"],"validator_aggregates":[["2","128"],true]}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeJSON.validator_aggregates of type []string",
		},
		{
			name:  "ValidatorAggregateEmpty",
			input: []byte(`{"validators":["2","128","4","61"],"validator_aggregates":[["2","128"],[]]}`),
			err:   "validator aggregate length cannot be 0",
		},
		{
			name:  "ValidatorAggregateIndexWrongType",
			input: []byte(`{"validators":["2","128","4","61"],"validator_aggregates":[["2","128"],["4",true]]}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeJSON.validator_aggregates of type string",
		},
		{
			name:  "ValidatorAggregateIndexInvalid",
			input: []byte(`{"validators":["2","128","4","61"],"validator_aggregates":[["2","128"],["4","invalid"]]}`),
			err:   "invalid value for validator aggregate: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"validators":["2","128","4","61"],"validator_aggregates":[["2","128"],["4","61"]]}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.SyncCommittee
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
