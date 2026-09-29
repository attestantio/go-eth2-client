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
	"math"
	"reflect"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"

	api "github.com/attestantio/go-eth2-client/api/v1"
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
			input: []byte(`{"slot":"4095940","block_root":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","builder_index":"12","block_hash":"invalid"}`),
			err:   "invalid value for block hash: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BlockHashShort",
			input: []byte(`{"slot":"4095940","block_root":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","builder_index":"12","block_hash":"0xe3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028"}`),
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

}

func TestExecutionPayloadEventRequiredFields(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{name: "BuilderIndexMissing", input: `{"slot":"10","block_hash":"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef","block_root":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf"}`, err: "builder index missing"},
		{name: "BlockHashMissing", input: `{"slot":"10","builder_index":"42","block_root":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf"}`, err: "block hash missing"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var event api.ExecutionPayloadEvent
			require.EqualError(t, json.Unmarshal([]byte(test.input), &event), test.err)
		})
	}
}

func TestExecutionPayloadGossipEventRequiredFields(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{name: "BuilderIndexMissing", input: `{"slot":"10","block_hash":"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef","block_root":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf"}`, err: "builder index missing"},
		{name: "BlockHashMissing", input: `{"slot":"10","builder_index":"42","block_root":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf"}`, err: "block hash missing"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var event api.ExecutionPayloadGossipEvent
			require.EqualError(t, json.Unmarshal([]byte(test.input), &event), test.err)
		})
	}
}

func TestExecutionPayloadEventBuilderIndexType(t *testing.T) {
	for _, event := range []any{api.ExecutionPayloadEvent{}, api.ExecutionPayloadGossipEvent{}} {
		field, ok := reflect.TypeOf(event).FieldByName("BuilderIndex")
		require.True(t, ok)
		require.Equal(t, reflect.TypeOf(gloas.BuilderIndex(0)), field.Type)
	}
}

func TestExecutionPayloadEventSelfBuild(t *testing.T) {
	const input = `{"slot":"16342","builder_index":"18446744073709551615","block_hash":"0x3e901234567890abcdef1234567890abcdef1234567890abcdef1234567890ab","block_root":"0x9ac61234567890abcdef1234567890abcdef1234567890abcdef1234567890ab","execution_optimistic":false}`
	for _, event := range []json.Unmarshaler{&api.ExecutionPayloadEvent{}, &api.ExecutionPayloadGossipEvent{}} {
		require.NoError(t, event.UnmarshalJSON([]byte(input)))
		field := reflect.ValueOf(event).Elem().FieldByName("BuilderIndex")
		require.Equal(t, gloas.BuilderIndexSelfBuild, field.Interface())
		require.Equal(t, uint64(math.MaxUint64), uint64(gloas.BuilderIndexSelfBuild))
	}
}

func TestExecutionPayloadEventsIgnoreNimbusStateRoot(t *testing.T) {
	const input = `{"slot":"10","builder_index":"42","block_hash":"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef","block_root":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","state_root":"0x0000000000000000000000000000000000000000000000000000000000000000"}`
	for _, event := range []json.Unmarshaler{&api.ExecutionPayloadEvent{}, &api.ExecutionPayloadGossipEvent{}} {
		require.NoError(t, event.UnmarshalJSON([]byte(input)))
		_, exists := reflect.TypeOf(event).Elem().FieldByName("StateRoot")
		require.False(t, exists)
		encoded, err := json.Marshal(event)
		require.NoError(t, err)
		require.NotContains(t, string(encoded), "state_root")
	}
}
