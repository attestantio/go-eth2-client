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

package auto_test

import (
	"context"
	"os"
	"testing"
	"time"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/auto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	tests := []struct {
		name    string
		address string
		err     string
		version string
	}{
		{
			name: "AddressMissing",
			err:  "problem with parameters: no address specified",
		},
		{
			name:    "Prysm",
			address: os.Getenv("PRYSMGRPC_ADDRESS"),
			version: "Prysm",
		},
		{
			name:    "Lighthouse",
			address: os.Getenv("LIGHTHOUSEHTTP_ADDRESS"),
			version: "Lighthouse",
		},
		{
			name:    "BadPort",
			address: "localhost:22",
			err:     "failed to connect to Ethereum 2 client with any known method",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, err := auto.New(context.Background(),
				auto.WithLogLevel(zerolog.Disabled),
				auto.WithTimeout(60*time.Second),
				auto.WithAddress(test.address),
			)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				version, err := service.(eth2client.NodeVersionProvider).NodeVersion(context.Background(), &api.NodeVersionOpts{})
				require.NoError(t, err)
				require.Contains(t, version, test.version)
			}
		})
	}
}
