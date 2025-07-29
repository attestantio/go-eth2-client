// Copyright © 2021 - 2024 Attestant Limited.
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
	"slices"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
)

// monitor monitors active and inactive clients, and moves them between
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
	activeClients := len(s.activeClients)
	clients = append(clients, s.activeClients...)
	clients = append(clients, s.inactiveClients...)
	s.clientsMu.RUnlock()

	// Ping each client to update its state.
	for _, client := range clients {
		// We actively recheck the connection state if we had no active clients at the start of the check, in an attempt to obtain
		// at least 1 active client.
		if activeClients == 0 {
			if httpClient, isHTTPClient := client.(*http.Service); isHTTPClient {
				httpClient.CheckConnectionState(ctx)
			}
		}
		switch {
		case client.IsSynced():
			s.activateClient(ctx, client)
		default:
			s.deactivateClient(ctx, client)
		}
	}
}

// deactivateClient marks a client as deactivated, moving it to the inactive list if not currently on it.
func (s *Service) deactivateClient(ctx context.Context, client consensusclient.Service) {
	log := zerolog.Ctx(ctx)

	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	activeClients := make([]consensusclient.Service, 0, len(s.activeClients)+len(s.inactiveClients))
	inactiveClients := s.inactiveClients
	for _, activeClient := range s.activeClients {
		if activeClient == client {
			inactiveClients = append(inactiveClients, activeClient)
			s.setProviderStateMetric(ctx, client.Address(), "inactive")
		} else {
			activeClients = append(activeClients, activeClient)
		}
	}
	if len(inactiveClients) != len(s.inactiveClients) {
		log.Trace().Str("client", client.Address()).
			Int("active", len(activeClients)).
			Int("inactive", len(inactiveClients)).
			Msg("Client deactivated")
	}

	s.activeClients = activeClients
	s.inactiveClients = inactiveClients
	s.setConnectionsMetric(ctx, len(s.activeClients), len(s.inactiveClients))
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
			s.setProviderStateMetric(ctx, client.Address(), "active")
		} else {
			inactiveClients = append(inactiveClients, inactiveClient)
		}
	}
	if len(inactiveClients) != len(s.inactiveClients) {
		log.Trace().Str("client", client.Address()).
			Int("active", len(activeClients)).
			Int("inactive", len(inactiveClients)).
			Msg("Client activated")
	}

	s.activeClients = activeClients
	s.inactiveClients = inactiveClients
	s.setConnectionsMetric(ctx, len(s.activeClients), len(s.inactiveClients))
}

// scoreClient returns client score.
func (s *Service) scoreClient(clientAddr string) int {
	return s.clientScores[clientAddr]
}

// penalizeClient client score.
func (s *Service) penalizeClient(clientAddr string) {
	s.clientScores[clientAddr]--
}

// callFunc is the definition for a call function.  It provides a generic return interface
// to allow the caller to unpick the results as it sees fit.
type callFunc func(ctx context.Context, client consensusclient.Service) (any, error)

// errHandlerFunc is the definition for an error handler function.  It looks at the error
// returned from the client, potentially rewrites it, and also states if the error should
// result in a provider failover.
type errHandlerFunc func(ctx context.Context, client consensusclient.Service, err error) (bool, error)

// doCall carries out a call on the active clients in turn until one succeeds.
func (s *Service) doCall(ctx context.Context, call callFunc, errHandler errHandlerFunc) (any, error) {
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
		return nil, errors.New("no clients to which to make call")
	}

	// Sort active clients in desc order (clients with the highest scores come first).
	// Deep-copy active clients so we can sort the underlying slice elements without
	// having to worry about concurrent readers/writers.
	activeClients = s.activeClientsCopy()
	s.clientScoresMu.RLock()
	slices.SortFunc(activeClients, func(a, b consensusclient.Service) int {
		aScore := s.scoreClient(a.Address())
		bScore := s.scoreClient(b.Address())
		if aScore < bScore {
			return 1
		}
		if aScore > bScore {
			return -1
		}

		return 0
	})
	s.clientScoresMu.RUnlock()

	var err error
	var res any
	for _, client := range activeClients {
		log := log.With().Str("client", client.Name()).Str("address", client.Address()).Logger()
		res, err = call(ctx, client)
		if err != nil {
			log.Trace().Err(err).Msg(fmt.Sprintf("Potentially failing over from client %s due to error", client.Address()))
			var apiErr *api.Error
			switch {
			case errors.As(err, &apiErr) && statusCodeFamily(apiErr.StatusCode) == 4:
				log.Trace().Err(err).Msg(fmt.Sprintf("Not failing over from client %s on user error", client.Address()))

				return res, err
			case errors.Is(err, context.Canceled):
				log.Trace().Msgf("Not failing over from client %s on canceled context", client.Address())

				return res, err
			}

			failover := true
			if errHandler != nil {
				failover, err = errHandler(ctx, client, err)
			}
			if failover {
				log.Debug().Err(err).Msg(fmt.Sprintf("Failing over from client %s on error", client.Address()))

				s.clientScoresMu.Lock()
				s.penalizeClient(client.Address())
				s.clientScoresMu.Unlock()

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

// activeClientsCopy returns deep-copy of activeClients slice.
func (s *Service) activeClientsCopy() []consensusclient.Service {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()

	result := make([]consensusclient.Service, 0, len(s.activeClients))
	for _, client := range s.activeClients {
		result = append(result, client)
	}
	return result
}

// providerInfo returns information on the provider.
// Currently this just returns the name of the service (lighthouse/teku/etc.).
func (*Service) providerInfo(ctx context.Context, provider consensusclient.Service) string {
	providerName := "<unknown>"
	nodeVersionProvider, isNodeVersionProvider := provider.(consensusclient.NodeVersionProvider)
	if isNodeVersionProvider {
		response, err := nodeVersionProvider.NodeVersion(ctx, &api.NodeVersionOpts{})
		if err == nil {
			switch {
			case strings.Contains(strings.ToLower(response.Data), "lighthouse"):
				providerName = "lighthouse"
			case strings.Contains(strings.ToLower(response.Data), "lodestar"):
				providerName = "lodestar"
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

func statusCodeFamily(status int) int {
	return status / 100
}
