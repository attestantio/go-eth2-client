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

	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// ValidatorBalancesByState provides all validator balancess for a given state.
// State can be a slot number or state root, or one of the special values "genesis", "head", "justified" or "finalized".
func (s *Service) ValidatorBalancesByState(ctx context.Context, stateID string) (map[uint64]uint64, error) {
	beaconChainClient := ethpb.NewBeaconChainClient(s.conn)
	if beaconChainClient == nil {
		return nil, errors.New("failed to obtain beacon chain client")
	}

	validatorBalancesReq := &ethpb.ListValidatorBalancesRequest{}

	_, _, slot, err := s.parseStateID(ctx, stateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain slot")
	}
	epoch := slot / (*s.slotsPerEpoch)
	if epoch == 0 {
		log.Trace().Msg("Fetching genesis validator balances")
		validatorBalancesReq.QueryFilter = &ethpb.ListValidatorBalancesRequest_Genesis{Genesis: true}
	} else {
		log.Trace().Uint64("epoch", epoch).Msg("Fetching epoch validator balances")
		validatorBalancesReq.QueryFilter = &ethpb.ListValidatorBalancesRequest_Epoch{Epoch: epoch}
	}

	// Start the results arrays with the minimum expected number of validators.
	balances := make(map[uint64]uint64, 16384)

	pageToken := ""
	for i := int32(0); ; i += s.maxPageSize {
		log.Trace().Msg("Calling ListValidatorBalances()")
		validatorBalancesReq.PageToken = pageToken
		validatorBalancesReq.PageSize = s.maxPageSize
		resp, err := beaconChainClient.ListValidatorBalances(ctx, validatorBalancesReq)
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain validator balances")
		}
		if len(resp.Balances) == 0 {
			break
		}

		for _, entry := range resp.Balances {
			balances[entry.Index] = entry.Balance
		}

		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	return balances, nil
}
