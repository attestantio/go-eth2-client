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
	"strings"

	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type attesterDutiesJSON struct {
	Data []*api.AttesterDuty `json:"data"`
}

// AttesterDuties obtains attester duties.
func (s *Service) AttesterDuties(ctx context.Context, epoch spec.Epoch, indices []spec.ValidatorIndex) ([]*api.AttesterDuty, error) {
	validatorIndices := make([]string, 0, len(indices))
	for i := range indices {
		validatorIndices = append(validatorIndices, fmt.Sprintf("%d", indices[i]))
	}

	url := fmt.Sprintf("/eth/v1/validator/duties/attester/%d?index=%s", epoch, strings.Join(validatorIndices, ","))
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request attester duties")
	}

	var resp attesterDutiesJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse attester duties response")
	}

	return resp.Data, nil
}
