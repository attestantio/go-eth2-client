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
	"github.com/stretchr/testify/require"
)

func TestVersionedExecutionPayloadEnvelopeAccessors(t *testing.T) {
	payload := &gloas.ExecutionPayload{BaseFeePerGas: uint256.NewInt(1)}
	executionRequests := &gloas.ExecutionRequests{}
	beaconBlockRoot := phase0.Root{0x01}
	parentBeaconBlockRoot := phase0.Root{0x02}
	populated := &spec.VersionedExecutionPayloadEnvelope{
		Version: spec.DataVersionGloas,
		Gloas: &gloas.ExecutionPayloadEnvelope{
			Payload:               payload,
			ExecutionRequests:     executionRequests,
			BuilderIndex:          3,
			BeaconBlockRoot:       beaconBlockRoot,
			ParentBeaconBlockRoot: parentBeaconBlockRoot,
		},
	}
	empty := &spec.VersionedExecutionPayloadEnvelope{Version: spec.DataVersionGloas}
	tests := []struct {
		name     string
		accessor func(*spec.VersionedExecutionPayloadEnvelope) (any, error)
		want     any
	}{
		{
			name: "Payload",
			accessor: func(v *spec.VersionedExecutionPayloadEnvelope) (any, error) {
				return v.Payload()
			},
			want: payload,
		},
		{
			name: "ExecutionRequests",
			accessor: func(v *spec.VersionedExecutionPayloadEnvelope) (any, error) {
				return v.ExecutionRequests()
			},
			want: &spec.VersionedExecutionRequests{
				Version: spec.DataVersionGloas,
				Gloas:   executionRequests,
			},
		},
		{
			name: "BuilderIndex",
			accessor: func(v *spec.VersionedExecutionPayloadEnvelope) (any, error) {
				return v.BuilderIndex()
			},
			want: gloas.BuilderIndex(3),
		},
		{
			name: "BeaconBlockRoot",
			accessor: func(v *spec.VersionedExecutionPayloadEnvelope) (any, error) {
				return v.BeaconBlockRoot()
			},
			want: beaconBlockRoot,
		},
		{
			name: "ParentBeaconBlockRoot",
			accessor: func(v *spec.VersionedExecutionPayloadEnvelope) (any, error) {
				return v.ParentBeaconBlockRoot()
			},
			want: parentBeaconBlockRoot,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.accessor(populated)
			require.NoError(t, err)
			require.Equal(t, test.want, got)

			_, err = test.accessor(empty)
			require.EqualError(t, err, "no gloas execution payload envelope")
		})
	}
}

func TestVersionedExecutionPayloadEnvelopeState(t *testing.T) {
	tests := []struct {
		name     string
		envelope *spec.VersionedExecutionPayloadEnvelope
		empty    bool
		expected string
	}{
		{
			name: "Gloas",
			envelope: &spec.VersionedExecutionPayloadEnvelope{
				Version: spec.DataVersionGloas,
				Gloas: &gloas.ExecutionPayloadEnvelope{
					Payload: &gloas.ExecutionPayload{BaseFeePerGas: uint256.NewInt(1)},
				},
			},
			expected: "payload:",
		},
		{
			name: "GloasNil",
			envelope: &spec.VersionedExecutionPayloadEnvelope{
				Version: spec.DataVersionGloas,
			},
			empty: true,
		},
		{
			name:     "UnknownVersion",
			envelope: &spec.VersionedExecutionPayloadEnvelope{},
			empty:    true,
			expected: "unknown version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.empty, test.envelope.IsEmpty())
			if test.expected == "" {
				require.Empty(t, test.envelope.String())
			} else {
				require.Contains(t, test.envelope.String(), test.expected)
			}
		})
	}
}
