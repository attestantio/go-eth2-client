// Copyright © 2020 Attestant Limited.
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

package altair_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	spec "github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/goccy/go-yaml"
	require "github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestPendingAttestationJSON(t *testing.T) {
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
			name:  "AggregationBitsMissing",
			input: []byte(`{"data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"inclusion_delay":"1","proposer_index":"2"}`),
			err:   "aggregation bits missing",
		},
		{
			name:  "AggregationBitsWrongType",
			input: []byte(`{"aggregation_bits":true,"data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"inclusion_delay":"1","proposer_index":"2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field pendingAttestationJSON.aggregation_bits of type string",
		},
		{
			name:  "AggregationBitsInvalid",
			input: []byte(`{"aggregation_bits":"invalid","data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"inclusion_delay":"1","proposer_index":"2"}`),
			err:   "invalid value for aggregation bits: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type altair.pendingAttestationJSON",
		},
		{
			name:  "DataMissing",
			input: []byte(`{"aggregation_bits":"0x010203","inclusion_delay":"1","proposer_index":"2"}`),
			err:   "data missing",
		},
		{
			name:  "DataWrongType",
			input: []byte(`{"aggregation_bits":"0x010203","data":true,"inclusion_delay":"1","proposer_index":"2"}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type altair.attestationDataJSON",
		},
		{
			name:  "InclusionDelayMissing",
			input: []byte(`{"aggregation_bits":"0x010203","data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"proposer_index":"2"}`),
			err:   "inclusion delay missing",
		},
		{
			name:  "InclusionDelayWrongType",
			input: []byte(`{"aggregation_bits":"0x010203","data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"inclusion_delay":true,"proposer_index":"2"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field pendingAttestationJSON.inclusion_delay of type string",
		},
		{
			name:  "InclusionDelayInvalid",
			input: []byte(`{"aggregation_bits":"0x010203","data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"inclusion_delay":"-1","proposer_index":"2"}`),
			err:   "invalid value for inclusion delay: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ProposerIndexMissing",
			input: []byte(`{"aggregation_bits":"0x010203","data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"inclusion_delay":"1"}`),
			err:   "proposer index missing",
		},
		{
			name:  "ProposerIndexWrongType",
			input: []byte(`{"aggregation_bits":"0x010203","data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"inclusion_delay":"1","proposer_index":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field pendingAttestationJSON.proposer_index of type string",
		},
		{
			name:  "ProposerIndexInvalid",
			input: []byte(`{"aggregation_bits":"0x010203","data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"inclusion_delay":"1","proposer_index":"-1"}`),
			err:   "invalid value for proposer index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"aggregation_bits":"0x010203","data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"inclusion_delay":"1","proposer_index":"2"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res spec.PendingAttestation
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

func TestPendingAttestationYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{aggregation_bits: '0x010203', data: {slot: 100, index: 1, beacon_block_root: '0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f', source: {epoch: 1, root: '0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f'}, target: {epoch: 2, root: '0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f'}}, inclusion_delay: 1, proposer_index: 2}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res spec.PendingAttestation
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestPendingAttestationSpec(t *testing.T) {
	if os.Getenv("ETH2_SPEC_TESTS_DIR") == "" {
		t.Skip("ETH2_SPEC_TESTS_DIR not suppplied, not running spec tests")
	}
	baseDir := filepath.Join(os.Getenv("ETH2_SPEC_TESTS_DIR"), "tests", "mainnet", "altair", "ssz_static", "PendingAttestation", "ssz_random")
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
				var res spec.PendingAttestation
				require.NoError(t, yaml.Unmarshal(specYAML, &res))

				specSSZ, err := ioutil.ReadFile(filepath.Join(path, "serialized.ssz"))
				require.NoError(t, err)

				ssz, err := res.MarshalSSZ()
				require.NoError(t, err)
				require.Equal(t, specSSZ, ssz)

				root, err := res.HashTreeRoot()
				require.NoError(t, err)
				rootsYAML, err := ioutil.ReadFile(filepath.Join(path, "roots.yaml"))
				require.NoError(t, err)
				require.Equal(t, string(rootsYAML), fmt.Sprintf("{root: '%#x'}\n", root))
			})
		}
		return nil
	}))
}
