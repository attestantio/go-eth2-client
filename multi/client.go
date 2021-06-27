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
	"strings"
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
			log.Trace().Str("client", client.Address()).Int("active", len(activeClients)).Int("inactive", len(inactiveClients)).Msg("Client deactivated")
			setProviderStateMetric(ctx, client.Address(), "inactive")
		} else {
			activeClients = append(activeClients, activeClient)
		}
	}

	s.activeClients = activeClients
	setProvidersMetric(ctx, "active", len(s.activeClients))
	s.inactiveClients = inactiveClients
	setProvidersMetric(ctx, "inactive", len(s.inactiveClients))
}

// activateClient activates a client, moving it to the active list if not currently on it.
func (s *Service) activateClient(ctx context.Context, client eth2client.Service) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	activeClients := s.activeClients
	inactiveClients := make([]eth2client.Service, 0, len(s.activeClients)+len(s.inactiveClients))
	for _, inactiveClient := range s.inactiveClients {
		if inactiveClient == client {
			activeClients = append(activeClients, inactiveClient)
			log.Trace().Str("client", client.Address()).Int("active", len(activeClients)).Int("inactive", len(inactiveClients)).Msg("Client activated")
			setProviderStateMetric(ctx, client.Address(), "active")
		} else {
			inactiveClients = append(inactiveClients, inactiveClient)
		}
	}

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

// callFunc is the definition for a call function.  It provides a generic return interface
// to allow the caller to unpick the results as it sees fit.
type callFunc func(ctx context.Context, client eth2client.Service) (interface{}, error)

// errHandlerFunc is the definition for an error handler function.  It looks at the error
// returned from the client, potentially rewrites it, and also states if the error should
// result in a provider failover.
type errHandlerFunc func(ctx context.Context, client eth2client.Service, err error) (bool, error)

// doCall carries out a call on the active clients in turn until one succeeds.
func (s *Service) doCall(ctx context.Context, call callFunc, errHandler errHandlerFunc) (interface{}, error) {
	// Grab local copy of active clients in case it is updated whilst we are using it.
	s.clientsMu.RLock()
	activeClients := s.activeClients
	s.clientsMu.RUnlock()

	var err error
	var res interface{}
	for _, client := range activeClients {
		res, err = call(ctx, client)
		if err != nil {
			failover := true
			if errHandler != nil {
				failover, err = errHandler(ctx, client, err)
			}

			if failover {
				// Failed with this client; try the next.
				s.deactivateClient(ctx, client)
				continue
			}

			// No failover required, return.
			return res, err
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

// providerInfo returns information on the provider.
// Currently this just returns the name of the service (lighthouse/teku/etc.).
func (s *Service) providerInfo(ctx context.Context, provider eth2client.Service) string {
	providerName := "<unknown>"
	nodeVersionProvider, isNodeVersionProvider := provider.(eth2client.NodeVersionProvider)
	if isNodeVersionProvider {
		nodeVersion, err := nodeVersionProvider.NodeVersion(ctx)
		if err == nil {
			switch {
			case strings.Contains(strings.ToLower(nodeVersion), "lighthouse"):
				providerName = "lighthouse"
			case strings.Contains(strings.ToLower(nodeVersion), "prysm"):
				providerName = "prysm"
			case strings.Contains(strings.ToLower(nodeVersion), "teku"):
				providerName = "teku"
			}
		}
	}

	return providerName
}
