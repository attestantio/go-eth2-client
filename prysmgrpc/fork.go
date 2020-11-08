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
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// Fork provides the fork at a given epoch.
// Prysm does not provide a method to obtain the current fork version, so provide the genesis fork version.
func (s *Service) Fork(ctx context.Context, stateID string) (*spec.Fork, error) {
	if s.genesisForkVersion == nil {
		conn := ethpb.NewBeaconChainClient(s.conn)
		opCtx, cancel := context.WithTimeout(ctx, s.timeout)
		config, err := conn.GetBeaconConfig(opCtx, &types.Empty{})
		cancel()
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain configuration")
		}

		val, exists := config.Config["GenesisForkVersion"]
		if !exists {
			return nil, errors.New("config did not provide GenesisForkVersion value")
		}
		s.genesisForkVersion, err = parseConfigByteArray(val)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert value %q for GenesisForkVersion", val)
		}
	}

	fork := &spec.Fork{
		Epoch: 0,
	}
	copy(fork.CurrentVersion[:], s.genesisForkVersion)
	copy(fork.PreviousVersion[:], s.genesisForkVersion)

	return fork, nil
}
