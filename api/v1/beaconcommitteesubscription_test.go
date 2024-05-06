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

func TestBeaconCommitteeSubscriptionJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.beaconCommitteeSubscriptionJSON",
		},
		{
			name:  "ValidatorIndexMissing",
			input: []byte(`{"slot":"1","committee_index":"2","committees_at_slot":"5","is_aggregator":true}`),
			err:   "validator index missing",
		},
		{
			name:  "ValidatorIndexWrongType",
			input: []byte(`{"validator_index":true,"slot":"1","committee_index":"2","committees_at_slot":"5","is_aggregator":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconCommitteeSubscriptionJSON.validator_index of type string",
		},
		{
			name:  "ValidatorIndexInvalid",
			input: []byte(`{"validator_index":"-1","slot":"1","committee_index":"2","committees_at_slot":"5","is_aggregator":true}`),
			err:   "invalid value for validator index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"validator_index":"10","committee_index":"2","committees_at_slot":"5","is_aggregator":true}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"validator_index":"10","slot":true,"committee_index":"2","committees_at_slot":"5","is_aggregator":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconCommitteeSubscriptionJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"validator_index":"10","slot":"-1","committee_index":"2","committees_at_slot":"5","is_aggregator":true}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "CommitteeIndexMissing",
			input: []byte(`{"validator_index":"10","slot":"1","committees_at_slot":"5","is_aggregator":true}`),
			err:   "committee index missing",
		},
		{
			name:  "CommitteeIndexWrongType",
			input: []byte(`{"validator_index":"10","slot":"1","committee_index":true,"committees_at_slot":"5","is_aggregator":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconCommitteeSubscriptionJSON.committee_index of type string",
		},
		{
			name:  "CommitteeIndexMissing",
			input: []byte(`{"validator_index":"10","slot":"1","committee_index":"-1","committees_at_slot":"5","is_aggregator":true}`),
			err:   "invalid value for committee index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "CommitteesAtSlotMissing",
			input: []byte(`{"validator_index":"10","slot":"1","committee_index":"2","is_aggregator":true}`),
			err:   "committees at slot missing",
		},
		{
			name:  "CommitteesAtSlotWrongType",
			input: []byte(`{"validator_index":"10","slot":"1","committee_index":"2","committees_at_slot":true,"is_aggregator":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconCommitteeSubscriptionJSON.committees_at_slot of type string",
		},
		{
			name:  "CommitteesAtSlotInvalid",
			input: []byte(`{"validator_index":"10","slot":"1","committee_index":"2","committees_at_slot":"-1","is_aggregator":true}`),
			err:   "invalid value for committees at slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "CommitteesAtSlotInvalid2",
			input: []byte(`{"validator_index":"10","slot":"1","committee_index":"2","committees_at_slot":"0","is_aggregator":true}`),
			err:   "committees at slot cannot be 0",
		},
		{
			name:  "Good",
			input: []byte(`{"validator_index":"10","slot":"1","committee_index":"2","committees_at_slot":"5","is_aggregator":true}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.BeaconCommitteeSubscription
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
