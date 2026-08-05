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

package mock

import (
	"context"
	"math/big"

	"github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/api"
	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
)

// EPBSProposal fetches an ePBS proposal for signing.
//
// The response follows opts.IncludePayload, since that is the axis a caller is
// most likely to be exercising: with it set the proposal carries the envelope,
// blobs and proofs needed to publish from anywhere, and without it only a block.
func (s *Service) EPBSProposal(ctx context.Context,
	opts *api.EPBSProposalOpts,
) (
	*api.Response[*api.VersionedEPBSProposal],
	error,
) {
	if s.EPBSProposalFunc != nil {
		return s.EPBSProposalFunc(ctx, opts)
	}

	block := mockGloasBeaconBlock(opts)

	proposal := &api.VersionedEPBSProposal{
		Version:        spec.DataVersionGloas,
		ConsensusValue: big.NewInt(1),
		ExecutionValue: big.NewInt(2),
	}

	if opts.IncludePayload != nil && *opts.IncludePayload {
		proposal.ExecutionPayloadIncluded = true
		proposal.GloasContents = &apiv1gloas.BlockContents{
			Block:                    block,
			ExecutionPayloadEnvelope: mockGloasExecutionPayloadEnvelope(),
			KZGProofs:                []deneb.KZGProof{},
			Blobs:                    []deneb.Blob{},
		}
	} else {
		proposal.Gloas = block
	}

	return &api.Response[*api.VersionedEPBSProposal]{
		Data:     proposal,
		Metadata: make(map[string]any),
	}, nil
}

// mockGloasBeaconBlock returns a block for the requested slot, echoing back the
// RANDAO reveal and graffiti so that a caller's own consistency checks pass.
func mockGloasBeaconBlock(opts *api.EPBSProposalOpts) *gloas.BeaconBlock {
	return &gloas.BeaconBlock{
		Slot: opts.Slot,
		Body: &gloas.BeaconBlockBody{
			RANDAOReveal: opts.RandaoReveal,
			Graffiti:     opts.Graffiti,
			ETH1Data: &phase0.ETH1Data{
				BlockHash: make([]byte, phase0.Hash32Length),
			},
			ProposerSlashings: []*phase0.ProposerSlashing{},
			AttesterSlashings: []*gloas.AttesterSlashing{},
			Attestations:      []*gloas.Attestation{},
			Deposits:          []*phase0.Deposit{},
			VoluntaryExits:    []*phase0.SignedVoluntaryExit{},
			SyncAggregate: &altair.SyncAggregate{
				SyncCommitteeBits: bitfield.NewBitvector512(),
			},
			BLSToExecutionChanges: []*capella.SignedBLSToExecutionChange{},
			SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{
				Message: &gloas.ExecutionPayloadBid{
					Slot:               opts.Slot,
					BlobKZGCommitments: []deneb.KZGCommitment{},
				},
			},
			PayloadAttestations:     []*gloas.PayloadAttestation{},
			ParentExecutionRequests: mockGloasExecutionRequests(),
		},
	}
}

// mockGloasExecutionPayloadEnvelope returns an envelope with every nilable field
// of its payload set.  BaseFeePerGas matters most: it is a *uint256.Int, so a
// consumer that reads it rather than the marshaled JSON dereferences nil.  The
// JSON marshaler itself substitutes a zero, which means an unset one would
// otherwise go unnoticed.
func mockGloasExecutionPayloadEnvelope() *gloas.ExecutionPayloadEnvelope {
	return &gloas.ExecutionPayloadEnvelope{
		Payload: &gloas.ExecutionPayload{
			BaseFeePerGas:   uint256.NewInt(1),
			ExtraData:       []byte{},
			Transactions:    []bellatrix.Transaction{},
			Withdrawals:     []*capella.Withdrawal{},
			BlockAccessList: []byte{},
		},
		ExecutionRequests: mockGloasExecutionRequests(),
	}
}

// mockGloasExecutionRequests returns execution requests with every list
// allocated, so that the container round-trips through JSON unchanged.
func mockGloasExecutionRequests() *gloas.ExecutionRequests {
	return &gloas.ExecutionRequests{
		Deposits:        []*electra.DepositRequest{},
		Withdrawals:     []*electra.WithdrawalRequest{},
		Consolidations:  []*electra.ConsolidationRequest{},
		BuilderDeposits: []*gloas.BuilderDepositRequest{},
		BuilderExits:    []*gloas.BuilderExitRequest{},
	}
}
