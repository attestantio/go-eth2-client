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
			err: "unsupported version",
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

// TestVersionedEPBSProposalProposerIndex verifies the proposer index is read
// from whichever arm the execution_payload_included flag selects.  The
// proposer signs the proposal together with this index, so returning the zero
// value for an absent block would let a caller sign with the wrong validator.
func TestVersionedEPBSProposalProposerIndex(t *testing.T) {
	tests := []struct {
		name     string
		proposal *api.VersionedEPBSProposal
		expected phase0.ValidatorIndex
		err      string
	}{
		{
			name: "PayloadExcluded",
			proposal: &api.VersionedEPBSProposal{
				Version: spec.DataVersionGloas,
				Gloas:   &gloas.BeaconBlock{ProposerIndex: 42},
			},
			expected: 42,
		},
		{
			name: "PayloadIncluded",
			proposal: &api.VersionedEPBSProposal{
				Version:                  spec.DataVersionGloas,
				ExecutionPayloadIncluded: true,
				GloasContents: &apiv1gloas.BlockContents{
					Block: &gloas.BeaconBlock{ProposerIndex: 43},
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
			err: "unsupported version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			index, err := test.proposal.ProposerIndex()
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, index)
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

// TestVersionedEPBSProposalParentRoot verifies the parent root is read from
// whichever arm the execution_payload_included flag selects, and that every
// way of arriving without a block is an error.  Each fixture also carries a
// distinct state root, so a fixed accessor that read the wrong field would
// fail rather than pass by coincidence.
func TestVersionedEPBSProposalParentRoot(t *testing.T) {
	tests := []struct {
		name     string
		proposal *api.VersionedEPBSProposal
		expected phase0.Root
		err      string
	}{
		{
			name: "PayloadExcluded",
			proposal: &api.VersionedEPBSProposal{
				Version: spec.DataVersionGloas,
				Gloas:   &gloas.BeaconBlock{ParentRoot: phase0.Root{0x01}, StateRoot: phase0.Root{0xaa}},
			},
			expected: phase0.Root{0x01},
		},
		{
			name: "PayloadIncluded",
			proposal: &api.VersionedEPBSProposal{
				Version:                  spec.DataVersionGloas,
				ExecutionPayloadIncluded: true,
				GloasContents: &apiv1gloas.BlockContents{
					Block: &gloas.BeaconBlock{ParentRoot: phase0.Root{0x02}, StateRoot: phase0.Root{0xbb}},
				},
			},
			expected: phase0.Root{0x02},
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
			err: "unsupported version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := test.proposal.ParentRoot()
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, root)
			}
		})
	}
}

// TestVersionedEPBSProposalStateRoot mirrors TestVersionedEPBSProposalParentRoot
// for the state root, with distinct fixture values so a parent/state field
// swap in the accessor would fail rather than pass by coincidence.
func TestVersionedEPBSProposalStateRoot(t *testing.T) {
	tests := []struct {
		name     string
		proposal *api.VersionedEPBSProposal
		expected phase0.Root
		err      string
	}{
		{
			name: "PayloadExcluded",
			proposal: &api.VersionedEPBSProposal{
				Version: spec.DataVersionGloas,
				Gloas:   &gloas.BeaconBlock{StateRoot: phase0.Root{0x03}, ParentRoot: phase0.Root{0xcc}},
			},
			expected: phase0.Root{0x03},
		},
		{
			name: "PayloadIncluded",
			proposal: &api.VersionedEPBSProposal{
				Version:                  spec.DataVersionGloas,
				ExecutionPayloadIncluded: true,
				GloasContents: &apiv1gloas.BlockContents{
					Block: &gloas.BeaconBlock{StateRoot: phase0.Root{0x04}, ParentRoot: phase0.Root{0xdd}},
				},
			},
			expected: phase0.Root{0x04},
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
			err: "unsupported version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := test.proposal.StateRoot()
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, root)
			}
		})
	}
}

