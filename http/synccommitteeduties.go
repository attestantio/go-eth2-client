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
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type syncCommitteeDutiesJSON struct {
	Data []*api.SyncCommitteeDuty `json:"data"`
}

// SyncCommitteeDuties obtains sync committee duties.
func (s *Service) SyncCommitteeDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*api.SyncCommitteeDuty, error) {
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
	url := fmt.Sprintf("/eth/v1/validator/duties/sync/%d", epoch)
	respBodyReader, err := s.post(ctx, url, &reqBodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request sync committee duties")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain sync committee duties")
	}

	var resp syncCommitteeDutiesJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse sync committee duties response")
	}

	return resp.Data, nil
}
