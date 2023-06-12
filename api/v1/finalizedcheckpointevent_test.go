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

func TestFinalizedCheckpointEventJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.finalizedCheckpointEventJSON",
		},
		{
			name:  "BlockMissing",
			input: []byte(`{"state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch":"2"}`),
			err:   "block missing",
		},
		{
			name:  "BlockWrongType",
			input: []byte(`{"block":true,"state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch":"2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field finalizedCheckpointEventJSON.block of type string",
		},
		{
			name:  "BlockInvalid",
			input: []byte(`{"block":"invalid","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch":"2"}`),
			err:   "invalid value for block: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BlockShort",
			input: []byte(`{"block":"0xe3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch":"2"}`),
			err:   "incorrect length 31 for block",
		},
		{
			name:  "BlockLong",
			input: []byte(`{"block":"0x9999e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch":"2"}`),
			err:   "incorrect length 33 for block",
		},
		{
			name:  "StateMissing",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","epoch":"2"}`),
			err:   "state missing",
		},
		{
			name:  "StateWrongType",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":true,"epoch":"2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field finalizedCheckpointEventJSON.state of type string",
		},
		{
			name:  "StateInvalid",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"invalid","epoch":"2"}`),
			err:   "invalid value for state: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "StateShort",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x9a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch":"2"}`),
			err:   "incorrect length 31 for state",
		},
		{
			name:  "StateLong",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x74749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch":"2"}`),
			err:   "incorrect length 33 for state",
		},
		{
			name:  "EpochMissing",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28"}`),
			err:   "epoch missing",
		},
		{
			name:  "EpochWrongType",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field finalizedCheckpointEventJSON.epoch of type string",
		},
		{
			name:  "EpochInvalid",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch":"-1"}`),
			err:   "invalid value for epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28","epoch":"2"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.FinalizedCheckpointEvent
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
