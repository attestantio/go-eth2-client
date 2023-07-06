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
	require "github.com/stretchr/testify/require"
)

func TestSyncCommitteeMessageJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type altair.syncCommitteeMessageJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"2","signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"slot":true,"beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"2","signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeMessageJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"2","signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "BeaconBlockRootMissing",
			input: []byte(`{"slot":"1","validator_index":"2","signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "beacon block root missing",
		},
		{
			name:  "BeaconBlockRootWrongType",
			input: []byte(`{"slot":"1","beacon_block_root":true,"validator_index":"2","signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeMessageJSON.beacon_block_root of type string",
		},
		{
			name:  "BeaconBlockRootShort",
			input: []byte(`{"slot":"1","beacon_block_root":"0xcd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"2","signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "incorrect length for beacon block root",
		},
		{
			name:  "BeaconBlockRootLong",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbabacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"2","signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "incorrect length for beacon block root",
		},
		{
			name:  "BeaconBlockRootInvalid",
			input: []byte(`{"slot":"1","beacon_block_root":"invalid","validator_index":"2","signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "invalid value for beacon block root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "ValidatorIndexMissing",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "validator index missing",
		},
		{
			name:  "ValidatorIndexWrongType",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":true,"signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeMessageJSON.validator_index of type string",
		},
		{
			name:  "ValidatorIndexInvalid",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"-1","signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "invalid value for validator index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "SignatureMissing",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"2"}`),
			err:   "signature missing",
		},
		{
			name:  "SignatureWrongType",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"2","signature":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeMessageJSON.signature of type string",
		},
		{
			name:  "SignatureShort",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"2","signature":"0xead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "incorrect length for signature",
		},
		{
			name:  "SignatureLong",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"2","signature":"0xb4b4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
			err:   "incorrect length for signature",
		},
		{
			name:  "SignatureInvalid",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"2","signature":"invalid"}`),
			err:   "invalid value for signature: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "Good",
			input: []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","validator_index":"2","signature":"0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.SyncCommitteeMessage
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

func TestSyncCommitteeMessageYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{slot: 1, beacon_block_root: '0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c', validator_index: 2, signature: '0xb4ead6da46dc0ce26343defc6f9607987ce0ecad5073e48c71f21d1a198cd68600a4c434dca26310460999c564885b6901c6f59ec3db84bd8e7adede27c5fdb270042a57d50415afe509c0c88edc5c611ca6f63bed63c88714ed56987ee3ca8f'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.SyncCommitteeMessage
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
