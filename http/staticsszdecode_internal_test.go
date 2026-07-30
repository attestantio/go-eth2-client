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

// Coverage for the static SSZ decode branch, i.e. the one a service takes when it has
// not opted in to WithCustomSpecSupport. That is the default, and so the path every
// mainnet-preset client actually runs, but no live test can reach it: the interim
// Glamsterdam devnet serves the minimal preset, whose bodies the generated
// mainnet-preset codecs reject outright (see ADR-0003). The bodies here are therefore
// marshalled in-process by the same generated codecs.
//
// Each test drives a decoder on a zero-value &Service{}, whose customSpecSupport is
// false: that branch reads nothing from the service, so no beacon node is involved.
//
// Note these tests are still gated by TestMain, which runs nothing unless HTTP_ADDRESS
// is set, even though they make no network calls of their own.

import (
	"context"
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"
	apiv1fulu "github.com/attestantio/go-eth2-client/api/v1/fulu"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/fulu"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// staticSlot is the slot every body below is built at. Any value that is not the zero
// value works; the point is that a skipped decode leaves 0 behind.
const staticSlot = phase0.Slot(11235813)

// mustMarshalSSZ marshals v with its own generated (static, mainnet-preset) codec.
func mustMarshalSSZ(t *testing.T, v interface{ MarshalSSZ() ([]byte, error) }) []byte {
	t.Helper()

	body, err := v.MarshalSSZ()
	require.NoError(t, err)

	return body
}

func TestBeaconStateFromSSZStaticCodec(t *testing.T) {
	ctx := context.Background()
	s := &Service{}

	tests := []struct {
		name    string
		version spec.DataVersion
		body    []byte
	}{
		{
			name:    "Phase0",
			version: spec.DataVersionPhase0,
			body:    mustMarshalSSZ(t, &phase0.BeaconState{Slot: staticSlot}),
		},
		{
			name:    "Altair",
			version: spec.DataVersionAltair,
			body:    mustMarshalSSZ(t, &altair.BeaconState{Slot: staticSlot}),
		},
		{
			name:    "Bellatrix",
			version: spec.DataVersionBellatrix,
			body:    mustMarshalSSZ(t, &bellatrix.BeaconState{Slot: staticSlot}),
		},
		{
			name:    "Capella",
			version: spec.DataVersionCapella,
			body:    mustMarshalSSZ(t, &capella.BeaconState{Slot: staticSlot}),
		},
		{
			name:    "Deneb",
			version: spec.DataVersionDeneb,
			body:    mustMarshalSSZ(t, &deneb.BeaconState{Slot: staticSlot}),
		},
		{
			name:    "Electra",
			version: spec.DataVersionElectra,
			body:    mustMarshalSSZ(t, &electra.BeaconState{Slot: staticSlot}),
		},
		{
			name:    "Fulu",
			version: spec.DataVersionFulu,
			body:    mustMarshalSSZ(t, &fulu.BeaconState{Slot: staticSlot}),
		},
		{
			name:    "Gloas",
			version: spec.DataVersionGloas,
			body:    mustMarshalSSZ(t, &gloas.BeaconState{Slot: staticSlot}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := &httpResponse{
				consensusVersion: test.version,
				body:             test.body,
			}

			response, err := s.beaconStateFromSSZ(ctx, res)
			require.NoError(t, err)
			require.Equal(t, test.version, response.Data.Version)

			// Slot() reads through the arm the version names, so it fails both when the
			// decode was skipped (slot 0) and when the body landed in another fork's arm.
			slot, err := response.Data.Slot()
			require.NoError(t, err)
			require.Equal(t, staticSlot, slot)
		})
	}
}

// staticGloasSignedBeaconBlock returns a gloas block at the slot these tests
// assert on.  Unlike the earlier forks' containers it has no encodable zero
// value — the sync aggregate and the execution payload bid are pointers the
// marshaler dereferences — so it reuses the fixture that already satisfies that,
// rather than restating it.
func staticGloasSignedBeaconBlock() *gloas.SignedBeaconBlock {
	block := validEPBSBeaconBlock()
	block.Slot = staticSlot

	return &gloas.SignedBeaconBlock{Message: block}
}

func TestSignedBeaconBlockFromSSZStaticCodec(t *testing.T) {
	ctx := context.Background()
	s := &Service{}

	tests := []struct {
		name    string
		version spec.DataVersion
		body    []byte
	}{
		{
			name:    "Phase0",
			version: spec.DataVersionPhase0,
			body:    mustMarshalSSZ(t, &phase0.SignedBeaconBlock{Message: &phase0.BeaconBlock{Slot: staticSlot}}),
		},
		{
			name:    "Altair",
			version: spec.DataVersionAltair,
			body:    mustMarshalSSZ(t, &altair.SignedBeaconBlock{Message: &altair.BeaconBlock{Slot: staticSlot}}),
		},
		{
			name:    "Bellatrix",
			version: spec.DataVersionBellatrix,
			body:    mustMarshalSSZ(t, &bellatrix.SignedBeaconBlock{Message: &bellatrix.BeaconBlock{Slot: staticSlot}}),
		},
		{
			name:    "Capella",
			version: spec.DataVersionCapella,
			body:    mustMarshalSSZ(t, &capella.SignedBeaconBlock{Message: &capella.BeaconBlock{Slot: staticSlot}}),
		},
		{
			name:    "Deneb",
			version: spec.DataVersionDeneb,
			body:    mustMarshalSSZ(t, &deneb.SignedBeaconBlock{Message: &deneb.BeaconBlock{Slot: staticSlot}}),
		},
		{
			name:    "Electra",
			version: spec.DataVersionElectra,
			body:    mustMarshalSSZ(t, &electra.SignedBeaconBlock{Message: &electra.BeaconBlock{Slot: staticSlot}}),
		},
		{
			// Fulu reuses electra's container, as VersionedSignedBeaconBlock.Fulu does.
			name:    "Fulu",
			version: spec.DataVersionFulu,
			body:    mustMarshalSSZ(t, &electra.SignedBeaconBlock{Message: &electra.BeaconBlock{Slot: staticSlot}}),
		},
		{
			// Gloas needs its body populated where the earlier forks do not: the
			// container's own marshaler dereferences the sync aggregate and the
			// execution payload bid, neither of which has a zero value it can
			// encode.
			name:    "Gloas",
			version: spec.DataVersionGloas,
			body:    mustMarshalSSZ(t, staticGloasSignedBeaconBlock()),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := &httpResponse{
				consensusVersion: test.version,
				body:             test.body,
			}

			response, err := s.signedBeaconBlockFromSSZ(ctx, res)
			require.NoError(t, err)
			require.Equal(t, test.version, response.Data.Version)

			slot, err := response.Data.Slot()
			require.NoError(t, err)
			require.Equal(t, staticSlot, slot)
		})
	}
}

