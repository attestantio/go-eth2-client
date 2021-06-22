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
	"time"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/pkg/errors"
)

// monitor monitors active and inactive connections, and moves them between
// lists accordingly.
func (s *Service) monitor(ctx context.Context) {
	log.Trace().Msg("Monitor starting")
	for {
		select {
		case <-ctx.Done():
			log.Trace().Msg("Context done; monitor stopping")
			return
		case <-time.After(30 * time.Second):
			// Fetch all clients.
			clients := make([]eth2client.Service, 0, len(s.activeClients)+len(s.inactiveClients))
			s.clientsMu.RLock()
			clients = append(clients, s.activeClients...)
			clients = append(clients, s.inactiveClients...)
			s.clientsMu.RUnlock()

			// Ping each client to update its state.
			for _, client := range clients {
				if ping(ctx, client) {
					s.activateClient(ctx, client)
				} else {
					s.deactivateClient(ctx, client)
				}
			}
		}
	}
}

// deactivateClient deactivates a client, moving it to the inactive list if not currently on it.
func (s *Service) deactivateClient(ctx context.Context, client eth2client.Service) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	activeClients := make([]eth2client.Service, 0, len(s.activeClients)+len(s.inactiveClients))
	inactiveClients := s.inactiveClients
	for _, activeClient := range s.activeClients {
		if activeClient == client {
			inactiveClients = append(inactiveClients, activeClient)
		} else {
			activeClients = append(activeClients, activeClient)
		}
	}

	log.Trace().Str("client", client.Address()).Int("active", len(activeClients)).Int("inactive", len(inactiveClients)).Msg("Client deactivated")
	setProviderStateMetric(ctx, client.Address(), "inactive")
	s.activeClients = activeClients
	setProvidersMetric(ctx, "active", len(s.activeClients))
	s.inactiveClients = inactiveClients
	setProvidersMetric(ctx, "inactive", len(s.inactiveClients))
}

// activateClient activates a client, moving it to the active list if not currently on it.
func (s *Service) activateClient(ctx context.Context, client eth2client.Service) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	activeClients := s.inactiveClients
	inactiveClients := make([]eth2client.Service, 0, len(s.activeClients)+len(s.inactiveClients))
	for _, inactiveClient := range s.inactiveClients {
		if inactiveClient == client {
			activeClients = append(activeClients, inactiveClient)
		} else {
			inactiveClients = append(inactiveClients, inactiveClient)
		}
	}

	setProviderStateMetric(ctx, client.Address(), "active")
	log.Trace().Str("client", client.Address()).Int("active", len(activeClients)).Int("inactive", len(inactiveClients)).Msg("Client activated")
	s.activeClients = activeClients
	setProvidersMetric(ctx, "active", len(s.activeClients))
	s.inactiveClients = inactiveClients
	setProvidersMetric(ctx, "inactive", len(s.inactiveClients))
}

// ping pings a client, returning true if it is ready to serve requests and
// false otherwise.
func ping(ctx context.Context, eth2Client eth2client.Service) bool {
	provider, isProvider := eth2Client.(eth2client.NodeSyncingProvider)
	if !isProvider {
		log.Debug().Str("provider", eth2Client.Address()).Msg("Client does not provide sync state")
		return false
	}

	syncState, err := provider.NodeSyncing(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to obtain sync state from node")
		return false
	}

	return !syncState.IsSyncing
}

// callFunc is the definition ofor a call function.  It provides a generic return interface
// to allow the caller to unpick the results as it sees fit.
type callFunc func(ctx context.Context, client eth2client.Service) (interface{}, error)

// doCall carries out a call on the active clients in turn until one succeeds.
func (s *Service) doCall(ctx context.Context, call callFunc) (interface{}, error) {
	// Grab local copy of active clients in case it is updated whilst we are using it.
	s.clientsMu.RLock()
	activeClients := s.activeClients
	s.clientsMu.RUnlock()

	var err error
	var res interface{}
	for _, client := range activeClients {
		res, err = call(ctx, client)
		if err != nil {
			// Failed with this client; try the next.
			s.deactivateClient(ctx, client)
			continue
		}
		if res == nil {
			// No response from this client; try the next.
			err = errors.New("empty response")
			continue
		}
		return res, nil
	}
	return nil, err
}
