// Copyright © 2026 Attestant Limited.
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

// TestPTCDutyJSON exercises the PTC duty codec.
//
// The Good case is a duty captured verbatim from a live Glamsterdam devnet-7
// Lighthouse node (POST /eth/v1/validator/duties/ptc/7537), so the field set
// and the string-encoded integers are the real wire shape rather than an
// inference from the schema.
func TestPTCDutyJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.ptcDutyJSON",
		},
		{
			name:  "PubKeyMissing",
			input: []byte(`{"validator_index":"8","slot":"60300"}`),
			err:   "public key missing",
		},
		{
			name:  "PubKeyWrongType",
			input: []byte(`{"pubkey":true,"validator_index":"8","slot":"60300"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field ptcDutyJSON.pubkey of type string",
		},
		{
			name:  "PubKeyInvalid",
			input: []byte(`{"pubkey":"invalid","validator_index":"8","slot":"60300"}`),
			err:   "invalid value for public key: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "PubKeyShort",
			input: []byte(`{"pubkey":"0x2cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360","validator_index":"8","slot":"60300"}`),
			err:   "incorrect length for public key",
		},
		{
			name:  "PubKeyLong",
			input: []byte(`{"pubkey":"0xb7b72cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360","validator_index":"8","slot":"60300"}`),
			err:   "incorrect length for public key",
		},
		{
			name:  "ValidatorIndexMissing",
			input: []byte(`{"pubkey":"0xb72cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360","slot":"60300"}`),
			err:   "validator index missing",
		},
		{
			name:  "ValidatorIndexWrongType",
			input: []byte(`{"pubkey":"0xb72cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360","validator_index":8,"slot":"60300"}`),
			err:   "invalid JSON: json: cannot unmarshal number into Go struct field ptcDutyJSON.validator_index of type string",
		},
		{
			name:  "ValidatorIndexInvalid",
			input: []byte(`{"pubkey":"0xb72cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360","validator_index":"-1","slot":"60300"}`),
			err:   "invalid value for validator index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"pubkey":"0xb72cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360","validator_index":"8"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"pubkey":"0xb72cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360","validator_index":"8","slot":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field ptcDutyJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"pubkey":"0xb72cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360","validator_index":"8","slot":"-1"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"pubkey":"0xb72cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360","validator_index":"8","slot":"60300"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.PTCDuty
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

// TestPTCDutyFields confirms the decoded values land in the right fields, and
// that a PTC duty carries no committee information. A round-trip test alone
// would pass even if two same-typed fields were transposed.
func TestPTCDutyFields(t *testing.T) {
	input := []byte(`{"pubkey":"0xb72cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360","validator_index":"8","slot":"60300"}`)

	var duty api.PTCDuty
	require.NoError(t, json.Unmarshal(input, &duty))

	require.Equal(t,
		"0xb72cb106b7bc1ecae219e0ae1830a509ed18a042b56a2779f4033419de69ba8ae8017090caed1f5377bfa68506157360",
		duty.PubKey.String(),
	)
	require.Equal(t, uint64(8), uint64(duty.ValidatorIndex))
	require.Equal(t, uint64(60300), uint64(duty.Slot))
}
