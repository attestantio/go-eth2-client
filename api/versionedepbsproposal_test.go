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

package api_test

import (
	"math/big"
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	require "github.com/stretchr/testify/require"
)

// TestVersionedEPBSProposalSlot verifies that the slot is read from whichever
// arm the execution_payload_included flag selects, and that every way of
// arriving without a block is an error rather than slot 0.  Reporting slot 0 for
// an absent block would let a caller sign a proposal for the wrong slot.
func TestVersionedEPBSProposalSlot(t *testing.T) {
	tests := []struct {
		name     string
		proposal *api.VersionedEPBSProposal
		expected phase0.Slot
		err      string
	}{
		{
			name: "PayloadExcluded",
			proposal: &api.VersionedEPBSProposal{
				Version: spec.DataVersionGloas,
				Gloas:   &gloas.BeaconBlock{Slot: 42},
			},
			expected: 42,
		},
		{
			name: "PayloadIncluded",
			proposal: &api.VersionedEPBSProposal{
				Version:                  spec.DataVersionGloas,
				ExecutionPayloadIncluded: true,
				GloasContents: &apiv1gloas.BlockContents{
					Block: &gloas.BeaconBlock{Slot: 43},
				},
			},
			expected: 43,
		},
		{
			name: "PayloadExcludedNilBlock",
			proposal: &api.VersionedEPBSProposal{
				Version: spec.DataVersionGloas,
			},
			err: "no gloas beacon block",
		},
		{
			name: "PayloadIncludedNilContents",
			proposal: &api.VersionedEPBSProposal{
				Version:                  spec.DataVersionGloas,
				ExecutionPayloadIncluded: true,
			},
			err: "no gloas block contents",
		},
		{
			// Contents present but its block absent: the flag says the payload
			// travelled with the block, so a missing block is still an error.
			name: "PayloadIncludedNilBlock",
			proposal: &api.VersionedEPBSProposal{
				Version:                  spec.DataVersionGloas,
				ExecutionPayloadIncluded: true,
				GloasContents:            &apiv1gloas.BlockContents{},
			},
			err: "no gloas beacon block",
		},
		{
			name: "PreGloas",
			proposal: &api.VersionedEPBSProposal{
				Version: spec.DataVersionFulu,
			},
			err: "no epbs proposal in fulu",
		},
		{
			name: "UnknownVersion",
			proposal: &api.VersionedEPBSProposal{
				Version: spec.DataVersion(99),
			},
			err: "unknown version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			slot, err := test.proposal.Slot()
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, slot)
			}
		})
	}
}

// TestVersionedEPBSProposalRandaoReveal verifies the RANDAO reveal is read
// through the same arm selection as the slot.  The block production endpoint
// compares this against the reveal it sent, so reading it from the wrong arm
// would make every proposal look inconsistent.
func TestVersionedEPBSProposalRandaoReveal(t *testing.T) {
	reveal := phase0.BLSSignature{0x01, 0x02}

	proposal := &api.VersionedEPBSProposal{
		Version:                  spec.DataVersionGloas,
		ExecutionPayloadIncluded: true,
		GloasContents: &apiv1gloas.BlockContents{
			Block: &gloas.BeaconBlock{
				Body: &gloas.BeaconBlockBody{RANDAOReveal: reveal},
			},
		},
	}

	got, err := proposal.RandaoReveal()
	require.NoError(t, err)
	require.Equal(t, reveal, got)

	// A block with no body cannot yield a reveal, and must not report the zero
	// signature as if it were one.
	_, err = (&api.VersionedEPBSProposal{
		Version: spec.DataVersionGloas,
		Gloas:   &gloas.BeaconBlock{},
	}).RandaoReveal()
	require.EqualError(t, err, "no gloas beacon block body")
}

