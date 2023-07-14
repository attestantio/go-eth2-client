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

func TestETH1DataJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type phase0.eth1DataJSON",
		},
		{
			name:  "DepositRootMissing",
			input: []byte(`{"deposit_count":"10","block_hash":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "deposit root missing",
		},
		{
			name:  "DepositRootWrongType",
			input: []byte(`{"deposit_root":true,"deposit_count":"10","block_hash":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field eth1DataJSON.deposit_root of type string",
		},
		{
			name:  "DepositRootInvalid",
			input: []byte(`{"deposit_root":"invalid","deposit_count":"10","block_hash":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "invalid value for deposit root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "DepositRootShort",
			input: []byte(`{"deposit_root":"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":"10","block_hash":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "incorrect length for deposit root",
		},
		{
			name:  "DepositRootLong",
			input: []byte(`{"deposit_root":"0x00000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":"10","block_hash":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "incorrect length for deposit root",
		},
		{
			name:  "DepositCountMissing",
			input: []byte(`{"deposit_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","block_hash":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "deposit count missing",
		},
		{
			name:  "DepositCountWrongType",
			input: []byte(`{"deposit_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":true,"block_hash":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field eth1DataJSON.deposit_count of type string",
		},
		{
			name:  "DepositCountInvalid",
			input: []byte(`{"deposit_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":"-1","block_hash":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "invalid value for deposit count: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BlockHashMissing",
			input: []byte(`{"deposit_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":"10"}`),
			err:   "block hash missing",
		},
		{
			name:  "BlockHashWrongType",
			input: []byte(`{"deposit_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":"10","block_hash":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field eth1DataJSON.block_hash of type string",
		},
		{
			name:  "BlockHashInvalid",
			input: []byte(`{"deposit_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":"10","block_hash":"invalid"}`),
			err:   "invalid value for block hash: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BlockHashShort",
			input: []byte(`{"deposit_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":"10","block_hash":"0x2122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "incorrect length for block hash",
		},
		{
			name:  "BlockHashLong",
			input: []byte(`{"deposit_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":"10","block_hash":"0x20202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "incorrect length for block hash",
		},
		{
			name:  "Good",
			input: []byte(`{"deposit_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":"10","block_hash":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.ETH1Data
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

func TestETH1DataYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{deposit_root: '0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f', deposit_count: 10, block_hash: '0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.ETH1Data
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