func TestBeaconBlockProposalFromSSZStaticCodec(t *testing.T) {
	ctx := context.Background()
	s := &Service{}

	// From bellatrix onwards each version decodes a blinded or a full body, so both
	// sides of that inner branch need a case of their own.
	tests := []struct {
		name    string
		version spec.DataVersion
		blinded bool
		body    []byte
	}{
		{
			name:    "Phase0",
			version: spec.DataVersionPhase0,
			body:    mustMarshalSSZ(t, &phase0.BeaconBlock{Slot: staticSlot}),
		},
		{
			name:    "Altair",
			version: spec.DataVersionAltair,
			body:    mustMarshalSSZ(t, &altair.BeaconBlock{Slot: staticSlot}),
		},
		{
			name:    "Bellatrix",
			version: spec.DataVersionBellatrix,
			body:    mustMarshalSSZ(t, &bellatrix.BeaconBlock{Slot: staticSlot}),
		},
		{
			name:    "BellatrixBlinded",
			version: spec.DataVersionBellatrix,
			blinded: true,
			body:    mustMarshalSSZ(t, &apiv1bellatrix.BlindedBeaconBlock{Slot: staticSlot}),
		},
		{
			name:    "Capella",
			version: spec.DataVersionCapella,
			body:    mustMarshalSSZ(t, &capella.BeaconBlock{Slot: staticSlot}),
		},
		{
			name:    "CapellaBlinded",
			version: spec.DataVersionCapella,
			blinded: true,
			body:    mustMarshalSSZ(t, &apiv1capella.BlindedBeaconBlock{Slot: staticSlot}),
		},
		{
			name:    "Deneb",
			version: spec.DataVersionDeneb,
			body:    mustMarshalSSZ(t, &apiv1deneb.BlockContents{Block: &deneb.BeaconBlock{Slot: staticSlot}}),
		},
		{
			name:    "DenebBlinded",
			version: spec.DataVersionDeneb,
			blinded: true,
			body:    mustMarshalSSZ(t, &apiv1deneb.BlindedBeaconBlock{Slot: staticSlot}),
		},
		{
			name:    "Electra",
			version: spec.DataVersionElectra,
			body:    mustMarshalSSZ(t, &apiv1electra.BlockContents{Block: &electra.BeaconBlock{Slot: staticSlot}}),
		},
		{
			name:    "ElectraBlinded",
			version: spec.DataVersionElectra,
			blinded: true,
			body:    mustMarshalSSZ(t, &apiv1electra.BlindedBeaconBlock{Slot: staticSlot}),
		},
		{
			name:    "Fulu",
			version: spec.DataVersionFulu,
			body:    mustMarshalSSZ(t, &apiv1fulu.BlockContents{Block: &electra.BeaconBlock{Slot: staticSlot}}),
		},
		{
			// The fulu blinded arm decodes electra's blinded block, as the decoder does.
			name:    "FuluBlinded",
			version: spec.DataVersionFulu,
			blinded: true,
			body:    mustMarshalSSZ(t, &apiv1electra.BlindedBeaconBlock{Slot: staticSlot}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// The decoder reads blindedness from this header, via
			// populateProposalDataFromHeaders, before it picks an arm.
			headers := map[string]string{}
			if test.blinded {
				headers["Eth-Execution-Payload-Blinded"] = "true"
			}

			res := &httpResponse{
				consensusVersion: test.version,
				headers:          headers,
				body:             test.body,
			}

			response, err := s.beaconBlockProposalFromSSZ(ctx, res)
			require.NoError(t, err)
			require.Equal(t, test.version, response.Data.Version)
			require.Equal(t, test.blinded, response.Data.Blinded)

			slot, err := response.Data.Slot()
			require.NoError(t, err)
			require.Equal(t, staticSlot, slot)
		})
	}
}

