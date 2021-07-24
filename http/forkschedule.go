// Copyright © 2020, 2021 Attestant Limited.
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

package http

import (
	"context"
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type forkScheduleJSON struct {
	Data []*phase0.Fork `json:"data"`
}

// ForkSchedule provides details of past and future changes in the chain's fork version.
func (s *Service) ForkSchedule(ctx context.Context) ([]*phase0.Fork, error) {
	if s.forkSchedule != nil {
		return s.forkSchedule, nil
	}

	s.forkScheduleMutex.Lock()
	defer s.forkScheduleMutex.Unlock()
	if s.forkSchedule != nil {
		// Someone else fetched this whilst we were waiting for the lock.
		return s.forkSchedule, nil
	}

	// Up to us to fetch the information.
	respBodyReader, err := s.get(ctx, "/eth/v1/config/fork_schedule")
	if err != nil {
		return nil, errors.Wrap(err, "failed to request fork schedule")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain fork schedule")
	}

	var resp forkScheduleJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse fork schedule")
	}
	s.forkSchedule = resp.Data
	return s.forkSchedule, nil
}
