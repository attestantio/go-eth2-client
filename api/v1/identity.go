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

package v1

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// NodeIdentity contains the node identity information.
type NodeIdentity struct {
	PeerID             string            `json:"peer_id"`
	Enr                string            `json:"enr"`
	P2PAddresses       []string          `json:"p2p_addresses"`
	DiscoveryAddresses []string          `json:"discovery_addresses"`
	Metadata           map[string]string `json:"metadata"`
}

type nodeIdentityJSON struct {
	PeerID             string            `json:"peer_id"`
	Enr                string            `json:"enr"`
	P2PAddresses       []string          `json:"p2p_addresses"`
	DiscoveryAddresses []string          `json:"discovery_addresses"`
	Metadata           map[string]string `json:"metadata"`
}

// MarshalJSON implements json.Marshaler.
func (n *NodeIdentity) MarshalJSON() ([]byte, error) {
	return json.Marshal(&nodeIdentityJSON{
		PeerID:             n.PeerID,
		Enr:                n.Enr,
		P2PAddresses:       n.P2PAddresses,
		DiscoveryAddresses: n.DiscoveryAddresses,
		Metadata:           n.Metadata,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *NodeIdentity) UnmarshalJSON(input []byte) error {
	var nodeIdentityJSON nodeIdentityJSON

	if err := json.Unmarshal(input, &nodeIdentityJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	n.PeerID = nodeIdentityJSON.PeerID
	n.Enr = nodeIdentityJSON.Enr
	n.P2PAddresses = nodeIdentityJSON.P2PAddresses
	n.DiscoveryAddresses = nodeIdentityJSON.DiscoveryAddresses
	n.Metadata = nodeIdentityJSON.Metadata

	return nil
}

func (n *NodeIdentity) String() string {
	data, err := json.Marshal(n)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
