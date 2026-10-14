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
	"encoding/json"
	"testing"

	"github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestVersionedAttestationMarshalJSON(t *testing.T) {
	data := &phase0.AttestationData{Source: &phase0.Checkpoint{}, Target: &phase0.Checkpoint{}}
	p0 := &phase0.Attestation{AggregationBits: bitfield.NewBitlist(8), Data: data}
	el := &electra.Attestation{AggregationBits: bitfield.NewBitlist(8), Data: data, CommitteeBits: bitfield.NewBitvector64()}
	gl := &gloas.Attestation{AggregationBits: bitfield.NewBitlist(8), Data: data, CommitteeBits: bitfield.NewBitvector64()}

	index := phase0.ValidatorIndex(7)
	tests := []struct {
		name        string
		attestation *spec.VersionedAttestation
		arm         any
		err         string
	}{
		{name: "Phase0", attestation: &spec.VersionedAttestation{Version: spec.DataVersionPhase0, ValidatorIndex: &index, Phase0: p0}, arm: p0},
		{name: "Altair", attestation: &spec.VersionedAttestation{Version: spec.DataVersionAltair, Altair: p0}, arm: p0},
		{name: "Bellatrix", attestation: &spec.VersionedAttestation{Version: spec.DataVersionBellatrix, Bellatrix: p0}, arm: p0},
		{name: "Capella", attestation: &spec.VersionedAttestation{Version: spec.DataVersionCapella, Capella: p0}, arm: p0},
		{name: "Deneb", attestation: &spec.VersionedAttestation{Version: spec.DataVersionDeneb, Deneb: p0}, arm: p0},
		{name: "Electra", attestation: &spec.VersionedAttestation{Version: spec.DataVersionElectra, Electra: el}, arm: el},
		{name: "Fulu", attestation: &spec.VersionedAttestation{Version: spec.DataVersionFulu, Fulu: el}, arm: el},
		{name: "Gloas", attestation: &spec.VersionedAttestation{Version: spec.DataVersionGloas, Gloas: gl}, arm: gl},
		{name: "NilArm", attestation: &spec.VersionedAttestation{Version: spec.DataVersionDeneb}, err: "no attestation"},
		{name: "Unknown", attestation: &spec.VersionedAttestation{}, err: "unknown version"},
		{name: "Unsupported", attestation: &spec.VersionedAttestation{Version: spec.DataVersion(1000)}, err: "unsupported version"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := json.Marshal(test.attestation)
			if test.err != "" {
				require.ErrorContains(t, err, test.err)

				return
			}
			require.NoError(t, err)

			// The output is the populated fork's attestation, with no wrapper or
			// per-fork keys, and not the ValidatorIndex, which is not part of it.
			expected, err := json.Marshal(test.arm)
			require.NoError(t, err)
			require.JSONEq(t, string(expected), string(res))
		})
	}
}

// Decoding infers only Phase0 or Electra from the attestation itself, as the wire
// object carries no fork, so Fulu and Gloas attestations decode as Electra.
func TestVersionedAttestationJSONForkInference(t *testing.T) {
	data := &phase0.AttestationData{Source: &phase0.Checkpoint{}, Target: &phase0.Checkpoint{}}
	gl := &gloas.Attestation{AggregationBits: bitfield.NewBitlist(8), Data: data, CommitteeBits: bitfield.NewBitvector64()}

	res, err := json.Marshal(&spec.VersionedAttestation{Version: spec.DataVersionGloas, Gloas: gl})
	require.NoError(t, err)

	var decoded spec.VersionedAttestation
	require.NoError(t, json.Unmarshal(res, &decoded))
	require.Equal(t, spec.DataVersionElectra, decoded.Version)
}
