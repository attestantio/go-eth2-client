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
	"sync"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
)

// Service handles multiple Ethereum 2 clients.
type Service struct {
	log zerolog.Logger

	clientsMu       sync.RWMutex
	activeClients   []consensusclient.Service
	inactiveClients []consensusclient.Service
}

// New creates a new Ethereum 2 client with multiple endpoints.
// The endpoints are periodiclaly checked to see if they are active,
// and requests will retry a different client if the currently active
// client fails to respond.
func New(ctx context.Context, params ...Parameter) (consensusclient.Service, error) {
	parameters, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	// Set logging.
	log := zerologger.With().Str("service", "client").Str("impl", "multi").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}
	ctx = log.WithContext(ctx)

	if parameters.monitor != nil {
		if err := registerMetrics(ctx, parameters.monitor); err != nil {
			return nil, errors.Wrap(err, "failed to register metrics")
		}
	}

	// Check the state of each client and put it in an active or inactive list, accordingly.
	activeClients := make([]consensusclient.Service, 0, len(parameters.clients))
	inactiveClients := make([]consensusclient.Service, 0, len(parameters.clients))
	for _, client := range parameters.clients {
		if ping(ctx, client) {
			activeClients = append(activeClients, client)
		} else {
			inactiveClients = append(inactiveClients, client)
		}
	}
	for _, address := range parameters.addresses {
		client, err := http.New(ctx,
			http.WithLogLevel(parameters.logLevel),
			http.WithTimeout(parameters.timeout),
			http.WithAddress(address),
			http.WithEnforceJSON(parameters.enforceJSON),
			http.WithExtraHeaders(parameters.extraHeaders),
		)
		if err != nil {
			log.Error().Str("provider", address).Msg("Provider not present; dropping from rotation")

			continue
		}
		if ping(ctx, client) {
			activeClients = append(activeClients, client)
			setProviderActiveMetric(ctx, client.Address(), "active")
		} else {
			inactiveClients = append(inactiveClients, client)
			setProviderActiveMetric(ctx, client.Address(), "inactive")
		}
	}
	if len(activeClients) == 0 {
		return nil, errors.New("No providers active, cannot proceed")
	}
	log.Trace().Int("active", len(activeClients)).Int("inactive", len(inactiveClients)).Msg("Initial providers")
	setProvidersMetric(ctx, "active", len(activeClients))
	setProvidersMetric(ctx, "inactive", len(inactiveClients))

	s := &Service{
		log:             log,
		activeClients:   activeClients,
		inactiveClients: inactiveClients,
	}

	// Kick off monitor.
	go s.monitor(ctx)

	return s, nil
}

// Name returns the name of the client implementation.
func (s *Service) Name() string {
	return "multi"
}

// Address returns the address of the client.
func (s *Service) Address() string {
	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()
	if len(s.activeClients) > 0 {
		return s.activeClients[0].Address()
	}

	return "none"
}
