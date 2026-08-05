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

package spec_test

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	require "github.com/stretchr/testify/require"
)

// TestVersionedSignedExecutionPayloadEnvelopeString verifies that stringifying the
// versioned wrapper does not panic when the innermost payload carries a nil
// BaseFeePerGas. String() delegates to the gloas envelope, whose YAML marshaler
// descends into ExecutionPayload.MarshalYAML, so the whole chain depends on that
// marshaler's nil guard.
func TestVersionedSignedExecutionPayloadEnvelopeString(t *testing.T) {
	tests := []struct {
		name     string
		envelope *spec.VersionedSignedExecutionPayloadEnvelope
		// expected is a substring the rendering must contain.
		expected string
	}{
		{
			// The zero value of DataVersion is DataVersionUnknown, not
			// DataVersionPhase0, so a zero-value wrapper takes the default arm.
			name:     "ZeroValue",
			envelope: &spec.VersionedSignedExecutionPayloadEnvelope{},
			expected: "unknown version",
		},
		{
			// The shape mock.Service.SignedExecutionPayloadEnvelope returns.
			name: "GloasNilBaseFeePerGas",
			envelope: &spec.VersionedSignedExecutionPayloadEnvelope{
				Version: spec.DataVersionGloas,
				Gloas: &gloas.SignedExecutionPayloadEnvelope{
					Message: &gloas.ExecutionPayloadEnvelope{
						Payload: &gloas.ExecutionPayload{},
					},
				},
			},
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
