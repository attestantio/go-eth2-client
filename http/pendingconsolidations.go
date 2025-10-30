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
	"github.com/attestantio/go-eth2-client/spec/electra"
	"go.opentelemetry.io/otel"
)

// PendingConsolidations returns the pending consolidations for a given state.
func (s *Service) PendingConsolidations(ctx context.Context,
	opts *api.PendingConsolidationsOpts,
) (
	*api.Response[[]*electra.PendingConsolidation],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "PendingConsolidations")
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

	endpoint := fmt.Sprintf("/eth/v1/beacon/states/%s/pending_consolidations", opts.State)
	query := ""

	resp, err := s.get(ctx, endpoint, query, &opts.Common, false)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request pending consolidations"), err)
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(resp.body), []*electra.PendingConsolidation{})
	if err != nil {
		return nil, err
	}

	return &api.Response[[]*electra.PendingConsolidation]{
		Data:     data,
		Metadata: metadata,
	}, nil
}
