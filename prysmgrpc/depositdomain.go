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

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// DepositDomain provides the deposit domain of the chain.
func (s *Service) DepositDomain(ctx context.Context) ([]byte, error) {
	if s.depositDomain == nil {
		client := ethpb.NewBeaconChainClient(s.conn)
		log.Trace().Msg("Fetching deposit domain")
		config, err := client.GetBeaconConfig(ctx, &types.Empty{})
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain configuration")
		}

		val, exists := config.Config["DomainDeposit"]
		if !exists {
			return nil, errors.New("config did not provide DomainDeposit value")
		}
		s.depositDomain, err = parseConfigByteArray(val)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert value %q for DomainDeposit", val)
		}
	}

	log.Trace().Str("domain", fmt.Sprintf("%#x", s.beaconAttesterDomain)).Msg("Returning deposit domain")
	return s.depositDomain, nil
}
