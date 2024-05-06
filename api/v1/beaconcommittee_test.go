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

func TestBeaconCommitteeJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.beaconCommitteeJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"index":"2","validators":["2","128","4","61"]}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"slot":true,"index":"2","validators":["2","128","4","61"]}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconCommitteeJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","index":"2","validators":["2","128","4","61"]}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "IndexMissing",
			input: []byte(`{"slot":"1","validators":["2","128","4","61"]}`),
			err:   "index missing",
		},
		{
			name:  "IndexWrongType",
			input: []byte(`{"slot":"1","index":true,"validators":["2","128","4","61"]}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconCommitteeJSON.index of type string",
		},
		{
			name:  "IndexInvalid",
			input: []byte(`{"slot":"1","index":"-1","validators":["2","128","4","61"]}`),
			err:   "invalid value for index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ValidatorsMissing",
			input: []byte(`{"slot":"1","index":"2"}`),
			err:   "validators missing",
		},
		{
			name:  "ValidatorsWrongType",
			input: []byte(`{"slot":"1","index":"2","validators":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconCommitteeJSON.validators of type []string",
		},
		{
			name:  "ValidatorsInvalid",
			input: []byte(`{"slot":"1","index":"2","validators":["-1","128","4","61"]}`),
			err:   "invalid value for validator: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ValidatorsEmpty",
			input: []byte(`{"slot":"1","index":"2","validators":[]}`),
			err:   "validators length cannot be 0",
		},
		{
			name:  "Good",
			input: []byte(`{"slot":"1","index":"2","validators":["2","128","4","61"]}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.BeaconCommittee
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
