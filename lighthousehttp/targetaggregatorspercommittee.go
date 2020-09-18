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
	"encoding/json"

	"github.com/pkg/errors"
)

// TargetAggregatorsPerCommittee provides the target number of aggregators for each attestation committee.
func (s *Service) TargetAggregatorsPerCommittee(ctx context.Context) (uint64, error) {
	if s.targetAggregatorsPerCommittee == nil {
		respBodyReader, err := s.get(ctx, "/spec")
		if err != nil {
			return 0, errors.Wrap(err, "failed to obtain configuration")
		}
		defer func() {
			if err := respBodyReader.Close(); err != nil {
				log.Warn().Err(err).Msg("Failed to close HTTP body")
			}
		}()
		var cfg map[string]interface{}
		if err := json.NewDecoder(respBodyReader).Decode(&cfg); err != nil {
			return 0, errors.Wrap(err, "failed to parse configuration")
		}
		val, exists := cfg["target_aggregators_per_committee"]
		if !exists {
			return 0, errors.New("failed to obtain target_aggregators_per_committee value")
		}
		floatVal, isNum := val.(float64)
		if !isNum {
			return 0, errors.New("target_aggregators_per_committee value not a number")
		}

		targetAggregatorsPerCommittee := uint64(floatVal)
		s.targetAggregatorsPerCommittee = &targetAggregatorsPerCommittee
	}
	return *s.targetAggregatorsPerCommittee, nil
}
