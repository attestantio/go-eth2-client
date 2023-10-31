// Copyright © 2020 - 2023 Attestant Limited.
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
	"bytes"
	"context"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// ForkSchedule provides details of past and future changes in the chain's fork version.
func (s *Service) ForkSchedule(ctx context.Context,
	opts *api.ForkScheduleOpts,
) (
	*api.Response[[]*phase0.Fork],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}

	s.forkScheduleMutex.RLock()
	if s.forkSchedule != nil {
		defer s.forkScheduleMutex.RUnlock()

		return &api.Response[[]*phase0.Fork]{
			Data:     s.forkSchedule,
			Metadata: make(map[string]any),
		}, nil
	}
	s.forkScheduleMutex.RUnlock()

	s.forkScheduleMutex.Lock()
	defer s.forkScheduleMutex.Unlock()
	if s.forkSchedule != nil {
		// Someone else fetched this whilst we were waiting for the lock.
		return &api.Response[[]*phase0.Fork]{
			Data:     s.forkSchedule,
			Metadata: make(map[string]any),
		}, nil
	}

	// Up to us to fetch the information.
	url := "/eth/v1/config/fork_schedule"
	httpResponse, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*phase0.Fork{})
	if err != nil {
		return nil, err
	}
	s.forkSchedule = data

	return &api.Response[[]*phase0.Fork]{
		Data:     s.forkSchedule,
		Metadata: metadata,
	}, nil
}
