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

// SignedBeaconBlockBySlot fetches a signed beacon block given its slot.
func (s *Service) SignedBeaconBlockBySlot(ctx context.Context, slot uint64) (*spec.SignedBeaconBlock, error) {
	conn := ethpb.NewBeaconChainClient(s.conn)

	req := &ethpb.ListBlocksRequest{}
	if slot == 0 {
		req.QueryFilter = &ethpb.ListBlocksRequest_Genesis{Genesis: true}
	} else {
		req.QueryFilter = &ethpb.ListBlocksRequest_Slot{Slot: slot}
	}
	resp, err := conn.ListBlocks(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "call to ListBlocks() failed")
	}
	if len(resp.BlockContainers) == 0 {
		return nil, nil
	}

	block := resp.BlockContainers[0].Block

	signedBeaconBlock := &spec.SignedBeaconBlock{
		Signature: block.Signature,
		Message: &spec.BeaconBlock{
			Slot:          block.Block.Slot,
			ProposerIndex: block.Block.ProposerIndex,
			ParentRoot:    block.Block.ParentRoot,
			StateRoot:     block.Block.StateRoot,
			Body: &spec.BeaconBlockBody{
				RANDAOReveal: block.Block.Body.RandaoReveal,
				ETH1Data: &spec.ETH1Data{
					DepositRoot:  block.Block.Body.Eth1Data.DepositRoot,
					DepositCount: block.Block.Body.Eth1Data.DepositCount,
					BlockHash:    block.Block.Body.Eth1Data.BlockHash,
				},
				Graffiti: block.Block.Body.Graffiti,
			},
		},
	}
	signedBeaconBlock.Message.Body.ProposerSlashings = make([]*spec.ProposerSlashing, len(block.Block.Body.ProposerSlashings))
	for i := range block.Block.Body.ProposerSlashings {
		signedBeaconBlock.Message.Body.ProposerSlashings[i] = &spec.ProposerSlashing{
			Header1: &spec.SignedBeaconBlockHeader{
				Message: &spec.BeaconBlockHeader{
					Slot:          block.Block.Body.ProposerSlashings[i].Header_1.Header.Slot,
					ProposerIndex: block.Block.Body.ProposerSlashings[i].Header_1.Header.ProposerIndex,
					ParentRoot:    block.Block.Body.ProposerSlashings[i].Header_1.Header.ParentRoot,
					StateRoot:     block.Block.Body.ProposerSlashings[i].Header_1.Header.StateRoot,
					BodyRoot:      block.Block.Body.ProposerSlashings[i].Header_1.Header.BodyRoot,
				},
				Signature: block.Block.Body.ProposerSlashings[i].Header_1.Signature,
			},
			Header2: &spec.SignedBeaconBlockHeader{
				Message: &spec.BeaconBlockHeader{
					Slot:          block.Block.Body.ProposerSlashings[i].Header_2.Header.Slot,
					ProposerIndex: block.Block.Body.ProposerSlashings[i].Header_2.Header.ProposerIndex,
					ParentRoot:    block.Block.Body.ProposerSlashings[i].Header_2.Header.ParentRoot,
					StateRoot:     block.Block.Body.ProposerSlashings[i].Header_2.Header.StateRoot,
					BodyRoot:      block.Block.Body.ProposerSlashings[i].Header_2.Header.BodyRoot,
				},
				Signature: block.Block.Body.ProposerSlashings[i].Header_2.Signature,
			},
		}
	}
	signedBeaconBlock.Message.Body.AttesterSlashings = make([]*spec.AttesterSlashing, len(block.Block.Body.AttesterSlashings))
	for i := range block.Block.Body.AttesterSlashings {
		signedBeaconBlock.Message.Body.AttesterSlashings[i] = &spec.AttesterSlashing{
			Attestation1: &spec.IndexedAttestation{
				AttestingIndices: block.Block.Body.AttesterSlashings[i].Attestation_1.AttestingIndices,
				Data: &spec.AttestationData{
					Slot:            block.Block.Body.AttesterSlashings[i].Attestation_1.Data.Slot,
					Index:           block.Block.Body.AttesterSlashings[i].Attestation_1.Data.CommitteeIndex,
					BeaconBlockRoot: block.Block.Body.AttesterSlashings[i].Attestation_1.Data.BeaconBlockRoot,
					Source: &spec.Checkpoint{
						Epoch: block.Block.Body.AttesterSlashings[i].Attestation_1.Data.Source.Epoch,
						Root:  block.Block.Body.AttesterSlashings[i].Attestation_1.Data.Source.Root,
					},
					Target: &spec.Checkpoint{
						Epoch: block.Block.Body.AttesterSlashings[i].Attestation_1.Data.Target.Epoch,
						Root:  block.Block.Body.AttesterSlashings[i].Attestation_1.Data.Target.Root,
					},
				},
				Signature: block.Block.Body.AttesterSlashings[i].Attestation_1.Signature,
			},
			Attestation2: &spec.IndexedAttestation{
				AttestingIndices: block.Block.Body.AttesterSlashings[i].Attestation_2.AttestingIndices,
				Data: &spec.AttestationData{
					Slot:            block.Block.Body.AttesterSlashings[i].Attestation_2.Data.Slot,
					Index:           block.Block.Body.AttesterSlashings[i].Attestation_2.Data.CommitteeIndex,
					BeaconBlockRoot: block.Block.Body.AttesterSlashings[i].Attestation_2.Data.BeaconBlockRoot,
					Source: &spec.Checkpoint{
						Epoch: block.Block.Body.AttesterSlashings[i].Attestation_2.Data.Source.Epoch,
						Root:  block.Block.Body.AttesterSlashings[i].Attestation_2.Data.Source.Root,
					},
					Target: &spec.Checkpoint{
						Epoch: block.Block.Body.AttesterSlashings[i].Attestation_2.Data.Target.Epoch,
						Root:  block.Block.Body.AttesterSlashings[i].Attestation_2.Data.Target.Root,
					},
				},
				Signature: block.Block.Body.AttesterSlashings[i].Attestation_2.Signature,
			},
		}
	}
	signedBeaconBlock.Message.Body.Attestations = make([]*spec.Attestation, len(block.Block.Body.Attestations))
	for i := range block.Block.Body.Attestations {
		signedBeaconBlock.Message.Body.Attestations[i] = &spec.Attestation{
			AggregationBits: block.Block.Body.Attestations[i].AggregationBits,
			Data: &spec.AttestationData{
				Slot:            block.Block.Body.Attestations[i].Data.Slot,
				Index:           block.Block.Body.Attestations[i].Data.CommitteeIndex,
				BeaconBlockRoot: block.Block.Body.Attestations[i].Data.BeaconBlockRoot,
				Source: &spec.Checkpoint{
					Epoch: block.Block.Body.Attestations[i].Data.Source.Epoch,
					Root:  block.Block.Body.Attestations[i].Data.Source.Root,
				},
				Target: &spec.Checkpoint{
					Epoch: block.Block.Body.Attestations[i].Data.Target.Epoch,
					Root:  block.Block.Body.Attestations[i].Data.Target.Root,
				},
			},
			Signature: block.Block.Body.Attestations[i].Signature,
		}
	}
	signedBeaconBlock.Message.Body.Deposits = make([]*spec.Deposit, len(block.Block.Body.Deposits))
	for i := range block.Block.Body.Deposits {
		signedBeaconBlock.Message.Body.Deposits[i] = &spec.Deposit{
			Proof: block.Block.Body.Deposits[i].Proof,
			Data: &spec.DepositData{
				PublicKey:             block.Block.Body.Deposits[i].Data.PublicKey,
				WithdrawalCredentials: block.Block.Body.Deposits[i].Data.WithdrawalCredentials,
				Amount:                block.Block.Body.Deposits[i].Data.Amount,
				Signature:             block.Block.Body.Deposits[i].Data.Signature,
			},
		}
	}
	signedBeaconBlock.Message.Body.VoluntaryExits = make([]*spec.SignedVoluntaryExit, len(block.Block.Body.VoluntaryExits))
	for i := range block.Block.Body.VoluntaryExits {
		signedBeaconBlock.Message.Body.VoluntaryExits[i] = &spec.SignedVoluntaryExit{
			Message: &spec.VoluntaryExit{
				Epoch:          block.Block.Body.VoluntaryExits[i].Exit.Epoch,
				ValidatorIndex: block.Block.Body.VoluntaryExits[i].Exit.ValidatorIndex,
			},
			Signature: block.Block.Body.VoluntaryExits[i].Signature,
		}
	}

	return signedBeaconBlock, nil
}
