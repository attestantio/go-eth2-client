// Copyright © 2021 - 2026 Attestant Limited.
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

package multi

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/rs/zerolog"
)

// Events feeds requested events with the given topics to the supplied handler.
//
// The caller must not mutate opts after calling this.  Each handler supplied is captured for the
// lifetime of the subscription, so replacing one afterwards has no effect.  The topics slice is
// retained by reference, and is read again whenever a client that was not synced at this point
// later comes into service, which can be long after this returns.
func (s *Service) Events(ctx context.Context,
	opts *api.EventsOpts,
) error {
	if opts == nil {
		return consensusclient.ErrNoOptions
	}

	// #nosec G404
	log := s.log.With().Str("id", fmt.Sprintf("%02x", rand.Int31())).Logger()

	// Because events are streams we treat them differently from all other calls.
	// We listen to all active clients, and only pass along events from the currently active provider.

	// Grab local copy of both active and inactive clients in case it is updated whilst we are using it.
	s.clientsMu.RLock()
	activeClients := s.activeClients
	inactiveClients := s.inactiveClients
	s.clientsMu.RUnlock()

	// Call all active clients immediately.
	for _, client := range activeClients {
		ah := newActiveHandler(s, log, client.Address(), opts)

		if err := client.(consensusclient.EventsProvider).Events(ctx, ah.clientOpts); err != nil {
			inactiveClients = append(inactiveClients, client)

			continue
		}

		log.Trace().Str("address", ah.address).Strs("topics", opts.Topics).Msg("Events handler active")
	}

	// Periodically try all inactive clients, quitting as they become active.
	for _, inactiveClient := range inactiveClients {
		ah := newActiveHandler(s, log, inactiveClient.Address(), opts)

		go func(c consensusclient.Service, ah *activeHandler) {
			for {
				provider, isProvider := c.(consensusclient.NodeSyncingProvider)
				if !isProvider {
					ah.log.Error().
						Str("address", ah.address).
						Strs("topics", opts.Topics).
						Msg("Not a node syncing provider")

					return
				}

				syncResponse, err := provider.NodeSyncing(ctx, &api.NodeSyncingOpts{})
				if err != nil {
					ah.log.Error().
						Str("address", ah.address).
						Strs("topics", opts.Topics).
						Err(err).
						Msg("Failed to obtain sync state from node")

					return
				}

				if !syncResponse.Data.IsSyncing {
					// Client is now synced, set up the events call.  This uses the same substituted
					// options as an initially-active client, so that events from it are subject to
					// the same active-address filtering.
					if err := c.(consensusclient.EventsProvider).Events(ctx, ah.clientOpts); err != nil {
						ah.log.Error().
							Str("address", ah.address).
							Strs("topics", opts.Topics).
							Err(err).
							Msg("Failed to set up events handler")
					}

					// Return either way.
					return
				}

				time.Sleep(5 * time.Second)
			}
		}(inactiveClient, ah)
	}

	return nil
}

type activeHandler struct {
	s       *Service
	log     zerolog.Logger
	address string

	// clientOpts are the options handed to the underlying client: the caller's topics and common
	// options, with each handler the caller supplied replaced by a wrapper that filters on the
	// active address before forwarding to it.  This is always a struct of its own and never the
	// caller's own options modified in place, because one handler is built per client: wrapping
	// in place would have the second client's wrapper wrap the first's, filtering an event
	// against two addresses at once.
	clientOpts *api.EventsOpts
}

