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
	"fmt"
	"testing"

	"github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestVersionedAggregateAndProofHashTreeRootGloas asserts that the Gloas arm's hash tree root is
// computed via the progressive-bitlist merkleization, by showing it differs from the root the
// electra container produces for the same logical aggregate.
//
// Under Gloas (consensus-spec v1.7.0-alpha.12, specs/gloas/beacon-chain.md) Attestation is a
// ProgressiveContainer whose aggregation_bits are a ProgressiveBitlist; electra's are a regular
// Bitlist, so the two merkleize to different hash tree roots.
func TestVersionedAggregateAndProofHashTreeRootGloas(t *testing.T) {
	electraAggregate, gloasAggregate := equivalentAggregates()

	gloasRoot, err := (&spec.VersionedAggregateAndProof{
		Version: spec.DataVersionGloas,
		Gloas:   gloasAggregate,
	}).HashTreeRoot()
	require.NoError(t, err)

	electraRoot, err := (&spec.VersionedAggregateAndProof{
		Version: spec.DataVersionElectra,
		Electra: electraAggregate,
	}).HashTreeRoot()
	require.NoError(t, err)

	require.NotEqual(t, electraRoot, gloasRoot,
		"the Gloas root must differ from electra's: aggregation bits merkleize progressively")

	// Pin the root itself, so that a change to the progressive merkleization is caught rather
	// than merely producing a different-but-still-not-electra value. This value is generated from
	// this repository's own codec, so it guards against drift; it is not independent verification
	// against the spec, which needs the consensus-spec ssz_static vectors.
	require.Equal(t, "0xc565330dfc3b575b20935a7f217a64119b90f86c0e987b023bbf6d94648334ce",
		fmt.Sprintf("%#x", gloasRoot),
		"the Gloas aggregate root changed: check go generate ./spec/gloas and the progressive merkleization")
}

// TestVersionedAggregateAndProofGloasAccessors covers every accessor's DataVersionGloas arm on the
// unsigned wrapper, asserting each reads the gloas container and reports a missing arm rather than
// panicking on a nil pointer.
//
// This wrapper has no callers inside this library (it exists for downstream signers, which is why
// its mistyped Gloas arm went unnoticed), so these are the only tests exercising its Gloas path.
func TestVersionedAggregateAndProofGloasAccessors(t *testing.T) {
	_, aggregate := equivalentAggregates()
	populated := &spec.VersionedAggregateAndProof{
		Version: spec.DataVersionGloas,
		Gloas:   aggregate,
	}
	empty := &spec.VersionedAggregateAndProof{Version: spec.DataVersionGloas}

	t.Run("AggregatorIndex", func(t *testing.T) {
		got, err := populated.AggregatorIndex()
		require.NoError(t, err)
		require.Equal(t, phase0.ValidatorIndex(7), got)

		_, err = empty.AggregatorIndex()
		require.EqualError(t, err, "no gloas aggregate and proof")
	})

	t.Run("SelectionProof", func(t *testing.T) {
		got, err := populated.SelectionProof()
		require.NoError(t, err)
		require.Equal(t, phase0.BLSSignature{0xbb}, got)

		_, err = empty.SelectionProof()
		require.EqualError(t, err, "no gloas aggregate and proof")
	})

	t.Run("HashTreeRoot", func(t *testing.T) {
		_, err := empty.HashTreeRoot()
		require.EqualError(t, err, "no gloas aggregate and proof")
	})

	t.Run("String", func(t *testing.T) {
		require.Contains(t, populated.String(), "aggregator_index")
		require.Empty(t, empty.String())
	})

	t.Run("IsEmpty", func(t *testing.T) {
		require.False(t, populated.IsEmpty())
		require.True(t, empty.IsEmpty())
	})
}

// equivalentAggregates returns an electra and a gloas aggregate and proof carrying identical
// logical content, so that any difference in their hash tree roots is attributable to
// merkleization alone.
//
// It, along with the aggregationBits helper below, is declared here and shared with
// versionedsignedaggregateandproof_test.go.
func equivalentAggregates() (*electra.AggregateAndProof, *gloas.AggregateAndProof) {
	data := &phase0.AttestationData{
		Slot:            42,
		Index:           1,
		BeaconBlockRoot: phase0.Root{0x01},
		Source:          &phase0.Checkpoint{Epoch: 1, Root: phase0.Root{0x02}},
		Target:          &phase0.Checkpoint{Epoch: 2, Root: phase0.Root{0x03}},
	}

	signature := phase0.BLSSignature{0xaa}
	selectionProof := phase0.BLSSignature{0xbb}

	return &electra.AggregateAndProof{
			AggregatorIndex: 7,
			Aggregate: &electra.Attestation{
				AggregationBits: aggregationBits(),
				Data:            data,
				Signature:       signature,
				CommitteeBits:   committeeBits(),
			},
			SelectionProof: selectionProof,
		}, &gloas.AggregateAndProof{
			AggregatorIndex: 7,
			Aggregate: &gloas.Attestation{
				AggregationBits: aggregationBits(),
				Data:            data,
				Signature:       signature,
				CommitteeBits:   committeeBits(),
			},
			SelectionProof: selectionProof,
		}
}

// aggregationBits returns a bitlist with two bits set, short enough that the progressive and
// regular merkleizations of it are both exercised on a single chunk.
func aggregationBits() bitfield.Bitlist {
	bits := bitfield.NewBitlist(8)
	bits.SetBitAt(0, true)
	bits.SetBitAt(3, true)

	return bits
}

// committeeBits returns a committee bitvector with a single committee selected.
func committeeBits() bitfield.Bitvector64 {
	bits := bitfield.NewBitvector64()
	bits.SetBitAt(0, true)

	return bits
}
