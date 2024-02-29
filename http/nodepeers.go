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
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// NodePeers obtains the peers of a node.
func (s *Service) NodePeers(ctx context.Context, opts *api.NodePeersOpts) (*api.Response[[]*apiv1.Peer], error) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}

	// all options are considered optional
	endpoint := "/eth/v1/node/peers"
	query := ""
	additionalFields := make([]string, 0, len(opts.State)+len(opts.Direction))

	for _, stateFilter := range opts.State {
		additionalFields = append(additionalFields, fmt.Sprintf("state=%s", stateFilter))
	}

	for _, directionFilter := range opts.Direction {
		additionalFields = append(additionalFields, fmt.Sprintf("direction=%s", directionFilter))
	}

	if len(additionalFields) > 0 {
		query = strings.Join(additionalFields, "&")
	}

	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common)
	if err != nil {
		return nil, err
	}

	if httpResponse.contentType != ContentTypeJSON {
		return nil, fmt.Errorf("unexpected content type %v (expected JSON)", httpResponse.contentType)
	}
	data, meta, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.Peer{})
	if err != nil {
		return nil, err
	}

	return &api.Response[[]*apiv1.Peer]{
		Data:     data,
		Metadata: meta,
	}, nil
}
