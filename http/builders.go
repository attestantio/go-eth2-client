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
	"slices"

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
// does not know the state.  Without filters the node returns every builder,
// which on a large registry can need a longer opts.Common.Timeout.
func (s *Service) Builders(ctx context.Context,
	opts *api.BuildersOpts,
) (
	*api.Response[map[gloas.BuilderIndex]*apiv1.Builder],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "Builders")
	defer span.End()

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	if opts.State == "" {
		return nil, errors.Join(errors.New("no state specified"), client.ErrInvalidOptions)
	}

	// Check the node after the options. A caller bug must not look like an
	// inactive node, or multi fails over and deactivates a healthy client.
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}

	// The spec requires unique ids and statuses; Nimbus rejects duplicates with 400.
	body := &buildersBody{}
	seenIDs := make(map[string]struct{})
	addID := func(id string) {
		if _, exists := seenIDs[id]; !exists {
			seenIDs[id] = struct{}{}
			body.IDs = append(body.IDs, id)
		}
	}

	for i := range opts.Indices {
		addID(fmt.Sprintf("%d", opts.Indices[i]))
	}

	for i := range opts.PubKeys {
		addID(opts.PubKeys[i].String())
	}

	for i := range opts.BuilderStatuses {
		if status := opts.BuilderStatuses[i].String(); !slices.Contains(body.Statuses, status) {
			body.Statuses = append(body.Statuses, status)
		}
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

	// Seeded with a nil slice so that a body with no data, or null data, is
	// distinguishable from an empty list of builders.
	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.Builder(nil))
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.New("no builders in response")
	}

	// Data is returned as an array but we want it as a map.
	mapData := make(map[gloas.BuilderIndex]*apiv1.Builder, len(data))
	for i, builder := range data {
		if builder == nil {
			return nil, fmt.Errorf("builder %d is null", i)
		}
		mapData[builder.Index] = builder
	}

	return &api.Response[map[gloas.BuilderIndex]*apiv1.Builder]{
		Data:     mapData,
		Metadata: metadata,
	}, nil
}
