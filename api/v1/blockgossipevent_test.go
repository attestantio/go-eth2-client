// Copyright © 2024 Attestant Limited.
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
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestBlockGossipEventJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.blockGossipEventJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"slot":true,"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blockGossipEventJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BlockMissing",
			input: []byte(`{"slot":"525277"}`),
			err:   "block missing",
		},
		{
			name:  "BlockWrongType",
			input: []byte(`{"slot":"525277","block":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blockGossipEventJSON.block of type string",
		},
		{
			name:  "BlockInvalid",
			input: []byte(`{"slot":"525277","block":"invalid"}`),
			err:   "invalid value for block: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BlockShort",
			input: []byte(`{"slot":"525277","block":"0xe3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`),
			err:   "incorrect length 31 for block",
		},
		{
			name:  "BlockLong",
			input: []byte(`{"slot":"525277","block":"0x9999e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`),
			err:   "incorrect length 33 for block",
		},
		{
			name:  "Good",
			input: []byte(`{"slot":"525277","block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.BlockGossipEvent
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
