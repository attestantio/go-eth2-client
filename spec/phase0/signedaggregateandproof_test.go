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

func TestSignedAggregateAndProofJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type phase0.signedAggregateAndProofJSON",
		},
		{
			name:  "MessageMissing",
			input: []byte(`{"signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "message missing",
		},
		{
			name:  "MessageWrongType",
			input: []byte(`{"message":true,"signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.aggregateAndProofJSON",
		},
		{
			name:  "MessageInvalid",
			input: []byte(`{"message":{},"signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "invalid JSON: aggregator index missing",
		},
		{
			name:  "SignatureMissing",
			input: []byte(`{"message":{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}}`),
			err:   "signature missing",
		},
		{
			name:  "SignatureWrongType",
			input: []byte(`{"message":{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"},"signature":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field signedAggregateAndProofJSON.signature of type string",
		},
		{
			name:  "SignatureInvalid",
			input: []byte(`{"message":{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"},"signature":"invalid"}`),
			err:   "invalid value for signature: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SignatureShort",
			input: []byte(`{"message":{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"},"signature":"0xead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "incorrect length for signature",
		},
		{
			name:  "SignatureLong",
			input: []byte(`{"message":{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"},"signature":"0xb4b4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "incorrect length for signature",
		},
		{
			name:  "Good",
			input: []byte(`{"message":{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"},"signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.SignedAggregateAndProof
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

func TestSignedAggregateAndProofYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{message: {aggregator_index: 402, aggregate: {aggregation_bits: '0xffffffff01', data: {slot: 66, index: 0, beacon_block_root: '0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100', source: {epoch: 0, root: '0x0000000000000000000000000000000000000000000000000000000000000000'}, target: {epoch: 2, root: '0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440'}}, signature: '0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d'}, selection_proof: '0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c'}, signature: '0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res phase0.SignedAggregateAndProof
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
