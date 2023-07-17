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

func TestBeaconBlockHeaderJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type phase0.beaconBlockHeaderJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"slot":true,"proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconBlockHeaderJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ProposerIndexMissing",
			input: []byte(`{"slot":"1","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "proposer index missing",
		},
		{
			name:  "ProposerWrongType",
			input: []byte(`{"slot":"1","proposer_index":true,"parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconBlockHeaderJSON.proposer_index of type string",
		},
		{
			name:  "ProposerInvalid",
			input: []byte(`{"slot":"1","proposer_index":"-1","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "invalid value for proposer index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ParentRootMissing",
			input: []byte(`{"slot":"1","proposer_index":"2","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "parent root missing",
		},
		{
			name:  "ParentRootWrongType",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":true,"state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconBlockHeaderJSON.parent_root of type string",
		},
		{
			name:  "ParentRootInvalid",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"invalid","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "invalid value for parent root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "ParentRootShort",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "incorrect length for parent root",
		},
		{
			name:  "ParentRootLong",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x00000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "incorrect length for parent root",
		},
		{
			name:  "StateRootMissing",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "state root missing",
		},
		{
			name:  "StateRootWrongType",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":true,"body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconBlockHeaderJSON.state_root of type string",
		},
		{
			name:  "StateRootInvalid",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"invalid","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "invalid value for state root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "StateRootShort",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x2122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "incorrect length for state root",
		},
		{
			name:  "StateRootLong",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x20202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "incorrect length for state root",
		},
		{
			name:  "BodyRootMissing",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"}`),
			err:   "body root missing",
		},
		{
			name:  "BodyRootWrongType",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconBlockHeaderJSON.body_root of type string",
		},
		{
			name:  "BodyRootInvalid",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"invalid"}`),
			err:   "invalid value for body root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "BodyRootShort",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x4142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "incorrect length for body root",
		},
		{
			name:  "BodyRootLong",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x40404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
			err:   "incorrect length for body root",
		},
		{
			name:  "Good",
			input: []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.BeaconBlockHeader
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

func TestBeaconBlockHeaderYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{slot: 13763605646674810374, proposer_index: 7704591909929553516, parent_root: '0x98f5655842c7d96474354de483dde443a6eb3180f102c425aae3b6a924332de1', state_root: '0x38cee1c1ff58231fe19959d16cb805af25618ca1175c39885a8a9a0f6f6560ea', body_root: '0x36122528bead53899c96cdcc4681c202802edd33a9f7c877c035f566e1fef59e'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.BeaconBlockHeader
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
