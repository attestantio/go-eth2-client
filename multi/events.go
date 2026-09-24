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
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"time"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/internal/eventdispatch"
	"github.com/rs/zerolog"
)

// deferredEventsClient is a client whose subscription can be retried until it is synced.
type deferredEventsClient interface {
	consensusclient.Service
	consensusclient.EventsProvider
	consensusclient.NodeSyncingProvider
}

var _ deferredEventsClient = (*Service)(nil)

// Events feeds requested events with the given topics to the supplied handler.  It returns an
// error if no client is subscribed or awaiting a retry, as no event would then arrive.
func (s *Service) Events(ctx context.Context,
	opts *api.EventsOpts,
) error {
	if err := http.ValidateEventsOpts(opts); err != nil {
		return err
	}

	// #nosec G404
	log := s.log.With().Str("id", fmt.Sprintf("%02x", rand.Int31())).Logger()

	// Because events are streams we treat them differently from all other calls.
	// We listen to all active clients, and only pass along events from the currently active provider.

	// Grab local copy of both active and inactive clients in case it is updated whilst we are using it.
	// The inactive list is cloned rather than shared, as active clients that fail to subscribe are
	// appended to it below, outside the lock, and the service's own list can have spare capacity
	// for that append to write into.
	s.clientsMu.RLock()
	activeClients := s.activeClients
	inactiveClients := slices.Clone(s.inactiveClients)
	s.clientsMu.RUnlock()

	// Call all active clients immediately.
	subscribed := 0
	for _, client := range activeClients {
		ah := newActiveHandler(s, log, client.Address(), opts)

		provider, isProvider := client.(consensusclient.EventsProvider)
		if !isProvider {
			ah.log.Error().Str("address", ah.address).Strs("topics", opts.Topics).Msg("Not an events provider")

			continue
		}

		if err := provider.Events(ctx, ah.clientOpts); err != nil {
			ah.log.Warn().Str("address", ah.address).Strs("topics", opts.Topics).Err(err).Msg("Failed to set up events handler")
			inactiveClients = append(inactiveClients, client)

			continue
		}

		subscribed++
		log.Trace().Str("address", ah.address).Strs("topics", opts.Topics).Msg("Events handler active")
	}

	deferredClients := make([]deferredEventsClient, 0, len(inactiveClients))
	for _, inactiveClient := range inactiveClients {
		deferredClient, isDeferrable := inactiveClient.(deferredEventsClient)
		if !isDeferrable {
			msg := "Not an events provider"
			if _, isEventsProvider := inactiveClient.(consensusclient.EventsProvider); isEventsProvider {
				msg = "Not a node syncing provider; cannot retry events subscription"
			}
			log.Error().Str("address", inactiveClient.Address()).Strs("topics", opts.Topics).Msg(msg)

			continue
		}

		deferredClients = append(deferredClients, deferredClient)
	}

	// With no client subscribed or to be retried, no event would ever arrive.
	if subscribed == 0 && len(deferredClients) == 0 {
		return errors.New("no client can provide events")
	}

	// Periodically try all inactive clients, quitting as they become active.  A failure to check
	// sync state or to subscribe is retried rather than final: the client can still be made the
	// active one later, and were it left unsubscribed its events would then never arrive.
	for _, deferredClient := range deferredClients {
		ah := newActiveHandler(s, log, deferredClient.Address(), opts)

		go func(c deferredEventsClient, ah *activeHandler) {
			for {
				syncResponse, err := c.NodeSyncing(ctx, &api.NodeSyncingOpts{})

				switch {
				case err != nil:
					ah.log.Warn().
						Str("address", ah.address).
						Strs("topics", ah.clientOpts.Topics).
						Err(err).
						Msg("Failed to obtain sync state from node; will retry")
				case !syncResponse.Data.IsSyncing:
					// Client is now synced, set up the events call.  This uses the same filtered
					// options as an initially-active client, so that events from it are subject to
					// the same active-address filtering.
					err := c.Events(ctx, ah.clientOpts)
					if err == nil {
						ah.log.Trace().Str("address", ah.address).Strs("topics", ah.clientOpts.Topics).Msg("Events handler active")

						return
					}

					ah.log.Warn().
						Str("address", ah.address).
						Strs("topics", ah.clientOpts.Topics).
						Err(err).
						Msg("Failed to set up events handler; will retry")
				default:
					// Still syncing; check again after the interval.
				}

				select {
				case <-ctx.Done():
					return
				case <-time.After(s.retryInterval()):
				}
			}
		}(deferredClient, ah)
	}

	return nil
}

// retryInterval returns how long Events waits between attempts to subscribe a deferred client.
// It is to be used in place of eventsRetryInterval, as a zero interval, as in a Service not built
// by New, would have the retry loop spin.
func (s *Service) retryInterval() time.Duration {
	if s.eventsRetryInterval <= 0 {
		return defaultEventsRetryInterval
	}

	return s.eventsRetryInterval
}

type activeHandler struct {
	s       *Service
	log     zerolog.Logger
	address string

	// clientOpts are the options handed to the underlying client: the caller's options filtered
	// on the active address.  One handler is built per client, each with options of its own.
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
	ah.clientOpts = eventdispatch.Filtered(opts, ah.forwards)

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
