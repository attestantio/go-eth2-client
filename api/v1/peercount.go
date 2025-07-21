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
	"strconv"

	"github.com/pkg/errors"
)

// PeerCount contains all the available information about node's peer count.
type PeerCount struct {
	Disconnected  uint64
	Connecting    uint64
	Connected     uint64
	Disconnecting uint64
}

type peerCountJSON struct {
	Disconnected  string `json:"disconnected"`
	Connecting    string `json:"connecting"`
	Connected     string `json:"connected"`
	Disconnecting string `json:"disconnecting"`
}

// MarshalJSON implements json.Marshaler.
func (p *PeerCount) MarshalJSON() ([]byte, error) {
	return json.Marshal(&peerCountJSON{
		Disconnected:  strconv.FormatUint(p.Disconnected, 10),
		Connecting:    strconv.FormatUint(p.Connecting, 10),
		Connected:     strconv.FormatUint(p.Connected, 10),
		Disconnecting: strconv.FormatUint(p.Disconnecting, 10),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PeerCount) UnmarshalJSON(input []byte) error {
	var err error

	var data peerCountJSON
	if err = json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.Disconnected == "" {
		return errors.New("disconnected missing")
	}
	disconnected, err := strconv.ParseUint(data.Disconnected, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for disconnected")
	}
	p.Disconnected = disconnected

	if data.Connecting == "" {
		return errors.New("connecting missing")
	}
	connecting, err := strconv.ParseUint(data.Connecting, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for connecting")
	}
	p.Connecting = connecting

	if data.Connected == "" {
		return errors.New("connected missing")
	}
	connected, err := strconv.ParseUint(data.Connected, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for connected")
	}
	p.Connected = connected

	if data.Disconnecting == "" {
		return errors.New("disconnecting missing")
	}
	disconnecting, err := strconv.ParseUint(data.Disconnecting, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for disconnecting")
	}
	p.Disconnecting = disconnecting

	return nil
}

func (p *PeerCount) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
