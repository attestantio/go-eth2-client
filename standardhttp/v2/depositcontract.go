// Copyright © 2021 Attestant Limited.
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

package v2

import (
	"context"
	"encoding/json"

	api "github.com/attestantio/go-eth2-client/api/v2"
	"github.com/pkg/errors"
)

type depositContractJSON struct {
	Data *api.DepositContract `json:"data"`
}

// DepositContract provides details of the Ethereum 1 deposit contract for the chain.
func (s *Service) DepositContract(ctx context.Context) (*api.DepositContract, error) {
	if s.depositContract != nil {
		return s.depositContract, nil
	}

	s.depositContractMutex.Lock()
	defer s.depositContractMutex.Unlock()
	if s.depositContract != nil {
		// Someone else fetched this whilst we were waiting for the lock.
		return s.depositContract, nil
	}

	// Up to us to fetch the information.
	respBodyReader, err := s.get(ctx, "/eth/v1/config/deposit_contract")
	if err != nil {
		return nil, errors.Wrap(err, "failed to request deposit contract")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain deposit contract")
	}

	var resp depositContractJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse deposit contract")
	}
	s.depositContract = resp.Data
	return s.depositContract, nil
}
