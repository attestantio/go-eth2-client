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
	"fmt"

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
	if len(validators) == 0 {
		return s.validators(ctx, stateID, true)
	}
	return s.validatorsByPubKeys(ctx, stateID, validators, true)
}

// ValidatorsWithoutBalance provides the validators, with their status, for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validators is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
func (s *Service) ValidatorsWithoutBalance(ctx context.Context, stateID string, validators []client.ValidatorIDProvider) (map[uint64]*api.Validator, error) {
	if len(validators) == 0 {
		return s.validators(ctx, stateID, false)
	}
	return s.validatorsByPubKeys(ctx, stateID, validators, false)
}

// validators returns all validators known by the client.
func (s *Service) validators(ctx context.Context, stateID string, includeBalances bool) (map[uint64]*api.Validator, error) {
	// The state ID could by dynamic ('head', 'finalized', etc.).  Becase we are making multiple calls and don't want to
	// fetch data from different states we resolve it to an epoch and use that.
	epoch, err := s.EpochFromStateID(ctx, stateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lock state ID")
	}

	conn := ethpb.NewBeaconChainClient(s.conn)
	if conn == nil {
		return nil, errors.New("failed to obtain beacon chain client")
	}

	validatorsReq := &ethpb.ListValidatorsRequest{}
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
		opCtx, cancel := context.WithTimeout(ctx, s.timeout)
		validatorsResp, err := conn.ListValidators(opCtx, validatorsReq)
		cancel()
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

		highest := uint64(0)
		for i := range res {
			if res[i].Index > highest {
				highest = res[i].Index
			}
		}

		if validatorsResp.NextPageToken == "" {
			// Means we're done.
			break
		}
		pageToken = validatorsResp.NextPageToken
	}

	if !includeBalances {
		return res, nil
	}

	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain slots per epoch")
	}
	balances, err := s.validatorBalances(ctx, fmt.Sprintf("%d", epoch*slotsPerEpoch))
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain validator balances")
	}
	for index, balance := range balances {
		if _, exists := res[index]; exists {
			res[index].Balance = balance
		}
	}

	return res, nil
}

// validatorsByPubKeys returns a subset of validators.
func (s *Service) validatorsByPubKeys(ctx context.Context, stateID string, validators []client.ValidatorIDProvider, includeBalances bool) (map[uint64]*api.Validator, error) {
	// The state ID could by dynamic ('head', 'finalized', etc.).  Becase we are making multiple calls and don't want to
	// fetch data from different states we resolve it to an epoch and use that.
	epoch, err := s.EpochFromStateID(ctx, stateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to lock state ID")
	}

	conn := ethpb.NewBeaconChainClient(s.conn)
	if conn == nil {
		return nil, errors.New("failed to obtain beacon chain client")
	}

	validatorsReq := &ethpb.ListValidatorsRequest{
		PageSize: s.maxPageSize,
	}
	validatorBalancesReq := &ethpb.ListValidatorBalancesRequest{
		PageSize: s.maxPageSize,
	}
	if epoch == 0 {
		log.Trace().Msg("Fetching genesis validators")
		validatorsReq.QueryFilter = &ethpb.ListValidatorsRequest_Genesis{Genesis: true}
		validatorBalancesReq.QueryFilter = &ethpb.ListValidatorBalancesRequest_Genesis{Genesis: true}
	} else {
		log.Trace().Uint64("epoch", epoch).Msg("Fetching epoch validators")
		validatorsReq.QueryFilter = &ethpb.ListValidatorsRequest_Epoch{Epoch: epoch}
		validatorBalancesReq.QueryFilter = &ethpb.ListValidatorBalancesRequest_Epoch{Epoch: epoch}
	}

	farFutureEpoch, err := s.FarFutureEpoch(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain far future epoch")
	}

	pubKeys := make([][]byte, 0, len(validators))
	for i := range validators {
		pubKey, err := validators[i].PubKey(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to obtain public key for validator")
		}
		pubKeys = append(pubKeys, pubKey)
	}

	// If we ask prysm for the balance of a validator that doesn't exist it errors, so
	// keep track of the validators that are known to prysm.
	known := make(map[[48]byte]bool)

	res := make(map[uint64]*api.Validator)
	for i := 0; i < len(pubKeys); i += int(s.maxPageSize) {
		lastIndex := i + int(s.maxPageSize)
		if lastIndex > len(pubKeys) {
			lastIndex = len(pubKeys)
		}

		validatorsReq.PublicKeys = pubKeys[i:lastIndex]

		log.Trace().Msg("Calling ListValidators()")
		opCtx, cancel := context.WithTimeout(ctx, s.timeout)
		validatorsResp, err := conn.ListValidators(opCtx, validatorsReq)
		cancel()
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

			// Add validator to known list.
			var pubKey [48]byte
			copy(pubKey[:], entry.Validator.PublicKey)
			known[pubKey] = true
		}
	}

	if !includeBalances {
		return res, nil
	}

	// If we ask prysm for the balance of a validator that doesn't exist it errors, so
	// reduce our validators to those that Prysm recognises.
	balancePubKeys := make([][]byte, 0, len(known))
	for k := range known {
		balancePubKeys = append(balancePubKeys, k[:])
	}

	// Fetch the balances
	for i := 0; i < len(balancePubKeys); i += int(s.maxPageSize) {
		lastIndex := i + int(s.maxPageSize)
		if lastIndex > len(balancePubKeys) {
			lastIndex = len(balancePubKeys)
		}

		validatorBalancesReq.PublicKeys = balancePubKeys[i:lastIndex]

		log.Trace().Msg("Calling ListValidatorBalances()")
		opCtx, cancel := context.WithTimeout(ctx, s.timeout)
		validatorBalancesResp, err := conn.ListValidatorBalances(opCtx, validatorBalancesReq)
		cancel()
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain validator balances")
		}
		if len(validatorBalancesResp.Balances) == 0 {
			break
		}
		for _, entry := range validatorBalancesResp.Balances {
			res[entry.Index].Balance = entry.Balance
		}
	}

	return res, nil
}
