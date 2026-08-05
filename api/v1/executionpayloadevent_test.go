// Copyright © 2026 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestExecutionPayloadEventJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.executionPayloadEventJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"builder_index":"12","block_hash":"0x1c3981b7439cd2dc53dca1a99122e1cacb36a13796d426d4c8a03ba745cb0c8b","block_root":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":false}`),
			err:   "slot missing",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","block_root":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BlockRootMissing",
			input: []byte(`{"slot":"4095940"}`),
			err:   "block root missing",
		},
		{
			name:  "BlockRootInvalid",
			input: []byte(`{"slot":"4095940","block_root":"invalid"}`),
			err:   "invalid value for block root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BlockRootShort",
			input: []byte(`{"slot":"4095940","block_root":"0xe3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`),
			err:   "incorrect length 31 for block root",
		},
		{
			name:  "BlockRootLong",
			input: []byte(`{"slot":"4095940","block_root":"0x9999e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`),
			err:   "incorrect length 33 for block root",
		},
		{
			name:  "BuilderIndexInvalid",
			input: []byte(`{"slot":"4095940","block_root":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","builder_index":"-1"}`),
			err:   "invalid value for builder index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BlockHashInvalid",
			input: []byte(`{"slot":"4095940","block_root":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","block_hash":"invalid"}`),
			err:   "invalid value for block hash: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BlockHashShort",
			input: []byte(`{"slot":"4095940","block_root":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","block_hash":"0xe3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`),
			err:   "incorrect length 31 for block hash",
		},
		{
			name:  "Good",
			input: []byte(`{"slot":"4095940","builder_index":"12","block_hash":"0x1c3981b7439cd2dc53dca1a99122e1cacb36a13796d426d4c8a03ba745cb0c8b","block_root":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","execution_optimistic":true}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.ExecutionPayloadEvent
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

	// Only slot and block_root are required; builder_index and block_hash are
	// optional and tolerated when absent. This input omits them, so it does not
	// round-trip byte-for-byte (MarshalJSON always emits all five fields);
	// assert on the parsed fields instead.
	t.Run("GoodOptionalFieldsAbsent", func(t *testing.T) {
		var res api.ExecutionPayloadEvent
		err := json.Unmarshal([]byte(`{"slot":"4095940","block_root":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`), &res)
		require.NoError(t, err)
		require.Equal(t, phase0.Slot(4095940), res.Slot)
		require.NotEqual(t, phase0.Root{}, res.BlockRoot)
		require.Equal(t, uint64(0), res.BuilderIndex)
		require.Equal(t, phase0.Hash32{}, res.BlockHash)
		require.False(t, res.ExecutionOptimistic)
	})
}
