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

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
)

// ForkChoice fetches all current fork choice context.
func (s *Service) ForkChoice(ctx context.Context) (*api.ForkChoice, error) {
	respBodyReader, err := s.get(ctx, "/eth/v1/debug/fork_choice")
	if err != nil {
		return nil, errors.Wrap(err, "failed to request fork choice")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain fork choice")
	}

	var forkChoice *api.ForkChoice
	if err := json.NewDecoder(respBodyReader).Decode(&forkChoice); err != nil {
		return nil, errors.Wrap(err, "failed to parse fork choice")
	}

	return forkChoice, nil
}
