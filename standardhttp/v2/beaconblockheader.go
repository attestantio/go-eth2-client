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

package v2

import (
	"context"
	"encoding/json"
	"fmt"

	api "github.com/attestantio/go-eth2-client/api/v2"
	"github.com/pkg/errors"
)

type beaconBlockHeaderJSON struct {
	Data *api.BeaconBlockHeader `json:"data"`
}

// BeaconBlockHeader provides the block header of a given block ID.
func (s *Service) BeaconBlockHeader(ctx context.Context, blockID string) (*api.BeaconBlockHeader, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/headers/%s", blockID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon block header")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain beacon block header")
	}

	var resp beaconBlockHeaderJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon block header")
	}

	return resp.Data, nil
}
