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
	"encoding/json"
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/api"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"
	apiv1fulu "github.com/attestantio/go-eth2-client/api/v1/fulu"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// gloasSignedProposal wraps a block in the shape the publish endpoint takes for
// gloas.
//
// Note this is a plain SignedBeaconBlock, not a contents wrapper as Deneb
// through Fulu use.  Post-Gloas the blobs travel in the execution payload
// envelope, which is published separately, so there is nothing to bundle with
// the block.
func gloasSignedProposal() *api.VersionedSignedProposal {
	return &api.VersionedSignedProposal{
		Version: spec.DataVersionGloas,
		Gloas: &gloas.SignedBeaconBlock{
			Message:   validEPBSBeaconBlock(),
			Signature: phase0.BLSSignature{0x0a, 0x0b, 0x0c},
		},
	}
}

// TestSubmitProposalDataStaticCodecParity pins two things about the default
// service, the one that has not opted in to WithCustomSpecSupport and so encodes
// against the compiled-in mainnet preset.
//
// The first is that routing the body through the dynamic codec left the wire
// format alone: for a type whose sizes are all static the dynamic codec
// delegates to the very generated codec this used to call, so every arm must
// still produce the bytes that codec produces, byte for byte.  The second is
// that each version reaches for the right container, which nothing else checks —
// the arms are eight hand-written lines and a wrong one still compiles.
func TestSubmitProposalDataStaticCodecParity(t *testing.T) {
	ctx := context.Background()
	s := &Service{}

	// Each fixture sits at a slot unique to its fork, so an arm reaching for the
	// wrong container encodes a different valid block rather than a nil one --
	// two empty containers could encode identically and let a mis-wire pass.
	proposal := api.VersionedSignedProposal{
		Phase0:    &phase0.SignedBeaconBlock{Message: &phase0.BeaconBlock{Slot: 100}},
		Altair:    &altair.SignedBeaconBlock{Message: &altair.BeaconBlock{Slot: 200}},
		Bellatrix: &bellatrix.SignedBeaconBlock{Message: &bellatrix.BeaconBlock{Slot: 300}},
		Capella:   &capella.SignedBeaconBlock{Message: &capella.BeaconBlock{Slot: 400}},
		Deneb: &apiv1deneb.SignedBlockContents{
			SignedBlock: &deneb.SignedBeaconBlock{Message: &deneb.BeaconBlock{Slot: 500}},
		},
		Electra: &apiv1electra.SignedBlockContents{
			SignedBlock: &electra.SignedBeaconBlock{Message: &electra.BeaconBlock{Slot: 600}},
		},
		Fulu: &apiv1fulu.SignedBlockContents{
			SignedBlock: &electra.SignedBeaconBlock{Message: &electra.BeaconBlock{Slot: 700}},
		},
		Gloas: &gloas.SignedBeaconBlock{Message: &gloas.BeaconBlock{Slot: 800}},
	}

	tests := []struct {
		name    string
		version spec.DataVersion
		want    []byte
	}{
		{name: "Phase0", version: spec.DataVersionPhase0, want: mustMarshalSSZ(t, proposal.Phase0)},
		{name: "Altair", version: spec.DataVersionAltair, want: mustMarshalSSZ(t, proposal.Altair)},
		{name: "Bellatrix", version: spec.DataVersionBellatrix, want: mustMarshalSSZ(t, proposal.Bellatrix)},
		{name: "Capella", version: spec.DataVersionCapella, want: mustMarshalSSZ(t, proposal.Capella)},
		{name: "Deneb", version: spec.DataVersionDeneb, want: mustMarshalSSZ(t, proposal.Deneb)},
		{name: "Electra", version: spec.DataVersionElectra, want: mustMarshalSSZ(t, proposal.Electra)},
		{name: "Fulu", version: spec.DataVersionFulu, want: mustMarshalSSZ(t, proposal.Fulu)},
		{name: "Gloas", version: spec.DataVersionGloas, want: mustMarshalSSZ(t, proposal.Gloas)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			versioned := proposal
			versioned.Version = test.version

			body, contentType, err := s.submitProposalData(ctx, &versioned)
			require.NoError(t, err)
			require.Equal(t, ContentTypeSSZ, contentType)
			require.Equal(t, test.want, body)
		})
	}
}

// TestSubmitProposalDataAtTheNodesPreset covers the request body a service that
// has opted in to WithCustomSpecSupport produces.
//
// A block body holds a preset-sized sync committee bitvector from Altair
// onwards — 512 bits at the mainnet preset, 32 at minimal — and SSZ addresses
// the variable-size fields that follow it by offsets, so the whole body shifts
// by 60 bytes between the two.  Encoding with the generated codec, which has
// the mainnet preset baked in, therefore produces a body the node cannot read.
// It does so silently: the codec pads a short bitvector out to the mainnet
// length and returns no error, so the only place the mismatch is observable is
// at a decoder holding the node's own spec.  That is what this asserts against.
func TestSubmitProposalDataAtTheNodesPreset(t *testing.T) {
	ctx := context.Background()
	custom := customSpecService(ctx, t)

	// The bitvector is the preset-derived field that shifts the offsets, so the
	// fixture has to be sized from the node's spec rather than from the
	// compiled-in preset to be a block that node would accept at all.
	specResponse, err := custom.Spec(ctx, &api.SpecOpts{})
	require.NoError(t, err)
	syncCommitteeSize, isUint64 := specResponse.Data["SYNC_COMMITTEE_SIZE"].(uint64)
	require.True(t, isUint64, "SYNC_COMMITTEE_SIZE missing from spec")

	proposal := gloasSignedProposal()
	proposal.Gloas.Message.Body.SyncAggregate.SyncCommitteeBits = make(bitfield.Bitvector512, syncCommitteeSize/8)

	body, contentType, err := custom.submitProposalData(ctx, proposal)
	require.NoError(t, err)
	require.Equal(t, ContentTypeSSZ, contentType)

	// Decode with a codec built from the node's spec -- the same one the node
	// uses, so a body it rejects is a body the node rejects.
	ds, err := custom.dynSSZForRequest(ctx)
	require.NoError(t, err)

	var decoded gloas.SignedBeaconBlock
	require.NoError(t, ds.UnmarshalSSZ(&decoded, body))
	require.Equal(t, phase0.Slot(60300), decoded.Message.Slot)
	require.Len(t, decoded.Message.Body.SyncAggregate.SyncCommitteeBits, int(syncCommitteeSize/8))
}

