package v1

import (
	"encoding/json"
	"fmt"
)

// NodePeer is the data regarding which peer a node has.
type NodePeer struct {
	PeerID             string `json:"peer_id"`
	Enr                string `json:"enr,omitempty"`
	LastSeenP2pAddress string `json:"last_seen_p2p_address"`
	State              string `json:"state"`
	Direction          string `json:"direction"`
}

// NodePeers is the response from a /eth/v1/node/peers endpoint.
type NodePeers struct {
	Data []NodePeer             `json:"data"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}

func (e *NodePeers) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}

func (e *NodePeer) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
