// Copyright © 2026 Attestant Limited.
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
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/gloas"
)

// Builders provides the builders, with their status, for the given options.
func (s *Service) Builders(ctx context.Context,
	opts *api.BuildersOpts,
) (
	*api.Response[map[gloas.BuilderIndex]*apiv1.Builder],
	error,
) {
	if s.BuildersFunc != nil {
		return s.BuildersFunc(ctx, opts)
	}

	return &api.Response[map[gloas.BuilderIndex]*apiv1.Builder]{
		Data:     map[gloas.BuilderIndex]*apiv1.Builder{},
		Metadata: make(map[string]any),
	}, nil
}
