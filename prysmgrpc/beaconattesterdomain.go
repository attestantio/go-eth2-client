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

// BeaconAttesterDomain provides the beacon attester domain of the chain.
func (s *Service) BeaconAttesterDomain(ctx context.Context) ([]byte, error) {
	if s.beaconAttesterDomain == nil {
		conn := ethpb.NewBeaconChainClient(s.conn)
		log.Trace().Msg("Fetching beacon attester domain")
		opCtx, cancel := context.WithTimeout(ctx, s.timeout)
		config, err := conn.GetBeaconConfig(opCtx, &types.Empty{})
		cancel()
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain configuration")
		}

		val, exists := config.Config["DomainBeaconAttester"]
		if !exists {
			return nil, errors.New("config did not provide DomainBeaconAttester value")
		}
		s.beaconAttesterDomain, err = parseConfigByteArray(val)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert value %q for DomainBeaconAttester", val)
		}
	}
	log.Trace().Str("domain", fmt.Sprintf("%#x", s.beaconAttesterDomain)).Msg("Returning beacon attester domain")
	return s.beaconAttesterDomain, nil
}
