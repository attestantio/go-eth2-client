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

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	require "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestSyncCommitteeDutyJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.syncCommitteeDutyJSON",
		},
		{
			name:  "PubKeyMissing",
			input: []byte(`{"validator_index":"1","validator_sync_committee_indices":["2","3","4"]}`),
			err:   "public key missing",
		},
		{
			name:  "PubKeyWrongType",
			input: []byte(`{"pubkey":true,"validator_index":"1","validator_sync_committee_indices":["2","3","4"]}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeDutyJSON.pubkey of type string",
		},
		{
			name:  "PubKeyInvalid",
			input: []byte(`{"pubkey":"invalid","validator_index":"1","validator_sync_committee_indices":["2","3","4"]}`),
			err:   "invalid value for public key: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "PubKeyShort",
			input: []byte(`{"pubkey":"0x9bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","validator_index":"1","validator_sync_committee_indices":["2","3","4"]}`),
			err:   "incorrect length for public key",
		},
		{
			name:  "ValidatorIndexMissing",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","validator_sync_committee_indices":["2","3","4"]}`),
			err:   "validator index missing",
		},
		{
			name:  "ValidatorIndexWrongType",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","validator_index":true,"validator_sync_committee_indices":["2","3","4"]}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeDutyJSON.validator_index of type string",
		},
		{
			name:  "ValidatorIndexInvalid",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","validator_index":"-1","validator_sync_committee_indices":["2","3","4"]}`),
			err:   "invalid value for validator index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ValidatorSyncCommitteeIndicesMissing",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","validator_index":"1"}`),
			err:   "validator sync committee indices missing",
		},
		{
			name:  "ValidatorSyncCommitteeIndicesWrongType",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","validator_index":"1","validator_sync_committee_indices":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeDutyJSON.validator_sync_committee_indices of type []string",
		},
		{
			name:  "ValidatorSyncCommitteeIndicesEmpty",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","validator_index":"1","validator_sync_committee_indices":[]}`),
			err:   "validator sync committee indices missing",
		},
		{
			name:  "ValidatorSyncCommitteeIndexWrongType",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","validator_index":"1","validator_sync_committee_indices":["2","3",true]}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeDutyJSON.validator_sync_committee_indices of type string",
		},
		{
			name:  "ValidatorSyncCommitteeIndexInvalid",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","validator_index":"1","validator_sync_committee_indices":["2","3","-1"]}`),
			err:   "invalid value for sync committee index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"pubkey":"0xb89bebc699769726a318c8e9971bd3171297c61aea4a6578a7a4f94b547dcba5bac16a89108b6b6a1fe3695d1a874a0b","validator_index":"1","validator_sync_committee_indices":["2","3","4"]}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.SyncCommitteeDuty
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
