package v1

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// Peer contains all the available information about a nodes peer.
type Peer struct {
	PeerID             string `json:"peer_id"`
	Enr                string `json:"enr,omitempty"`
	LastSeenP2PAddress string `json:"last_seen_p2p_address"`
	State              string `json:"state"`
	Direction          string `json:"direction"`
}

type peerJSON struct {
	PeerID             string `json:"peer_id"`
	Enr                string `json:"enr,omitempty"`
	LastSeenP2PAddress string `json:"last_seen_p2p_address"`
	State              string `json:"state"`
	Direction          string `json:"direction"`
}

// validPeerDirections are all the accepted options for peer direction.
var validPeerDirections = map[string]int{"inbound": 1, "outbound": 1}

// validPeerStates are all the accepted options for peer states.
var validPeerStates = map[string]int{"connected": 1, "connecting": 1, "disconnected": 1, "disconnecting": 1}

// MarshalJSON implements json.Marshaler.
func (p *Peer) MarshalJSON() ([]byte, error) {
	// make sure we have valid peer states and directions
	_, exists := validPeerDirections[p.Direction]
	if !exists {
		return nil, fmt.Errorf("invalid value for peer direction: %s", p.Direction)
	}
	_, exists = validPeerStates[p.State]
	if !exists {
		return nil, fmt.Errorf("invalid value for peer state: %s", p.State)
	}

	return json.Marshal(&peerJSON{
		PeerID:             p.PeerID,
		Enr:                p.Enr,
		LastSeenP2PAddress: p.LastSeenP2PAddress,
		State:              p.State,
		Direction:          p.Direction,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *Peer) UnmarshalJSON(input []byte) error {
	var peerJSON peerJSON

	if err := json.Unmarshal(input, &peerJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	_, ok := validPeerStates[peerJSON.State]
	if !ok {
		return fmt.Errorf("invalid value for peer state: %s", peerJSON.State)
	}
	p.State = peerJSON.State
	_, ok = validPeerDirections[peerJSON.Direction]
	if !ok {
		return fmt.Errorf("invalid value for peer direction: %s", peerJSON.Direction)
	}
	p.Direction = peerJSON.Direction
	p.Enr = peerJSON.Enr
	p.PeerID = peerJSON.PeerID
	p.LastSeenP2PAddress = peerJSON.LastSeenP2PAddress

	return nil
}

func (p *Peer) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
