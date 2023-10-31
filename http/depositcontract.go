// Copyright © 2020 - 2023 Attestant Limited.
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

package http

import (
	"bytes"
	"context"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
)

// DepositContract provides details of the execution deposit contract for the chain.
func (s *Service) DepositContract(ctx context.Context,
	opts *api.DepositContractOpts,
) (
	*api.Response[*apiv1.DepositContract],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}

	s.depositContractMutex.RLock()
	if s.depositContract != nil {
		defer s.depositContractMutex.RUnlock()

		return &api.Response[*apiv1.DepositContract]{
			Data:     s.depositContract,
			Metadata: map[string]any{},
		}, nil
	}
	s.depositContractMutex.RUnlock()

	s.depositContractMutex.Lock()
	defer s.depositContractMutex.Unlock()
	if s.depositContract != nil {
		// Someone else fetched this whilst we were waiting for the lock.
		return &api.Response[*apiv1.DepositContract]{
			Data:     s.depositContract,
			Metadata: map[string]any{},
		}, nil
	}

	// Up to us to fetch the information.
	url := "/eth/v1/config/deposit_contract"
	httpResponse, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), apiv1.DepositContract{})
	if err != nil {
		return nil, err
	}

	s.depositContract = &data

	return &api.Response[*apiv1.DepositContract]{
		Data:     s.depositContract,
		Metadata: metadata,
	}, nil
}
