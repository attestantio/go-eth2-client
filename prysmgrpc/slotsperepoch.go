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
	"strconv"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// SlotsPerEpoch provides the number of slots per epoch of the chain.
func (s *Service) SlotsPerEpoch(ctx context.Context) (uint64, error) {
	if s.slotsPerEpoch == nil {
		conn := ethpb.NewBeaconChainClient(s.conn)
		opCtx, cancel := context.WithTimeout(ctx, s.timeout)
		config, err := conn.GetBeaconConfig(opCtx, &types.Empty{})
		cancel()
		if err != nil {
			return 0, errors.Wrap(err, "failed to obtain configuration")
		}

		val, exists := config.Config["SlotsPerEpoch"]
		if !exists {
			return 0, errors.New("failed to obtain SlotsPerEpoch value")
		}
		uintVal, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse SlotsPerEpoch value")
		}
		s.slotsPerEpoch = &uintVal
	}
	return *s.slotsPerEpoch, nil
}
