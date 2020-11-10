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

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// BeaconProposerDomain provides the beacon proposer domain of the chain.
func (s *Service) BeaconProposerDomain(ctx context.Context) (spec.DomainType, error) {
	if s.beaconProposerDomain == nil {
		conn := ethpb.NewBeaconChainClient(s.conn)
		log.Trace().Msg("Fetching beacon proposer domain")
		opCtx, cancel := context.WithTimeout(ctx, s.timeout)
		config, err := conn.GetBeaconConfig(opCtx, &types.Empty{})
		cancel()
		if err != nil {
			return spec.DomainType{}, errors.Wrap(err, "failed to obtain configuration")
		}

		val, exists := config.Config["DomainBeaconProposer"]
		if !exists {
			return spec.DomainType{}, errors.New("config did not provide DomainBeaconProposer value")
		}
		tmp, err := parseConfigByteArray(val)
		if err != nil {
			return spec.DomainType{}, errors.Wrapf(err, "failed to convert value %q for DomainBeaconProposer", val)
		}
		var domainType spec.DomainType
		copy(domainType[:], tmp)
		s.beaconProposerDomain = &domainType
	}

	log.Trace().Str("domain", fmt.Sprintf("%#x", s.beaconAttesterDomain)).Msg("Returning beacon proposer domain")
	return *s.beaconProposerDomain, nil
}