// TestVersionedEPBSProposalBodyRoot verifies BodyRoot returns the hash tree
// root of the block's body, not the block itself: dirk's proposer domain
// signs over the body root, so returning the block root instead would sign
// gibberish.  The happy-path cases assert the result matches an
// independently-computed body HTR and differs from the block's own HTR,
// rather than merely asserting no error.
func TestVersionedEPBSProposalBodyRoot(t *testing.T) {
	newBlock := func(reveal byte) *gloas.BeaconBlock {
		return &gloas.BeaconBlock{
			Slot: 1,
			Body: &gloas.BeaconBlockBody{RANDAOReveal: phase0.BLSSignature{reveal}},
		}
	}

	t.Run("PayloadExcluded", func(t *testing.T) {
		block := newBlock(0x01)
		wantBodyRoot, err := block.Body.HashTreeRoot()
		require.NoError(t, err)
		bodyRoot := phase0.Root(wantBodyRoot)
		blockRoot, err := block.HashTreeRoot()
		require.NoError(t, err)
		require.NotEqual(t, bodyRoot, phase0.Root(blockRoot))

		got, err := (&api.VersionedEPBSProposal{
			Version:             spec.DataVersionGloas,
			Gloas:               block,
			BeaconBlockBodyRoot: &bodyRoot,
		}).BodyRoot()
		require.NoError(t, err)
		require.NotEqual(t, phase0.Root{}, got)
		require.Equal(t, bodyRoot, got)
		require.NotEqual(t, phase0.Root(blockRoot), got)
	})

	t.Run("PayloadIncluded", func(t *testing.T) {
		block := newBlock(0x02)
		wantBodyRoot, err := block.Body.HashTreeRoot()
		require.NoError(t, err)
		bodyRoot := phase0.Root(wantBodyRoot)
		blockRoot, err := block.HashTreeRoot()
		require.NoError(t, err)
		require.NotEqual(t, bodyRoot, phase0.Root(blockRoot))

		got, err := (&api.VersionedEPBSProposal{
			Version:                  spec.DataVersionGloas,
			ExecutionPayloadIncluded: true,
			GloasContents:            &apiv1gloas.BlockContents{Block: block},
			BeaconBlockBodyRoot:      &bodyRoot,
		}).BodyRoot()
		require.NoError(t, err)
		require.NotEqual(t, phase0.Root{}, got)
		require.Equal(t, bodyRoot, got)
		require.NotEqual(t, phase0.Root(blockRoot), got)
	})

	// A block with no body cannot yield a body root, and must not report the
	// zero root as if it were one.
	t.Run("NilBody", func(t *testing.T) {
		_, err := (&api.VersionedEPBSProposal{
			Version: spec.DataVersionGloas,
			Gloas:   &gloas.BeaconBlock{},
		}).BodyRoot()
		require.EqualError(t, err, "no gloas beacon block body")
	})

	// Generated HashTreeRoot methods inline mainnet preset sizes, so a body
	// present without a retained BeaconBlockBodyRoot must not fall back to
	// computing one: that fallback is the exact bug this field exists to
	// avoid.
	t.Run("BeaconBlockBodyRootUnset", func(t *testing.T) {
		_, err := (&api.VersionedEPBSProposal{
			Version: spec.DataVersionGloas,
			Gloas:   newBlock(0x03),
		}).BodyRoot()
		require.EqualError(t, err, "no beacon block body root")
	})

	errTests := []struct {
		name     string
		proposal *api.VersionedEPBSProposal
		err      string
	}{
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
			err: "unsupported version",
		},
	}

	for _, test := range errTests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.proposal.BodyRoot()
			require.EqualError(t, err, test.err)
		})
	}
}

// TestVersionedEPBSProposalRoot verifies Root derives the block root by
// composing a phase0.BeaconBlockHeader from the block's own fields and the
// retained BeaconBlockBodyRoot, rather than by calling the block's own
// generated HashTreeRoot: that generated method recurses into the body and
// inherits the same mainnet-baked sizes that make Body.HashTreeRoot() wrong
// on other presets, so it could disagree with BodyRoot() on a custom preset.
// Each case's expected root is composed independently of the code under
// test, so a Root that mixed up a field would fail rather than pass by
// coincidence.
func TestVersionedEPBSProposalRoot(t *testing.T) {
	newBlock := func(reveal byte) *gloas.BeaconBlock {
		return &gloas.BeaconBlock{
			Slot:          10,
			ProposerIndex: 5,
			ParentRoot:    phase0.Root{0x11},
			StateRoot:     phase0.Root{0x22},
			Body:          &gloas.BeaconBlockBody{RANDAOReveal: phase0.BLSSignature{reveal}},
		}
	}

	wantRootFor := func(t *testing.T, block *gloas.BeaconBlock, bodyRoot phase0.Root) phase0.Root {
		t.Helper()

		root, err := (&phase0.BeaconBlockHeader{
			Slot:          block.Slot,
			ProposerIndex: block.ProposerIndex,
			ParentRoot:    block.ParentRoot,
			StateRoot:     block.StateRoot,
			BodyRoot:      bodyRoot,
		}).HashTreeRoot()
		require.NoError(t, err)

		return phase0.Root(root)
	}

	t.Run("PayloadExcluded", func(t *testing.T) {
		block := newBlock(0x07)
		htr, err := block.Body.HashTreeRoot()
		require.NoError(t, err)
		bodyRoot := phase0.Root(htr)

		got, err := (&api.VersionedEPBSProposal{
			Version:             spec.DataVersionGloas,
			Gloas:               block,
			BeaconBlockBodyRoot: &bodyRoot,
		}).Root()
		require.NoError(t, err)
		require.NotEqual(t, phase0.Root{}, got)
		require.Equal(t, wantRootFor(t, block, bodyRoot), got)
	})

	t.Run("PayloadIncluded", func(t *testing.T) {
		block := newBlock(0x08)
		htr, err := block.Body.HashTreeRoot()
		require.NoError(t, err)
		bodyRoot := phase0.Root(htr)

		got, err := (&api.VersionedEPBSProposal{
			Version:                  spec.DataVersionGloas,
			ExecutionPayloadIncluded: true,
			GloasContents:            &apiv1gloas.BlockContents{Block: block},
			BeaconBlockBodyRoot:      &bodyRoot,
		}).Root()
		require.NoError(t, err)
		require.NotEqual(t, phase0.Root{}, got)
		require.Equal(t, wantRootFor(t, block, bodyRoot), got)
	})

	// Root must propagate BodyRoot's error rather than deriving a root from a
	// zero body root: a proposal missing its retained root is not a proposal
	// with an all-zero body.
	t.Run("BeaconBlockBodyRootUnset", func(t *testing.T) {
		_, err := (&api.VersionedEPBSProposal{
			Version: spec.DataVersionGloas,
			Gloas:   newBlock(0x09),
		}).Root()
		require.EqualError(t, err, "no beacon block body root")
	})

	t.Run("NoBlock", func(t *testing.T) {
		_, err := (&api.VersionedEPBSProposal{
			Version: spec.DataVersionGloas,
		}).Root()
		require.EqualError(t, err, "no gloas beacon block")
	})
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
