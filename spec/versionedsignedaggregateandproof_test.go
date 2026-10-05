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
	"github.com/stretchr/testify/require"
)

// TestVersionedSignedAggregateAndProofGloasFetchToSubmit walks a Gloas aggregate through the
// sequence a validator actually performs, asserting the types line up at every hop with no
// conversion:
//
//	fetch:     AggregateAttestation returns a VersionedAttestation whose Gloas arm is a
//	           *gloas.Attestation
//	aggregate: that attestation is placed directly into a VersionedAggregateAndProof
//	sign:      HashTreeRoot on that wrapper yields the signing root
//	submit:    the signed result is wrapped in a VersionedSignedAggregateAndProof
//
// The hop that matters is placing the fetched *gloas.Attestation into the aggregate: while the
// Gloas arms were electra-typed this step did not compile, and the aggregate had to be rebuilt as
// an electra value, which is what made the signing root wrong. Each assignment below is therefore
// itself the assertion; the accessor checks confirm the wrappers resolve their Gloas arms.
func TestVersionedSignedAggregateAndProofGloasFetchToSubmit(t *testing.T) {
	_, aggregate := equivalentAggregates()

	// Fetch: the shape AggregateAttestation returns for a Gloas node.
	fetched := &spec.VersionedAttestation{
		Version: spec.DataVersionGloas,
		Gloas:   aggregate.Aggregate,
	}

	// The fetched wrapper resolves its Gloas arm through the versioned accessors.
	fetchedBits, err := fetched.AggregationBits()
	require.NoError(t, err)
	require.Equal(t, aggregationBits(), fetchedBits)

	// Aggregate: the fetched attestation is carried straight through, not converted.
	unsigned := &spec.VersionedAggregateAndProof{
		Version: spec.DataVersionGloas,
		Gloas: &gloas.AggregateAndProof{
			AggregatorIndex: 7,
			Aggregate:       fetched.Gloas,
			SelectionProof:  phase0.BLSSignature{0xbb},
		},
	}

	// Sign: the root the validator signs comes from the unsigned wrapper.
	_, err = unsigned.HashTreeRoot()
	require.NoError(t, err)

	// Submit: the signed aggregate keeps the same message.
	signed := &spec.VersionedSignedAggregateAndProof{
		Version: spec.DataVersionGloas,
		Gloas: &gloas.SignedAggregateAndProof{
			Message:   unsigned.Gloas,
			Signature: phase0.BLSSignature{0xcc},
		},
	}

	// The signed wrapper reports the values that came off the wire, through its accessors.
	aggregatorIndex, err := signed.AggregatorIndex()
	require.NoError(t, err)
	require.Equal(t, phase0.ValidatorIndex(7), aggregatorIndex)

	slot, err := signed.Slot()
	require.NoError(t, err)
	require.Equal(t, phase0.Slot(42), slot)
}

// TestVersionedSignedAggregateAndProofGloasAccessors covers every accessor's DataVersionGloas arm,
// asserting each reads the gloas container rather than erroring or returning a zero value, and that
// each reports a missing arm rather than panicking on a nil pointer.
func TestVersionedSignedAggregateAndProofGloasAccessors(t *testing.T) {
	_, aggregate := equivalentAggregates()
	populated := &spec.VersionedSignedAggregateAndProof{
		Version: spec.DataVersionGloas,
		Gloas: &gloas.SignedAggregateAndProof{
			Message:   aggregate,
			Signature: phase0.BLSSignature{0xcc},
		},
	}
	empty := &spec.VersionedSignedAggregateAndProof{Version: spec.DataVersionGloas}

	t.Run("AggregatorIndex", func(t *testing.T) {
		got, err := populated.AggregatorIndex()
		require.NoError(t, err)
		require.Equal(t, phase0.ValidatorIndex(7), got)

		_, err = empty.AggregatorIndex()
		require.EqualError(t, err, "no gloas signed aggregate and proof")
	})

	t.Run("SelectionProof", func(t *testing.T) {
		got, err := populated.SelectionProof()
		require.NoError(t, err)
		require.Equal(t, phase0.BLSSignature{0xbb}, got)

		_, err = empty.SelectionProof()
		require.EqualError(t, err, "no gloas signed aggregate and proof")
	})

	t.Run("Signature", func(t *testing.T) {
		got, err := populated.Signature()
		require.NoError(t, err)
		require.Equal(t, phase0.BLSSignature{0xcc}, got)

		_, err = empty.Signature()
		require.EqualError(t, err, "no gloas signed aggregate and proof")
	})

	t.Run("Slot", func(t *testing.T) {
		got, err := populated.Slot()
		require.NoError(t, err)
		require.Equal(t, phase0.Slot(42), got)

		_, err = empty.Slot()
		require.EqualError(t, err, "no gloas signed aggregate and proof")
	})

	t.Run("String", func(t *testing.T) {
		// The gloas container's own String, reached through the arm, reports its fields.
		require.Contains(t, populated.String(), "aggregator_index")
		require.Empty(t, empty.String())
	})

	t.Run("IsEmpty", func(t *testing.T) {
		require.False(t, populated.IsEmpty())
		require.True(t, empty.IsEmpty())
	})
}
