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
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
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

func TestVersionedSignedExecutionPayloadEnvelopeUnknownVersion(t *testing.T) {
	_, err := (&spec.VersionedSignedExecutionPayloadEnvelope{
		Version: spec.DataVersion(255),
	}).Payload()
	require.EqualError(t, err, "unknown version")
}

func TestVersionedSignedExecutionPayloadEnvelopeAccessors(t *testing.T) {
	payload := &gloas.ExecutionPayload{BaseFeePerGas: uint256.NewInt(1)}
	executionRequests := &gloas.ExecutionRequests{}
	beaconBlockRoot := phase0.Root{0x01}
	parentBeaconBlockRoot := phase0.Root{0x02}
	message := &gloas.ExecutionPayloadEnvelope{
		Payload:               payload,
		ExecutionRequests:     executionRequests,
		BuilderIndex:          3,
		BeaconBlockRoot:       beaconBlockRoot,
		ParentBeaconBlockRoot: parentBeaconBlockRoot,
	}
	signature := phase0.BLSSignature{0x04}
	populated := &spec.VersionedSignedExecutionPayloadEnvelope{
		Version: spec.DataVersionGloas,
		Gloas: &gloas.SignedExecutionPayloadEnvelope{
			Message:   message,
			Signature: signature,
		},
	}
	gloasNoMessage := &spec.VersionedSignedExecutionPayloadEnvelope{
		Version: spec.DataVersionGloas,
		Gloas:   &gloas.SignedExecutionPayloadEnvelope{},
	}
	gloasNil := &spec.VersionedSignedExecutionPayloadEnvelope{Version: spec.DataVersionGloas}
	tests := []struct {
		name     string
		accessor func(*spec.VersionedSignedExecutionPayloadEnvelope) (any, error)
		want     any
		missing  *spec.VersionedSignedExecutionPayloadEnvelope
	}{
		{
			name: "Message",
			accessor: func(v *spec.VersionedSignedExecutionPayloadEnvelope) (any, error) {
				return v.Message()
			},
			want:    message,
			missing: gloasNoMessage,
		},
		{
			name: "Payload",
			accessor: func(v *spec.VersionedSignedExecutionPayloadEnvelope) (any, error) {
				return v.Payload()
			},
			want:    payload,
			missing: gloasNoMessage,
		},
		{
			name: "ExecutionRequests",
			accessor: func(v *spec.VersionedSignedExecutionPayloadEnvelope) (any, error) {
				return v.ExecutionRequests()
			},
			want: &spec.VersionedExecutionRequests{
				Version: spec.DataVersionGloas,
				Gloas:   executionRequests,
			},
			missing: gloasNoMessage,
		},
		{
			name: "BuilderIndex",
			accessor: func(v *spec.VersionedSignedExecutionPayloadEnvelope) (any, error) {
				return v.BuilderIndex()
			},
			want:    gloas.BuilderIndex(3),
			missing: gloasNoMessage,
		},
		{
			name: "BeaconBlockRoot",
			accessor: func(v *spec.VersionedSignedExecutionPayloadEnvelope) (any, error) {
				return v.BeaconBlockRoot()
			},
			want:    beaconBlockRoot,
			missing: gloasNoMessage,
		},
		{
			name: "ParentBeaconBlockRoot",
			accessor: func(v *spec.VersionedSignedExecutionPayloadEnvelope) (any, error) {
				return v.ParentBeaconBlockRoot()
			},
			want:    parentBeaconBlockRoot,
			missing: gloasNoMessage,
		},
		{
			name: "Signature",
			accessor: func(v *spec.VersionedSignedExecutionPayloadEnvelope) (any, error) {
				return v.Signature()
			},
			want:    signature,
			missing: gloasNil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.accessor(populated)
			require.NoError(t, err)
			require.Equal(t, test.want, got)

			_, err = test.accessor(test.missing)
			require.EqualError(t, err, "no gloas signed execution payload envelope")
		})
	}
}

func TestVersionedSignedExecutionPayloadEnvelopeIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		envelope *spec.VersionedSignedExecutionPayloadEnvelope
		expected bool
	}{
		{
			name: "Gloas",
			envelope: &spec.VersionedSignedExecutionPayloadEnvelope{
				Version: spec.DataVersionGloas,
				Gloas:   &gloas.SignedExecutionPayloadEnvelope{},
			},
		},
		{
			name: "GloasNil",
			envelope: &spec.VersionedSignedExecutionPayloadEnvelope{
				Version: spec.DataVersionGloas,
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, test.envelope.IsEmpty())
		})
	}
}
