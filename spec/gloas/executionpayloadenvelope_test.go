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
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	require "github.com/stretchr/testify/require"
)

// TestExecutionPayloadEnvelopeString verifies that stringifying an envelope does not
// panic when its payload carries a nil BaseFeePerGas. The envelope's YAML marshaler
// embeds the payload, so String() descends into ExecutionPayload.MarshalYAML and
// depends on that marshaler's nil guard — this is the shape
// mock.Service.SignedExecutionPayloadEnvelope builds, which made logging a mocked or
// fetched envelope crash the caller.
func TestExecutionPayloadEnvelopeString(t *testing.T) {
	tests := []struct {
		name     string
		envelope *gloas.ExecutionPayloadEnvelope
		// expected is a substring the rendering must contain.
		expected string
	}{
		{
			// A nil payload is rendered as null by the YAML encoder, so this case
			// never reaches the payload marshaler.
			name:     "ZeroValue",
			envelope: &gloas.ExecutionPayloadEnvelope{},
			expected: "payload: null",
		},
		{
			name:     "NilBaseFeePerGas",
			envelope: &gloas.ExecutionPayloadEnvelope{Payload: &gloas.ExecutionPayload{}},
			expected: `base_fee_per_gas: '0'`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var str string
			require.NotPanics(t, func() {
				str = test.envelope.String()
			})
			require.Contains(t, str, test.expected)
		})
	}
}
