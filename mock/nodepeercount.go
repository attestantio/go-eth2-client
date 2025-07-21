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

package mock

import (
	"context"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// NodePeerCount provides the peer count of the node.
func (s *Service) NodePeerCount(ctx context.Context,
	opts *api.NodePeerCountOpts,
) (
	*api.Response[*apiv1.PeerCount],
	error,
) {
	if s.NodePeerCountFunc != nil {
		return s.NodePeerCountFunc(ctx, opts)
	}

	return &api.Response[*apiv1.PeerCount]{
		Data: &apiv1.PeerCount{
			Disconnected:  1,
			Connecting:    2,
			Connected:     3,
			Disconnecting: 4,
		},
	}, nil
}
