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

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	require "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestEvent(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.eventJSON",
		},
		{
			name:  "TopicMissing",
			input: []byte(`{"data":{}}`),
			err:   "topic missing",
		},
		{
			name:  "TopicWrongType",
			input: []byte(`{"topic":[],"data":{}}`),
			err:   "invalid JSON: json: cannot unmarshal array into Go struct field eventJSON.topic of type string",
		},
		{
			name:  "TopicUnsupported",
			input: []byte(`{"topic":"foo","data":{"block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","current_duty_dependent_root":"0x92c6b763f610d5941d2041906007bf9449d37772aacf0483a76275ac27c096b4","epoch_transition":false,"previous_duty_dependent_root":"0xa692c095bbca3eeaf99eeabada78874c028c02b176ccf691f3e8fa075d67f5c6","slot":"231192","state":"0x61099b2c1dee0104c93ce0e14e5f5fc4b6faceff4cb863278d055bdfb73b7dc7"}}`),
			err:   "unsupported event topic foo",
		},
		{
			name:  "DataMissing",
			input: []byte(`{"topic":"head"}`),
			err:   "data missing",
		},
		{
			name:  "GoodAttestation",
			input: []byte(`{"topic":"attestation","data":{"aggregation_bits":"0x010203","data":{"beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","index":"1","slot":"100","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}}`),
		},
		{
			name:  "GoodBlock",
			input: []byte(`{"topic":"block","data":{"block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","slot":"1"}}`),
		},
		{
			name:  "GoodChainReorg",
			input: []byte(`{"topic":"chain_reorg","data":{"depth":"2","epoch":"16405","new_head_block":"0xa3fe14d8d749318359aa3790d3588a23e12ea3b02bd879fbfbf04c3a66770df7","new_head_state":"0x4ab800aaa51c14c786fe7e924abd1355aa2ac2e0434d7cb5ae568720ed1bf522","old_head_block":"0x2ffc0a5b75de20f2a12853dff3e09b263e7c3cb19515134cba756b28e5ba25ee","old_head_state":"0x97cc0a37b77fbac6fa140f330c92521ddcd5b1dfefeef99d86996a51f1993d60","slot":"524986"}}`),
		},
		{
			name:  "GoodFinalizedCheckpoint",
			input: []byte(`{"topic":"finalized_checkpoint","data":{"block":"0x99e3f24aab3dd084045a0c927a33b8463eb5c7b17eeadfecdcf4e4badf7b6028","epoch":"2","state":"0x749a95b1355828b758864ea601c007e69aabed7b34a0f2084c43c26242f77e28"}}`),
		},
		{
			name:  "GoodHead",
			input: []byte(`{"topic":"head","data":{"block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","current_duty_dependent_root":"0x92c6b763f610d5941d2041906007bf9449d37772aacf0483a76275ac27c096b4","epoch_transition":false,"previous_duty_dependent_root":"0xa692c095bbca3eeaf99eeabada78874c028c02b176ccf691f3e8fa075d67f5c6","slot":"231192","state":"0x61099b2c1dee0104c93ce0e14e5f5fc4b6faceff4cb863278d055bdfb73b7dc7"}}`),
		},
		{
			name:  "GoodVoluntaryExit",
			input: []byte(`{"topic":"voluntary_exit","data":{"message":{"epoch":"1","validator_index":"2"},"signature":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`),
		},
		{
			name:  "GoodContributionAndProof",
			input: []byte(`{"topic":"contribution_and_proof","data":{"message":{"aggregator_index":"6568","contribution":{"aggregation_bits":"0x3f7f7f9fbffd9fddaf77fff7fffffdff","beacon_block_root":"0x3471a569ed74fb13f6638d7b759cd17c8ed08045d4668ae635349cc5f4dd2a75","signature":"0xa7260b90db427b85806cdaaecef08146a02c8c450aae96245be862f67fe6f54fefc7fdf1d4adfafad0164f5bbc0ceb65197b29e25b1dd9efd44ba8c390e95bd966b5dd97bf877a0ce277c757b68643054238659932348185775dc36d036b38da","slot":"45566","subcommittee_index":"2"},"selection_proof":"0x8c28b4b2f304f957735986e89ed3e429e007592e854d2c9a794333d5dfa05505412d70d0ba91e9fe3453816b01cd846415f82c864e7337e4796101ac9d2e351f2f8172d1d9061fd212f353ecf0ffd9dd17da42598adeae2046e5a74cbcb43474"},"signature":"0xb992ac86e1bbd6e2d1b7d18e8467aa435fecf583f5d13739db99b8d343093177caf167010b480c4e10f858f84cd05a1704c7ae3a253b1c454e5ebeefb7c35f7b8b51ba7aba0019cbc92d5bd8e9bcb61608a2ef47ce0a024b7b497ac9e813620f"}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.Event
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
