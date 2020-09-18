// Copyright © 2020 Attestant Limited.
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

package phase0

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ETH1Data provides information about the state of Ethereum 1 as viewed by the
// Ethereum 2 chain.
type ETH1Data struct {
	DepositRoot  []byte `ssz-size:"32"`
	DepositCount uint64
	BlockHash    []byte `ssz-size:"32"`
}

// eth1DataJSON is the spec representation of the struct.
type eth1DataJSON struct {
	DepositRoot  string `json:"deposit_root"`
	DepositCount string `json:"deposit_count"`
	BlockHash    string `json:"block_hash"`
}

// MarshalJSON implements json.Marshaler.
func (e *ETH1Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(&eth1DataJSON{
		DepositRoot:  fmt.Sprintf("%#x", e.DepositRoot),
		DepositCount: fmt.Sprintf("%d", e.DepositCount),
		BlockHash:    fmt.Sprintf("%#x", e.BlockHash),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ETH1Data) UnmarshalJSON(input []byte) error {
	var err error

	var eth1DataJSON eth1DataJSON
	if err = json.Unmarshal(input, &eth1DataJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if eth1DataJSON.DepositRoot == "" {
		return errors.New("deposit root missing")
	}
	if e.DepositRoot, err = hex.DecodeString(strings.TrimPrefix(eth1DataJSON.DepositRoot, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for deposit root")
	}
	if len(e.DepositRoot) != rootLength {
		return errors.New("incorrect length for deposit root")
	}
	if eth1DataJSON.DepositCount == "" {
		return errors.New("deposit count missing")
	}
	if e.DepositCount, err = strconv.ParseUint(eth1DataJSON.DepositCount, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for deposit count")
	}
	if eth1DataJSON.BlockHash == "" {
		return errors.New("block hash missing")
	}
	if e.BlockHash, err = hex.DecodeString(strings.TrimPrefix(eth1DataJSON.BlockHash, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for block hash")
	}
	if len(e.BlockHash) != hashLength {
		return errors.New("incorrect length for block hash")
	}

	return nil
}

// String returns a string version of the structure.
func (e *ETH1Data) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
