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

func TestAggregateAndProofJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type phase0.aggregateAndProofJSON",
		},
		{
			name:  "AggregatorIndexMissing",
			input: []byte(`{"aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "aggregator index missing",
		},
		{
			name:  "AggregatorIndexWrongType",
			input: []byte(`{"aggregator_index":true,"aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field aggregateAndProofJSON.aggregator_index of type string",
		},
		{
			name:  "AggregatorIndexInvalid",
			input: []byte(`{"aggregator_index":"-1","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid value for aggregator index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "AggregateMissing",
			input: []byte(`{"aggregator_index":"402","selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "aggregate missing",
		},
		{
			name:  "AggregateWrongType",
			input: []byte(`{"aggregator_index":"402","aggregate":true,"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.attestationJSON",
		},
		{
			name:  "AggregateInvalid",
			input: []byte(`{"aggregator_index":"402","aggregate":{},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid JSON: aggregation bits missing",
		},
		{
			name:  "SelectionProofMissing",
			input: []byte(`{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"}}`),
			err:   "selection proof missing",
		},
		{
			name:  "SelectionProofWrongType",
			input: []byte(`{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field aggregateAndProofJSON.selection_proof of type string",
		},
		{
			name:  "SelectionProofInvalid",
			input: []byte(`{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"invalid"}`),
			err:   "invalid value for selection proof: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SelectionProofShort",
			input: []byte(`{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "incorrect length for selection proof",
		},
		{
			name:  "SelectionProofLong",
			input: []byte(`{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "incorrect length for selection proof",
		},
		{
			name:  "Good",
			input: []byte(`{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.AggregateAndProof
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

func TestAggregateAndProofYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{aggregator_index: 402, aggregate: {aggregation_bits: '0xffffffff01', data: {slot: 66, index: 0, beacon_block_root: '0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100', source: {epoch: 0, root: '0x0000000000000000000000000000000000000000000000000000000000000000'}, target: {epoch: 2, root: '0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440'}}, signature: '0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d'}, selection_proof: '0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.AggregateAndProof
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
