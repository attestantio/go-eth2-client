// Copyright © 2020, 2021 Attestant Limited.
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

package phase0_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestForkJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type phase0.forkJSON",
		},
		{
			name:  "PreviousVersionMissing",
			input: []byte(`{"current_version":"0x00000002","epoch":"3"}`),
			err:   "previous version missing",
		},
		{
			name:  "PreviousVersionWrongType",
			input: []byte(`{"previous_version":true,"current_version":"0x00000002","epoch":"3"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field forkJSON.previous_version of type string",
		},
		{
			name:  "PreviousVersionInvalid",
			input: []byte(`{"previous_version":"invalid","current_version":"0x00000002","epoch":"3"}`),
			err:   "invalid value for previous version: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "PreviousVersionShort",
			input: []byte(`{"previous_version":"0x000001","current_version":"0x00000002","epoch":"3"}`),
			err:   "incorrect length for previous version",
		},
		{
			name:  "PreviousVersionLong",
			input: []byte(`{"previous_version":"0x0000000001","current_version":"0x00000002","epoch":"3"}`),
			err:   "incorrect length for previous version",
		},
		{
			name:  "CurrentVersionMissing",
			input: []byte(`{"previous_version":"0x00000001","epoch":"3"}`),
			err:   "current version missing",
		},
		{
			name:  "CurrentVersionWrongType",
			input: []byte(`{"previous_version":"0x00000001","current_version":true,"epoch":"3"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field forkJSON.current_version of type string",
		},
		{
			name:  "CurrentVersionInvalid",
			input: []byte(`{"previous_version":"0x00000001","current_version":"invalid","epoch":"3"}`),
			err:   "invalid value for current version: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "CurrentVersionShort",
			input: []byte(`{"previous_version":"0x00000001","current_version":"0x000002","epoch":"3"}`),
			err:   "incorrect length for current version",
		},
		{
			name:  "CurrentVersionLong",
			input: []byte(`{"previous_version":"0x00000001","current_version":"0x0000000002","epoch":"3"}`),
			err:   "incorrect length for current version",
		},
		{
			name:  "EpochMissing",
			input: []byte(`{"previous_version":"0x00000001","current_version":"0x00000002"}`),
			err:   "epoch missing",
		},
		{
			name:  "EpochWrongType",
			input: []byte(`{"previous_version":"0x00000001","current_version":"0x00000002","epoch":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field forkJSON.epoch of type string",
		},
		{
			name:  "EpochInvalid",
			input: []byte(`{"previous_version":"0x00000001","current_version":"0x00000002","epoch":"-1"}`),
			err:   "invalid value for epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"previous_version":"0x00000001","current_version":"0x00000002","epoch":"3"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.Fork
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestForkYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{previous_version: '0x00000001', current_version: '0x00000002', epoch: 3}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.Fork
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(rt), res.String())
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}
