// Copyright © 2020 Attestant Limited.
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

package prysmgrpc

import (
	"context"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// SubmitBeaconBlock submits a beacon block.
func (s *Service) SubmitBeaconBlock(ctx context.Context, block *spec.SignedBeaconBlock) error {
	proposal := &ethpb.SignedBeaconBlock{
		Block: &ethpb.BeaconBlock{
			Slot:          uint64(block.Message.Slot),
			ProposerIndex: uint64(block.Message.ProposerIndex),
			ParentRoot:    block.Message.ParentRoot[:],
			StateRoot:     block.Message.StateRoot[:],
			Body: &ethpb.BeaconBlockBody{
				RandaoReveal: block.Message.Body.RANDAOReveal[:],
				Eth1Data: &ethpb.Eth1Data{
					DepositRoot:  block.Message.Body.ETH1Data.DepositRoot[:],
					DepositCount: block.Message.Body.ETH1Data.DepositCount,
					BlockHash:    block.Message.Body.ETH1Data.BlockHash,
				},
				Graffiti: block.Message.Body.Graffiti,
			},
		},
		Signature: block.Signature[:],
	}
	// Shorthand for references below.
	body := block.Message.Body
	proposal.Block.Body.ProposerSlashings = make([]*ethpb.ProposerSlashing, len(body.ProposerSlashings))
	for i := range body.ProposerSlashings {
		proposal.Block.Body.ProposerSlashings[i] = &ethpb.ProposerSlashing{
			Header_1: &ethpb.SignedBeaconBlockHeader{
				Header: &ethpb.BeaconBlockHeader{
					Slot:          uint64(body.ProposerSlashings[i].SignedHeader1.Message.Slot),
					ProposerIndex: uint64(body.ProposerSlashings[i].SignedHeader1.Message.ProposerIndex),
					ParentRoot:    body.ProposerSlashings[i].SignedHeader1.Message.ParentRoot[:],
					StateRoot:     body.ProposerSlashings[i].SignedHeader1.Message.StateRoot[:],
					BodyRoot:      body.ProposerSlashings[i].SignedHeader1.Message.BodyRoot[:],
				},
				Signature: body.ProposerSlashings[i].SignedHeader1.Signature[:],
			},
			Header_2: &ethpb.SignedBeaconBlockHeader{
				Header: &ethpb.BeaconBlockHeader{
					Slot:          uint64(body.ProposerSlashings[i].SignedHeader2.Message.Slot),
					ProposerIndex: uint64(body.ProposerSlashings[i].SignedHeader2.Message.ProposerIndex),
					ParentRoot:    body.ProposerSlashings[i].SignedHeader2.Message.ParentRoot[:],
					StateRoot:     body.ProposerSlashings[i].SignedHeader2.Message.StateRoot[:],
					BodyRoot:      body.ProposerSlashings[i].SignedHeader2.Message.BodyRoot[:],
				},
				Signature: body.ProposerSlashings[i].SignedHeader2.Signature[:],
			},
		}
	}

	proposal.Block.Body.AttesterSlashings = make([]*ethpb.AttesterSlashing, len(body.AttesterSlashings))
	for i := range body.AttesterSlashings {
		proposal.Block.Body.AttesterSlashings[i] = &ethpb.AttesterSlashing{
			Attestation_1: &ethpb.IndexedAttestation{
				AttestingIndices: body.AttesterSlashings[i].Attestation1.AttestingIndices,
				Data: &ethpb.AttestationData{
					Slot:            uint64(body.AttesterSlashings[i].Attestation1.Data.Slot),
					CommitteeIndex:  uint64(body.AttesterSlashings[i].Attestation1.Data.Index),
					BeaconBlockRoot: body.AttesterSlashings[i].Attestation1.Data.BeaconBlockRoot[:],
					Source: &ethpb.Checkpoint{
						Epoch: uint64(body.AttesterSlashings[i].Attestation1.Data.Source.Epoch),
						Root:  body.AttesterSlashings[i].Attestation1.Data.Source.Root[:],
					},
					Target: &ethpb.Checkpoint{
						Epoch: uint64(body.AttesterSlashings[i].Attestation1.Data.Target.Epoch),
						Root:  body.AttesterSlashings[i].Attestation1.Data.Target.Root[:],
					},
				},
				Signature: body.AttesterSlashings[i].Attestation1.Signature[:],
			},
			Attestation_2: &ethpb.IndexedAttestation{
				AttestingIndices: body.AttesterSlashings[i].Attestation2.AttestingIndices,
				Data: &ethpb.AttestationData{
					Slot:            uint64(body.AttesterSlashings[i].Attestation2.Data.Slot),
					CommitteeIndex:  uint64(body.AttesterSlashings[i].Attestation2.Data.Index),
					BeaconBlockRoot: body.AttesterSlashings[i].Attestation2.Data.BeaconBlockRoot[:],
					Source: &ethpb.Checkpoint{
						Epoch: uint64(body.AttesterSlashings[i].Attestation2.Data.Source.Epoch),
						Root:  body.AttesterSlashings[i].Attestation2.Data.Source.Root[:],
					},
					Target: &ethpb.Checkpoint{
						Epoch: uint64(body.AttesterSlashings[i].Attestation2.Data.Target.Epoch),
						Root:  body.AttesterSlashings[i].Attestation2.Data.Target.Root[:],
					},
				},
				Signature: body.AttesterSlashings[i].Attestation2.Signature[:],
			},
		}
	}

	proposal.Block.Body.Attestations = make([]*ethpb.Attestation, len(body.Attestations))
	for i := range body.Attestations {
		proposal.Block.Body.Attestations[i] = &ethpb.Attestation{
			AggregationBits: body.Attestations[i].AggregationBits,
			Data: &ethpb.AttestationData{
				Slot:            uint64(body.Attestations[i].Data.Slot),
				CommitteeIndex:  uint64(body.Attestations[i].Data.Index),
				BeaconBlockRoot: body.Attestations[i].Data.BeaconBlockRoot[:],
				Source: &ethpb.Checkpoint{
					Epoch: uint64(body.Attestations[i].Data.Source.Epoch),
					Root:  body.Attestations[i].Data.Source.Root[:],
				},
				Target: &ethpb.Checkpoint{
					Epoch: uint64(body.Attestations[i].Data.Target.Epoch),
					Root:  body.Attestations[i].Data.Target.Root[:],
				},
			},
			Signature: body.Attestations[i].Signature[:],
		}
	}

	proposal.Block.Body.Deposits = make([]*ethpb.Deposit, len(body.Deposits))
	for i := range body.Deposits {
		proposal.Block.Body.Deposits[i] = &ethpb.Deposit{
			Proof: body.Deposits[i].Proof,
			Data: &ethpb.Deposit_Data{
				PublicKey:             body.Deposits[i].Data.PublicKey[:],
				WithdrawalCredentials: body.Deposits[i].Data.WithdrawalCredentials,
				Amount:                uint64(body.Deposits[i].Data.Amount),
				Signature:             body.Deposits[i].Data.Signature[:],
			},
		}
	}

	proposal.Block.Body.VoluntaryExits = make([]*ethpb.SignedVoluntaryExit, len(body.VoluntaryExits))
	for i := range body.VoluntaryExits {
		proposal.Block.Body.VoluntaryExits[i] = &ethpb.SignedVoluntaryExit{
			Exit: &ethpb.VoluntaryExit{
				Epoch:          uint64(body.VoluntaryExits[i].Message.Epoch),
				ValidatorIndex: uint64(body.VoluntaryExits[i].Message.ValidatorIndex),
			},
			Signature: body.VoluntaryExits[i].Signature[:],
		}
	}

	conn := ethpb.NewBeaconNodeValidatorClient(s.conn)
	log.Trace().Msg("Calling ProposeBlock()")
	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	_, err := conn.ProposeBlock(opCtx, proposal)
	cancel()

	if err != nil {
		return errors.Wrap(err, "failed to submit beacon block")
	}
	return nil
}
