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

// SubmitAttestations submits attestations.
// Prysm does not provide the ability to submit attestations natively, so send individually.
func (s *Service) SubmitAttestations(ctx context.Context, attestations []*spec.Attestation) error {
	var anyErr error
	for i := range attestations {
		if err := s.submitAttestation(ctx, attestations[i]); err != nil {
			anyErr = err
		}
	}

	return anyErr
}

// submitAttestation submits an attestation.
func (s *Service) submitAttestation(ctx context.Context, attestation *spec.Attestation) error {
	prysmAttestation := &ethpb.Attestation{
		AggregationBits: attestation.AggregationBits,
		Data: &ethpb.AttestationData{
			Slot:            uint64(attestation.Data.Slot),
			CommitteeIndex:  uint64(attestation.Data.Index),
			BeaconBlockRoot: attestation.Data.BeaconBlockRoot[:],
			Source: &ethpb.Checkpoint{
				Epoch: uint64(attestation.Data.Source.Epoch),
				Root:  attestation.Data.Source.Root[:],
			},
			Target: &ethpb.Checkpoint{
				Epoch: uint64(attestation.Data.Target.Epoch),
				Root:  attestation.Data.Target.Root[:],
			},
		},
		Signature: attestation.Signature[:],
	}

	conn := ethpb.NewBeaconNodeValidatorClient(s.conn)
	log.Trace().Msg("Calling ProposeAttestation()")
	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	_, err := conn.ProposeAttestation(opCtx, prysmAttestation)
	cancel()
	if err != nil {
		return errors.Wrap(err, "failed to submit attestation")
	}

	return nil
}
