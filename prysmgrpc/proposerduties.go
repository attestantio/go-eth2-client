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

	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// ProposerDuties obtains proposer duties.
func (s *Service) ProposerDuties(ctx context.Context, epoch spec.Epoch, indices []spec.ValidatorIndex) ([]*api.ProposerDuty, error) {
	conn := ethpb.NewBeaconNodeValidatorClient(s.conn)

	pubKeys := make([][]byte, 0, len(indices))
	if len(indices) == 0 {
		// Prysm requires we send it a list of validators, so fetch them.
		prysmValidators, err := s.Validators(ctx, "head", nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain validators")
		}
		log.Trace().Int("validators", len(prysmValidators)).Msg("Obtained validators")

		for _, prysmValidator := range prysmValidators {
			pubKeys = append(pubKeys, prysmValidator.Validator.PublicKey[:])
		}
	} else {
		// Convert provided indices to pubkeys.
		validatorPubKeys, err := s.indicesToPubKeys(ctx, indices)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert indices to public keys")
		}
		for i := range validatorPubKeys {
			pubKeys = append(pubKeys, validatorPubKeys[i][:])
		}
	}

	req := &ethpb.DutiesRequest{
		Epoch:      uint64(epoch),
		PublicKeys: pubKeys,
	}
	log.Trace().Msg("Calling GetDuties()")

	opCtx, cancel := context.WithTimeout(ctx, s.timeout)
	resp, err := conn.GetDuties(opCtx, req)
	cancel()
	if err != nil {
		return nil, errors.Wrap(err, "call to GetDuties() failed")
	}

	proposerDuties := make([]*api.ProposerDuty, 0)
	index := 0
	for _, duty := range resp.CurrentEpochDuties {
		for _, slot := range duty.ProposerSlots {
			log.Trace().Uint64("slot", slot).Uint64("validator_index", duty.ValidatorIndex).Msg("Received proposer duty")
			proposerDuties = append(proposerDuties, &api.ProposerDuty{
				Slot:           spec.Slot(slot),
				ValidatorIndex: spec.ValidatorIndex(duty.ValidatorIndex),
			})
			copy(proposerDuties[index].PubKey[:], duty.PublicKey)
			index++
		}
	}

	return proposerDuties, nil
}
