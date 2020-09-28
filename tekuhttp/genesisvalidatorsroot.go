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

package tekuhttp

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type genesisValidatorsRootJSON struct {
	GenesisValidatorsRoot string `json:"genesis_validators_root"`
}

// GenesisValidatorsRoot provides the genesis validators root of the chain.
func (s *Service) GenesisValidatorsRoot(ctx context.Context) ([]byte, error) {
	if s.genesisValidatorsRoot == nil {
		slot, err := s.CurrentSlot(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain current slot")
		}
		// Go back by 1 to ensure that the state is present.
		if slot > 0 {
			slot--
		}
		respBodyReader, err := s.get(ctx, fmt.Sprintf("/beacon/state?slot=%d", slot))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to obtain beacon state for slot %d", slot))
		}

		var genesisValidatorsRootJSON *genesisValidatorsRootJSON
		if err := json.NewDecoder(respBodyReader).Decode(&genesisValidatorsRootJSON); err != nil {
			return nil, errors.Wrap(err, "failed to parse beacon state")
		}
		genesisValidatorsRootHex := strings.TrimSuffix(strings.TrimPrefix(genesisValidatorsRootJSON.GenesisValidatorsRoot, `"`), `"`)
		genesisValidatorsRoot, err := hex.DecodeString(strings.TrimPrefix(genesisValidatorsRootHex, "0x"))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse genesis validators root")
		}
		s.genesisValidatorsRoot = genesisValidatorsRoot
	}
	return s.genesisValidatorsRoot, nil
}
