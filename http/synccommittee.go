// Copyright © 2021 Attestant Limited.
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

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
)

type syncCommitteeJSON struct {
	Data *api.SyncCommittee `json:"data"`
}

// SyncCommittee fetches the sync committee for the epoch at the given state.
func (s *Service) SyncCommittee(ctx context.Context, stateID string) (*api.SyncCommittee, error) {
	url := fmt.Sprintf("/eth/v1/beacon/states/%s/sync_committees", stateID)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		log.Trace().Str("url", url).Err(err).Msg("Request failed")
		return nil, errors.Wrap(err, "failed to request sync committee")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain sync committee")
	}

	var resp syncCommitteeJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse sync committee")
	}

	return resp.Data, nil
}
