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

package gloas_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/holiman/uint256"
	require "github.com/stretchr/testify/require"
)

// TestExecutionPayloadMarshalJSON verifies that marshaling does not panic on a
// zero-value payload, whose BaseFeePerGas (*uint256.Int) is nil — mirroring the
// nil guard the SSZ marshaler already applies.
func TestExecutionPayloadMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		payload  *gloas.ExecutionPayload
		expected string
	}{
		{
			name:     "NilBaseFeePerGas",
			payload:  &gloas.ExecutionPayload{},
			expected: `"base_fee_per_gas":"0"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var data []byte
			var err error
			require.NotPanics(t, func() {
				data, err = json.Marshal(test.payload)
			})
			require.NoError(t, err)
			require.Contains(t, string(data), test.expected)
		})
	}
}

// TestExecutionPayloadUnmarshalJSON verifies that unmarshaling does not panic on
// a transaction element too short to carry the "0x" framing, whose decoded
// length (len-4)/2 would otherwise underflow to a negative make() size.
func TestExecutionPayloadUnmarshalJSON(t *testing.T) {
	// A valid payload marshaled to JSON gives a well-formed prefix up to the
	// (empty) transactions array; each case replaces it with a single raw
	// element too short to carry the "0x" framing, reaching the decode guard.
	// The names count the raw JSON element's bytes (e.g. `[1]` is one byte).
	base, err := json.Marshal(&gloas.ExecutionPayload{BaseFeePerGas: uint256.NewInt(0)})
	require.NoError(t, err)

	tests := []struct {
		name         string
		transactions string
		err          string
	}{
		{
			name:         "OneByteRawElement",
			transactions: `[1]`,
			err:          "transaction 0: missing or malformed",
		},
		{
			name:         "TwoByteRawElement",
			transactions: `[12]`,
			err:          "transaction 0: missing or malformed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := bytes.Replace(base, []byte(`"transactions":[]`), []byte(`"transactions":`+test.transactions), 1)
			var payload gloas.ExecutionPayload
			var err error
			require.NotPanics(t, func() {
				err = json.Unmarshal(input, &payload)
			})
			require.EqualError(t, err, test.err)
		})
	}
}
