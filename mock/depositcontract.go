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

package mock

import (
	"context"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// DepositContract provides details of the execution layer deposit contract for the chain.
func (s *Service) DepositContract(_ context.Context) (*api.Response[*apiv1.DepositContract], error) {
	return &api.Response[*apiv1.DepositContract]{
		Data:     &apiv1.DepositContract{},
		Metadata: make(map[string]any),
	}, nil
}
