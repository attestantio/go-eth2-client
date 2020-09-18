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
	"encoding/json"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// BeaconBlockProposal fetches a proposed beacon block for signing.
func (s *Service) BeaconBlockProposal(ctx context.Context, slot uint64, randaoReveal []byte, graffiti []byte) (*spec.BeaconBlock, error) {
	conn := ethpb.NewBeaconNodeValidatorClient(s.conn)

	// Graffiti should be 32 bytes.
	fixedGraffiti := make([]byte, 32)
	copy(fixedGraffiti, graffiti)

	req := &ethpb.BlockRequest{
		Slot:         slot,
		RandaoReveal: randaoReveal,
		Graffiti:     fixedGraffiti,
	}

	if e := log.Trace(); e.Enabled() {
		jsonData, err := json.Marshal(req)
		if err == nil {
			log.Trace().Str("data", string(jsonData)).Msg("Calling GetBlock()")
		}
	}
	resp, err := conn.GetBlock(ctx, &ethpb.BlockRequest{
		Slot:         slot,
		RandaoReveal: randaoReveal,
		Graffiti:     fixedGraffiti,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain beacon block data")
	}

	block := &spec.BeaconBlock{
		Slot:          resp.Slot,
		ProposerIndex: resp.ProposerIndex,
		ParentRoot:    resp.ParentRoot,
		StateRoot:     resp.StateRoot,
		Body: &spec.BeaconBlockBody{
			RANDAOReveal: randaoReveal,
			ETH1Data: &spec.ETH1Data{
				DepositRoot:  resp.Body.Eth1Data.DepositRoot,
				DepositCount: resp.Body.Eth1Data.DepositCount,
				BlockHash:    resp.Body.Eth1Data.BlockHash,
			},
			Graffiti: fixedGraffiti,
		},
	}
	block.Body.ProposerSlashings = make([]*spec.ProposerSlashing, len(resp.Body.ProposerSlashings))
	for i := range resp.Body.ProposerSlashings {
		block.Body.ProposerSlashings[i] = &spec.ProposerSlashing{
			Header1: &spec.SignedBeaconBlockHeader{
				Message: &spec.BeaconBlockHeader{
					Slot:          resp.Body.ProposerSlashings[i].Header_1.Header.Slot,
					ProposerIndex: resp.Body.ProposerSlashings[i].Header_1.Header.ProposerIndex,
					ParentRoot:    resp.Body.ProposerSlashings[i].Header_1.Header.ParentRoot,
					StateRoot:     resp.Body.ProposerSlashings[i].Header_1.Header.StateRoot,
					BodyRoot:      resp.Body.ProposerSlashings[i].Header_1.Header.BodyRoot,
				},
				Signature: resp.Body.ProposerSlashings[i].Header_1.Signature,
			},
			Header2: &spec.SignedBeaconBlockHeader{
				Message: &spec.BeaconBlockHeader{
					Slot:          resp.Body.ProposerSlashings[i].Header_2.Header.Slot,
					ProposerIndex: resp.Body.ProposerSlashings[i].Header_2.Header.ProposerIndex,
					ParentRoot:    resp.Body.ProposerSlashings[i].Header_2.Header.ParentRoot,
					StateRoot:     resp.Body.ProposerSlashings[i].Header_2.Header.StateRoot,
					BodyRoot:      resp.Body.ProposerSlashings[i].Header_2.Header.BodyRoot,
				},
				Signature: resp.Body.ProposerSlashings[i].Header_2.Signature,
			},
		}
	}
	block.Body.AttesterSlashings = make([]*spec.AttesterSlashing, len(resp.Body.AttesterSlashings))
	for i := range resp.Body.AttesterSlashings {
		block.Body.AttesterSlashings[i] = &spec.AttesterSlashing{
			Attestation1: &spec.IndexedAttestation{
				AttestingIndices: resp.Body.AttesterSlashings[i].Attestation_1.AttestingIndices,
				Data: &spec.AttestationData{
					Slot:            resp.Body.AttesterSlashings[i].Attestation_1.Data.Slot,
					Index:  resp.Body.AttesterSlashings[i].Attestation_1.Data.CommitteeIndex,
					BeaconBlockRoot: resp.Body.AttesterSlashings[i].Attestation_1.Data.BeaconBlockRoot,
					Source: &spec.Checkpoint{
						Epoch: resp.Body.AttesterSlashings[i].Attestation_1.Data.Source.Epoch,
						Root:  resp.Body.AttesterSlashings[i].Attestation_1.Data.Source.Root,
					},
					Target: &spec.Checkpoint{
						Epoch: resp.Body.AttesterSlashings[i].Attestation_1.Data.Target.Epoch,
						Root:  resp.Body.AttesterSlashings[i].Attestation_1.Data.Target.Root,
					},
				},
				Signature: resp.Body.AttesterSlashings[i].Attestation_1.Signature,
			},
			Attestation2: &spec.IndexedAttestation{
				AttestingIndices: resp.Body.AttesterSlashings[i].Attestation_2.AttestingIndices,
				Data: &spec.AttestationData{
					Slot:            resp.Body.AttesterSlashings[i].Attestation_2.Data.Slot,
					Index:  resp.Body.AttesterSlashings[i].Attestation_2.Data.CommitteeIndex,
					BeaconBlockRoot: resp.Body.AttesterSlashings[i].Attestation_2.Data.BeaconBlockRoot,
					Source: &spec.Checkpoint{
						Epoch: resp.Body.AttesterSlashings[i].Attestation_2.Data.Source.Epoch,
						Root:  resp.Body.AttesterSlashings[i].Attestation_2.Data.Source.Root,
					},
					Target: &spec.Checkpoint{
						Epoch: resp.Body.AttesterSlashings[i].Attestation_2.Data.Target.Epoch,
						Root:  resp.Body.AttesterSlashings[i].Attestation_2.Data.Target.Root,
					},
				},
				Signature: resp.Body.AttesterSlashings[i].Attestation_2.Signature,
			},
		}
	}
	block.Body.Attestations = make([]*spec.Attestation, len(resp.Body.Attestations))
	for i := range resp.Body.Attestations {
		block.Body.Attestations[i] = &spec.Attestation{
			AggregationBits: resp.Body.Attestations[i].AggregationBits,
			Data: &spec.AttestationData{
				Slot:            resp.Body.Attestations[i].Data.Slot,
				Index:  resp.Body.Attestations[i].Data.CommitteeIndex,
				BeaconBlockRoot: resp.Body.Attestations[i].Data.BeaconBlockRoot,
				Source: &spec.Checkpoint{
					Epoch: resp.Body.Attestations[i].Data.Source.Epoch,
					Root:  resp.Body.Attestations[i].Data.Source.Root,
				},
				Target: &spec.Checkpoint{
					Epoch: resp.Body.Attestations[i].Data.Target.Epoch,
					Root:  resp.Body.Attestations[i].Data.Target.Root,
				},
			},
			Signature: resp.Body.Attestations[i].Signature,
		}
	}
	block.Body.Deposits = make([]*spec.Deposit, len(resp.Body.Deposits))
	for i := range resp.Body.Deposits {
		block.Body.Deposits[i] = &spec.Deposit{
			Proof: resp.Body.Deposits[i].Proof,
			Data: &spec.DepositData{
				PublicKey:             resp.Body.Deposits[i].Data.PublicKey,
				WithdrawalCredentials: resp.Body.Deposits[i].Data.WithdrawalCredentials,
				Amount:                resp.Body.Deposits[i].Data.Amount,
				Signature:             resp.Body.Deposits[i].Data.Signature,
			},
		}
	}
	block.Body.VoluntaryExits = make([]*spec.SignedVoluntaryExit, len(resp.Body.VoluntaryExits))
	for i := range resp.Body.VoluntaryExits {
		block.Body.VoluntaryExits[i] = &spec.SignedVoluntaryExit{
			Message: &spec.VoluntaryExit{
				Epoch:          resp.Body.VoluntaryExits[i].Exit.Epoch,
				ValidatorIndex: resp.Body.VoluntaryExits[i].Exit.ValidatorIndex,
			},
			Signature: resp.Body.VoluntaryExits[i].Signature,
		}
	}

	return block, nil
}
