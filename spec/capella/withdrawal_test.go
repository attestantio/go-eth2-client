// Copyright © 2022 Attestant Limited.
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
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/goccy/go-yaml"
	"github.com/golang/snappy"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestWithdrawalJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type capella.withdrawalJSON",
		},
		{
			name:  "IndexMissing",
			input: []byte(`{"validator_index":"3","address":"0x000102030405060708090a0b0c0d0e0f10111213","amount":"1000000000000000000"}`),
			err:   "index missing",
		},
		{
			name:  "IndexWrongType",
			input: []byte(`{"index":true,"validator_index":"3","address":"0x000102030405060708090a0b0c0d0e0f10111213","amount":"1000000000000000000"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field withdrawalJSON.index of type string",
		},
		{
			name:  "IndexInvalid",
			input: []byte(`{"index":"true","validator_index":"3","address":"0x000102030405060708090a0b0c0d0e0f10111213","amount":"1000000000000000000"}`),
			err:   "invalid value for index: strconv.ParseUint: parsing \"true\": invalid syntax",
		},
		{
			name:  "AddressMissing",
			input: []byte(`{"index":"2","validator_index":"3","amount":"1000000000000000000"}`),
			err:   "address missing",
		},
		{
			name:  "AddressWrongType",
			input: []byte(`{"index":"2","validator_index":"3","address":true,"amount":"1000000000000000000"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field withdrawalJSON.address of type string",
		},
		{
			name:  "AddressInvalid",
			input: []byte(`{"index":"2","validator_index":"3","address":"invalid","amount":"1000000000000000000"}`),
			err:   "invalid value for address: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "AddressWrongLength",
			input: []byte(`{"index":"2","validator_index":"3","address":"0x0102030405060708090a0b0c0d0e0f10111213","amount":"1000000000000000000"}`),
			err:   "incorrect length for address",
		},
		{
			name:  "AmountMissing",
			input: []byte(`{"index":"2","validator_index":"3","address":"0x000102030405060708090a0b0c0d0e0f10111213"}`),
			err:   "amount missing",
		},
		{
			name:  "AmountWrongType",
			input: []byte(`{"index":"2","validator_index":"3","address":"0x000102030405060708090a0b0c0d0e0f10111213","amount":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field withdrawalJSON.amount of type string",
		},
		{
			name:  "AmountInvalid",
			input: []byte(`{"index":"2","validator_index":"3","address":"0x000102030405060708090a0b0c0d0e0f10111213","amount":"true"}`),
			err:   "invalid value for amount: strconv.ParseUint: parsing \"true\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"index":"2","validator_index":"3","address":"0x000102030405060708090a0b0c0d0e0f10111213","amount":"1000000000000000000"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.Withdrawal
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

func TestWithdrawalYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{index: 2, validator_index: 3, address: '0x000102030405060708090a0b0c0d0e0f10111213', amount: 1000000000000000000}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res capella.Withdrawal
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

func TestWithdrawalSpec(t *testing.T) {
	if os.Getenv("ETH2_SPEC_TESTS_DIR") == "" {
		t.Skip("ETH2_SPEC_TESTS_DIR not suppplied, not running spec tests")
	}
	baseDir := filepath.Join(os.Getenv("ETH2_SPEC_TESTS_DIR"), "tests", "mainnet", "capella", "ssz_static", "Withdrawal", "ssz_random")
	require.NoError(t, filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if path == baseDir {
			// Only interested in subdirectories.
			return nil
		}
		require.NoError(t, err)
		if info.IsDir() {
			t.Run(info.Name(), func(t *testing.T) {
				specYAML, err := os.ReadFile(filepath.Join(path, "value.yaml"))
				require.NoError(t, err)
				var res capella.Withdrawal
				require.NoError(t, yaml.Unmarshal(specYAML, &res))

				compressedSpecSSZ, err := os.ReadFile(filepath.Join(path, "serialized.ssz_snappy"))
				require.NoError(t, err)
				var specSSZ []byte
				specSSZ, err = snappy.Decode(specSSZ, compressedSpecSSZ)
				require.NoError(t, err)

				// Ensure this matches the expected hash tree root.
				ssz, err := res.MarshalSSZ()
				require.NoError(t, err)
				require.Equal(t, specSSZ, ssz)
			})
		}
		return nil
	}))
}
