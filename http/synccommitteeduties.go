// Copyright © 2021, 2023 Attestant Limited.
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
	"fmt"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
)

// SyncCommitteeDuties obtains sync committee duties.
func (s *Service) SyncCommitteeDuties(ctx context.Context,
	opts *api.SyncCommitteeDutiesOpts,
) (
	*api.Response[[]*apiv1.SyncCommitteeDuty],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if len(opts.Indices) == 0 {
		return nil, errors.New("no validator indices specified")
	}

	var reqBodyReader bytes.Buffer
	if _, err := reqBodyReader.WriteString(`[`); err != nil {
		return nil, errors.Wrap(err, "failed to write validator index array start")
	}
	for i := range opts.Indices {
		if _, err := reqBodyReader.WriteString(fmt.Sprintf(`"%d"`, opts.Indices[i])); err != nil {
			return nil, errors.Wrap(err, "failed to write index")
		}
		if i != len(opts.Indices)-1 {
			if _, err := reqBodyReader.WriteString(`,`); err != nil {
				return nil, errors.Wrap(err, "failed to write separator")
			}
		}
	}
	if _, err := reqBodyReader.WriteString(`]`); err != nil {
		return nil, errors.Wrap(err, "failed to write end of validator index array")
	}

	url := fmt.Sprintf("/eth/v1/validator/duties/sync/%d", opts.Epoch)
	respBodyReader, err := s.post(ctx, url, &reqBodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request sync committee duties")
	}

	data, metadata, err := decodeJSONResponse(respBodyReader, []*apiv1.SyncCommitteeDuty{})
	if err != nil {
		return nil, err
	}

	return &api.Response[[]*apiv1.SyncCommitteeDuty]{
		Metadata: metadata,
		Data:     data,
	}, nil
}
