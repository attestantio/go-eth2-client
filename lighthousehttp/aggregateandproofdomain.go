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

// AggregateAndProofDomain provides the aggregate and proof domain of the chain.
func (s *Service) AggregateAndProofDomain(ctx context.Context) ([]byte, error) {
	if s.aggregateAndProofDomain == nil {
		respBodyReader, err := s.get(ctx, "/spec")
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain configuration")
		}
		defer func() {
			if err := respBodyReader.Close(); err != nil {
				log.Warn().Err(err).Msg("Failed to close HTTP body")
			}
		}()
		var cfg map[string]interface{}
		if err := json.NewDecoder(respBodyReader).Decode(&cfg); err != nil {
			return nil, errors.Wrap(err, "failed to parse configuration")
		}
		val, exists := cfg["domain_aggregate_and_proof"]
		if !exists {
			return nil, errors.New("failed to obtain domain_aggregate_and_proof value")
		}
		floatVal, isNum := val.(float64)
		if !isNum {
			return nil, errors.New("domain_aggregate_and_proof value not a number")
		}
		intVal := uint32(floatVal)
		var domain [4]byte
		binary.LittleEndian.PutUint32(domain[:], intVal)
		s.aggregateAndProofDomain = domain[:]
	}
	return s.aggregateAndProofDomain, nil
}
