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
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// Validators provides the validators, with their balance and status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validators is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
func (s *Service) Validators(ctx context.Context, stateID string, validators []client.ValidatorIDProvider) (map[uint64]*api.Validator, error) {
	beaconChainClient := ethpb.NewBeaconChainClient(s.conn)
	if beaconChainClient == nil {
		return nil, errors.New("failed to obtain beacon chain client")
	}

	validatorsReq := &ethpb.ListValidatorsRequest{}

	epoch, err := s.epochFromStateID(ctx, stateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain epoch from state ID")
	}
	if epoch == 0 {
		log.Trace().Msg("Fetching genesis validators")
		validatorsReq.QueryFilter = &ethpb.ListValidatorsRequest_Genesis{Genesis: true}
	} else {
		log.Trace().Uint64("epoch", epoch).Msg("Fetching epoch validators")
		validatorsReq.QueryFilter = &ethpb.ListValidatorsRequest_Epoch{Epoch: epoch}
	}
	farFutureEpoch, err := s.FarFutureEpoch(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain far future epoch")
	}

	res := make(map[uint64]*api.Validator)

	pageToken := ""
	for i := int32(0); ; i += s.maxPageSize {
		log.Trace().Msg("Calling ListValidators()")
		validatorsReq.PageToken = pageToken
		validatorsReq.PageSize = s.maxPageSize
		validatorsResp, err := beaconChainClient.ListValidators(ctx, validatorsReq)
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain validators")
		}
		if len(validatorsResp.ValidatorList) == 0 {
			break
		}

		for _, entry := range validatorsResp.ValidatorList {
			validator := &spec.Validator{
				PublicKey:                  entry.Validator.PublicKey,
				WithdrawalCredentials:      entry.Validator.WithdrawalCredentials,
				EffectiveBalance:           entry.Validator.EffectiveBalance,
				Slashed:                    entry.Validator.Slashed,
				ActivationEligibilityEpoch: entry.Validator.ActivationEligibilityEpoch,
				ActivationEpoch:            entry.Validator.ActivationEpoch,
				ExitEpoch:                  entry.Validator.ExitEpoch,
				WithdrawableEpoch:          entry.Validator.WithdrawableEpoch,
			}
			res[entry.Index] = &api.Validator{
				Index:     entry.Index,
				State:     api.ValidatorToState(validator, epoch, farFutureEpoch),
				Validator: validator,
			}
		}

		if validatorsResp.NextPageToken == "" {
			break
		}
		pageToken = validatorsResp.NextPageToken
	}

	balances, err := s.ValidatorBalances(ctx, stateID, validators)
	if err != nil {
		return nil, errors.New("failed to obtain validator balances")
	}
	for index, balance := range balances {
		res[index].Balance = balance
	}

	return res, nil
}
