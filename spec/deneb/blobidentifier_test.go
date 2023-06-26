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

package deneb_test

import (
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/goccy/go-yaml"
	require "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestBlobIdentifierJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type map[string]json.RawMessage",
		},
		{
			name:  "BlockRootMissing",
			input: []byte(`{"index":"17189299593882149153"}`),
			err:   "block_root: missing",
		},
		{
			name:  "BlockRootWrongType",
			input: []byte(`{"block_root":true,"index":"17189299593882149153"}`),
			err:   "block_root: invalid prefix",
		},
		{
			name:  "BlockRootInvalid",
			input: []byte(`{"block_root":"true","index":"17189299593882149153"}`),
			err:   "block_root: invalid prefix",
		},
		{
			name:  "BlockRootIncorrectLength",
			input: []byte(`{"block_root":"0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec7","index":"17189299593882149153"}`),
			err:   "block_root: incorrect length",
		},
		{
			name:  "IndexMissing",
			input: []byte(`{"block_root":"0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec799"}`),
			err:   "index: missing",
		},
		{
			name:  "IndexWrongType",
			input: []byte(`{"block_root":"0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec799","index":true}`),
			err:   "index: invalid prefix",
		},
		{
			name:  "IndexInvalid",
			input: []byte(`{"block_root":"0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec799","index":"-1"}`),
			err:   "index: invalid value -1: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"block_root":"0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec799","index":"17189299593882149153"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res deneb.BlobIdentifier
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

func TestBlobIdentifierYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{block_root: '0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec799', index: 17189299593882149153}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res deneb.BlobIdentifier
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, testYAMLFormat([]byte(res.String())), testYAMLFormat(rt))
				assert.Equal(t, testYAMLFormat(test.input), testYAMLFormat(rt))
			}
		})
	}
}
