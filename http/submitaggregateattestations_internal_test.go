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

package http

import (
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestCreateUnversionedAggregates verifies that every version arm unwraps the field matching
// its version for submission, that mixed versions are rejected, and that an unrecognised
// version does not fall through silently.
func TestCreateUnversionedAggregates(t *testing.T) {
	phase0Aggregate := func(index phase0.ValidatorIndex) *phase0.SignedAggregateAndProof {
		return &phase0.SignedAggregateAndProof{
			Message: &phase0.AggregateAndProof{AggregatorIndex: index},
		}
	}
	electraAggregate := func(index phase0.ValidatorIndex) *electra.SignedAggregateAndProof {
		return &electra.SignedAggregateAndProof{
			Message: &electra.AggregateAndProof{AggregatorIndex: index},
		}
	}
	gloasAggregate := func(index phase0.ValidatorIndex) *gloas.SignedAggregateAndProof {
		return &gloas.SignedAggregateAndProof{
			Message: &gloas.AggregateAndProof{AggregatorIndex: index},
		}
	}

	// A distinct aggregator index per fork, so forwarding the wrong field produces a legible diff.
	phase0Value := phase0Aggregate(1)
	altairValue := phase0Aggregate(2)
	bellatrixValue := phase0Aggregate(3)
	capellaValue := phase0Aggregate(4)
	denebValue := phase0Aggregate(5)
	electraValue := electraAggregate(6)
	fuluValue := electraAggregate(7)
	gloasValue := gloasAggregate(8)

	first, second, third := gloasAggregate(42), gloasAggregate(43), gloasAggregate(44)

	tests := []struct {
		name               string
		aggregateAndProofs []*spec.VersionedSignedAggregateAndProof
		expected           []any
		expectedErr        string
	}{
		{
			name: "Phase0",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionPhase0, Phase0: phase0Value},
			},
			expected: []any{phase0Value},
		},
		{
			name: "Phase0NilAggregate",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionPhase0, Phase0: nil},
			},
			expectedErr: "nil phase0 aggregate and proof supplied\ninvalid options",
		},
		{
			name: "Altair",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionAltair, Altair: altairValue},
			},
			expected: []any{altairValue},
		},
		{
			name: "AltairNilAggregate",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionAltair, Altair: nil},
			},
			expectedErr: "nil altair aggregate and proof supplied\ninvalid options",
		},
		{
			name: "Bellatrix",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionBellatrix, Bellatrix: bellatrixValue},
			},
			expected: []any{bellatrixValue},
		},
		{
			name: "BellatrixNilAggregate",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionBellatrix, Bellatrix: nil},
			},
			expectedErr: "nil bellatrix aggregate and proof supplied\ninvalid options",
		},
		{
			name: "Capella",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionCapella, Capella: capellaValue},
			},
			expected: []any{capellaValue},
		},
		{
			name: "CapellaNilAggregate",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionCapella, Capella: nil},
			},
			expectedErr: "nil capella aggregate and proof supplied\ninvalid options",
		},
		{
			name: "Deneb",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionDeneb, Deneb: denebValue},
			},
			expected: []any{denebValue},
		},
		{
			name: "DenebNilAggregate",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionDeneb, Deneb: nil},
			},
			expectedErr: "nil deneb aggregate and proof supplied\ninvalid options",
		},
		{
			name: "Electra",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionElectra, Electra: electraValue},
			},
			expected: []any{electraValue},
		},
		{
			name: "ElectraNilAggregate",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionElectra, Electra: nil},
			},
			expectedErr: "nil electra aggregate and proof supplied\ninvalid options",
		},
		{
			name: "Fulu",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionFulu, Fulu: fuluValue},
			},
			expected: []any{fuluValue},
		},
		{
			name: "FuluNilAggregate",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionFulu, Fulu: nil},
			},
			expectedErr: "nil fulu aggregate and proof supplied\ninvalid options",
		},
		{
			name: "Gloas",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionGloas, Gloas: gloasValue},
			},
			expected: []any{gloasValue},
		},
		{
			name: "GloasSeveralAggregates",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionGloas, Gloas: first},
				{Version: spec.DataVersionGloas, Gloas: second},
				{Version: spec.DataVersionGloas, Gloas: third},
			},
			expected: []any{first, second, third},
		},
		{
			name: "GloasNilAggregate",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionGloas},
			},
			expectedErr: "nil gloas aggregate and proof supplied\ninvalid options",
		},
		{
			name: "NilVersionedAggregate",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				nil,
			},
			expectedErr: "nil aggregate and proof version supplied\ninvalid options",
		},
		{
			name: "MixedVersions",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersionElectra, Electra: electraValue},
				{Version: spec.DataVersionGloas, Gloas: gloasValue},
			},
			expectedErr: "aggregate and proofs must all be of the same version\ninvalid options",
		},
		{
			name: "UnknownVersion",
			aggregateAndProofs: []*spec.VersionedSignedAggregateAndProof{
				{Version: spec.DataVersion(999)},
			},
			expectedErr: "unknown aggregate and proof version\ninvalid options",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			unversioned, err := createUnversionedAggregates(test.aggregateAndProofs)
			if test.expectedErr != "" {
				require.ErrorIs(t, err, client.ErrInvalidOptions)
				require.EqualError(t, err, test.expectedErr)

				return
			}
			require.NoError(t, err)
			require.Equal(t, test.expected, unversioned)
		})
	}
}
