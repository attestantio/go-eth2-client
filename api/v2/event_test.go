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

package v2_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v2"
	require "github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestEvent(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v2.eventJSON",
		},
		{
			name:  "TopicMissing",
			input: []byte(`{"data":{}}`),
			err:   "topic missing",
		},
		{
			name:  "TopicWrongType",
			input: []byte(`{"topic":[],"data":{}}`),
			err:   "invalid JSON: json: cannot unmarshal array into Go struct field eventJSON.topic of type string",
		},
		{
			name:  "TopicUnsupported",
			input: []byte(`{"topic":"foo","data":{"block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","current_duty_dependent_root":"0x92c6b763f610d5941d2041906007bf9449d37772aacf0483a76275ac27c096b4","epoch_transition":false,"previous_duty_dependent_root":"0xa692c095bbca3eeaf99eeabada78874c028c02b176ccf691f3e8fa075d67f5c6","slot":"231192","state":"0x61099b2c1dee0104c93ce0e14e5f5fc4b6faceff4cb863278d055bdfb73b7dc7"}}`),
			err:   "unsupported event topic foo",
		},
		{
			name:  "DataMissing",
			input: []byte(`{"topic":"head"}`),
			err:   "data missing",
		},
		{
			name:  "Good",
			input: []byte(`{"topic":"head","data":{"block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","current_duty_dependent_root":"0x92c6b763f610d5941d2041906007bf9449d37772aacf0483a76275ac27c096b4","epoch_transition":false,"previous_duty_dependent_root":"0xa692c095bbca3eeaf99eeabada78874c028c02b176ccf691f3e8fa075d67f5c6","slot":"231192","state":"0x61099b2c1dee0104c93ce0e14e5f5fc4b6faceff4cb863278d055bdfb73b7dc7"}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.Event
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
