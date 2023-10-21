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

package v1

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// DepositContract represents the details of the Ethereum 1 deposit contract for a chain.
type DepositContract struct {
	ChainID uint64
	Address []byte
}

// depositContractJSON is the standard API representation of the struct.
type depositContractJSON struct {
	ChainID string `json:"chain_id"`
	Address string `json:"address"`
}

// MarshalJSON implements json.Marshaler.
func (d *DepositContract) MarshalJSON() ([]byte, error) {
	return json.Marshal(&depositContractJSON{
		ChainID: strconv.FormatUint(d.ChainID, 10),
		Address: fmt.Sprintf("%#x", d.Address),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *DepositContract) UnmarshalJSON(input []byte) error {
	var err error

	var depositContractJSON depositContractJSON
	if err = json.Unmarshal(input, &depositContractJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if depositContractJSON.ChainID == "" {
		return errors.New("chain ID missing")
	}
	if d.ChainID, err = strconv.ParseUint(depositContractJSON.ChainID, 10, 64); err != nil {
		return errors.Wrap(err, "invalid value for chain ID")
	}
	if depositContractJSON.Address == "" {
		return errors.New("address missing")
	}
	if d.Address, err = hex.DecodeString(strings.TrimPrefix(depositContractJSON.Address, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for address")
	}
	if len(d.Address) != eth1AddressLength {
		return fmt.Errorf("incorrect length %d for address", len(d.Address))
	}

	return nil
}

// String returns a string version of the structure.
func (d *DepositContract) String() string {
	data, err := json.Marshal(d)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
