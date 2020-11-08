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

// SubmitAggregateAttestations submits aggregate attestations.
func (s *Service) SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*spec.SignedAggregateAndProof) error {
	conn := ethpb.NewBeaconNodeValidatorClient(s.conn)
	for _, aggregateAndProof := range aggregateAndProofs {
		prysmAggregateAndProof := &ethpb.SignedAggregateAttestationAndProof{
			Message: &ethpb.AggregateAttestationAndProof{
				AggregatorIndex: uint64(aggregateAndProof.Message.AggregatorIndex),
				Aggregate: &ethpb.Attestation{
					AggregationBits: aggregateAndProof.Message.Aggregate.AggregationBits,
					Data: &ethpb.AttestationData{
						Slot:            uint64(aggregateAndProof.Message.Aggregate.Data.Slot),
						CommitteeIndex:  uint64(aggregateAndProof.Message.Aggregate.Data.Index),
						BeaconBlockRoot: aggregateAndProof.Message.Aggregate.Data.BeaconBlockRoot[:],
						Source: &ethpb.Checkpoint{
							Epoch: uint64(aggregateAndProof.Message.Aggregate.Data.Source.Epoch),
							Root:  aggregateAndProof.Message.Aggregate.Data.Source.Root[:],
						},
						Target: &ethpb.Checkpoint{
							Epoch: uint64(aggregateAndProof.Message.Aggregate.Data.Target.Epoch),
							Root:  aggregateAndProof.Message.Aggregate.Data.Target.Root[:],
						},
					},
					Signature: aggregateAndProof.Message.Aggregate.Signature[:],
				},
				SelectionProof: aggregateAndProof.Message.SelectionProof[:],
			},
			Signature: aggregateAndProof.Signature[:],
		}

		log.Trace().Msg("Calling ProposeSignedAggregateSelectionProof()")
		opCtx, cancel := context.WithTimeout(ctx, s.timeout)
		_, err := conn.SubmitSignedAggregateSelectionProof(opCtx, &ethpb.SignedAggregateSubmitRequest{
			SignedAggregateAndProof: prysmAggregateAndProof,
		})
		cancel()
		if err != nil {
			return errors.Wrap(err, "failed to submit signed aggregate attestation")
		}
	}

	return nil
}
