// Copyright © 2021, 2024 Attestant Limited.
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
	"errors"
	"fmt"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// SyncCommittee fetches the sync committee for epoch at the given state.
func (s *Service) SyncCommittee(ctx context.Context,
	opts *api.SyncCommitteeOpts,
) (
	*api.Response[*apiv1.SyncCommittee],
	error,
) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	if opts.State == "" {
		return nil, errors.Join(errors.New("no state specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/eth/v1/beacon/states/%s/sync_committees", opts.State)

	query := ""
	if opts.Epoch != nil {
		query = fmt.Sprintf("epoch=%d", *opts.Epoch)
	}

	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common, false)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), apiv1.SyncCommittee{})
	if err != nil {
		return nil, err
	}

	return &api.Response[*apiv1.SyncCommittee]{
		Metadata: metadata,
		Data:     &data,
	}, nil
}
