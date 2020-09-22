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

package tekuhttp

import (
	"context"
	"encoding/json"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// Fork provides the fork of the chain at a given epoch.
func (s *Service) Fork(ctx context.Context, stateID string) (*spec.Fork, error) {
	respBodyReader, err := s.get(ctx, "/node/fork")
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain current fork")
	}
	defer func() {
		if err := respBodyReader.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close HTTP body")
		}
	}()

	var fork *spec.Fork
	if err := json.NewDecoder(respBodyReader).Decode(&fork); err != nil {
		return nil, errors.Wrap(err, "failed to parse fork")
	}

	// There is no way to obtain the fork version at a given epoch.  The only
	// check we can make here is to ensure that the epoch of the fork is before
	// the requested epoch.
	epoch, err := s.EpochFromStateID(ctx, stateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain epoch from state ID")
	}
	if epoch < fork.Epoch {
		return nil, errors.New("incorrect fork obtained")
	}

	return fork, nil
}