// newActiveHandler creates a handler that filters the events of the client at the given address,
// forwarding those from the currently active client to the caller's handlers.
func newActiveHandler(s *Service, log zerolog.Logger, address string, opts *api.EventsOpts) *activeHandler {
	ah := &activeHandler{
		s:       s,
		log:     log,
		address: address,
	}

	sub := &api.EventsOpts{
		Common: opts.Common,
		Topics: opts.Topics,
	}
	ah.clientOpts = sub

	// These need no nil check of their own: substitute leaves a handler the caller did not
	// supply nil, for the reason given on it.
	sub.Handler = substituteGeneric(ah, opts.Handler)
	sub.AttestationHandler = substitute(ah, "attestation", opts.AttestationHandler)
	sub.AttesterSlashingHandler = substitute(ah, "attester_slashing", opts.AttesterSlashingHandler)
	sub.BlobSidecarHandler = substitute(ah, "blob_sidecar", opts.BlobSidecarHandler)
	sub.BlockHandler = substitute(ah, "block", opts.BlockHandler)
	sub.BlockGossipHandler = substitute(ah, "block_gossip", opts.BlockGossipHandler)
	sub.BLSToExecutionChangeHandler = substitute(ah, "bls_to_execution_change", opts.BLSToExecutionChangeHandler)
	sub.ChainReorgHandler = substitute(ah, "chain_reorg", opts.ChainReorgHandler)
	sub.ContributionAndProofHandler = substitute(ah, "contribution_and_proof", opts.ContributionAndProofHandler)
	sub.DataColumnSidecarHandler = substitute(ah, "data_column_sidecar", opts.DataColumnSidecarHandler)
	sub.ExecutionPayloadHandler = substitute(ah, "execution_payload", opts.ExecutionPayloadHandler)
	sub.ExecutionPayloadAvailableHandler = substitute(ah, "execution_payload_available", opts.ExecutionPayloadAvailableHandler)
	sub.ExecutionPayloadBidHandler = substitute(ah, "execution_payload_bid", opts.ExecutionPayloadBidHandler)
	sub.ExecutionPayloadGossipHandler = substitute(ah, "execution_payload_gossip", opts.ExecutionPayloadGossipHandler)
	sub.FastConfirmationHandler = substitute(ah, "fast_confirmation", opts.FastConfirmationHandler)
	sub.FinalizedCheckpointHandler = substitute(ah, "finalized_checkpoint", opts.FinalizedCheckpointHandler)
	sub.HeadHandler = substitute(ah, "head", opts.HeadHandler)
	sub.PayloadAttestationMessageHandler = substitute(ah, "payload_attestation_message", opts.PayloadAttestationMessageHandler)
	sub.PayloadAttributesHandler = substitute(ah, "payload_attributes", opts.PayloadAttributesHandler)
	sub.ProposerPreferencesHandler = substitute(ah, "proposer_preferences", opts.ProposerPreferencesHandler)
	sub.ProposerSlashingHandler = substitute(ah, "proposer_slashing", opts.ProposerSlashingHandler)
	sub.SingleAttestationHandler = substitute(ah, "single_attestation", opts.SingleAttestationHandler)
	sub.VoluntaryExitHandler = substitute(ah, "voluntary_exit", opts.VoluntaryExitHandler)

	return ah
}

// forwards reports whether an event just received from this handler's client should be passed on
// to the caller.  We only forward events from the currently active provider.  If we did not do
// this then we could end up with inconsistent results, for example a client may receive a `head`
// event and a subsequent call to fetch the head block end up with an earlier block.
func (h *activeHandler) forwards(topic string) bool {
	forwarding := h.s.Address() == h.address

	h.log.Trace().
		Str("address", h.address).
		Str("topic", topic).
		Bool("forwarding", forwarding).
		Msg("Event received")

	return forwarding
}

// substitute wraps one of the caller's topic handlers so that an event reaches it only when this
// handler's client is the currently active one.  The wrapper is what the underlying client is
// given in place of the caller's own handler.  A nil handler yields a nil wrapper, leaving the
// field unset, because clients fall back to the generic handler for any topic whose specific
// handler is nil and substituting one the caller did not supply would starve that fallback.
//
// topic is diagnostic only.  It names the handler in the trace log and has no bearing on whether
// the event is forwarded, so a wrong one here mislabels a log line rather than misrouting a call.
func substitute[T any](h *activeHandler, topic string, handler func(context.Context, T)) func(context.Context, T) {
	if handler == nil {
		return nil
	}

	return func(ctx context.Context, data T) {
		if !h.forwards(topic) {
			return
		}

		handler(ctx, data)
	}
}

// substituteGeneric is substitute for the caller's generic handler, which carries its own topic
// and takes no context, and so does not fit substitute's shape.
func substituteGeneric(h *activeHandler, handler api.EventHandlerFunc) api.EventHandlerFunc {
	if handler == nil {
		return nil
	}

	return func(event *apiv1.Event) {
		if !h.forwards(event.Topic) {
			return
		}

		handler(event)
	}
}
