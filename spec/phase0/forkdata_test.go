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

func TestForkDataJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type phase0.forkDataJSON",
		},
		{
			name:  "CurrentVersionMissing",
			input: []byte(`{"genesis_validators_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"}`),
			err:   "current version missing",
		},
		{
			name:  "CurrentVersionWrongType",
			input: []byte(`{"current_version":true,"genesis_validators_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field forkDataJSON.current_version of type string",
		},
		{
			name:  "CurrentVersionInvalid",
			input: []byte(`{"current_version":"invalid","genesis_validators_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"}`),
			err:   "invalid value for current version: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "CurrentVersionShort",
			input: []byte(`{"current_version":"0x000002","genesis_validators_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"}`),
			err:   "incorrect length for current version",
		},
		{
			name:  "CurrentVersionLong",
			input: []byte(`{"current_version":"0x0000000002","genesis_validators_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"}`),
			err:   "incorrect length for current version",
		},
		{
			name:  "GenesisValidatorsRootMissing",
			input: []byte(`{"current_version":"0x00000002"}`),
			err:   "genesis validators root missing",
		},
		{
			name:  "GenesisValidatorsRootWrongType",
			input: []byte(`{"current_version":"0x00000002","genesis_validators_root":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field forkDataJSON.genesis_validators_root of type string",
		},
		{
			name:  "GenesisValidatorsRootInvalid",
			input: []byte(`{"current_version":"0x00000002","genesis_validators_root":"invalid"}`),
			err:   "invalid value for genesis validators root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "GenesisValidatorsRootShort",
			input: []byte(`{"current_version":"0x00000002","genesis_validators_root":"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"}`),
			err:   "incorrect length for genesis validators root",
		},
		{
			name:  "GenesisValidatorsRootLong",
			input: []byte(`{"current_version":"0x00000002","genesis_validators_root":"0x00000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"}`),
			err:   "incorrect length for genesis validators root",
		},
		{
			name:  "Good",
			input: []byte(`{"current_version":"0x00000002","genesis_validators_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.ForkData
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

func TestForkDataYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{current_version: '0x00000002', genesis_validators_root: '0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.ForkData
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
