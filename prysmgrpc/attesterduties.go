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

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// AttesterDuties obtains attester duties.
func (s *Service) AttesterDuties(ctx context.Context, epoch uint64, validators []client.ValidatorIDProvider) ([]*api.AttesterDuty, error) {
	conn := ethpb.NewBeaconNodeValidatorClient(s.conn)

	pubKeys := make([][]byte, 0, len(validators))
	for i := range validators {
		pubKey, err := validators[i].PubKey(ctx)
		if err != nil {
			// Warn but do not exit as we want to obtain as many proposers as possible.
			log.Warn().Err(err).Msg("Failed to obtain public key for validator")
			continue
		}
		pubKeys = append(pubKeys, pubKey)
	}
	req := &ethpb.DutiesRequest{
		Epoch:      epoch,
		PublicKeys: pubKeys,
	}
	if e := log.Trace(); e.Enabled() {
		jsonData, err := json.Marshal(req)
		if err == nil {
			log.Trace().Str("req", string(jsonData)).Msg("Calling GetDuties()")
		}
	}
	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	resp, err := conn.GetDuties(opCtx, req)
	cancel()
	if err != nil {
		return nil, errors.Wrap(err, "call to GetDuties() failed")
	}

	duties := make([]*api.AttesterDuty, 0, len(resp.CurrentEpochDuties))
	for _, duty := range resp.CurrentEpochDuties {
		validatorCommitteeIndex := 0
		for i := range duty.Committee {
			if duty.Committee[i] == duty.ValidatorIndex {
				validatorCommitteeIndex = i
				break
			}
		}
		duties = append(duties, &api.AttesterDuty{
			Slot:                    duty.AttesterSlot,
			ValidatorIndex:          duty.ValidatorIndex,
			CommitteeIndex:          duty.CommitteeIndex,
			ValidatorCommitteeIndex: uint64(validatorCommitteeIndex),
			CommitteeLength:         uint64(len(duty.Committee)),
		})
	}

	if e := log.Trace(); e.Enabled() {
		jsonData, err := json.Marshal(duties)
		if err == nil {
			log.Trace().Str("data", string(jsonData)).Msg("Returning attester duties")
		}
	}
	return duties, nil
}
