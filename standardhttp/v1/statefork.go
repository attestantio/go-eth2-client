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

package v1

import (
	"context"
	"encoding/json"
	"fmt"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type stateForkJSON struct {
	Data *spec.Fork `json:"data"`
}

// StateFork provides the fork given a state ID.
func (s *Service) StateFork(ctx context.Context, stateID string) (*spec.Fork, error) {
	if stateID == "" {
		return nil, errors.New("no state ID specified")
	}

	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/states/%s/fork", stateID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request state fork")
	}

	var stateForkJSON stateForkJSON
	if err := json.NewDecoder(respBodyReader).Decode(&stateForkJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse state fork")
	}
	if stateForkJSON.Data == nil {
		return nil, errors.New("no state fork returned")
	}

	return stateForkJSON.Data, nil
}