// TestVersionedEPBSProposalContents verifies that the three publish-side
// accessors are reachable only when the execution payload actually travelled
// with the block.  When it did not, the caller must fetch the envelope from the
// producing node instead, so returning empty values here would silently produce
// an unpublishable proposal.
func TestVersionedEPBSProposalContents(t *testing.T) {
	envelope := &gloas.ExecutionPayloadEnvelope{BuilderIndex: 7}
	included := &api.VersionedEPBSProposal{
		Version:                  spec.DataVersionGloas,
		ExecutionPayloadIncluded: true,
		GloasContents: &apiv1gloas.BlockContents{
			Block:                    &gloas.BeaconBlock{Slot: 1},
			ExecutionPayloadEnvelope: envelope,
			KZGProofs:                []deneb.KZGProof{{0x01}},
			Blobs:                    []deneb.Blob{{}},
		},
	}

	excluded := &api.VersionedEPBSProposal{
		Version: spec.DataVersionGloas,
		Gloas:   &gloas.BeaconBlock{Slot: 1},
	}

	const excludedErr = "no block contents; the execution payload was not included"

	t.Run("EnvelopeIncluded", func(t *testing.T) {
		got, err := included.ExecutionPayloadEnvelope()
		require.NoError(t, err)
		require.Equal(t, envelope, got)
	})

	t.Run("EnvelopeExcluded", func(t *testing.T) {
		_, err := excluded.ExecutionPayloadEnvelope()
		require.EqualError(t, err, excludedErr)
	})

	t.Run("KZGProofsIncluded", func(t *testing.T) {
		got, err := included.KZGProofs()
		require.NoError(t, err)
		require.Equal(t, []deneb.KZGProof{{0x01}}, got)
	})

	t.Run("KZGProofsExcluded", func(t *testing.T) {
		_, err := excluded.KZGProofs()
		require.EqualError(t, err, excludedErr)
	})

	t.Run("BlobsIncluded", func(t *testing.T) {
		got, err := included.Blobs()
		require.NoError(t, err)
		// Assert the length rather than emptiness: require.Empty is satisfied by
		// a nil return, so an accessor that read nothing would pass.
		require.Len(t, got, 1)
	})

	t.Run("BlobsExcluded", func(t *testing.T) {
		_, err := excluded.Blobs()
		require.EqualError(t, err, excludedErr)
	})

	t.Run("EnvelopeIncludedButAbsent", func(t *testing.T) {
		_, err := (&api.VersionedEPBSProposal{
			Version:                  spec.DataVersionGloas,
			ExecutionPayloadIncluded: true,
			GloasContents:            &apiv1gloas.BlockContents{Block: &gloas.BeaconBlock{}},
		}).ExecutionPayloadEnvelope()
		require.EqualError(t, err, "no gloas execution payload envelope")
	})
}

// TestVersionedEPBSProposalValue verifies the two reward components are summed
// nil-safely.  Both are populated from response headers that a node may omit.
func TestVersionedEPBSProposalValue(t *testing.T) {
	tests := []struct {
		name      string
		consensus *big.Int
		execution *big.Int
		expected  *big.Int
	}{
		{name: "Both", consensus: big.NewInt(3), execution: big.NewInt(4), expected: big.NewInt(7)},
		{name: "NeitherSet", expected: big.NewInt(0)},
		{name: "ConsensusOnly", consensus: big.NewInt(5), expected: big.NewInt(5)},
		{name: "ExecutionOnly", execution: big.NewInt(6), expected: big.NewInt(6)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proposal := &api.VersionedEPBSProposal{
				ConsensusValue: test.consensus,
				ExecutionValue: test.execution,
			}

			require.Equal(t, test.expected, proposal.Value())
		})
	}
}

// TestVersionedEPBSProposalIsEmpty verifies emptiness tracks the arms rather
// than the flag: either arm populated means the proposal carries data.
func TestVersionedEPBSProposalIsEmpty(t *testing.T) {
	require.True(t, (&api.VersionedEPBSProposal{Version: spec.DataVersionGloas}).IsEmpty())
	require.False(t, (&api.VersionedEPBSProposal{
		Version: spec.DataVersionGloas,
		Gloas:   &gloas.BeaconBlock{},
	}).IsEmpty())
	require.False(t, (&api.VersionedEPBSProposal{
		Version:       spec.DataVersionGloas,
		GloasContents: &apiv1gloas.BlockContents{},
	}).IsEmpty())
}
