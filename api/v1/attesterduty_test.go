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

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	require "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestAttesterDutyJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.attesterDutyJSON",
		},
		{
			name:  "PublicKeyMissing",
			input: []byte(`{"slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "public key missing",
		},
		{
			name:  "PublicKeyWrongType",
			input: []byte(`{"pubkey":[true],"slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "invalid JSON: json: cannot unmarshal array into Go struct field attesterDutyJSON.pubkey of type string",
		},
		{
			name:  "PublicKeyInvalid",
			input: []byte(`{"pubkey":"invalid","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "invalid value for public key: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "PublicKeyShort",
			input: []byte(`{"pubkey":"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "incorrect length for public key",
		},
		{
			name:  "PublicKeyLong",
			input: []byte(`{"pubkey":"0x00000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "incorrect length for public key",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":true,"validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field attesterDutyJSON.slot of type string",
		},
		{
			name:  "SlotInvalidValue",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"-1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "ValidatorIndexMissing",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "validator index missing",
		},
		{
			name:  "ValidatorIndexWrongType",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":true,"committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field attesterDutyJSON.validator_index of type string",
		},
		{
			name:  "ValidatorIndexInvalid",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"-1","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "invalid value for validator index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "CommitteeIndexMissing",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "committee index missing",
		},
		{
			name:  "CommitteeIndexWrongType",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":true,"committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field attesterDutyJSON.committee_index of type string",
		},
		{
			name:  "CommitteeIndexInvalid",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"-1","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "invalid value for committee index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "CommitteeLengthMissing",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "committee length missing",
		},
		{
			name:  "CommitteeLengthWrongType",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":true,"committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field attesterDutyJSON.committee_length of type string",
		},
		{
			name:  "CommitteeLengthInvalid",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"-1","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "invalid value for committee length: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "CommitteeLengthRejectZero",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"0","committees_at_slot":"4","validator_committee_index":"61"}`),
			err:   "committee length cannot be 0",
		},
		{
			name:  "CommitteesAtSlotMissing",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","validator_committee_index":"61"}`),
			err:   "committees at slot missing",
		},
		{
			name:  "CommitteesAtSlotWrongType",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":true,"validator_committee_index":"61"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field attesterDutyJSON.committees_at_slot of type string",
		},
		{
			name:  "CommitteesAtSlotInvalid",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"invalid","validator_committee_index":"61"}`),
			err:   "invalid value for committees at slot: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "CommitteesAtSlotZero",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"0","validator_committee_index":"61"}`),
			err:   "committees at slot cannot be 0",
		},
		{
			name:  "ValidatorCommitteeIndexMissing",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4"}`),
			err:   "validator committee index missing",
		},
		{
			name:  "ValidatorCommitteeIndexWrongType",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field attesterDutyJSON.validator_committee_index of type string",
		},
		{
			name:  "ValidatorCommitteeIndexInvalid",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"-1"}`),
			err:   "invalid value for validator committee index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.AttesterDuty
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
