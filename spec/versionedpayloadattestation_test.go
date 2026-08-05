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

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestVersionedPayloadAttestationAccessors(t *testing.T) {
	bits := bitfield.NewBitvector512()
	bits.SetBitAt(3, true)

	populated := &spec.VersionedPayloadAttestation{
		Version: spec.DataVersionGloas,
		Gloas: &gloas.PayloadAttestation{
			AggregationBits: bits,
			Data:            testPayloadAttestationData(),
			Signature:       phase0.BLSSignature{0xaa},
		},
	}
	nilArm := &spec.VersionedPayloadAttestation{Version: spec.DataVersionGloas}
	preGloas := &spec.VersionedPayloadAttestation{Version: spec.DataVersionFulu}
	unknown := &spec.VersionedPayloadAttestation{Version: spec.DataVersion(99)}

	t.Run("Data", func(t *testing.T) {
		data, err := populated.Data()
		require.NoError(t, err)
		require.Equal(t, phase0.Slot(60300), data.Slot)

		_, err = nilArm.Data()
		require.EqualError(t, err, "no gloas payload attestation")

		_, err = preGloas.Data()
		require.EqualError(t, err, "no payload attestation in fulu")

		_, err = unknown.Data()
		require.EqualError(t, err, "unknown version")
	})

	t.Run("AggregationBits", func(t *testing.T) {
		aggregationBits, err := populated.AggregationBits()
		require.NoError(t, err)
		require.Equal(t, []int{3}, aggregationBits.BitIndices())

		// Read participation through the container's accessor, not the bitvector's
		// mainnet-only BitAt, which this mainnet-width fixture would hide.
		attested, err := populated.Gloas.AttestedAt(3)
		require.NoError(t, err)
		require.True(t, attested)

		attested, err = populated.Gloas.AttestedAt(4)
		require.NoError(t, err)
		require.False(t, attested)

		_, err = nilArm.AggregationBits()
		require.EqualError(t, err, "no gloas payload attestation")
	})

	t.Run("Signature", func(t *testing.T) {
		sig, err := populated.Signature()
		require.NoError(t, err)
		require.Equal(t, phase0.BLSSignature{0xaa}, sig)

		_, err = nilArm.Signature()
		require.EqualError(t, err, "no gloas payload attestation")
	})

	t.Run("IsEmpty", func(t *testing.T) {
		require.False(t, populated.IsEmpty())
		require.True(t, nilArm.IsEmpty())
	})

	t.Run("String", func(t *testing.T) {
		require.Contains(t, populated.String(), "60300")
		require.Empty(t, nilArm.String())
		require.Equal(t, "unknown version", unknown.String())
	})

	// A populated attestation whose inner Data is nil must not hand back a
	// usable-looking zero value: the aggregate carries signatures over that
	// datum, so an empty one silently misrepresents what was signed.
	t.Run("NilData", func(t *testing.T) {
		partial := &spec.VersionedPayloadAttestation{
			Version: spec.DataVersionGloas,
			Gloas:   &gloas.PayloadAttestation{AggregationBits: bits},
		}

		_, err := partial.Data()
		require.EqualError(t, err, "no gloas payload attestation data")
	})
}
