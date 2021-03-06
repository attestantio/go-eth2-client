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

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// GenesisValidatorsRoot provides the genesis validators root of the chain.
func (s *Service) GenesisValidatorsRoot(ctx context.Context) ([]byte, error) {
	if s.genesisValidatorsRoot == nil {
		conn := ethpb.NewNodeClient(s.conn)
		opCtx, cancel := context.WithTimeout(ctx, s.timeout)
		res, err := conn.GetGenesis(opCtx, &types.Empty{})
		cancel()
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain genesis validators root")
		}
		s.genesisValidatorsRoot = res.GenesisValidatorsRoot
	}
	return s.genesisValidatorsRoot, nil
}
