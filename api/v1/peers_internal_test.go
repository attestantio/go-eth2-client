package v1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestNodePeerJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "GoodNoENR",
			input: []byte(`{"peer_id":"16Uiu2HAm7ukVy4XugqVShYbLih4H2jBJjYevevznBZaHsmd1FM96","last_seen_p2p_address":"/ip4/10.0.20.8/tcp/43402","state":"connected","direction":"inbound"}`),
		},
		{
			name:  "GoodWithENR",
			input: []byte(`{"peer_id":"16Uiu2HAmTJgqKuVcN1QReyWzwELRkfWCjLAfBSu3KxuBuWFvvaLX","enr":"enr:-MS4QExfvXqHhj-nqAqkg1Sn55uV7YgpRtlImGCvMJkrkbnLDo8sGhecAGid9B3NjXzN3UtGxpOOUqHZVcEDQxkniwoBh2F0dG5ldHOIAAAAAAAAAACEZXRoMpAEg2rFBAAGZgIAAAAAAAAAgmlkgnY0gmlwhAoAFBWEcXVpY4IjKYlzZWNwMjU2azGhA9mr4bIskWVeMt0dEn4IlQJhOFgOqgR9V3gkHTl1lTioiHN5bmNuZXRzAIN0Y3CCIyiDdWRwgiMo","last_seen_p2p_address":"/ip4/10.0.20.21/udp/9001/quic-v1/p2p/16Uiu2HAmTJgqKuVcN1QReyWzwELRkfWCjLAfBSu3KxuBuWFvvaLX","state":"connected","direction":"outbound"}`),
		},
		{
			name:  "BadDirection",
			input: []byte(`{"peer_id":"16Uiu2HAm7ukVy4XugqVShYbLih4H2jBJjYevevznBZaHsmd1FM96","last_seen_p2p_address":"/ip4/10.0.20.8/tcp/43402","state":"connected","direction":"backwards"}`),
			err:   "invalid value for peer direction: backwards",
		},
		{
			name:  "BadState",
			input: []byte(`{"peer_id":"16Uiu2HAm7ukVy4XugqVShYbLih4H2jBJjYevevznBZaHsmd1FM96","last_seen_p2p_address":"/ip4/10.0.20.8/tcp/43402","state":"tightly-coupled","direction":"inbound"}`),
			err:   "invalid value for peer state: tightly-coupled",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.peerJSON",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res Peer
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}
