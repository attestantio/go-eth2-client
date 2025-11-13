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
			name: "Default options (current slot and all committee indices)",
			opts: &api.AttestationPoolOpts{},
			assert: func(t *testing.T, response *api.Response[[]*spec.VersionedAttestation]) {
				require.NotNil(t, response)
				require.Greater(t, len(response.Data), 0, "Beacon node probably returned no attestations. Try again.")

				committeeIndices := make(map[int]bool)
				for _, attestation := range response.Data {
					committeeBits, err := attestation.CommitteeBits()
					require.NoError(t, err)
					for _, committeeIndex := range committeeBits.BitIndices() {
						committeeIndices[committeeIndex] = true
					}
					_, err = attestation.Data()
					require.NoError(t, err)
				}
				require.Greater(t, len(committeeIndices), 0, "Beacon node returned attestations, we should have at least one committee index.")
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

func TestAttestationPoolCommitteeIndexSet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	genesisResponse, err := service.(client.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
	require.NoError(t, err)
	slotDuration, err := service.(client.SlotDurationProvider).SlotDuration(ctx)
	require.NoError(t, err)
	currentSlot := uint64(time.Since(genesisResponse.Data.GenesisTime).Seconds()) / uint64(slotDuration.Seconds())

	// Collect all committee indices that have attestations to avoid testing invalid indices
	committeeIndices := collectCommitteeIndicesWithAttestations(ctx, t, service)
	require.Greater(t, len(committeeIndices), 0, "Beacon node should have returned attestations for at least one committee")

	// Helper function to find attestations for any committee index
	findAttestationsForCommittee := func(opts *api.AttestationPoolOpts) (*api.Response[[]*spec.VersionedAttestation], error) {
		for committeeIndex := range committeeIndices {
			optsCopy := *opts // Create a copy to avoid modifying the original
			optsCopy.CommitteeIndex = &committeeIndex
			response, err := service.(client.AttestationPoolProvider).AttestationPool(ctx, &optsCopy)
			if err == nil && response != nil && len(response.Data) > 0 {
				return response, nil
			}
		}
		return nil, errors.New("no attestations found for any committee index")
	}

	tests := []struct {
		name string
		opts *api.AttestationPoolOpts
	}{
		{
			name: "Committee Index Set, Slot specified",
			opts: &api.AttestationPoolOpts{
				Slot: slotptr(currentSlot - 1), // Previous slot to increase chance of finding attestations
			},
		},
		{
			name: "Committee Index Set, Slot not specified",
			opts: &api.AttestationPoolOpts{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := findAttestationsForCommittee(test.opts)
			require.NoError(t, err)
			require.NotNil(t, response)
			require.Greater(t, len(response.Data), 0, "Beacon node should have returned attestations")
		})
	}
}

// collectCommitteeIndicesWithAttestations fetches all attestations and returns the set of committee indices that have attestations
func collectCommitteeIndicesWithAttestations(ctx context.Context, t *testing.T, service client.Service) map[phase0.CommitteeIndex]bool {
	t.Helper()
	response, err := service.(client.AttestationPoolProvider).AttestationPool(ctx, &api.AttestationPoolOpts{})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Greater(t, len(response.Data), 0, "Beacon node should have returned some attestations")

	indices := make(map[phase0.CommitteeIndex]bool)
	for _, attestation := range response.Data {
		committeeBits, err := attestation.CommitteeBits()
		require.NoError(t, err)
		for _, committeeIndex := range committeeBits.BitIndices() {
			indices[phase0.CommitteeIndex(committeeIndex)] = true
		}
	}
	return indices
}
