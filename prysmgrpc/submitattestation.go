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

// SubmitAttestation submits an attestation.
func (s *Service) SubmitAttestation(ctx context.Context, attestation *spec.Attestation) error {
	prysmAttestation := &ethpb.Attestation{
		AggregationBits: attestation.AggregationBits,
		Data: &ethpb.AttestationData{
			Slot:            attestation.Data.Slot,
			CommitteeIndex:  attestation.Data.Index,
			BeaconBlockRoot: attestation.Data.BeaconBlockRoot,
			Source: &ethpb.Checkpoint{
				Epoch: attestation.Data.Source.Epoch,
				Root:  attestation.Data.Source.Root,
			},
			Target: &ethpb.Checkpoint{
				Epoch: attestation.Data.Target.Epoch,
				Root:  attestation.Data.Target.Root,
			},
		},
		Signature: attestation.Signature,
	}

	client := ethpb.NewBeaconNodeValidatorClient(s.conn)
	log.Trace().Msg("Calling ProposeAttestation()")
	_, err := client.ProposeAttestation(ctx, prysmAttestation)
	if err != nil {
		return errors.Wrap(err, "failed to submit attestation")
	}

	return nil
}
