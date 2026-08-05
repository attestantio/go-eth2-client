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

// testPayloadAttestationData is the datum a PTC validator signs: the block it
// is voting on, plus the two timeliness flags.
//
// The two flags are deliberately set apart.  Both are bool and both mean
// "something about this payload is fine", so an accessor wired to the other
// one's field is invisible to any test built on a fixture that sets them
// alike — and getting it wrong makes a validator vote the opposite of what it
// observed.
func testPayloadAttestationData() *gloas.PayloadAttestationData {
	return &gloas.PayloadAttestationData{
		BeaconBlockRoot:   phase0.Root{0x01, 0x02, 0x03},
		Slot:              60300,
		PayloadPresent:    true,
		BlobDataAvailable: false,
	}
}

func TestVersionedPayloadAttestationDataAccessors(t *testing.T) {
	populated := &spec.VersionedPayloadAttestationData{
		Version: spec.DataVersionGloas,
		Gloas:   testPayloadAttestationData(),
	}
	nilArm := &spec.VersionedPayloadAttestationData{Version: spec.DataVersionGloas}
	preGloas := &spec.VersionedPayloadAttestationData{Version: spec.DataVersionElectra}
	unknown := &spec.VersionedPayloadAttestationData{Version: spec.DataVersion(99)}

	t.Run("Slot", func(t *testing.T) {
		slot, err := populated.Slot()
		require.NoError(t, err)
		require.Equal(t, phase0.Slot(60300), slot)

		_, err = nilArm.Slot()
		require.EqualError(t, err, "no gloas payload attestation data")

		_, err = preGloas.Slot()
		require.EqualError(t, err, "no payload attestation data in electra")

		_, err = unknown.Slot()
		require.EqualError(t, err, "unknown version")
	})

	t.Run("BeaconBlockRoot", func(t *testing.T) {
		root, err := populated.BeaconBlockRoot()
		require.NoError(t, err)
		require.Equal(t, phase0.Root{0x01, 0x02, 0x03}, root)

		_, err = nilArm.BeaconBlockRoot()
		require.EqualError(t, err, "no gloas payload attestation data")

		_, err = preGloas.BeaconBlockRoot()
		require.EqualError(t, err, "no payload attestation data in electra")
	})

	t.Run("PayloadPresent", func(t *testing.T) {
		present, err := populated.PayloadPresent()
		require.NoError(t, err)
		require.True(t, present)

		_, err = nilArm.PayloadPresent()
		require.EqualError(t, err, "no gloas payload attestation data")
	})

	t.Run("BlobDataAvailable", func(t *testing.T) {
		available, err := populated.BlobDataAvailable()
		require.NoError(t, err)
		require.False(t, available)

		_, err = nilArm.BlobDataAvailable()
		require.EqualError(t, err, "no gloas payload attestation data")
	})

	t.Run("IsEmpty", func(t *testing.T) {
		require.False(t, populated.IsEmpty())
		require.True(t, nilArm.IsEmpty())
		require.True(t, preGloas.IsEmpty())
	})

	t.Run("String", func(t *testing.T) {
		require.Contains(t, populated.String(), "60300")
		require.Empty(t, nilArm.String())
		require.Empty(t, preGloas.String())
		require.Equal(t, "unknown version", unknown.String())
	})
}
