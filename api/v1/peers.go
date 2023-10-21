package v1

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Peers struct {
	Peers []Peer
}

type peersJSON struct {
	Data []Peer                 `json:"data"`
	Meta map[string]interface{} `json:"meta,omitempty"`
}

type Peer struct {
	PeerID             string `json:"peer_id"`
	Enr                string `json:"enr,omitempty"`
	LastSeenP2PAddress string `json:"last_seen_p2p_address"`
	State              string `json:"state"`
	Direction          string `json:"direction"`
}

func (p *Peers) MarshalJSON() ([]byte, error) {
	meta := make(map[string]interface{})
	meta["count"] = len(p.Peers)

	return json.Marshal(&peersJSON{
		Data: p.Peers,
		Meta: meta,
	})
}

func (p *Peers) UnmarshalJSON(input []byte) error {
	var err error

	var peersJSON peersJSON
	if err = json.Unmarshal(input, &peersJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	p.Peers = peersJSON.Data

	return nil
}

func (p *Peers) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}

func (e *Peer) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
