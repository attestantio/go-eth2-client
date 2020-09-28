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

package lighthousehttp

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

// GenesisValidatorsRoot provides the genesis validators root of the chain.
func (s *Service) GenesisValidatorsRoot(ctx context.Context) ([]byte, error) {
	if s.genesisValidatorsRoot == nil {
		respBodyReader, err := s.get(ctx, "/beacon/genesis_validators_root")
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain genesis validators root")
		}

		genesisValidatorsRootBytes, err := ioutil.ReadAll(respBodyReader)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read genesis validators root")
		}
		genesisValidatorsRootHex := strings.TrimSuffix(strings.TrimPrefix(string(genesisValidatorsRootBytes), `"`), `"`)
		genesisValidatorsRoot, err := hex.DecodeString(strings.TrimPrefix(genesisValidatorsRootHex, "0x"))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse genesis validators root")
		}
		s.genesisValidatorsRoot = genesisValidatorsRoot
	}
	return s.genesisValidatorsRoot, nil
}
