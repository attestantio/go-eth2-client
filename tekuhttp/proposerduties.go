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
	"fmt"

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
)

type proposerDutiesJSON struct {
	Data []*api.ProposerDuty `json:"data"`
}

// ProposerDuties obtains proposer duties.
func (s *Service) ProposerDuties(ctx context.Context, epoch uint64, validators []client.ValidatorIDProvider) ([]*api.ProposerDuty, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/validator/duties/proposer/%d", epoch))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request proposer duties")
	}
	defer func() {
		if err := respBodyReader.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close HTTP body")
		}
	}()

	var resp proposerDutiesJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse proposer duties response")
	}

	return resp.Data, nil
}
