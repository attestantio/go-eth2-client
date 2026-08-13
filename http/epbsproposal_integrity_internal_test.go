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
	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestAssertIncludedEPBSProposalEnvelopeMatchesBlock(t *testing.T) {
	block := &gloas.BeaconBlock{
		Body: &gloas.BeaconBlockBody{
			ETH1Data:                  &phase0.ETH1Data{},
			SyncAggregate:             &altair.SyncAggregate{SyncCommitteeBits: bitfield.NewBitvector512()},
			SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{}},
			ParentExecutionRequests:   &gloas.ExecutionRequests{},
		},
	}
	contents := &apiv1gloas.BlockContents{
		Block: block,
		ExecutionPayloadEnvelope: &gloas.ExecutionPayloadEnvelope{
			Payload:           &gloas.ExecutionPayload{BaseFeePerGas: uint256.NewInt(1)},
			ExecutionRequests: &gloas.ExecutionRequests{},
		},
	}

	bodyRoot, err := block.Body.HashTreeRoot()
	require.NoError(t, err)
	root := phase0.Root(bodyRoot)

	proposal := &api.VersionedEPBSProposal{
		Version:                  spec.DataVersionGloas,
		ExecutionPayloadIncluded: true,
		GloasContents:            contents,
		BeaconBlockBodyRoot:      &root,
	}

	blockRoot, err := proposal.Root()
	require.NoError(t, err)
	contents.ExecutionPayloadEnvelope.BeaconBlockRoot = blockRoot
	require.NoError(t, assertIncludedEPBSProposalEnvelopeMatchesBlock(proposal))

	contents.ExecutionPayloadEnvelope.BeaconBlockRoot = phase0.Root{0xff}
	err = assertIncludedEPBSProposalEnvelopeMatchesBlock(proposal)
	require.ErrorIs(t, err, client.ErrInconsistentResult)

	// A proposal whose body root was never retained must not fall back to
	// the generated (and potentially wrong-preset) block hash: the guard
	// must error rather than pass.
	contents.ExecutionPayloadEnvelope.BeaconBlockRoot = blockRoot
	proposal.BeaconBlockBodyRoot = nil
	err = assertIncludedEPBSProposalEnvelopeMatchesBlock(proposal)
	require.Error(t, err)
	require.NotErrorIs(t, err, client.ErrInconsistentResult)
}
