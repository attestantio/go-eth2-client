// Copyright © 2021 Attestant Limited.
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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestSyncCommitteeContribution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	// Needed to fetch current epoch.
	genesisResponse, err := service.(client.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
	require.NoError(t, err)
	slotDuration, err := service.(client.SlotDurationProvider).SlotDuration(ctx)
	require.NoError(t, err)

	tests := []struct {
		name string
		opts *api.SyncCommitteeContributionOpts
	}{
		{
			name: "Current",
			opts: &api.SyncCommitteeContributionOpts{
				Slot:              0,
				SubcommitteeIndex: 0,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.opts.Slot == 0 {
				test.opts.Slot = phase0.Slot(uint64(time.Since(genesisResponse.Data.GenesisTime).Seconds()) / uint64(slotDuration.Seconds()))
			}
			rootResponse, err := service.(client.BeaconBlockRootProvider).BeaconBlockRoot(ctx, &api.BeaconBlockRootOpts{Block: "head"})
			require.NoError(t, err)
			test.opts.BeaconBlockRoot = *rootResponse.Data
			response, err := service.(client.SyncCommitteeContributionProvider).SyncCommitteeContribution(ctx, test.opts)
			// Possible that the node is not aggregating sync committee messages...
			if err != nil {
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Equal(t, 404, apiErr.StatusCode)
				}
			} else {
				require.NotNil(t, response.Data)
			}
		})
	}
}
