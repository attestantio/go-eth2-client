// Copyright © 2023 Attestant Limited.
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

package capella_test

import (
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/capella"
	require "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestHistoricalSummaryJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		output []byte
		err    string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type capella.historicalSummaryJSON",
		},
		{
			name:  "BlockSummaryRootMissing",
			input: []byte(`{"state_summary_root":"0xea760203509bdde017a506b12c825976d12b04db7bce9eca9e1ed007056a3f36"}`),
			err:   "block summary root missing",
		},
		{
			name:  "BlockSummaryRootWrongType",
			input: []byte(`{"block_summary_root":true,"state_summary_root":"0xea760203509bdde017a506b12c825976d12b04db7bce9eca9e1ed007056a3f36"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field historicalSummaryJSON.block_summary_root of type string",
		},
		{
			name:  "BlockSummaryRootInvalid",
			input: []byte(`{"block_summary_root":"true","state_summary_root":"0xea760203509bdde017a506b12c825976d12b04db7bce9eca9e1ed007056a3f36"}`),
			err:   "invalid value for block summary root: encoding/hex: invalid byte: U+0074 't'",
		},
		{
			name:  "BlockSummaryRootWrongLength",
			input: []byte(`{"block_summary_root":"0x6e230e6eceb8f3db582777b1500b8b31b9d268339e7b32bba8d6f1311b211d","state_summary_root":"0xea760203509bdde017a506b12c825976d12b04db7bce9eca9e1ed007056a3f36"}`),
			err:   "incorrect length for block summary root",
		},
		{
			name:  "StateSummaryRootMissing",
			input: []byte(`{"block_summary_root":"0x3d6e230e6eceb8f3db582777b1500b8b31b9d268339e7b32bba8d6f1311b211d"}`),
			err:   "state summary root missing",
		},
		{
			name:  "StateSummaryRootWrongType",
			input: []byte(`{"block_summary_root":"0x3d6e230e6eceb8f3db582777b1500b8b31b9d268339e7b32bba8d6f1311b211d","state_summary_root":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field historicalSummaryJSON.state_summary_root of type string",
		},
		{
			name:  "StateSummaryRootInvalid",
			input: []byte(`{"block_summary_root":"0x3d6e230e6eceb8f3db582777b1500b8b31b9d268339e7b32bba8d6f1311b211d","state_summary_root":"true"}`),
			err:   "invalid value for state summary root: encoding/hex: invalid byte: U+0074 't'",
		},
		{
			name:  "StateSummaryRootWrongLength",
			input: []byte(`{"block_summary_root":"0x3d6e230e6eceb8f3db582777b1500b8b31b9d268339e7b32bba8d6f1311b211d","state_summary_root":"0x760203509bdde017a506b12c825976d12b04db7bce9eca9e1ed007056a3f36"}`),
			err:   "incorrect length for state summary root",
		},
		{
			name:  "Good",
			input: []byte(`{"block_summary_root":"0x3d6e230e6eceb8f3db582777b1500b8b31b9d268339e7b32bba8d6f1311b211d","state_summary_root":"0xea760203509bdde017a506b12c825976d12b04db7bce9eca9e1ed007056a3f36"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.HistoricalSummary
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				if len(test.output) > 0 {
					assert.Equal(t, string(test.output), string(rt))
				} else {
					assert.Equal(t, string(test.input), string(rt))
				}
			}
		})
	}
}