// staticBlob returns a blob whose first and last bytes are set, so that a decoded blob
// can be told apart from the 128KiB of zeroes a freshly allocated one holds.
func staticBlob() *deneb.Blob {
	blob := new(deneb.Blob)
	blob[0] = 0xab
	blob[len(blob)-1] = 0xcd

	return blob
}

func TestBlobsFromSSZStaticCodec(t *testing.T) {
	ctx := context.Background()
	s := &Service{}

	t.Run("Blobs", func(t *testing.T) {
		blobs := apiv1.Blobs{staticBlob()}
		res := &httpResponse{
			consensusVersion: spec.DataVersionFulu,
			body:             mustMarshalSSZ(t, &blobs),
		}

		response, err := s.blobsFromSSZ(ctx, res)
		require.NoError(t, err)
		require.Len(t, response.Data, 1)
		require.Equal(t, byte(0xab), response.Data[0][0])
		require.Equal(t, byte(0xcd), response.Data[0][len(response.Data[0])-1])
	})

	// An empty body is a valid "no blobs for this block" response, and returns before
	// the decoder reaches either branch.
	t.Run("EmptyBody", func(t *testing.T) {
		res := &httpResponse{consensusVersion: spec.DataVersionFulu}

		response, err := s.blobsFromSSZ(ctx, res)
		require.NoError(t, err)
		require.Empty(t, response.Data)
	})
}

func TestBlobSidecarsFromSSZStaticCodec(t *testing.T) {
	ctx := context.Background()
	s := &Service{}

	t.Run("BlobSidecars", func(t *testing.T) {
		sidecars := &api.BlobSidecars{
			Sidecars: []*deneb.BlobSidecar{{Index: 3, Blob: *staticBlob()}},
		}
		res := &httpResponse{
			consensusVersion: spec.DataVersionDeneb,
			body:             mustMarshalSSZ(t, sidecars),
		}

		response, err := s.blobSidecarsFromSSZ(ctx, res)
		require.NoError(t, err)
		require.Len(t, response.Data, 1)
		require.Equal(t, deneb.BlobIndex(3), response.Data[0].Index)
		require.Equal(t, byte(0xab), response.Data[0].Blob[0])
	})

	t.Run("EmptyBody", func(t *testing.T) {
		res := &httpResponse{consensusVersion: spec.DataVersionDeneb}

		response, err := s.blobSidecarsFromSSZ(ctx, res)
		require.NoError(t, err)
		require.Empty(t, response.Data)
	})
}
