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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/golang/snappy"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestVoluntaryExitJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type phase0.voluntaryExitJSON",
		},
		{
			name:  "EpochMissing",
			input: []byte(`{"validator_index":"2"}`),
			err:   "epoch missing",
		},
		{
			name:  "EpochWrongType",
			input: []byte(`{"epoch":true,"validator_index":"2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field voluntaryExitJSON.epoch of type string",
		},
		{
			name:  "EpochInvalid",
			input: []byte(`{"epoch":"-1","validator_index":"2"}`),
			err:   "invalid value for epoch: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ValidatorIndexMissing",
			input: []byte(`{"epoch":"1"}`),
			err:   "validator index missing",
		},
		{
			name:  "ValidatorIndexWrongType",
			input: []byte(`{"epoch":"1","validator_index":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field voluntaryExitJSON.validator_index of type string",
		},
		{
			name:  "ValidatorIndexInvalid",
			input: []byte(`{"epoch":"1","validator_index":"-1"}`),
			err:   "invalid value for validator index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"epoch":"1","validator_index":"2"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.VoluntaryExit
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

func TestVoluntaryExitYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{epoch: 1, validator_index: 2}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.VoluntaryExit
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

func TestVoluntaryExitSpec(t *testing.T) {
	if os.Getenv("ETH2_SPEC_TESTS_DIR") == "" {
		t.Skip("ETH2_SPEC_TESTS_DIR not suppplied, not running spec tests")
	}
	baseDir := filepath.Join(os.Getenv("ETH2_SPEC_TESTS_DIR"), "tests", "mainnet", "phase0", "ssz_static", "VoluntaryExit", "ssz_random")
	require.NoError(t, filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if path == baseDir {
			// Only interested in subdirectories.
			return nil
		}
		require.NoError(t, err)
		if info.IsDir() {
			t.Run(info.Name(), func(t *testing.T) {
				specYAML, err := ioutil.ReadFile(filepath.Join(path, "value.yaml"))
				require.NoError(t, err)
				var res phase0.VoluntaryExit
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
