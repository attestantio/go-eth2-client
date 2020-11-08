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
	"strconv"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// SignedBeaconBlock fetches a signed beacon block given a block ID.
func (s *Service) SignedBeaconBlock(ctx context.Context, blockID string) (*spec.SignedBeaconBlock, error) {
	slot, err := strconv.ParseUint(blockID, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "invalid block ID")
	}
	return s.SignedBeaconBlockBySlot(ctx, slot)
}

// SignedBeaconBlockBySlot fetches a signed beacon block given its slot.
func (s *Service) SignedBeaconBlockBySlot(ctx context.Context, slot uint64) (*spec.SignedBeaconBlock, error) {
	conn := ethpb.NewBeaconChainClient(s.conn)

	req := &ethpb.ListBlocksRequest{}
	if slot == 0 {
		req.QueryFilter = &ethpb.ListBlocksRequest_Genesis{Genesis: true}
	} else {
		req.QueryFilter = &ethpb.ListBlocksRequest_Slot{Slot: slot}
	}
	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	resp, err := conn.ListBlocks(opCtx, req)
	cancel()
	if err != nil {
		return nil, errors.Wrap(err, "call to ListBlocks() failed")
	}
	if len(resp.BlockContainers) == 0 {
		return nil, nil
	}

	block := resp.BlockContainers[0].Block

	signedBeaconBlock := &spec.SignedBeaconBlock{
		Message: &spec.BeaconBlock{
			Slot:          spec.Slot(block.Block.Slot),
			ProposerIndex: spec.ValidatorIndex(block.Block.ProposerIndex),
			Body: &spec.BeaconBlockBody{
				ETH1Data: &spec.ETH1Data{
					DepositCount: block.Block.Body.Eth1Data.DepositCount,
					BlockHash:    block.Block.Body.Eth1Data.BlockHash,
				},
				Graffiti: block.Block.Body.Graffiti,
			},
		},
	}
	copy(signedBeaconBlock.Signature[:], block.Signature)
	copy(signedBeaconBlock.Message.ParentRoot[:], block.Block.ParentRoot)
	copy(signedBeaconBlock.Message.StateRoot[:], block.Block.StateRoot)
	copy(signedBeaconBlock.Message.Body.RANDAOReveal[:], block.Block.Body.RandaoReveal)
	copy(signedBeaconBlock.Message.Body.ETH1Data.DepositRoot[:], block.Block.Body.Eth1Data.DepositRoot)
	signedBeaconBlock.Message.Body.ProposerSlashings = make([]*spec.ProposerSlashing, len(block.Block.Body.ProposerSlashings))
	for i := range block.Block.Body.ProposerSlashings {
		signedBeaconBlock.Message.Body.ProposerSlashings[i] = &spec.ProposerSlashing{
			SignedHeader1: &spec.SignedBeaconBlockHeader{
				Message: &spec.BeaconBlockHeader{
					Slot:          spec.Slot(block.Block.Body.ProposerSlashings[i].Header_1.Header.Slot),
					ProposerIndex: spec.ValidatorIndex(block.Block.Body.ProposerSlashings[i].Header_1.Header.ProposerIndex),
				},
			},
			SignedHeader2: &spec.SignedBeaconBlockHeader{
				Message: &spec.BeaconBlockHeader{
					Slot:          spec.Slot(block.Block.Body.ProposerSlashings[i].Header_2.Header.Slot),
					ProposerIndex: spec.ValidatorIndex(block.Block.Body.ProposerSlashings[i].Header_2.Header.ProposerIndex),
				},
			},
		}
		copy(signedBeaconBlock.Message.Body.ProposerSlashings[i].SignedHeader1.Message.ParentRoot[:], block.Block.Body.ProposerSlashings[i].Header_1.Header.ParentRoot)
		copy(signedBeaconBlock.Message.Body.ProposerSlashings[i].SignedHeader1.Message.StateRoot[:], block.Block.Body.ProposerSlashings[i].Header_1.Header.StateRoot)
		copy(signedBeaconBlock.Message.Body.ProposerSlashings[i].SignedHeader1.Message.BodyRoot[:], block.Block.Body.ProposerSlashings[i].Header_1.Header.BodyRoot)
		copy(signedBeaconBlock.Message.Body.ProposerSlashings[i].SignedHeader1.Signature[:], block.Block.Body.ProposerSlashings[i].Header_1.Signature)
		copy(signedBeaconBlock.Message.Body.ProposerSlashings[i].SignedHeader2.Message.ParentRoot[:], block.Block.Body.ProposerSlashings[i].Header_2.Header.ParentRoot)
		copy(signedBeaconBlock.Message.Body.ProposerSlashings[i].SignedHeader2.Message.StateRoot[:], block.Block.Body.ProposerSlashings[i].Header_2.Header.StateRoot)
		copy(signedBeaconBlock.Message.Body.ProposerSlashings[i].SignedHeader2.Message.BodyRoot[:], block.Block.Body.ProposerSlashings[i].Header_2.Header.BodyRoot)
		copy(signedBeaconBlock.Message.Body.ProposerSlashings[i].SignedHeader2.Signature[:], block.Block.Body.ProposerSlashings[i].Header_2.Signature)
	}
	signedBeaconBlock.Message.Body.AttesterSlashings = make([]*spec.AttesterSlashing, len(block.Block.Body.AttesterSlashings))
	for i := range block.Block.Body.AttesterSlashings {
		signedBeaconBlock.Message.Body.AttesterSlashings[i] = &spec.AttesterSlashing{
			Attestation1: &spec.IndexedAttestation{
				AttestingIndices: block.Block.Body.AttesterSlashings[i].Attestation_1.AttestingIndices,
				Data: &spec.AttestationData{
					Slot:  spec.Slot(block.Block.Body.AttesterSlashings[i].Attestation_1.Data.Slot),
					Index: spec.CommitteeIndex(block.Block.Body.AttesterSlashings[i].Attestation_1.Data.CommitteeIndex),
					Source: &spec.Checkpoint{
						Epoch: spec.Epoch(block.Block.Body.AttesterSlashings[i].Attestation_1.Data.Source.Epoch),
					},
					Target: &spec.Checkpoint{
						Epoch: spec.Epoch(block.Block.Body.AttesterSlashings[i].Attestation_1.Data.Target.Epoch),
					},
				},
			},
			Attestation2: &spec.IndexedAttestation{
				AttestingIndices: block.Block.Body.AttesterSlashings[i].Attestation_2.AttestingIndices,
				Data: &spec.AttestationData{
					Slot:  spec.Slot(block.Block.Body.AttesterSlashings[i].Attestation_2.Data.Slot),
					Index: spec.CommitteeIndex(block.Block.Body.AttesterSlashings[i].Attestation_2.Data.CommitteeIndex),
					Source: &spec.Checkpoint{
						Epoch: spec.Epoch(block.Block.Body.AttesterSlashings[i].Attestation_2.Data.Source.Epoch),
					},
					Target: &spec.Checkpoint{
						Epoch: spec.Epoch(block.Block.Body.AttesterSlashings[i].Attestation_2.Data.Target.Epoch),
					},
				},
			},
		}
		copy(signedBeaconBlock.Message.Body.AttesterSlashings[i].Attestation1.Data.BeaconBlockRoot[:], block.Block.Body.AttesterSlashings[i].Attestation_1.Data.BeaconBlockRoot)
		copy(signedBeaconBlock.Message.Body.AttesterSlashings[i].Attestation1.Data.Source.Root[:], block.Block.Body.AttesterSlashings[i].Attestation_1.Data.Source.Root)
		copy(signedBeaconBlock.Message.Body.AttesterSlashings[i].Attestation1.Data.Target.Root[:], block.Block.Body.AttesterSlashings[i].Attestation_1.Data.Target.Root)
		copy(signedBeaconBlock.Message.Body.AttesterSlashings[i].Attestation1.Signature[:], block.Block.Body.AttesterSlashings[i].Attestation_1.Signature)
		copy(signedBeaconBlock.Message.Body.AttesterSlashings[i].Attestation2.Data.BeaconBlockRoot[:], block.Block.Body.AttesterSlashings[i].Attestation_2.Data.BeaconBlockRoot)
		copy(signedBeaconBlock.Message.Body.AttesterSlashings[i].Attestation2.Data.Source.Root[:], block.Block.Body.AttesterSlashings[i].Attestation_2.Data.Source.Root)
		copy(signedBeaconBlock.Message.Body.AttesterSlashings[i].Attestation2.Data.Target.Root[:], block.Block.Body.AttesterSlashings[i].Attestation_2.Data.Target.Root)
		copy(signedBeaconBlock.Message.Body.AttesterSlashings[i].Attestation2.Signature[:], block.Block.Body.AttesterSlashings[i].Attestation_2.Signature)
	}
	signedBeaconBlock.Message.Body.Attestations = make([]*spec.Attestation, len(block.Block.Body.Attestations))
	for i := range block.Block.Body.Attestations {
		signedBeaconBlock.Message.Body.Attestations[i] = &spec.Attestation{
			AggregationBits: block.Block.Body.Attestations[i].AggregationBits,
			Data: &spec.AttestationData{
				Slot:  spec.Slot(block.Block.Body.Attestations[i].Data.Slot),
				Index: spec.CommitteeIndex(block.Block.Body.Attestations[i].Data.CommitteeIndex),
				Source: &spec.Checkpoint{
					Epoch: spec.Epoch(block.Block.Body.Attestations[i].Data.Source.Epoch),
				},
				Target: &spec.Checkpoint{
					Epoch: spec.Epoch(block.Block.Body.Attestations[i].Data.Target.Epoch),
				},
			},
		}
		copy(signedBeaconBlock.Message.Body.Attestations[i].Data.BeaconBlockRoot[:], block.Block.Body.Attestations[i].Data.BeaconBlockRoot)
		copy(signedBeaconBlock.Message.Body.Attestations[i].Data.Source.Root[:], block.Block.Body.Attestations[i].Data.Source.Root)
		copy(signedBeaconBlock.Message.Body.Attestations[i].Data.Target.Root[:], block.Block.Body.Attestations[i].Data.Target.Root)
		copy(signedBeaconBlock.Message.Body.Attestations[i].Signature[:], block.Block.Body.Attestations[i].Signature)
	}
	signedBeaconBlock.Message.Body.Deposits = make([]*spec.Deposit, len(block.Block.Body.Deposits))
	for i := range block.Block.Body.Deposits {
		signedBeaconBlock.Message.Body.Deposits[i] = &spec.Deposit{
			Proof: block.Block.Body.Deposits[i].Proof,
			Data: &spec.DepositData{
				WithdrawalCredentials: block.Block.Body.Deposits[i].Data.WithdrawalCredentials,
				Amount:                spec.Gwei(block.Block.Body.Deposits[i].Data.Amount),
			},
		}
		copy(signedBeaconBlock.Message.Body.Deposits[i].Data.PublicKey[:], block.Block.Body.Deposits[i].Data.PublicKey)
		copy(signedBeaconBlock.Message.Body.Deposits[i].Data.Signature[:], block.Block.Body.Deposits[i].Data.Signature)
	}
	signedBeaconBlock.Message.Body.VoluntaryExits = make([]*spec.SignedVoluntaryExit, len(block.Block.Body.VoluntaryExits))
	for i := range block.Block.Body.VoluntaryExits {
		signedBeaconBlock.Message.Body.VoluntaryExits[i] = &spec.SignedVoluntaryExit{
			Message: &spec.VoluntaryExit{
				Epoch:          spec.Epoch(block.Block.Body.VoluntaryExits[i].Exit.Epoch),
				ValidatorIndex: spec.ValidatorIndex(block.Block.Body.VoluntaryExits[i].Exit.ValidatorIndex),
			},
		}
		copy(signedBeaconBlock.Message.Body.VoluntaryExits[i].Signature[:], block.Block.Body.VoluntaryExits[i].Signature)
	}

	return signedBeaconBlock, nil
}
