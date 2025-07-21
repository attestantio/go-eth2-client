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
	"testing"

	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestNodePeerCountJSON(t *testing.T) {
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
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.peerCountJSON",
		},
		{
			name:  "DisconnectedMissing",
			input: []byte(`{"connecting":"2","connected":"3","disconnecting":"4"}`),
			err:   "disconnected missing",
		},
		{
			name:  "DisconnectedWrongType",
			input: []byte(`{"disconnected":true,"connecting":"2","connected":"3","disconnecting":"4"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field peerCountJSON.disconnected of type string",
		},
		{
			name:  "DisconnectedInvalid",
			input: []byte(`{"disconnected":"invalid","connecting":"2","connected":"3","disconnecting":"4"}`),
			err:   "invalid value for disconnected: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "ConnectingMissing",
			input: []byte(`{"disconnected":"1","connected":"3","disconnecting":"4"}`),
			err:   "connecting missing",
		},
		{
			name:  "ConnectingWrongType",
			input: []byte(`{"disconnected":"1","connecting":true,"connected":"3","disconnecting":"4"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field peerCountJSON.connecting of type string",
		},
		{
			name:  "ConnectingInvalid",
			input: []byte(`{"disconnected":"1","connecting":"invalid","connected":"3","disconnecting":"4"}`),
			err:   "invalid value for connecting: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "ConnectedMissing",
			input: []byte(`{"disconnected":"1","connecting":"2","disconnecting":"4"}`),
			err:   "connected missing",
		},
		{
			name:  "ConnectedWrongType",
			input: []byte(`{"disconnected":"1","connecting":true,"connected":true,"disconnecting":"4"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field peerCountJSON.connecting of type string",
		},
		{
			name:  "ConnectedInvalid",
			input: []byte(`{"disconnected":"1","connecting":"2","connected":"invalid","disconnecting":"4"}`),
			err:   "invalid value for connected: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "DisconnectingMissing",
			input: []byte(`{"disconnected":"1","connecting":"2","connected":"3"}`),
			err:   "disconnecting missing",
		},
		{
			name:  "DisconnectingWrongType",
			input: []byte(`{"disconnected":"1","connecting":"2","connected":"3","disconnecting":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field peerCountJSON.disconnecting of type string",
		},
		{
			name:  "DisconnectingInvalid",
			input: []byte(`{"disconnected":"1","connecting":"2","connected":"3","disconnecting":"invalid"}`),
			err:   "invalid value for disconnecting: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"disconnected":"1","connecting":"2","connected":"3","disconnecting":"4"}`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res PeerCount
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
