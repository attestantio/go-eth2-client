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
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type beaconStateJSON struct {
	Data *phase0.BeaconState `json:"data"`
}

// BeaconState fetches a beacon state.
// N.B if the requested beacon state is not available this will return nil without an error.
func (s *Service) BeaconState(ctx context.Context, stateID string) (*phase0.BeaconState, error) {
	url := fmt.Sprintf("/eth/v1/debug/beacon/states/%s", stateID)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		log.Trace().Str("url", url).Err(err).Msg("Request failed")
		return nil, errors.Wrap(err, "failed to request beacon state")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var resp beaconStateJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon state")
	}

	return resp.Data, nil
}
