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
	"fmt"
	"math"
	"slices"
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

// The aggregation bits of a payload attestation are Bitvector[PTC_SIZE], so their
// width is preset-dependent: 512 bits (64 bytes) at mainnet, 16 bits (2 bytes) on
// the minimal preset the validation devnet runs.  Both widths below are valid,
// correctly-decoded values.
const (
	minimalPTCBytes = 2
	mainnetPTCBytes = 64
)

// TestPayloadAttestationAttestedAtPresets reads every position of a correctly-sized
// bitvector at both the mainnet and minimal widths, including the highest valid
// index, which is where a mainnet-fixed width check bites.
//
// The fixtures are hand-written bytes, never built with SetAttestedAt: deriving them
// from the code under test would make the assertion circular.  BitIndices and Count
// are the length-tolerant members of bitfield's own API, so they serve as an
// independent oracle that the byte and bit ordering here matches go-bitfield's.
func TestPayloadAttestationAttestedAtPresets(t *testing.T) {
	mainnetBits := make(bitfield.Bitvector512, mainnetPTCBytes)
	mainnetBits[0] = 0x09  // positions 0 and 3
	mainnetBits[63] = 0x80 // position 511, the highest valid at mainnet

	tests := []struct {
		name string
		bits bitfield.Bitvector512
		size uint64
		// attesting lists, ascending, the positions the fixture bytes set.
		attesting []int
	}{
		{
			name:      "MinimalPreset",
			bits:      bitfield.Bitvector512{0x09, 0x80}, // positions 0, 3 and 15
			size:      minimalPTCBytes * 8,
			attesting: []int{0, 3, 15},
		},
		{
			name:      "MainnetPreset",
			bits:      mainnetBits,
			size:      mainnetPTCBytes * 8,
			attesting: []int{0, 3, 511},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			attestation := &gloas.PayloadAttestation{AggregationBits: test.bits}

			for position := range test.size {
				attested, err := attestation.AttestedAt(position)
				require.NoError(t, err)
				require.Equal(t, slices.Contains(test.attesting, int(position)), attested,
					"position %d", position)
			}

			// Out of range is an error rather than a false: the whole defect being
			// fixed here is that a failed read is indistinguishable from an unset bit.
			_, err := attestation.AttestedAt(test.size)
			require.ErrorContains(t, err, "out of range")

			require.Equal(t, test.attesting, test.bits.BitIndices())
			require.Equal(t, uint64(len(test.attesting)), test.bits.Count())
		})
	}
}

// TestPayloadAttestationPTCSize covers the committee width, which bitfield's own
// Len() reports as the mainnet 512 whatever the bitvector actually holds.
func TestPayloadAttestationPTCSize(t *testing.T) {
	tests := []struct {
		name     string
		bits     bitfield.Bitvector512
		expected uint64
	}{
		{
			name:     "MinimalPreset",
			bits:     make(bitfield.Bitvector512, minimalPTCBytes),
			expected: 16,
		},
		{
			name:     "MainnetPreset",
			bits:     make(bitfield.Bitvector512, mainnetPTCBytes),
			expected: 512,
		},
		{
			// An absent bitvector covers no committee at all.  Reporting 512 is what
			// makes a zero-valued field look like a full mainnet committee.
			name:     "Absent",
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			attestation := &gloas.PayloadAttestation{AggregationBits: test.bits}
			require.Equal(t, test.expected, attestation.PTCSize())
		})
	}

	// The premise this whole file exists for, pinned deliberately once: bitfield's own
	// Len() is the mainnet constant at every width, including a nil bitvector covering
	// no committee at all.  Should this ever fail, go-bitfield has grown preset-aware
	// accessors and these accessors are worth revisiting.
	require.Equal(t, uint64(512), bitfield.Bitvector512(nil).Len())
	require.Equal(t, uint64(512), bitfield.Bitvector512{0x00, 0x00}.Len())
}

