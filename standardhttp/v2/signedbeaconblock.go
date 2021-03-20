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

	spec "github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/pkg/errors"
)

type signedBeaconBlockJSON struct {
	Data *spec.SignedBeaconBlock `json:"data"`
}

// SignedBeaconBlock fetches a signed beacon block given a block ID.
// N.B if a signed beacon block for the block ID is not available this will return nil without an error.
func (s *Service) SignedBeaconBlock(ctx context.Context, blockID string) (*spec.SignedBeaconBlock, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/blocks/%s", blockID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request signed beacon block")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var resp signedBeaconBlockJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse signed beacon block")
	}

	return resp.Data, nil
}
