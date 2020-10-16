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

package v1

import (
	"context"
	"encoding/json"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type forkScheduleJSON struct {
	Data []*spec.Fork `json:"data"`
}

// ForkSchedule provides details of past and future changes in the chain's fork version.
func (s *Service) ForkSchedule(ctx context.Context) ([]*spec.Fork, error) {
	if s.forkSchedule == nil {
		respBodyReader, err := s.get(ctx, "/eth/v1/config/fork_schedule")
		if err != nil {
			return nil, errors.Wrap(err, "failed to request fork schedule")
		}

		var forkScheduleJSON forkScheduleJSON
		if err := json.NewDecoder(respBodyReader).Decode(&forkScheduleJSON); err != nil {
			return nil, errors.Wrap(err, "failed to parse fork schedule")
		}

		s.forkSchedule = forkScheduleJSON.Data
	}

	forkSchedule := make([]*spec.Fork, len(s.forkSchedule))
	for i := range s.forkSchedule {
		forkSchedule[i] = &spec.Fork{
			PreviousVersion: s.forkSchedule[i].PreviousVersion,
			CurrentVersion:  s.forkSchedule[i].CurrentVersion,
			Epoch:           s.forkSchedule[i].Epoch,
		}
	}

	return forkSchedule, nil
}
