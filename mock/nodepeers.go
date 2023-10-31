// Copyright © 2021, 2023 Attestant Limited.
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

// NodePeers provides the peers of the node.
func (s *Service) NodePeers(_ context.Context, _ *api.NodePeersOpts) (*api.Response[[]*apiv1.Peer], error) {
	return &api.Response[[]*apiv1.Peer]{
		Data: []*apiv1.Peer{{
			PeerID:             "MOCK16Uiu2HAm7ukVy4XugqVShYbLih4H2jBJjYevevznBZaHsmd1FM96",
			LastSeenP2PAddress: "/ip4/10.0.20.8/tcp/43402",
			State:              "connected",
			Direction:          "outbound",
		}},
	}, nil
}
