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

package testclients

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// NetworkName returns the network name (Mainnet, Sepolia, Holesky, or Hoodi) for the given service.
// It identifies the network by querying the genesis validators root and comparing it against known values.
// Returns "Unknown" if the network cannot be identified or if genesis information is unavailable.
func NetworkName(ctx context.Context, service any) string {
	// Known genesis validators roots for different networks
	knownNetworks := map[string]string{
		"0x4b363db94e286120d76eb905340fdd4e54bfe9f06bf33ff6cf5ad27f511bfe95": "mainnet",
		"0xd8ea171f3c94aea21ebc42a1ed61052acf3f9209c00e4efbaaddac09ed9b8078": "sepolia",
		"0x212f13fc4df078b6cb7db228f1c8307566dcecf900867401a92023d7ba99cb5f": "hoodi",
	}

	// Type assert to GenesisProvider interface
	genesisProvider, ok := service.(interface {
		Genesis(ctx context.Context, opts *api.GenesisOpts) (*api.Response[*apiv1.Genesis], error)
	})
	if !ok {
		return "Unknown"
	}

	// Query genesis information
	response, err := genesisProvider.Genesis(ctx, &api.GenesisOpts{})
	if err != nil {
		return "Unknown"
	}

	if response == nil || response.Data == nil {
		return "Unknown"
	}

	// Format the genesis validators root as hex string
	rootHex := strings.ToLower(hex.EncodeToString(response.Data.GenesisValidatorsRoot[:]))
	rootHexWithPrefix := "0x" + rootHex

	// Look up network name
	if networkName, exists := knownNetworks[rootHexWithPrefix]; exists {
		return networkName
	}

	return "Unknown"
}
