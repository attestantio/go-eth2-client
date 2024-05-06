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

func TestSyncCommitteeSubscriptionJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.syncCommitteeSubscriptionJSON",
		},
		{
			name:  "ValidatorIndexMissing",
			input: []byte(`{"sync_committee_indices":["2"],"until_epoch":"5"}`),
			err:   "validator index missing",
		},
		{
			name:  "ValidatorIndexWrongType",
			input: []byte(`{"validator_index":true,"sync_committee_indices":["2"],"until_epoch":"5"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeSubscriptionJSON.validator_index of type string",
		},
		{
			name:  "ValidatorIndexInvalid",
			input: []byte(`{"validator_index":"invalid","sync_committee_indices":["2"],"until_epoch":"5"}`),
			err:   "invalid value for validator index: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "SyncCommitteeIndicesMissing",
			input: []byte(`{"validator_index":"10","until_epoch":"5"}`),
			err:   "sync committee indices missing",
		},
		{
			name:  "SyncCommitteeIndicesWrongType",
			input: []byte(`{"validator_index":"10","sync_committee_indices":true,"until_epoch":"5"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeSubscriptionJSON.sync_committee_indices of type []string",
		},
		{
			name:  "SyncCommitteeIndexWrongType",
			input: []byte(`{"validator_index":"10","sync_committee_indices":[true],"until_epoch":"5"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeSubscriptionJSON.sync_committee_indices of type string",
		},
		{
			name:  "SyncCommitteeIndexInvalid",
			input: []byte(`{"validator_index":"10","sync_committee_indices":["invalid"],"until_epoch":"5"}`),
			err:   "invalid value for sync committee index: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "UntilEpochMissing",
			input: []byte(`{"validator_index":"10","sync_committee_indices":["2"]}`),
			err:   "until epoch missing",
		},
		{
			name:  "UntilEpochWrongType",
			input: []byte(`{"validator_index":"10","sync_committee_indices":["2"],"until_epoch":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncCommitteeSubscriptionJSON.until_epoch of type string",
		},
		{
			name:  "UntilEpochInvalid",
			input: []byte(`{"validator_index":"10","sync_committee_indices":["2"],"until_epoch":"invalid"}`),
			err:   "invalid value for until epoch: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"validator_index":"10","sync_committee_indices":["2"],"until_epoch":"5"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.SyncCommitteeSubscription
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
