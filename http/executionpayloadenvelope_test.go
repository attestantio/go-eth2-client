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
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestExecutionPayloadEnvelope exercises the validator-side envelope retrieval
// against a live node.  This is how a caller that asked for a proposal with the
// payload excluded gets the envelope it needs to publish, so the node only holds
// one: the envelope it built for the slot it is proposing.
func TestExecutionPayloadEnvelope(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)
	provider := service.(client.ExecutionPayloadEnvelopeProvider)

	t.Run("NilOpts", func(t *testing.T) {
		_, err := provider.ExecutionPayloadEnvelope(ctx, nil)
		require.ErrorIs(t, err, client.ErrNoOptions)
	})

	t.Run("NoSlot", func(t *testing.T) {
		_, err := provider.ExecutionPayloadEnvelope(ctx, &api.ExecutionPayloadEnvelopeOpts{
			BeaconBlockRoot: phase0.Root{0x01},
		})
		require.ErrorIs(t, err, client.ErrInvalidOptions)
		require.ErrorContains(t, err, "no slot specified")
	})

	// The block root is what makes the lookup re-org resistant, so an unset one
	// would ask the node to hand over whatever it happens to hold.
	t.Run("NoBeaconBlockRoot", func(t *testing.T) {
		_, err := provider.ExecutionPayloadEnvelope(ctx, &api.ExecutionPayloadEnvelopeOpts{Slot: 1})
		require.ErrorIs(t, err, client.ErrInvalidOptions)
		require.ErrorContains(t, err, "no beacon block root specified")
	})

	// A node caches only the current slot's locally built envelope, so a slot
	// already behind it holds nothing.  That is a defined outcome rather than a
	// failure, and it surfaces as a sentinel: reported as a raw transport error,
	// a caller could not tell "this node did not build that block" — the re-org
	// and wrong-node cases, which are recoverable — apart from a real fault.
	t.Run("PastSlotIsSentinel", func(t *testing.T) {
		// Gated although it passes on a pre-Gloas node, because it passes there
		// for the wrong reason: that node has no such endpoint, and its 404 maps
		// to the same sentinel this asserts.  What is meant to be under test is a
		// Gloas node reporting a cache miss, which needs a gloas head.
		requireOnGloas(ctx, t, service)

		// Far enough back to be certain nothing is cached, and still comfortably
		// after the fork, which the node rejects separately.
		slot := headSlot(ctx, t, service) - 100

		_, err := provider.ExecutionPayloadEnvelope(ctx, &api.ExecutionPayloadEnvelopeOpts{
			Slot:            slot,
			BeaconBlockRoot: phase0.Root{0x01},
		})
		require.ErrorIs(t, err, client.ErrNoExecutionPayloadEnvelope)
	})

	// A pre-Gloas slot is a different refusal, and it must not be folded into the
	// sentinel: asking for an envelope before the fork that introduced envelopes
	// is a caller mistake, not a cache miss to be retried elsewhere.
	t.Run("PreGloasSlotIsNotSentinel", func(t *testing.T) {
		// The weaker gate, and deliberately so: this is a test of pre-Gloas
		// behaviour, whose only precondition is that an endpoint exists to do the
		// refusing.  Gating it on the head's fork instead would keep it dark on a
		// mainnet-lineage node forever, when it should start running there the day
		// that node's client learns Gloas — years before the fork activates.
		requireKnowsGloas(ctx, t, service)

		_, err := provider.ExecutionPayloadEnvelope(ctx, &api.ExecutionPayloadEnvelopeOpts{
			Slot:            1,
			BeaconBlockRoot: phase0.Root{0x01},
		})
		require.Error(t, err)
		require.NotErrorIs(t, err, client.ErrNoExecutionPayloadEnvelope)

		var apiErr *api.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, 400, apiErr.StatusCode)
	})
}
