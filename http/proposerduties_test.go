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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	eth2http "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestProposerDutiesV2Transport(t *testing.T) {
	ctx := context.Background()
	currentEpochRoot := phase0.Root{0xaa}
	nextEpochRoot := phase0.Root{0xbb}
	headRoot := phase0.Root{0xcc}
	legacyRoot := phase0.Root{0xdd}
	pubkey := phase0.BLSPubKey{0xee}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/eth/v1/node/version":
			fmt.Fprint(w, `{"data":{"version":"stub"}}`)
		case "/eth/v1/node/syncing":
			fmt.Fprint(w, `{"data":{"head_slot":"64","sync_distance":"0","is_syncing":false,"is_optimistic":false}}`)
		case "/eth/v1/config/spec":
			fmt.Fprint(w, `{"data":{"SLOTS_PER_EPOCH":"32"}}`)
		case "/eth/v2/validator/duties/proposer/2":
			if r.Method != http.MethodGet {
				t.Errorf("unexpected request method %s", r.Method)
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			fmt.Fprintf(w, `{"dependent_root":"%s","execution_optimistic":true,"data":[{"pubkey":"%s","validator_index":"7","slot":"65"}]}`,
				currentEpochRoot.String(), pubkey.String())
		case "/eth/v2/validator/duties/proposer/3":
			if r.Method != http.MethodGet {
				t.Errorf("unexpected request method %s", r.Method)
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			fmt.Fprintf(w, `{"dependent_root":"%s","execution_optimistic":true,"data":[{"pubkey":"%s","validator_index":"7","slot":"97"}]}`,
				nextEpochRoot.String(), pubkey.String())
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	service, err := eth2http.New(ctx,
		eth2http.WithAddress(server.URL),
		eth2http.WithLogLevel(zerolog.Disabled),
	)
	require.NoError(t, err)

	provider, supported := service.(client.ProposerDutiesV2Provider)
	require.True(t, supported)

	tests := []struct {
		name  string
		epoch phase0.Epoch
		slot  phase0.Slot
		root  phase0.Root
	}{
		{
			name:  "CurrentEpoch",
			epoch: 2,
			slot:  65,
			root:  currentEpochRoot,
		},
		{
			name:  "NextEpoch",
			epoch: 3,
			slot:  97,
			root:  nextEpochRoot,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := provider.ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{Epoch: test.epoch})
			require.NoError(t, err)
			require.Equal(t, []*apiv1.ProposerDuty{{
				PubKey:         pubkey,
				ValidatorIndex: 7,
				Slot:           test.slot,
			}}, response.Data)
			require.Equal(t, test.root, response.Metadata["dependent_root"])
			require.NotEqual(t, phase0.Root{}, response.Metadata["dependent_root"])
			require.NotEqual(t, headRoot, response.Metadata["dependent_root"])
			require.NotEqual(t, legacyRoot, response.Metadata["dependent_root"])
			require.Equal(t, true, response.Metadata["execution_optimistic"])
		})
	}
}

func TestProposerDuties(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

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
