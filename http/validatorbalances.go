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
	"encoding/json"
	"errors"
	"fmt"

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

	endpoint := fmt.Sprintf("/eth/v1/beacon/states/%s/validator_balances", opts.State)
	query := ""

	body := make([]string, 0, len(opts.Indices)+len(opts.PubKeys))
	for i := range opts.Indices {
		body = append(body, fmt.Sprintf("%d", opts.Indices[i]))
	}
	for i := range opts.PubKeys {
		body = append(body, opts.PubKeys[i].String())
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal request data"), err)
	}

	httpResponse, err := s.post(ctx, endpoint, query, &opts.Common, bytes.NewReader(data), ContentTypeJSON, map[string]string{})
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

func (*Service) validatorBalancesFromJSON(_ context.Context,
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
