// Copyright © 2021, 2022 Attestant Limited.
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
	"time"

	"github.com/attestantio/go-eth2-client/api"
)

// Spec provides the spec information of the chain.
// This returns various useful values.
func (s *Service) Spec(_ context.Context, _ *api.SpecOpts) (*api.Response[map[string]any], error) {
	data := map[string]any{
		"SECONDS_PER_SLOT": 12 * time.Second,
		"SLOTS_PER_EPOCH":  uint64(32),
	}

	return &api.Response[map[string]any]{
		Data:     data,
		Metadata: make(map[string]any),
	}, nil
}
