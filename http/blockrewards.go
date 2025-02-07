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
	"errors"
	"fmt"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"go.opentelemetry.io/otel"
)

// BlockRewards provides rewards for proposing a block.
func (s *Service) BlockRewards(ctx context.Context,
	opts *api.BlockRewardsOpts,
) (
	*api.Response[*apiv1.BlockRewards],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "BlockRewards")
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

	endpoint := fmt.Sprintf("/eth/v1/beacon/rewards/blocks/%s", opts.Block)
	query := ""

	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common, false)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request block rewards"), err)
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), &apiv1.BlockRewards{})
	if err != nil {
		return nil, err
	}

	return &api.Response[*apiv1.BlockRewards]{
		Data:     data,
		Metadata: metadata,
	}, nil
}
