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

package mock

import (
	"context"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Fork fetches fork information for the given state.
func (s *Service) Fork(ctx context.Context, _ *api.ForkOpts) (*api.Response[*phase0.Fork], error) {
	fork, err := s.forkAtEpoch(ctx, 1)
	if err != nil {
		return nil, err
	}

	return &api.Response[*phase0.Fork]{
		Data:     fork,
		Metadata: make(map[string]any),
	}, nil
}
