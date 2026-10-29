// Copyright © 2026 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"go.opentelemetry.io/otel"
)

// buildersBody is always sent as a JSON object, even when empty, as Prysm
// rejects an array body.
type buildersBody struct {
	IDs      []string `json:"ids,omitempty"`
	Statuses []string `json:"statuses,omitempty"`
}

// Builders provides the builders, with their status, for the given options.
// A non-2xx response returns an error wrapping an *api.Error with the node's
// status code.  A 404 means either the node does not serve the endpoint or it
// does not know the state.
func (s *Service) Builders(ctx context.Context,
	opts *api.BuildersOpts,
) (
	*api.Response[map[gloas.BuilderIndex]*apiv1.Builder],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "Builders")
	defer span.End()

	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	if opts.State == "" {
		return nil, errors.Join(errors.New("no state specified"), client.ErrInvalidOptions)
	}

	body := &buildersBody{}
	for i := range opts.Indices {
		body.IDs = append(body.IDs, fmt.Sprintf("%d", opts.Indices[i]))
	}

	for i := range opts.PubKeys {
		body.IDs = append(body.IDs, opts.PubKeys[i].String())
	}

	for i := range opts.BuilderStatuses {
		body.Statuses = append(body.Statuses, opts.BuilderStatuses[i].String())
	}

	reqData, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal request data"), err)
	}

	endpoint := fmt.Sprintf("/eth/v1/beacon/states/%s/builders", opts.State)

	httpResponse, err := s.post(ctx, endpoint, "", &opts.Common, bytes.NewReader(reqData), ContentTypeJSON, map[string]string{})
	if err != nil {
		return nil, errors.Join(errors.New("failed to request builders"), err)
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.Builder{})
	if err != nil {
		return nil, err
	}

	// Data is returned as an array but we want it as a map.
	mapData := make(map[gloas.BuilderIndex]*apiv1.Builder, len(data))
	for _, builder := range data {
		mapData[builder.Index] = builder
	}

	return &api.Response[map[gloas.BuilderIndex]*apiv1.Builder]{
		Data:     mapData,
		Metadata: metadata,
	}, nil
}
