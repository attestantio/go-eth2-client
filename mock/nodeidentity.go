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

// NodeIdentity provides the identity information of the node.
func (s *Service) NodeIdentity(ctx context.Context,
	opts *api.NodeIdentityOpts,
) (
	*api.Response[*apiv1.NodeIdentity],
	error,
) {
	if s.NodeIdentityFunc != nil {
		return s.NodeIdentityFunc(ctx, opts)
	}

	return &api.Response[*apiv1.NodeIdentity]{
		Data: &apiv1.NodeIdentity{
			PeerID:             "16Uiu2HAmMockPeerID",
			Enr:                "enr:-mock-enr-value",
			P2PAddresses:       []string{"/ip4/127.0.0.1/tcp/9000"},
			DiscoveryAddresses: []string{"/ip4/127.0.0.1/udp/9000"},
			Metadata: map[string]string{
				"seq_number": "1",
				"attnets":    "0xffffffffffffffff",
			},
		},
	}, nil
}
