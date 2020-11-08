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

// PrysmValidatorBalances provides the validator balances for a given state.
// stateID can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
// validators is a list of validators to restrict the returned values.  If no validators are supplied no filter will be applied.
func (s *Service) PrysmValidatorBalances(ctx context.Context, stateID string, validatorPubKeys []spec.BLSPubKey) (map[spec.ValidatorIndex]spec.Gwei, error) {
	if len(validatorPubKeys) == 0 {
		return s.validatorBalances(ctx, stateID)
	}
	return s.validatorBalancesByPubKeys(ctx, stateID, validatorPubKeys)
}

func (s *Service) validatorBalances(ctx context.Context, stateID string) (map[spec.ValidatorIndex]spec.Gwei, error) {
	conn := ethpb.NewBeaconChainClient(s.conn)
	if conn == nil {
		return nil, errors.New("failed to obtain beacon chain client")
	}

	validatorBalancesReq := &ethpb.ListValidatorBalancesRequest{
		PageSize: s.maxPageSize,
	}

	epoch, err := s.EpochFromStateID(ctx, stateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain epoch from state ID")
	}
	if epoch == 0 {
		log.Trace().Msg("Fetching genesis validator balances")
		validatorBalancesReq.QueryFilter = &ethpb.ListValidatorBalancesRequest_Genesis{Genesis: true}
	} else {
		log.Trace().Uint64("epoch", uint64(epoch)).Msg("Fetching epoch validator balances")
		validatorBalancesReq.QueryFilter = &ethpb.ListValidatorBalancesRequest_Epoch{Epoch: uint64(epoch)}
	}

	res := make(map[spec.ValidatorIndex]spec.Gwei)

	pageToken := ""
	for i := int32(0); ; i += s.maxPageSize {
		log.Trace().Msg("Calling ListValidators()")
		validatorBalancesReq.PageToken = pageToken
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
			res[spec.ValidatorIndex(entry.Index)] = spec.Gwei(entry.Balance)
		}

		if validatorBalancesResp.NextPageToken == "" {
			break
		}
		pageToken = validatorBalancesResp.NextPageToken
	}

	return res, nil
}

// validatorsByPubKeys returns a subset of validator balances.
func (s *Service) validatorBalancesByPubKeys(ctx context.Context, stateID string, validatorPubKeys []spec.BLSPubKey) (map[spec.ValidatorIndex]spec.Gwei, error) {
	conn := ethpb.NewBeaconChainClient(s.conn)
	if conn == nil {
		return nil, errors.New("failed to obtain beacon chain client")
	}

	validatorBalancesReq := &ethpb.ListValidatorBalancesRequest{
		PageSize: s.maxPageSize,
	}

	epoch, err := s.EpochFromStateID(ctx, stateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain epoch from state ID")
	}
	if epoch == 0 {
		log.Trace().Msg("Fetching genesis validator balances")
		validatorBalancesReq.QueryFilter = &ethpb.ListValidatorBalancesRequest_Genesis{Genesis: true}
	} else {
		log.Trace().Uint64("epoch", uint64(epoch)).Msg("Fetching epoch validator balances")
		validatorBalancesReq.QueryFilter = &ethpb.ListValidatorBalancesRequest_Epoch{Epoch: uint64(epoch)}
	}

	pubKeys := make([][]byte, len(validatorPubKeys))
	for i := range validatorPubKeys {
		pubKeys[i] = validatorPubKeys[i][:]
	}

	res := make(map[spec.ValidatorIndex]spec.Gwei)
	for i := 0; i < len(validatorPubKeys); i += int(s.maxPageSize) {
		lastIndex := i + int(s.maxPageSize)
		if lastIndex > len(validatorPubKeys) {
			lastIndex = len(validatorPubKeys)
		}
		validatorBalancesReq.PublicKeys = pubKeys[i:lastIndex]

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
			res[spec.ValidatorIndex(entry.Index)] = spec.Gwei(entry.Balance)
		}
	}

	return res, nil
}
