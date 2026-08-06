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
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestVersionedSignedProposalExecutionBlockHash verifies that a Gloas proposal
// reports the execution block hash its proposer committed to.  Post-Gloas that
// hash lives in the execution payload bid rather than in a payload or a payload
// header, so a wrapper that recognised every earlier fork but not this one would
// refuse every single post-ePBS proposal.
//
// The absent-data cases matter for a subtler reason: they must fail as missing
// data, not as an unsupported version.  A caller distinguishes "this node gave
// me an incomplete proposal" from "this library cannot read this fork" by that
// error alone, and the second answer is wrong for a fork the wrapper carries a
// field for.
func TestVersionedSignedProposalExecutionBlockHash(t *testing.T) {
	blockHash := phase0.Hash32{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
	}

	tests := []struct {
		name     string
		proposal *api.VersionedSignedProposal
		expected phase0.Hash32
		err      error
	}{
		{
			name: "Gloas",
			proposal: &api.VersionedSignedProposal{
				Version: spec.DataVersionGloas,
				Gloas: &gloas.SignedBeaconBlock{
					Message: &gloas.BeaconBlock{
						Body: &gloas.BeaconBlockBody{
							SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{
								Message: &gloas.ExecutionPayloadBid{
									BlockHash: blockHash,
								},
							},
						},
					},
				},
			},
			expected: blockHash,
		},
		{
			// The bid is the only place a Gloas proposal carries an execution
			// block hash, so its absence is the case that must not report the
			// zero hash as a successful answer.
			name: "GloasNoBid",
			proposal: &api.VersionedSignedProposal{
				Version: spec.DataVersionGloas,
				Gloas: &gloas.SignedBeaconBlock{
					Message: &gloas.BeaconBlock{
						Body: &gloas.BeaconBlockBody{},
					},
				},
			},
			err: api.ErrDataMissing,
		},
		{
			// A signed bid whose message is absent carries a signature over
			// nothing; the hash has to be read from the message, not the wrapper.
			name: "GloasNoBidMessage",
			proposal: &api.VersionedSignedProposal{
				Version: spec.DataVersionGloas,
				Gloas: &gloas.SignedBeaconBlock{
					Message: &gloas.BeaconBlock{
						Body: &gloas.BeaconBlockBody{
							SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{},
						},
					},
				},
			},
			err: api.ErrDataMissing,
		},
		{
			name: "GloasNoBody",
			proposal: &api.VersionedSignedProposal{
				Version: spec.DataVersionGloas,
				Gloas: &gloas.SignedBeaconBlock{
					Message: &gloas.BeaconBlock{},
				},
			},
			err: api.ErrDataMissing,
		},
		{
			name: "GloasNoMessage",
			proposal: &api.VersionedSignedProposal{
				Version: spec.DataVersionGloas,
				Gloas:   &gloas.SignedBeaconBlock{},
			},
			err: api.ErrDataMissing,
		},
		{
			name: "GloasNoBlock",
			proposal: &api.VersionedSignedProposal{
				Version: spec.DataVersionGloas,
			},
			err: api.ErrDataMissing,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.proposal.ExecutionBlockHash()
			if test.err != nil {
				require.ErrorIs(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, got)
			}
		})
	}
}

func TestVersionedSignedProposalGloasBlindedState(t *testing.T) {
	tests := []struct {
		name     string
		proposal *api.VersionedSignedProposal
		err      string
	}{
		{
			name: "Unblinded",
			proposal: &api.VersionedSignedProposal{
				Version: spec.DataVersionGloas,
				Gloas:   &gloas.SignedBeaconBlock{},
			},
		},
		{
			name: "Blinded",
			proposal: &api.VersionedSignedProposal{
				Version: spec.DataVersionGloas,
				Blinded: true,
				Gloas:   &gloas.SignedBeaconBlock{},
			},
			err: "gloas proposals are never blinded",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.proposal.AssertPresent()
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
