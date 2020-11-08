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

// PrysmAggregateAttestation fetches the aggregate attestation given an attestation.
func (s *Service) PrysmAggregateAttestation(ctx context.Context, attestation *spec.Attestation, validatorPubKey spec.BLSPubKey, slotSignature spec.BLSSignature) (*spec.Attestation, error) {
	conn := ethpb.NewBeaconNodeValidatorClient(s.conn)
	log.Trace().Msg("Calling SubmitAggregateSelectionProof()")
	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	resp, err := conn.SubmitAggregateSelectionProof(opCtx, &ethpb.AggregateSelectionRequest{
		Slot:           uint64(attestation.Data.Slot),
		CommitteeIndex: uint64(attestation.Data.Index),
		PublicKey:      validatorPubKey[:],
		SlotSignature:  slotSignature[:],
	})
	cancel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain attestation data")
	}

	aggregateAttestation := &spec.Attestation{
		AggregationBits: resp.AggregateAndProof.Aggregate.AggregationBits,
		Data: &spec.AttestationData{
			Slot:  spec.Slot(resp.AggregateAndProof.Aggregate.Data.Slot),
			Index: spec.CommitteeIndex(resp.AggregateAndProof.Aggregate.Data.CommitteeIndex),
			Source: &spec.Checkpoint{
				Epoch: spec.Epoch(resp.AggregateAndProof.Aggregate.Data.Source.Epoch),
			},
			Target: &spec.Checkpoint{
				Epoch: spec.Epoch(resp.AggregateAndProof.Aggregate.Data.Target.Epoch),
			},
		},
	}
	copy(aggregateAttestation.Data.BeaconBlockRoot[:], resp.AggregateAndProof.Aggregate.Data.BeaconBlockRoot)
	copy(aggregateAttestation.Data.Source.Root[:], resp.AggregateAndProof.Aggregate.Data.Source.Root)
	copy(aggregateAttestation.Data.Target.Root[:], resp.AggregateAndProof.Aggregate.Data.Target.Root)
	copy(aggregateAttestation.Signature[:], resp.AggregateAndProof.Aggregate.Signature)

	if spec.Slot(resp.AggregateAndProof.Aggregate.Data.Slot) != attestation.Data.Slot {
		return nil, errors.New("aggregate attestation data returned for incorrect slot")
	}
	if spec.CommitteeIndex(resp.AggregateAndProof.Aggregate.Data.CommitteeIndex) != attestation.Data.Index {
		return nil, errors.New("aggregate attestation data returned for incorrect committee index")
	}

	if e := log.Trace(); e.Enabled() {
		jsonData, err := json.Marshal(aggregateAttestation)
		if err == nil {
			log.Trace().Str("data", string(jsonData)).Msg("Returning aggregate attestation")
		}
	}
	return aggregateAttestation, nil
}
