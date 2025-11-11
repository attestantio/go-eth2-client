// Copyright © 2021, 2023 Attestant Limited.
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
	"errors"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func slotptr(slot uint64) *phase0.Slot {
	res := phase0.Slot(slot)
	return &res
}

func TestAttestationPool(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	// Need to fetch current slot for attestation pools.
	genesisResponse, err := service.(client.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
	require.NoError(t, err)
	slotDuration, err := service.(client.SlotDurationProvider).SlotDuration(ctx)
	require.NoError(t, err)
	currentSlot := uint64(time.Since(genesisResponse.Data.GenesisTime).Seconds()) / uint64(slotDuration.Seconds())
	t.Logf("currentSlot: %d", currentSlot)
	committeeIndex := phase0.CommitteeIndex(0)

	tests := []struct {
		name    string
		opts    *api.AttestationPoolOpts
		assert  func(t *testing.T, response *api.Response[[]*spec.VersionedAttestation])
		err     string
		errCode int
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name: "Empty",
			opts: &api.AttestationPoolOpts{
				Slot: slotptr(1),
			},
			assert: func(t *testing.T, response *api.Response[[]*spec.VersionedAttestation]) {
				require.NotNil(t, response)
				require.Equal(t, len(response.Data), 0)
			},
		},
		{
			name: "Previous Slot",
			opts: &api.AttestationPoolOpts{
				Slot: slotptr(currentSlot - 1), // Get the previous slot to decrease the chance of getting an empty pool
			},
			assert: func(t *testing.T, response *api.Response[[]*spec.VersionedAttestation]) {
				require.NotNil(t, response)
				require.Greater(t, len(response.Data), 0, "Beacon node probably returned no attestations. Try again.")
				data, err := response.Data[0].Data()
				require.NoError(t, err)
				require.Equal(t, uint64(data.Slot), currentSlot-1)
			},
		},
		{
			name: "Previous Slot, Committee Index 0",
			opts: &api.AttestationPoolOpts{
				Slot:           slotptr(currentSlot - 1), // Get the previous slot to decrease the chance of getting an empty pool
				CommitteeIndex: &committeeIndex,
			},
			assert: func(t *testing.T, response *api.Response[[]*spec.VersionedAttestation]) {
				require.NotNil(t, response)
				require.Greater(t, len(response.Data), 0, "Beacon node probably returned no attestations. Try again.")
				data, err := response.Data[0].Data()
				require.NoError(t, err)
				require.Equal(t, data.Index, committeeIndex)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.AttestationPoolProvider).AttestationPool(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			case test.errCode != 0:
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Equal(t, test.errCode, apiErr.StatusCode)
				}
			default:
				require.NoError(t, err)
				require.NotNil(t, response)
				if test.assert != nil {
					test.assert(t, response)
				}
			}
		})
	}
}
