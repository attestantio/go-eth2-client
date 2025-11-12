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

package testclients_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNetworkName(t *testing.T) {
	if os.Getenv("HTTP_ADDRESS") == "" {
		t.Skip("HTTP_ADDRESS not set")
	}

	if logLevel := os.Getenv("HTTP_DEBUG_LOG_ENABLED"); strings.ToLower(logLevel) == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create an HTTP service
	service, err := http.New(ctx,
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
		http.WithAllowDelayedStart(true),
	)
	require.NoError(t, err)

	// Get the network name using the utility function
	network := testclients.NetworkName(ctx, service)

	// The network name should be one of the known networks
	require.Contains(t, []string{"mainnet", "sepolia", "hoodi", "unknown"}, network)

	t.Logf("Connected to network: %s", network)
}
