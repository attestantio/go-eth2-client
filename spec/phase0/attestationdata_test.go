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

package phase0_test

import (
	"encoding/json"
	"testing"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	require "github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestAttestationDataJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type phase0.attestationDataJSON",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"slot":true,"index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field attestationDataJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "IndexWrongType",
			input: []byte(`{"slot":"100","index":true,"beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field attestationDataJSON.index of type string",
		},
		{
			name:  "IndexInvalid",
			input: []byte(`{"slot":"100","index":"-1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
			err:   "invalid value for index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BeaconBlockRootWrongType",
			input: []byte(`{"slot":"100","index":"1","beacon_block_root":true,"source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field attestationDataJSON.beacon_block_root of type string",
		},
		{
			name:  "BeaconBlockRootInvalid",
			input: []byte(`{"slot":"100","index":"1","beacon_block_root":"invalid","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
			err:   "invalid value for beacon block root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BeaconBlockRootShort",
			input: []byte(`{"slot":"100","index":"1","beacon_block_root":"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
			err:   "incorrect length for beacon block root",
		},
		{
			name:  "BeaconBlockRootLong",
			input: []byte(`{"slot":"100","index":"1","beacon_block_root":"0x00000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
			err:   "incorrect length for beacon block root",
		},
		{
			name:  "SourceMissing",
			input: []byte(`{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
			err:   "source missing",
		},
		{
			name:  "SourceInvalid",
			input: []byte(`{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":true,"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.checkpointJSON",
		},
		{
			name:  "TargetMissing",
			input: []byte(`{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}}`),
			err:   "target missing",
		},
		{
			name:  "TargetInvalid",
			input: []byte(`{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":true}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.checkpointJSON",
		},
		{
			name:  "Good",
			input: []byte(`{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res spec.AttestationData
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
