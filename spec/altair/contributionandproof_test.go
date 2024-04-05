// Copyright © 2021 Attestant Limited.
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
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContributionAndProofJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type altair.contributionAndProofJSON",
		},
		{
			name:  "AggregatorIndexMissing",
			input: []byte(`{"contribution":{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "aggregator index missing",
		},
		{
			name:  "AggregatorIndexWrongType",
			input: []byte(`{"aggregator_index":true,"contribution":{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field contributionAndProofJSON.aggregator_index of type string",
		},
		{
			name:  "AggregatorIndexInvalid",
			input: []byte(`{"aggregator_index":"-1","contribution":{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid value for aggregator index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ContributionMissing",
			input: []byte(`{"aggregator_index":"402","selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "contribution missing",
		},
		{
			name:  "ContributionWrongType",
			input: []byte(`{"aggregator_index":"402","contribution":true,"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type altair.syncCommitteeContributionJSON",
		},
		{
			name:  "ContributionInvalid",
			input: []byte(`{"aggregator_index":"402","contribution":{},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid JSON: slot missing",
		},
		{
			name:  "SelectionProofMissing",
			input: []byte(`{"aggregator_index":"402","contribution":{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}}`),
			err:   "selection proof missing",
		},
		{
			name:  "SelectionProofWrongType",
			input: []byte(`{"aggregator_index":"402","contribution":{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"},"selection_proof":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field contributionAndProofJSON.selection_proof of type string",
		},
		{
			name:  "SelectionProofInvalid",
			input: []byte(`{"aggregator_index":"402","contribution":{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"},"selection_proof":"invalid"}`),
			err:   "invalid value for selection proof: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SelectionProofShort",
			input: []byte(`{"aggregator_index":"402","contribution":{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"},"selection_proof":"0x5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "incorrect length for selection proof",
		},
		{
			name:  "SelectionProofLong",
			input: []byte(`{"aggregator_index":"402","contribution":{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"},"selection_proof":"0x8b8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "incorrect length for selection proof",
		},
		{
			name:  "Good",
			input: []byte(`{"aggregator_index":"402","contribution":{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x0004000000000000000000000000000001","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.ContributionAndProof
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

func TestContributionAndProofYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{aggregator_index: 402, contribution: {slot: 1, beacon_block_root: '0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c', subcommittee_index: 3, aggregation_bits: '0x0004000000000000000000000000000001', signature: '0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c'}, selection_proof: '0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.ContributionAndProof
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, res.String(), string(rt))
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}
