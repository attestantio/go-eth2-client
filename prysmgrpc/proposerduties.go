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

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// ProposerDuties obtains proposer duties.
func (s *Service) ProposerDuties(ctx context.Context, epoch uint64, validators []client.ValidatorIDProvider) ([]*api.ProposerDuty, error) {
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
	log.Trace().Msg("Calling GetDuties()")
	resp, err := conn.GetDuties(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "call to GetDuties() failed")
	}

	proposerDuties := make([]*api.ProposerDuty, 0)
	for _, duty := range resp.CurrentEpochDuties {
		for _, slot := range duty.ProposerSlots {
			log.Trace().Uint64("slot", slot).Uint64("validator_index", duty.ValidatorIndex).Msg("Received proposer duty")
			proposerDuties = append(proposerDuties, &api.ProposerDuty{
				Slot:           slot,
				ValidatorIndex: duty.ValidatorIndex,
			})
		}
	}

	return proposerDuties, nil
}
