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
func (s *Service) BeaconBlockProposal(ctx context.Context, slot spec.Slot, randaoReveal spec.BLSSignature, graffiti []byte) (*spec.BeaconBlock, error) {
	conn := ethpb.NewBeaconNodeValidatorClient(s.conn)

	// Graffiti should be 32 bytes.
	fixedGraffiti := make([]byte, 32)
	copy(fixedGraffiti, graffiti)

	req := &ethpb.BlockRequest{
		Slot:         uint64(slot),
		RandaoReveal: randaoReveal[:],
		Graffiti:     fixedGraffiti,
	}

	if e := log.Trace(); e.Enabled() {
		jsonData, err := json.Marshal(req)
		if err == nil {
			log.Trace().Str("data", string(jsonData)).Msg("Calling GetBlock()")
		}
	}
	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	resp, err := conn.GetBlock(opCtx, &ethpb.BlockRequest{
		Slot:         uint64(slot),
		RandaoReveal: randaoReveal[:],
		Graffiti:     fixedGraffiti,
	})
	cancel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain beacon block data")
	}

	block := &spec.BeaconBlock{
		Slot:          spec.Slot(resp.Slot),
		ProposerIndex: spec.ValidatorIndex(resp.ProposerIndex),
		Body: &spec.BeaconBlockBody{
			ETH1Data: &spec.ETH1Data{
				DepositCount: resp.Body.Eth1Data.DepositCount,
				BlockHash:    resp.Body.Eth1Data.BlockHash,
			},
			Graffiti: fixedGraffiti,
		},
	}
	copy(block.ParentRoot[:], resp.ParentRoot)
	copy(block.StateRoot[:], resp.StateRoot)
	copy(block.Body.RANDAOReveal[:], randaoReveal[:])
	copy(block.Body.ETH1Data.DepositRoot[:], resp.Body.Eth1Data.DepositRoot)
	block.Body.ProposerSlashings = make([]*spec.ProposerSlashing, len(resp.Body.ProposerSlashings))
	for i := range resp.Body.ProposerSlashings {
		block.Body.ProposerSlashings[i] = &spec.ProposerSlashing{
			SignedHeader1: &spec.SignedBeaconBlockHeader{
				Message: &spec.BeaconBlockHeader{
					Slot:          spec.Slot(resp.Body.ProposerSlashings[i].Header_1.Header.Slot),
					ProposerIndex: spec.ValidatorIndex(resp.Body.ProposerSlashings[i].Header_1.Header.ProposerIndex),
				},
			},
			SignedHeader2: &spec.SignedBeaconBlockHeader{
				Message: &spec.BeaconBlockHeader{
					Slot:          spec.Slot(resp.Body.ProposerSlashings[i].Header_2.Header.Slot),
					ProposerIndex: spec.ValidatorIndex(resp.Body.ProposerSlashings[i].Header_2.Header.ProposerIndex),
				},
			},
		}
		copy(block.Body.ProposerSlashings[i].SignedHeader1.Message.ParentRoot[:], resp.Body.ProposerSlashings[i].Header_1.Header.ParentRoot)
		copy(block.Body.ProposerSlashings[i].SignedHeader1.Message.StateRoot[:], resp.Body.ProposerSlashings[i].Header_1.Header.StateRoot)
		copy(block.Body.ProposerSlashings[i].SignedHeader1.Message.BodyRoot[:], resp.Body.ProposerSlashings[i].Header_1.Header.BodyRoot)
		copy(block.Body.ProposerSlashings[i].SignedHeader1.Signature[:], resp.Body.ProposerSlashings[i].Header_1.Signature)
		copy(block.Body.ProposerSlashings[i].SignedHeader2.Message.ParentRoot[:], resp.Body.ProposerSlashings[i].Header_2.Header.ParentRoot)
		copy(block.Body.ProposerSlashings[i].SignedHeader2.Message.StateRoot[:], resp.Body.ProposerSlashings[i].Header_2.Header.StateRoot)
		copy(block.Body.ProposerSlashings[i].SignedHeader2.Message.BodyRoot[:], resp.Body.ProposerSlashings[i].Header_2.Header.BodyRoot)
		copy(block.Body.ProposerSlashings[i].SignedHeader2.Signature[:], resp.Body.ProposerSlashings[i].Header_2.Signature)
	}
	block.Body.AttesterSlashings = make([]*spec.AttesterSlashing, len(resp.Body.AttesterSlashings))
	for i := range resp.Body.AttesterSlashings {
		block.Body.AttesterSlashings[i] = &spec.AttesterSlashing{
			Attestation1: &spec.IndexedAttestation{
				AttestingIndices: resp.Body.AttesterSlashings[i].Attestation_1.AttestingIndices,
				Data: &spec.AttestationData{
					Slot:  spec.Slot(resp.Body.AttesterSlashings[i].Attestation_1.Data.Slot),
					Index: spec.CommitteeIndex(resp.Body.AttesterSlashings[i].Attestation_1.Data.CommitteeIndex),
					Source: &spec.Checkpoint{
						Epoch: spec.Epoch(resp.Body.AttesterSlashings[i].Attestation_1.Data.Source.Epoch),
					},
					Target: &spec.Checkpoint{
						Epoch: spec.Epoch(resp.Body.AttesterSlashings[i].Attestation_1.Data.Target.Epoch),
					},
				},
			},
			Attestation2: &spec.IndexedAttestation{
				AttestingIndices: resp.Body.AttesterSlashings[i].Attestation_2.AttestingIndices,
				Data: &spec.AttestationData{
					Slot:  spec.Slot(resp.Body.AttesterSlashings[i].Attestation_2.Data.Slot),
					Index: spec.CommitteeIndex(resp.Body.AttesterSlashings[i].Attestation_2.Data.CommitteeIndex),
					Source: &spec.Checkpoint{
						Epoch: spec.Epoch(resp.Body.AttesterSlashings[i].Attestation_2.Data.Source.Epoch),
					},
					Target: &spec.Checkpoint{
						Epoch: spec.Epoch(resp.Body.AttesterSlashings[i].Attestation_2.Data.Target.Epoch),
					},
				},
			},
		}
		copy(block.Body.AttesterSlashings[i].Attestation1.Data.BeaconBlockRoot[:], resp.Body.AttesterSlashings[i].Attestation_1.Data.BeaconBlockRoot)
		copy(block.Body.AttesterSlashings[i].Attestation1.Data.Source.Root[:], resp.Body.AttesterSlashings[i].Attestation_1.Data.Source.Root)
		copy(block.Body.AttesterSlashings[i].Attestation1.Data.Target.Root[:], resp.Body.AttesterSlashings[i].Attestation_1.Data.Target.Root)
		copy(block.Body.AttesterSlashings[i].Attestation1.Signature[:], resp.Body.AttesterSlashings[i].Attestation_1.Signature)
		copy(block.Body.AttesterSlashings[i].Attestation2.Data.BeaconBlockRoot[:], resp.Body.AttesterSlashings[i].Attestation_2.Data.BeaconBlockRoot)
		copy(block.Body.AttesterSlashings[i].Attestation2.Data.Source.Root[:], resp.Body.AttesterSlashings[i].Attestation_2.Data.Source.Root)
		copy(block.Body.AttesterSlashings[i].Attestation2.Data.Target.Root[:], resp.Body.AttesterSlashings[i].Attestation_2.Data.Target.Root)
		copy(block.Body.AttesterSlashings[i].Attestation2.Signature[:], resp.Body.AttesterSlashings[i].Attestation_2.Signature)
	}
	block.Body.Attestations = make([]*spec.Attestation, len(resp.Body.Attestations))
	for i := range resp.Body.Attestations {
		block.Body.Attestations[i] = &spec.Attestation{
			AggregationBits: resp.Body.Attestations[i].AggregationBits,
			Data: &spec.AttestationData{
				Slot:  spec.Slot(resp.Body.Attestations[i].Data.Slot),
				Index: spec.CommitteeIndex(resp.Body.Attestations[i].Data.CommitteeIndex),
				Source: &spec.Checkpoint{
					Epoch: spec.Epoch(resp.Body.Attestations[i].Data.Source.Epoch),
				},
				Target: &spec.Checkpoint{
					Epoch: spec.Epoch(resp.Body.Attestations[i].Data.Target.Epoch),
				},
			},
		}
		copy(block.Body.Attestations[i].Data.BeaconBlockRoot[:], resp.Body.Attestations[i].Data.BeaconBlockRoot)
		copy(block.Body.Attestations[i].Data.Source.Root[:], resp.Body.Attestations[i].Data.Source.Root)
		copy(block.Body.Attestations[i].Data.Target.Root[:], resp.Body.Attestations[i].Data.Target.Root)
		copy(block.Body.Attestations[i].Signature[:], resp.Body.Attestations[i].Signature)
	}
	block.Body.Deposits = make([]*spec.Deposit, len(resp.Body.Deposits))
	for i := range resp.Body.Deposits {
		block.Body.Deposits[i] = &spec.Deposit{
			Proof: resp.Body.Deposits[i].Proof,
			Data: &spec.DepositData{
				WithdrawalCredentials: resp.Body.Deposits[i].Data.WithdrawalCredentials,
				Amount:                spec.Gwei(resp.Body.Deposits[i].Data.Amount),
			},
		}
		copy(block.Body.Deposits[i].Data.PublicKey[:], resp.Body.Deposits[i].Data.PublicKey)
		copy(block.Body.Deposits[i].Data.Signature[:], resp.Body.Deposits[i].Data.Signature)
	}
	block.Body.VoluntaryExits = make([]*spec.SignedVoluntaryExit, len(resp.Body.VoluntaryExits))
	for i := range resp.Body.VoluntaryExits {
		block.Body.VoluntaryExits[i] = &spec.SignedVoluntaryExit{
			Message: &spec.VoluntaryExit{
				Epoch:          spec.Epoch(resp.Body.VoluntaryExits[i].Exit.Epoch),
				ValidatorIndex: spec.ValidatorIndex(resp.Body.VoluntaryExits[i].Exit.ValidatorIndex),
			},
		}
		copy(block.Body.VoluntaryExits[i].Signature[:], resp.Body.VoluntaryExits[i].Signature)
	}

	return block, nil
}
