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
	"context"
	"testing"

	"github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestAttestationPoolFromJSONGloas(t *testing.T) {
	poolResponse, err := new(Service).attestationPoolFromJSON(context.Background(),
		&api.AttestationPoolOpts{},
		&httpResponse{
			consensusVersion: spec.DataVersionGloas,
			body:             []byte(`{"version":"gloas","data":[` + gloasAttestationJSON + `]}`),
		},
	)
	require.NoError(t, err)
	require.Len(t, poolResponse.Data, 1)
	require.Equal(t, spec.DataVersionGloas, poolResponse.Data[0].Version)
	require.NotNil(t, poolResponse.Data[0].Gloas)

	aggregate, _, err := decodeAggregateAttestation(&httpResponse{
		consensusVersion: spec.DataVersionGloas,
		body:             []byte(`{"version":"gloas","data":` + gloasAttestationJSON + `}`),
	})
	require.NoError(t, err)

	poolRoot, err := poolResponse.Data[0].HashTreeRoot()
	require.NoError(t, err)
	aggregateRoot, err := aggregate.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, aggregateRoot, poolRoot)
}

func TestVerifyAttestationPoolGloas(t *testing.T) {
	attestation := func(slot phase0.Slot, committeeIndex phase0.CommitteeIndex) *spec.VersionedAttestation {
		committeeBits := bitfield.NewBitvector64()
		committeeBits.SetBitAt(uint64(committeeIndex), true)

		return &spec.VersionedAttestation{
			Version: spec.DataVersionGloas,
			Gloas: &gloas.Attestation{
				Data:          &phase0.AttestationData{Slot: slot},
				CommitteeBits: committeeBits,
			},
		}
	}

	matchingSlot := phase0.Slot(84434)
	otherSlot := phase0.Slot(84435)
	matchingCommitteeIndex := phase0.CommitteeIndex(3)
	otherCommitteeIndex := phase0.CommitteeIndex(4)
	matchingAttestation := attestation(matchingSlot, matchingCommitteeIndex)

	tests := []struct {
		name string
		opts *api.AttestationPoolOpts
		err  string
	}{
		{
			name: "NoFilters",
			opts: &api.AttestationPoolOpts{},
		},
		{
			name: "SlotAndCommitteeIndexMatch",
			opts: &api.AttestationPoolOpts{
				Slot:           &matchingSlot,
				CommitteeIndex: &matchingCommitteeIndex,
			},
		},
		{
			name: "SlotDiffers",
			opts: &api.AttestationPoolOpts{
				Slot: &otherSlot,
			},
			err: "attestation data not for requested slot",
		},
		{
			name: "CommitteeIndexDiffers",
			opts: &api.AttestationPoolOpts{
				CommitteeIndex: &otherCommitteeIndex,
			},
			err: "attestation data not for requested committee index",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := verifyAttestationPool(test.opts, []*spec.VersionedAttestation{matchingAttestation})
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
