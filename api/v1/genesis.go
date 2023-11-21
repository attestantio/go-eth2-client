// Copyright © 2020, 2021 Attestant Limited.
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
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// Genesis provides information about the genesis of a chain.
type Genesis struct {
	GenesisTime           time.Time
	GenesisValidatorsRoot phase0.Root
	GenesisForkVersion    phase0.Version
}

// genesisJSON is the spec representation of the struct.
type genesisJSON struct {
	GenesisTime           string `json:"genesis_time"`
	GenesisValidatorsRoot string `json:"genesis_validators_root"`
	GenesisForkVersion    string `json:"genesis_fork_version"`
}

// MarshalJSON implements json.Marshaler.
func (g *Genesis) MarshalJSON() ([]byte, error) {
	return json.Marshal(&genesisJSON{
		GenesisTime:           strconv.FormatInt(g.GenesisTime.Unix(), 10),
		GenesisValidatorsRoot: fmt.Sprintf("%#x", g.GenesisValidatorsRoot),
		GenesisForkVersion:    fmt.Sprintf("%#x", g.GenesisForkVersion),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (g *Genesis) UnmarshalJSON(input []byte) error {
	var err error

	var genesisJSON genesisJSON
	if err = json.Unmarshal(input, &genesisJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if genesisJSON.GenesisTime == "" {
		return errors.New("genesis time missing")
	}
	genesisTime, err := strconv.ParseInt(genesisJSON.GenesisTime, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for genesis time")
	}
	g.GenesisTime = time.Unix(genesisTime, 0)

	if genesisJSON.GenesisValidatorsRoot == "" {
		return errors.New("genesis validators root missing")
	}
	genesisValidatorsRoot, err := hex.DecodeString(strings.TrimPrefix(genesisJSON.GenesisValidatorsRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for genesis validators root")
	}
	if len(genesisValidatorsRoot) != rootLength {
		return fmt.Errorf("incorrect length %d for genesis validators root", len(genesisValidatorsRoot))
	}
	copy(g.GenesisValidatorsRoot[:], genesisValidatorsRoot)

	if genesisJSON.GenesisForkVersion == "" {
		return errors.New("genesis fork version missing")
	}
	genesisForkVersion, err := hex.DecodeString(strings.TrimPrefix(genesisJSON.GenesisForkVersion, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for genesis fork version")
	}
	if len(genesisForkVersion) != forkLength {
		return fmt.Errorf("incorrect length %d for genesis fork version", len(genesisForkVersion))
	}
	copy(g.GenesisForkVersion[:], genesisForkVersion)

	return nil
}

// String returns a string version of the structure.
func (g *Genesis) String() string {
	data, err := json.Marshal(g)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
