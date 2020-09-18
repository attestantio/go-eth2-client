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

// AttestationData obtains attestation data for a slot.
func (s *Service) AttestationData(ctx context.Context, slot uint64, committeeIndex uint64) (*spec.AttestationData, error) {
	beaconNodeClient := ethpb.NewBeaconNodeValidatorClient(s.conn)
	log.Trace().Msg("Calling GetAttestationData()")
	resp, err := beaconNodeClient.GetAttestationData(ctx, &ethpb.AttestationDataRequest{
		Slot:           slot,
		CommitteeIndex: committeeIndex,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain attestation data")
	}

	if resp.Slot != slot {
		return nil, errors.New("attestation data returned for incorrect slot")
	}
	if resp.CommitteeIndex != committeeIndex {
		return nil, errors.New("attestation data returned for incorrect committee index")
	}

	attestationData := &spec.AttestationData{
		Slot:            resp.Slot,
		Index:           resp.CommitteeIndex,
		BeaconBlockRoot: resp.BeaconBlockRoot,
		Source: &spec.Checkpoint{
			Epoch: resp.Source.Epoch,
			Root:  resp.Source.Root,
		},
		Target: &spec.Checkpoint{
			Epoch: resp.Target.Epoch,
			Root:  resp.Target.Root,
		},
	}

	if e := log.Trace(); e.Enabled() {
		jsonData, err := json.Marshal(attestationData)
		if err == nil {
			log.Trace().Str("attestation_data", string(jsonData)).Msg("Attestation data")
		}
	}

	return attestationData, nil
}
