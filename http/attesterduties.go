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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type attesterDutiesJSON struct {
	Data []*api.AttesterDuty `json:"data"`
}

// AttesterDuties obtains attester duties.
func (s *Service) AttesterDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*api.AttesterDuty, error) {
	// Try a POST request.
	var reqBodyReader bytes.Buffer
	if _, err := reqBodyReader.WriteString(`[`); err != nil {
		return nil, errors.Wrap(err, "failed to write validator index array start")
	}
	for i := range validatorIndices {
		if _, err := reqBodyReader.WriteString(fmt.Sprintf(`"%d"`, validatorIndices[i])); err != nil {
			return nil, errors.Wrap(err, "failed to write index")
		}
		if i != len(validatorIndices)-1 {
			if _, err := reqBodyReader.WriteString(`,`); err != nil {
				return nil, errors.Wrap(err, "failed to write separator")
			}
		}
	}
	if _, err := reqBodyReader.WriteString(`]`); err != nil {
		return nil, errors.Wrap(err, "failed to write end of validator index array")
	}
	url := fmt.Sprintf("/eth/v1/validator/duties/attester/%d", epoch)
	respBodyReader, err := s.post(ctx, url, &reqBodyReader)
	if err != nil {
		// Didn't work.  Try a GET request.
		indices := make([]string, len(validatorIndices))
		for i := range validatorIndices {
			indices[i] = fmt.Sprintf("%d", validatorIndices[i])
		}
		url := fmt.Sprintf("/eth/v1/validator/duties/attester/%d?index=%s", epoch, strings.Join(indices, ","))
		respBodyReader, err = s.get(ctx, url)
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to request attester duties")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain attester duties")
	}

	var resp attesterDutiesJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse attester duties response")
	}

	return resp.Data, nil
}
