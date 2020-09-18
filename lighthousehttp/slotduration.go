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
	"time"

	"github.com/pkg/errors"
)

// SlotDuration provides the duration of a slot of the chain.
func (s *Service) SlotDuration(ctx context.Context) (time.Duration, error) {
	if s.slotDuration == nil {
		respBodyReader, err := s.get(ctx, "/spec")
		if err != nil {
			return time.Duration(0), errors.Wrap(err, "failed to obtain configuration")
		}
		defer func() {
			if err := respBodyReader.Close(); err != nil {
				log.Warn().Err(err).Msg("Failed to close HTTP body")
			}
		}()

		var cfg map[string]interface{}
		if err := json.NewDecoder(respBodyReader).Decode(&cfg); err != nil {
			return time.Duration(0), errors.Wrap(err, "failed to parse configuration")
		}
		val, exists := cfg["milliseconds_per_slot"]
		if !exists {
			return time.Duration(0), errors.New("failed to obtain milliseconds_per_slot value")
		}
		floatVal, isNum := val.(float64)
		if !isNum {
			return time.Duration(0), errors.New("milliseconds_per_slot value not a number")
		}

		slotDuration := time.Duration(uint64(floatVal)) * time.Millisecond
		s.slotDuration = &slotDuration
	}
	return *s.slotDuration, nil
}
