// Copyright © 2021, 2022 Attestant Limited.
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

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// monitor monitors active and inactive connections, and moves them between
// lists accordingly.
func (s *Service) monitor(ctx context.Context) {
	log := s.log.With().Logger()
	ctx = log.WithContext(ctx)

	log.Trace().Msg("Monitor starting")
	for {
		select {
		case <-ctx.Done():
			log.Trace().Msg("Context done; monitor stopping")

			return
		case <-time.After(30 * time.Second):
			s.recheck(ctx)
		}
	}
}

// recheck checks clients to update their state.
func (s *Service) recheck(ctx context.Context) {
	// Fetch all clients.
	clients := make([]consensusclient.Service, 0, len(s.activeClients)+len(s.inactiveClients))
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

// deactivateClient deactivates a client, moving it to the inactive list if not currently on it.
func (s *Service) deactivateClient(ctx context.Context, client consensusclient.Service) {
	log := zerolog.Ctx(ctx)

	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	activeClients := make([]consensusclient.Service, 0, len(s.activeClients)+len(s.inactiveClients))
	inactiveClients := s.inactiveClients
	for _, activeClient := range s.activeClients {
		if activeClient == client {
			inactiveClients = append(inactiveClients, activeClient)
			setProviderActiveMetric(ctx, client.Address(), "inactive")
		} else {
			activeClients = append(activeClients, activeClient)
		}
	}
	if len(inactiveClients) != len(s.inactiveClients) {
		log.Trace().Str("client", client.Address()).Int("active", len(activeClients)).Int("inactive", len(inactiveClients)).Msg("Client deactivated")
	}

	s.activeClients = activeClients
	setProvidersMetric(ctx, "active", len(s.activeClients))
	s.inactiveClients = inactiveClients
	setProvidersMetric(ctx, "inactive", len(s.inactiveClients))
}

// activateClient activates a client, moving it to the active list if not currently on it.
func (s *Service) activateClient(ctx context.Context, client consensusclient.Service) {
	log := zerolog.Ctx(ctx)

	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	activeClients := s.activeClients
	inactiveClients := make([]consensusclient.Service, 0, len(s.activeClients)+len(s.inactiveClients))
	for _, inactiveClient := range s.inactiveClients {
		if inactiveClient == client {
			activeClients = append(activeClients, inactiveClient)
			setProviderActiveMetric(ctx, client.Address(), "active")
		} else {
			inactiveClients = append(inactiveClients, inactiveClient)
		}
	}
	if len(inactiveClients) != len(s.inactiveClients) {
		log.Trace().Str("client", client.Address()).Int("active", len(activeClients)).Int("inactive", len(inactiveClients)).Msg("Client activated")
	}

	s.activeClients = activeClients
	setProvidersMetric(ctx, "active", len(s.activeClients))
	s.inactiveClients = inactiveClients
	setProvidersMetric(ctx, "inactive", len(s.inactiveClients))
}

// ping pings a client, returning true if it is ready to serve requests and
// false otherwise.
func ping(ctx context.Context, client consensusclient.Service) bool {
	log := zerolog.Ctx(ctx)

	provider, isProvider := client.(consensusclient.NodeSyncingProvider)
	if !isProvider {
		log.Debug().Str("provider", client.Address()).Msg("Client does not provide sync state")

		return false
	}

	response, err := provider.NodeSyncing(ctx, &api.NodeSyncingOpts{})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to obtain sync state from node")

		return false
	}

	return (!response.Data.IsSyncing) || (response.Data.HeadSlot == 0 && response.Data.SyncDistance == 0)
}

// callFunc is the definition for a call function.  It provides a generic return interface
// to allow the caller to unpick the results as it sees fit.
type callFunc func(ctx context.Context, client consensusclient.Service) (interface{}, error)

// errHandlerFunc is the definition for an error handler function.  It looks at the error
// returned from the client, potentially rewrites it, and also states if the error should
// result in a provider failover.
type errHandlerFunc func(ctx context.Context, client consensusclient.Service, err error) (bool, error)

// doCall carries out a call on the active clients in turn until one succeeds.
func (s *Service) doCall(ctx context.Context, call callFunc, errHandler errHandlerFunc) (interface{}, error) {
	log := s.log.With().Logger()
	ctx = log.WithContext(ctx)

	// Grab local copy of active clients in case it is updated whilst we are using it.
	s.clientsMu.RLock()
	activeClients := s.activeClients
	s.clientsMu.RUnlock()

	if len(activeClients) == 0 {
		// There are no active clients; attempt to re-enable the inactive clients.
		s.recheck(ctx)
		s.clientsMu.RLock()
		activeClients = s.activeClients
		s.clientsMu.RUnlock()
	}

	if len(activeClients) == 0 {
		return nil, errors.New("no active clients to which to make call")
	}

	var err error
	var res interface{}
	for _, client := range activeClients {
		res, err = call(ctx, client)
		if err != nil {
			log.Trace().Err(err).Msg("Potentially deactivating client due to error")
			var apiErr *api.Error
			if errors.As(err, &apiErr) && apiErr.StatusCode/100 == 4 {
				log.Trace().Str("client", client.Name()).Str("address", client.Address()).Err(err).Msg("Not deactivating client on user error")

				return res, err
			}
			if errors.Is(err, context.Canceled) {
				log.Trace().Str("client", client.Name()).Str("address", client.Address()).Msg("Not deactivating client on canceled context")

				return res, err
			}

			failover := true
			if errHandler != nil {
				failover, err = errHandler(ctx, client, err)
			}

			if failover {
				log.Debug().Str("client", client.Name()).Str("address", client.Address()).Err(err).Msg("Deactivating client on error")
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
func (s *Service) providerInfo(ctx context.Context, provider consensusclient.Service) string {
	providerName := "<unknown>"
	nodeVersionProvider, isNodeVersionProvider := provider.(consensusclient.NodeVersionProvider)
	if isNodeVersionProvider {
		response, err := nodeVersionProvider.NodeVersion(ctx, &api.NodeVersionOpts{})
		if err == nil {
			switch {
			case strings.Contains(strings.ToLower(response.Data), "lighthouse"):
				providerName = "lighthouse"
			case strings.Contains(strings.ToLower(response.Data), "prysm"):
				providerName = "prysm"
			case strings.Contains(strings.ToLower(response.Data), "teku"):
				providerName = "teku"
			case strings.Contains(strings.ToLower(response.Data), "nimbus"):
				providerName = "nimbus"
			}
		}
	}

	return providerName
}
