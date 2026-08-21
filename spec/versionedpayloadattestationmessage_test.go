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

func TestVersionedPayloadAttestationMessageAccessors(t *testing.T) {
	populated := &spec.VersionedPayloadAttestationMessage{
		Version: spec.DataVersionGloas,
		Gloas: &gloas.PayloadAttestationMessage{
			ValidatorIndex: 8,
			Data:           testPayloadAttestationData(),
			Signature:      phase0.BLSSignature{0xbb},
		},
	}
	nilArm := &spec.VersionedPayloadAttestationMessage{Version: spec.DataVersionGloas}
	preGloas := &spec.VersionedPayloadAttestationMessage{Version: spec.DataVersionDeneb}
	unknown := &spec.VersionedPayloadAttestationMessage{Version: spec.DataVersion(99)}

	t.Run("ValidatorIndex", func(t *testing.T) {
		index, err := populated.ValidatorIndex()
		require.NoError(t, err)
		require.Equal(t, phase0.ValidatorIndex(8), index)

		_, err = nilArm.ValidatorIndex()
		require.EqualError(t, err, "no gloas payload attestation message")

		_, err = preGloas.ValidatorIndex()
		require.EqualError(t, err, "no payload attestation message in deneb")

		_, err = unknown.ValidatorIndex()
		require.EqualError(t, err, "unknown version")
	})

	t.Run("Data", func(t *testing.T) {
		data, err := populated.Data()
		require.NoError(t, err)
		require.Equal(t, phase0.Slot(60300), data.Slot)

		_, err = nilArm.Data()
		require.EqualError(t, err, "no gloas payload attestation message")
	})

	t.Run("Signature", func(t *testing.T) {
		sig, err := populated.Signature()
		require.NoError(t, err)
		require.Equal(t, phase0.BLSSignature{0xbb}, sig)

		_, err = nilArm.Signature()
		require.EqualError(t, err, "no gloas payload attestation message")
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

	// The message is what gets submitted to the node; a nil Data would be
	// serialised as a null datum rather than rejected.
	t.Run("NilData", func(t *testing.T) {
		partial := &spec.VersionedPayloadAttestationMessage{
			Version: spec.DataVersionGloas,
			Gloas:   &gloas.PayloadAttestationMessage{ValidatorIndex: 8},
		}

		_, err := partial.Data()
		require.EqualError(t, err, "no gloas payload attestation data")
	})
}
