// Copyright © 2020 - 2024 Attestant Limited.
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
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// ValidatorBalances provides the validator balances for the given options.
func (s *Service) ValidatorBalances(ctx context.Context,
	opts *api.ValidatorBalancesOpts,
) (
	*api.Response[map[phase0.ValidatorIndex]phase0.Gwei],
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

	if len(opts.Indices) > s.indexChunkSize(ctx) {
		return s.chunkedValidatorBalances(ctx, opts)
	}

	endpoint := fmt.Sprintf("/eth/v1/beacon/states/%s/validator_balances", opts.State)
	query := ""
	if len(opts.Indices) > 0 {
		ids := make([]string, len(opts.Indices))
		for i := range opts.Indices {
			ids[i] = fmt.Sprintf("%d", opts.Indices[i])
		}
		query = "id=" + strings.Join(ids, ",")
	}

	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common)
	if err != nil {
		return nil, err
	}

	switch httpResponse.contentType {
	case ContentTypeJSON:
		return s.validatorBalancesFromJSON(ctx, httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}
}

// chunkedValidatorBalances obtains the validator balances a chunk at a time.
func (s *Service) chunkedValidatorBalances(ctx context.Context,
	opts *api.ValidatorBalancesOpts,
) (
	*api.Response[map[phase0.ValidatorIndex]phase0.Gwei],
	error,
) {
	data := make(map[phase0.ValidatorIndex]phase0.Gwei)
	metadata := make(map[string]any)
	chunkSize := s.indexChunkSize(ctx)
	for i := 0; i < len(opts.Indices); i += chunkSize {
		chunkStart := i
		chunkEnd := i + chunkSize
		if len(opts.Indices) < chunkEnd {
			chunkEnd = len(opts.Indices)
		}
		chunkOpts := &api.ValidatorBalancesOpts{
			State:   opts.State,
			Indices: opts.Indices[chunkStart:chunkEnd],
		}
		chunkResponse, err := s.ValidatorBalances(ctx, chunkOpts)
		if err != nil {
			return nil, errors.Join(errors.New("failed to obtain chunk"), err)
		}
		for k, v := range chunkResponse.Data {
			data[k] = v
		}
		for k, v := range chunkResponse.Metadata {
			metadata[k] = v
		}
	}

	return &api.Response[map[phase0.ValidatorIndex]phase0.Gwei]{
		Data:     data,
		Metadata: metadata,
	}, nil
}

func (s *Service) validatorBalancesFromJSON(_ context.Context,
	httpResponse *httpResponse,
) (
	*api.Response[map[phase0.ValidatorIndex]phase0.Gwei],
	error,
) {
	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.ValidatorBalance{})
	if err != nil {
		return nil, err
	}

	response := &api.Response[map[phase0.ValidatorIndex]phase0.Gwei]{
		Data:     make(map[phase0.ValidatorIndex]phase0.Gwei),
		Metadata: metadata,
	}

	for _, datum := range data {
		response.Data[datum.Index] = datum.Balance
	}

	return response, nil
}