// TestPayloadAttestationSetAttestedAt covers the write direction at both presets.
// bitfield offers no length-tolerant setter at all, so without this accessor
// building a minimal-preset aggregate is not expressible against this type.
func TestPayloadAttestationSetAttestedAt(t *testing.T) {
	tests := []struct {
		name  string
		bytes int
		size  uint64
	}{
		{name: "MinimalPreset", bytes: minimalPTCBytes, size: minimalPTCBytes * 8},
		{name: "MainnetPreset", bytes: mainnetPTCBytes, size: mainnetPTCBytes * 8},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			attestation := &gloas.PayloadAttestation{
				AggregationBits: make(bitfield.Bitvector512, test.bytes),
			}

			// Aggregate two members, one of them the highest valid position.
			require.NoError(t, attestation.SetAttestedAt(0, true))
			require.NoError(t, attestation.SetAttestedAt(test.size-1, true))

			// The length-tolerant family is what every other reader of these bits
			// uses, so its agreement is what proves the write landed in the right
			// place rather than merely in a place this file's reader agrees with.
			require.Equal(t, []int{0, int(test.size - 1)}, attestation.AggregationBits.BitIndices())
			require.Equal(t, uint64(2), attestation.AggregationBits.Count())

			// Hand-computed bytes: position 0 is bit 0 of the first byte, and the
			// highest position is bit 7 of the last.
			require.Equal(t, byte(0x01), attestation.AggregationBits[0])
			require.Equal(t, byte(0x80), attestation.AggregationBits[test.bytes-1])

			attested, err := attestation.AttestedAt(test.size - 1)
			require.NoError(t, err)
			require.True(t, attested)

			// Clearing must clear, and must leave its neighbours alone.
			require.NoError(t, attestation.SetAttestedAt(test.size-1, false))
			attested, err = attestation.AttestedAt(test.size - 1)
			require.NoError(t, err)
			require.False(t, attested)
			require.Equal(t, uint64(1), attestation.AggregationBits.Count())

			// Out of range must error rather than no-op: a silently dropped write
			// is how SetBitAt loses an aggregate.
			require.ErrorContains(t, attestation.SetAttestedAt(test.size, true), "out of range")
			require.Equal(t, uint64(1), attestation.AggregationBits.Count())
		})
	}

	// The exact message, pinned once.
	attestation := &gloas.PayloadAttestation{
		AggregationBits: make(bitfield.Bitvector512, minimalPTCBytes),
	}
	require.EqualError(t, attestation.SetAttestedAt(16, true),
		"payload timeliness committee index 16 out of range for committee of 16")
}

// TestPayloadAttestationZeroValueAccessors covers the zero-valued container, whose
// bitvector is nil.  Indexing one without a bound check panics, which is a failure
// mode this package has recorded before for a hand-written accessor.
func TestPayloadAttestationZeroValueAccessors(t *testing.T) {
	attestation := &gloas.PayloadAttestation{}

	require.Zero(t, attestation.PTCSize())

	_, err := attestation.AttestedAt(0)
	require.EqualError(t, err, "payload timeliness committee index 0 out of range for committee of 0")

	require.EqualError(t, attestation.SetAttestedAt(0, true),
		"payload timeliness committee index 0 out of range for committee of 0")
}

// TestPayloadAttestationArbitraryWidths pins the property the accessors' safety rests
// on: the width of AggregationBits is whatever a beacon node sent, since nothing on
// the JSON path checks it against the chain's PTC_SIZE.  So every width must be
// handled without panicking -- widths no preset produces, and one far larger than
// mainnet, included.
func TestPayloadAttestationArbitraryWidths(t *testing.T) {
	for _, width := range []int{0, 1, 2, 3, 7, 63, 64, 65, 1000} {
		t.Run(fmt.Sprintf("%dBytes", width), func(t *testing.T) {
			attestation := &gloas.PayloadAttestation{
				AggregationBits: make(bitfield.Bitvector512, width),
			}
			size := uint64(width) * 8
			require.Equal(t, size, attestation.PTCSize())

			// Every position the reported size admits round-trips, and the first
			// position beyond it is refused in both directions.
			for position := range size {
				require.NoError(t, attestation.SetAttestedAt(position, true))
				attested, err := attestation.AttestedAt(position)
				require.NoError(t, err)
				require.True(t, attested)
			}
			require.Equal(t, size, attestation.PTCSize())

			_, err := attestation.AttestedAt(size)
			require.ErrorContains(t, err, "out of range")
			require.ErrorContains(t, attestation.SetAttestedAt(size, true), "out of range")

			// The largest index representable cannot be made to index anything.
			_, err = attestation.AttestedAt(math.MaxUint64)
			require.ErrorContains(t, err, "out of range")
			require.ErrorContains(t, attestation.SetAttestedAt(math.MaxUint64, true), "out of range")
		})
	}
}
