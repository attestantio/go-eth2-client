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

type finalityJSON struct {
	Data *api.Finality `json:"data"`
}

// Finality provides the finality given a state ID.
func (s *Service) Finality(ctx context.Context, stateID string) (*api.Finality, error) {
	if stateID == "" {
		return nil, errors.New("no state ID specified")
	}

	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/states/%s/finality_checkpoints", stateID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request finality checkpoints")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain finality checkpoints")
	}

	var finalityJSON finalityJSON
	if err := json.NewDecoder(respBodyReader).Decode(&finalityJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse finality")
	}
	if finalityJSON.Data == nil {
		return nil, errors.New("no finality returned")
	}

	return finalityJSON.Data, nil
}
