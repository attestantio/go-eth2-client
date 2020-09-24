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

package lighthousehttp

import (
	"context"
	"encoding/binary"
	"encoding/json"

	"github.com/pkg/errors"
)

// BeaconProposerDomain provides the beacon proposer domain of the chain.
func (s *Service) BeaconProposerDomain(ctx context.Context) ([]byte, error) {
	if s.beaconProposerDomain == nil {
		respBodyReader, cancel, err := s.get(ctx, "/spec")
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain configuration")
		}
		defer cancel()

		var cfg map[string]interface{}
		if err := json.NewDecoder(respBodyReader).Decode(&cfg); err != nil {
			return nil, errors.Wrap(err, "failed to parse configuration")
		}
		val, exists := cfg["domain_beacon_proposer"]
		if !exists {
			return nil, errors.New("failed to obtain domain_beacon_proposer value")
		}
		floatVal, isNum := val.(float64)
		if !isNum {
			return nil, errors.New("domain_beacon_proposer value not a number")
		}
		intVal := uint32(floatVal)
		var domain [4]byte
		binary.LittleEndian.PutUint32(domain[:], intVal)
		s.beaconProposerDomain = domain[:]
	}
	return s.beaconProposerDomain, nil
}