// TestSubmitProposalGloasMarshaling covers the gloas arm of the request-body
// switch, on both content types.  Without that arm the version falls through to
// "unknown proposal version", so the block never leaves the process and no
// server-side test can tell why.
func TestSubmitProposalGloasMarshaling(t *testing.T) {
	ctx := context.Background()

	// Both services are zero-valued but for the flag under test, so neither
	// reaches a beacon node: the SSZ one has customSpecSupport false, which is
	// the branch of dynSSZForRequest that reads nothing from the service.
	s := &Service{}
	jsonService := &Service{enforceJSON: true}

	t.Run("JSON", func(t *testing.T) {
		body, contentType, err := jsonService.submitProposalData(ctx, gloasSignedProposal())
		require.NoError(t, err)
		require.Equal(t, ContentTypeJSON, contentType)

		// The body is the signed block itself, unwrapped, so it must carry the
		// block's own keys rather than a contents envelope's.
		var decoded map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(body, &decoded))
		require.Contains(t, decoded, "message")
		require.Contains(t, decoded, "signature")
		require.NotContains(t, decoded, "signed_block")
	})

	t.Run("SSZ", func(t *testing.T) {
		body, contentType, err := s.submitProposalData(ctx, gloasSignedProposal())
		require.NoError(t, err)
		require.Equal(t, ContentTypeSSZ, contentType)
		require.NotEmpty(t, body)

		// Round-trips back to the same block, which is what tells us the arm
		// marshaled the right container.
		var decoded gloas.SignedBeaconBlock
		require.NoError(t, decoded.UnmarshalSSZ(body))
		require.Equal(t, phase0.Slot(60300), decoded.Message.Slot)
	})

	// A version claiming gloas with no gloas arm populated must be refused
	// before anything is marshaled: the switch runs on Version alone, so an
	// unchecked nil arm marshals as JSON null, and as a zero-valued block on the
	// SSZ path, where the generated codecs treat a nil container as a zero one.
	//
	// Only the SSZ service is exercised here, and in the two subtests below: the
	// presence check runs ahead of the content-type branch, so both services
	// travel the same code to reach this error.
	t.Run("MissingArm", func(t *testing.T) {
		proposal := &api.VersionedSignedProposal{Version: spec.DataVersionGloas}

		require.ErrorContains(t, proposal.AssertPresent(), "gloas proposal not present")

		_, _, err := s.submitProposalData(ctx, proposal)
		require.ErrorContains(t, err, "gloas proposal not present")
	})

	// There is no blinded proposal post-Gloas, so a proposal marked blinded is a
	// caller that believes it is withholding a payload the block never carried.
	// Ignoring the flag would publish a plain block while the caller thought it
	// had published a blinded one, so it is refused instead.
	t.Run("MarkedBlinded", func(t *testing.T) {
		proposal := gloasSignedProposal()
		proposal.Blinded = true

		require.ErrorContains(t, proposal.AssertPresent(), "gloas proposals are never blinded")

		_, _, err := s.submitProposalData(ctx, proposal)
		require.ErrorContains(t, err, "gloas proposals are never blinded")
	})

	// An unrecognised version has to be refused before encoding as well, and it
	// is the presence check rather than the switch that does it.
	t.Run("UnknownVersion", func(t *testing.T) {
		proposal := &api.VersionedSignedProposal{Version: spec.DataVersion(99)}

		_, _, err := s.submitProposalData(ctx, proposal)
		require.ErrorContains(t, err, "unsupported version")
	})

	// The wrapper's accessors have to know about the arm as well as its
	// marshalers do.  Nothing in this repository calls them — the submitter
	// itself does not — so a missing arm compiles, passes, and surfaces only in a
	// consumer that inspects a proposal the same way it does for every other
	// fork, and then reports the fork as unsupported while holding a valid block.
	t.Run("Accessors", func(t *testing.T) {
		proposal := gloasSignedProposal()

		slot, err := proposal.Slot()
		require.NoError(t, err)
		require.Equal(t, phase0.Slot(60300), slot)

		proposerIndex, err := proposal.ProposerIndex()
		require.NoError(t, err)
		require.Equal(t, phase0.ValidatorIndex(403), proposerIndex)

		require.NotEmpty(t, proposal.String())
		require.Contains(t, proposal.String(), "60300")
	})

	// The accessors funnel through a presence check, so a version claiming gloas
	// with an empty arm must report the data as missing rather than the fork as
	// unsupported.
	t.Run("AccessorsWithMissingArm", func(t *testing.T) {
		proposal := &api.VersionedSignedProposal{Version: spec.DataVersionGloas}

		_, err := proposal.Slot()
		require.ErrorIs(t, err, api.ErrDataMissing)

		_, err = proposal.ProposerIndex()
		require.ErrorIs(t, err, api.ErrDataMissing)

		require.Empty(t, proposal.String())
	})
}
