// Copyright © 2025 Attestant Limited.
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
	"errors"
	"fmt"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// ValidatorLiveness provides the liveness data to the given validators.
func (s *Service) ValidatorLiveness(
	ctx context.Context,
	opts *api.ValidatorLivenessOpts,
) (
	*api.Response[[]*apiv1.ValidatorLiveness],
	error,
) {
	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	if len(opts.Indices) == 0 {
		return nil, errors.Join(errors.New("no validator indices specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/eth/v1/validator/liveness/%d", opts.Epoch)
	query := ""

	reqData, err := json.Marshal(opts.Indices)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal validator indices: %w", err)
	}

	httpResponse, err := s.post(ctx,
		endpoint,
		query,
		&opts.Common,
		bytes.NewReader(reqData),
		ContentTypeJSON,
		map[string]string{},
	)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request validator liveness"), err)
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.ValidatorLiveness{})
	if err != nil {
		return nil, errors.Join(errors.New("failed to decode validator liveness response"), err)
	}

	return &api.Response[[]*apiv1.ValidatorLiveness]{
		Data:     data,
		Metadata: metadata,
	}, nil
}
