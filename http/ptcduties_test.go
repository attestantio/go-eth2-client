// Copyright © 2026 Attestant Limited.
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

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// headSlot returns the slot of the node's current head.  Every payload
// timeliness committee test anchors on it for the reason given on headEpoch
// below.
func headSlot(ctx context.Context, t *testing.T, service client.Service) phase0.Slot {
	t.Helper()

	syncState, err := service.(client.NodeSyncingProvider).NodeSyncing(ctx, &api.NodeSyncingOpts{})
	require.NoError(t, err)

	return syncState.Data.HeadSlot
}

// headEpoch returns the epoch of the node's current head.
//
// PTC duties are derived from the beacon state, so the epoch has to be one the
// node can actually build a state for.  Deriving it from the head rather than
// from the wall clock keeps the test usable on a devnet whose chain has
// stalled, and is equally correct on a live chain, where the two coincide.
func headEpoch(ctx context.Context, t *testing.T, service client.Service) (phase0.Epoch, uint64) {
	t.Helper()

	specResponse, err := service.(client.SpecProvider).Spec(ctx, &api.SpecOpts{})
	require.NoError(t, err)

	slotsPerEpoch, isCorrectType := specResponse.Data["SLOTS_PER_EPOCH"].(uint64)
	require.True(t, isCorrectType)
	require.NotZero(t, slotsPerEpoch)

	return phase0.Epoch(uint64(headSlot(ctx, t, service)) / slotsPerEpoch), slotsPerEpoch
}

func TestPTCDuties(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	epoch, slotsPerEpoch := headEpoch(ctx, t, service)

	// Request every validator on the devnet so that the payload timeliness
	// committee for at least one slot in the epoch is covered.
	indices := make([]phase0.ValidatorIndex, 256)
	for i := range indices {
		indices[i] = phase0.ValidatorIndex(i)
	}

	t.Run("NilOpts", func(t *testing.T) {
		_, err := service.(client.PTCDutiesProvider).PTCDuties(ctx, nil)
		require.ErrorIs(t, err, client.ErrNoOptions)
	})

	t.Run("NoIndices", func(t *testing.T) {
		_, err := service.(client.PTCDutiesProvider).PTCDuties(ctx, &api.PTCDutiesOpts{
			Epoch: epoch,
		})
		require.ErrorContains(t, err, "no validator indices specified")
		require.True(t, errors.Is(err, client.ErrInvalidOptions))
	})

	// Only this subtest needs the fork: the two above are refused by the client
	// before the network is reached, and they cover that refusal on any node.
	// Duties come from the beacon state, so a pre-Gloas state has no committee to
	// return — the node answers 500 IncorrectStateVariant rather than 404.
	t.Run("Good", func(t *testing.T) {
		requireOnGloas(ctx, t, service)

		response, err := service.(client.PTCDutiesProvider).PTCDuties(ctx, &api.PTCDutiesOpts{
			Epoch:   epoch,
			Indices: indices,
		})
		require.NoError(t, err)
		require.NotNil(t, response)
		require.NotEmpty(t, response.Data, "no PTC duties returned for the head epoch")

		startSlot := phase0.Slot(uint64(epoch) * slotsPerEpoch)
		endSlot := startSlot + phase0.Slot(slotsPerEpoch) - 1

		for _, duty := range response.Data {
			require.NotNil(t, duty)
			require.GreaterOrEqual(t, duty.Slot, startSlot)
			require.LessOrEqual(t, duty.Slot, endSlot)
			require.NotEqual(t, phase0.BLSPubKey{}, duty.PubKey, "duty carries a zero public key")
		}

		// The dependent root identifies the state the duties were computed
		// from; a validator needs it to detect a reorg that changes them.
		require.Contains(t, response.Metadata, "dependent_root")
	})
}
