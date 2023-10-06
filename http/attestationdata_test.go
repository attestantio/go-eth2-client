// Copyright © 2020 - 2023 Attestant Limited.
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

package http_test

import (
	"context"
	"os"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestAttestationData(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name    string
		opts    *api.AttestationDataOpts
		err     string
		errCode int
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name: "Good",
			opts: &api.AttestationDataOpts{},
		},
	}

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	// Need to fetch current slot for attestation data.
	genesis, err := service.(client.GenesisProvider).Genesis(ctx)
	require.NoError(t, err)
	slotDuration, err := service.(client.SlotDurationProvider).SlotDuration(ctx)
	require.NoError(t, err)

	for _, test := range tests {
		if test.opts != nil && test.opts.Slot == 0 {
			test.opts.Slot = phase0.Slot(uint64(time.Since(genesis.GenesisTime).Seconds()) / uint64(slotDuration.Seconds()))
		}
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.AttestationDataProvider).AttestationData(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			case test.errCode != 0:
				require.Equal(t, test.errCode, err.(api.Error).StatusCode)
			default:
				require.NoError(t, err)
				require.NotNil(t, response)
			}
		})
	}
}
