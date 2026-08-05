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

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestCreateUnversionedAttestationsGloas verifies that a Gloas-versioned attestation is
// converted to the single-attestation submission form via gloas.Attestation.ToSingleAttestation,
// mirroring the Electra/Fulu arms; without the arm the switch hits its default and errors.
func TestCreateUnversionedAttestationsGloas(t *testing.T) {
	committeeBits := bitfield.NewBitvector64()
	committeeBits.SetBitAt(2, true)
	aggregationBits := bitfield.NewBitlist(8)
	aggregationBits.SetBitAt(1, true)

	attestationData := &phase0.AttestationData{
		Slot:   9,
		Index:  0,
		Source: &phase0.Checkpoint{},
		Target: &phase0.Checkpoint{},
	}
	attestation := &gloas.Attestation{
		AggregationBits: aggregationBits,
		Data:            attestationData,
		Signature:       phase0.BLSSignature{0x0a},
		CommitteeBits:   committeeBits,
	}
	validatorIndex := phase0.ValidatorIndex(7)
	attestations := []*spec.VersionedAttestation{
		{Version: spec.DataVersionGloas, Gloas: attestation, ValidatorIndex: &validatorIndex},
	}

	s := &Service{}
	unversioned, err := s.createUnversionedAttestations(attestations)
	require.NoError(t, err)
	require.Len(t, unversioned, 1)

	single, ok := unversioned[0].(*electra.SingleAttestation)
	require.True(t, ok)
	require.Equal(t, phase0.CommitteeIndex(2), single.CommitteeIndex)
	require.Equal(t, validatorIndex, single.AttesterIndex)
	require.Equal(t, attestationData, single.Data)
}
