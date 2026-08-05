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
	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestAssertIncludedEPBSProposalEnvelopeMatchesBlock(t *testing.T) {
	contents := &apiv1gloas.BlockContents{
		Block: &gloas.BeaconBlock{
			Body: &gloas.BeaconBlockBody{
				ETH1Data:                  &phase0.ETH1Data{},
				SyncAggregate:             &altair.SyncAggregate{SyncCommitteeBits: bitfield.NewBitvector512()},
				SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{}},
				ParentExecutionRequests:   &gloas.ExecutionRequests{},
			},
		},
		ExecutionPayloadEnvelope: &gloas.ExecutionPayloadEnvelope{
			Payload:           &gloas.ExecutionPayload{BaseFeePerGas: uint256.NewInt(1)},
			ExecutionRequests: &gloas.ExecutionRequests{},
		},
	}

	blockRoot, err := contents.Block.HashTreeRoot()
	require.NoError(t, err)
	contents.ExecutionPayloadEnvelope.BeaconBlockRoot = blockRoot
	require.NoError(t, assertIncludedEPBSProposalEnvelopeMatchesBlock(contents))

	contents.ExecutionPayloadEnvelope.BeaconBlockRoot = phase0.Root{0xff}
	err = assertIncludedEPBSProposalEnvelopeMatchesBlock(contents)
	require.ErrorIs(t, err, client.ErrInconsistentResult)
}
