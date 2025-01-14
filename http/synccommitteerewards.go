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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// SyncCommitteeRewards provides rewards to the given validators for being members of a sync committee.
func (s *Service) SyncCommitteeRewards(ctx context.Context,
	opts *api.SyncCommitteeRewardsOpts,
) (
	*api.Response[[]*apiv1.SyncCommitteeReward],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "SyncCommitteeRewards")
	defer span.End()

	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}
	if opts.Block == "" {
		return nil, errors.Join(errors.New("no block specified"), client.ErrInvalidOptions)
	}
	span.SetAttributes(attribute.Int("validators", len(opts.Indices)+len(opts.PubKeys)))

	endpoint := fmt.Sprintf("/eth/v1/beacon/rewards/sync_committee/%s", opts.Block)
	query := ""

	body := make([]string, 0, len(opts.Indices)+len(opts.PubKeys))

	for i := range opts.Indices {
		body = append(body, fmt.Sprintf("%d", opts.Indices[i]))
	}
	for i := range opts.PubKeys {
		body = append(body, opts.PubKeys[i].String())
	}

	reqData, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal request data"), err)
	}

	httpResponse, err := s.post(ctx, endpoint, query, &opts.Common, bytes.NewReader(reqData), ContentTypeJSON, map[string]string{})
	if err != nil {
		return nil, errors.Join(errors.New("failed to request sync committee rewards"), err)
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.SyncCommitteeReward{})
	if err != nil {
		return nil, err
	}

	return &api.Response[[]*apiv1.SyncCommitteeReward]{
		Data:     data,
		Metadata: metadata,
	}, nil
}
