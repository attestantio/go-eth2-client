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

var (
	_ deferredEventsClient = (*Service)(nil)
	// A client that is not a deferredEventsClient is never retried.  This covers the HTTP
	// clients WithAddresses creates; those supplied through WithClients are checked when Events
	// is called.
	_ deferredEventsClient = (*http.Service)(nil)
)

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
	// The copies are only read.  Clients to retry go on a list of Events' own, as the service's
	// lists can have spare capacity, into which appending here, outside the lock, would write.
	s.clientsMu.RLock()
	activeClients := s.activeClients
	inactiveClients := s.inactiveClients
	s.clientsMu.RUnlock()

	// Call all active clients immediately.  Those that fail are retried with the inactive clients,
	// but only after an interval, as they have just been tried.
	subscribed := 0
	retries := make([]eventsRetry, 0, len(activeClients)+len(inactiveClients))
	for _, client := range activeClients {
		var ah *activeHandler
		var err error
		if provider, isProvider := client.(consensusclient.EventsProvider); isProvider {
			ah = newActiveHandler(s, log, client.Address(), opts)
			err = provider.Events(ctx, ah.clientOpts)
			if err == nil {
				subscribed++
				log.Trace().Str("address", ah.address).Strs("topics", opts.Topics).Msg("Events handler active")

				continue
			}
		}

		if deferredClient, isRetryable := retryableEventsClient(log, client, opts.Topics, err); isRetryable {
			retries = append(retries, eventsRetry{client: deferredClient, handler: ah, wait: true})
		}
	}

	for _, inactiveClient := range inactiveClients {
		if deferredClient, isRetryable := retryableEventsClient(log, inactiveClient, opts.Topics, nil); isRetryable {
			retries = append(retries, eventsRetry{
				client:  deferredClient,
				handler: newActiveHandler(s, log, deferredClient.Address(), opts),
			})
		}
	}

	// With no client subscribed or to be retried, no event would ever arrive.
	if subscribed == 0 && len(retries) == 0 {
		return errors.New("no client can provide events")
	}

	// Periodically try each client to retry, quitting as it becomes active.  A failure to check
	// sync state or to subscribe is retried rather than final: the client can still be made the
	// active one later, and were it left unsubscribed its events would then never arrive.
	for _, retry := range retries {
		go func(c deferredEventsClient, ah *activeHandler, wait bool) {
			for {
				if wait {
					select {
					case <-ctx.Done():
						return
					case <-time.After(s.retryInterval()):
					}
				}
				wait = true

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
			}
		}(retry.client, retry.handler, retry.wait)
	}

	return nil
}

// eventsRetry is a client whose subscription Events retries.
type eventsRetry struct {
	client deferredEventsClient
	// handler forwards the client's events while it is the active one.
	handler *activeHandler
	// wait is set if the client has just failed to subscribe, so is to wait an interval before
	// it is first retried.
	wait bool
}

// retryableEventsClient reports whether the subscription of a client that is not subscribed can
// be retried, logging once for the client why it is being retried or dropped.  subscribeErr is
// the error from an attempt to subscribe the client, or nil if none was made.
func retryableEventsClient(log zerolog.Logger,
	client consensusclient.Service,
	topics []string,
	subscribeErr error,
) (
	deferredEventsClient,
	bool,
) {
	deferredClient, isDeferrable := client.(deferredEventsClient)
	_, isProvider := client.(consensusclient.EventsProvider)

	level := zerolog.ErrorLevel
	var msg string
	switch {
	case !isProvider:
		msg = "Not an events provider"
	case isDeferrable && subscribeErr == nil:
		return deferredClient, true
	case isDeferrable:
		level, msg = zerolog.WarnLevel, "Failed to set up events handler; will retry"
	case subscribeErr != nil:
		msg = "Failed to set up events handler; not a node syncing provider, so will not retry"
	default:
		msg = "Not a node syncing provider; not subscribing to events"
	}
	log.WithLevel(level).Str("address", client.Address()).Strs("topics", topics).Err(subscribeErr).Msg(msg)

	return deferredClient, isDeferrable
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
