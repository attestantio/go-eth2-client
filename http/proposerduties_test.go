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
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestProposerDuties(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := http.New(ctx,
		http.WithLogLevel(zerolog.TraceLevel),
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	// Needed to fetch current epoch.
	genesisResponse, err := service.(client.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
	require.NoError(t, err)
	slotDuration, err := service.(client.SlotDurationProvider).SlotDuration(ctx)
	require.NoError(t, err)
	slotsPerEpoch, err := service.(client.SlotsPerEpochProvider).SlotsPerEpoch(ctx)
	require.NoError(t, err)

	tests := []struct {
		name             string
		opts             *api.ProposerDutiesOpts
		validatorIndices []phase0.ValidatorIndex
		expected         int
		err              string
	}{
		{
			name:     "Epoch",
			opts:     &api.ProposerDutiesOpts{Epoch: 0},
			expected: int(slotsPerEpoch - 1),
		},
		{
			name:     "Old",
			opts:     &api.ProposerDutiesOpts{Epoch: 1},
			expected: int(slotsPerEpoch),
			err:      "GET failed with status 404",
		},
		{
			name:     "Current",
			opts:     &api.ProposerDutiesOpts{Epoch: phase0.Epoch(uint64(time.Since(genesisResponse.Data.GenesisTime).Seconds()) / (uint64(slotDuration.Seconds()) * slotsPerEpoch))},
			expected: int(slotsPerEpoch),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.ProposerDutiesProvider).ProposerDuties(ctx, test.opts)
			if test.err != "" {
				require.ErrorContains(t, err, test.err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, response)
				require.NotNil(t, response.Data)
				require.NotNil(t, response.Metadata)
			}
		})
	}
}
