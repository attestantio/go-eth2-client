// Copyright © 2020, 2023 Attestant Limited.
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
	"context"
	"strings"

	"github.com/attestantio/go-eth2-client/api"
)

// NodeClient provides the client for the node.
func (s *Service) NodeClient(ctx context.Context) (*api.Response[string], error) {
	response, err := s.NodeVersion(ctx, &api.NodeVersionOpts{})
	if err != nil {
		return nil, err
	}

	nodeVersion := strings.ToLower(response.Data)

	var client string
	switch {
	case strings.HasPrefix(nodeVersion, "lighthouse"):
		client = "lighthouse"
	case strings.HasPrefix(nodeVersion, "nimbus"):
		client = "nimbus"
	case strings.HasPrefix(nodeVersion, "prysm"):
		client = "prysm"
	case strings.HasPrefix(nodeVersion, "teku"):
		client = "teku"
	default:
		client = nodeVersion
	}

	return &api.Response[string]{
		Data:     client,
		Metadata: make(map[string]any),
	}, nil
}
