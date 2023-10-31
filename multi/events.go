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
func (s *Service) Events(ctx context.Context,
	topics []string,
	handler consensusclient.EventHandlerFunc,
) error {
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
		ah := &activeHandler{
			s:       s,
			log:     log.With().Logger(),
			address: client.Address(),
			handler: handler,
		}
		if err := client.(consensusclient.EventsProvider).Events(ctx, topics, ah.handleEvent); err != nil {
			inactiveClients = append(inactiveClients, client)

			continue
		}
		log.Trace().Str("address", ah.address).Strs("topics", topics).Msg("Events handler active")
	}

	// Periodically try all inactive clients, quitting as they become active.
	for _, inactiveClient := range inactiveClients {
		ah := &activeHandler{
			s:       s,
			log:     log.With().Logger(),
			address: inactiveClient.Address(),
			handler: handler,
		}
		go func(c consensusclient.Service, ah *activeHandler) {
			for {
				provider, isProvider := c.(consensusclient.NodeSyncingProvider)
				if !isProvider {
					ah.log.Error().Str("address", ah.address).Strs("topics", topics).Msg("Not a node syncing provider")

					return
				}
				syncResponse, err := provider.NodeSyncing(ctx, &api.NodeSyncingOpts{})
				if err != nil {
					ah.log.Error().Str("address", ah.address).Strs("topics", topics).Err(err).Msg("Failed to obtain sync state from node")

					return
				}
				if !syncResponse.Data.IsSyncing {
					// Client is now synced, set up the events call.
					if err := c.(consensusclient.EventsProvider).Events(ctx, topics, ah.handleEvent); err != nil {
						ah.log.Error().Str("address", ah.address).Strs("topics", topics).Err(err).Msg("Failed to set up events handler")
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
	handler consensusclient.EventHandlerFunc
}

func (h *activeHandler) handleEvent(event *apiv1.Event) {
	h.log.Trace().Str("address", h.address).Str("topic", event.Topic).Msg("Event received")
	// We only forward events from the currently active provider.  If we did not do this then we could end up with
	// inconsistent results, for example a client may receive a `head` event and a subsequent call to fetch the head
	// block end up with an earlier block.
	if h.s.Address() == h.address {
		h.log.Trace().Str("address", h.address).Str("topic", event.Topic).Msg("Forwarding due to primary active address")
		h.handler(event)
	}
}
