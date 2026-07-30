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

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
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

// TestSubmitProposalGloasMarshaling covers the gloas arm of both request-body
// encoders.  Without an arm each falls through to "unknown proposal version", so
// the block never leaves the process and no server-side test can tell why.
func TestSubmitProposalGloasMarshaling(t *testing.T) {
	ctx := context.Background()
	s := &Service{}

	t.Run("JSON", func(t *testing.T) {
		body, err := s.submitProposalJSON(ctx, gloasSignedProposal())
		require.NoError(t, err)

		// The body is the signed block itself, unwrapped, so it must carry the
		// block's own keys rather than a contents envelope's.
		var decoded map[string]json.RawMessage
		require.NoError(t, json.Unmarshal(body, &decoded))
		require.Contains(t, decoded, "message")
		require.Contains(t, decoded, "signature")
		require.NotContains(t, decoded, "signed_block")
	})

	t.Run("SSZ", func(t *testing.T) {
		body, err := s.submitProposalSSZ(ctx, gloasSignedProposal())
		require.NoError(t, err)
		require.NotEmpty(t, body)

		// Round-trips back to the same block, which is what tells us the arm
		// marshaled the right container.
		var decoded gloas.SignedBeaconBlock
		require.NoError(t, decoded.UnmarshalSSZ(body))
		require.Equal(t, phase0.Slot(60300), decoded.Message.Slot)
	})

	// A version claiming gloas with no gloas arm populated must be refused
	// before anything is marshaled: the encoders switch on Version, so an
	// unchecked nil arm marshals as JSON null or panics on the SSZ path.
	t.Run("MissingArm", func(t *testing.T) {
		proposal := &api.VersionedSignedProposal{Version: spec.DataVersionGloas}

		require.ErrorContains(t, proposal.AssertPresent(), "gloas proposal not present")

		_, err := s.submitProposalJSON(ctx, proposal)
		require.ErrorContains(t, err, "gloas proposal not present")

		_, err = s.submitProposalSSZ(ctx, proposal)
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

		_, err := s.submitProposalJSON(ctx, proposal)
		require.ErrorContains(t, err, "gloas proposals are never blinded")

		_, err = s.submitProposalSSZ(ctx, proposal)
		require.ErrorContains(t, err, "gloas proposals are never blinded")
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
