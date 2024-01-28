// Copyright © 2020, 2021 Attestant Limited.
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

package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
)

// Service is an Ethereum 2 client service.
type Service struct {
	// log is a service-wide logger.
	log zerolog.Logger

	base    *url.URL
	address string
	client  *http.Client
	timeout time.Duration

	// Various information from the node that does not change during the
	// lifetime of a beacon node.
	genesis              *apiv1.Genesis
	genesisMutex         sync.RWMutex
	spec                 map[string]interface{}
	specMutex            sync.RWMutex
	depositContract      *apiv1.DepositContract
	depositContractMutex sync.RWMutex
	forkSchedule         []*phase0.Fork
	forkScheduleMutex    sync.RWMutex
	nodeVersion          string
	nodeVersionMutex     sync.RWMutex

	// User-specified chunk sizes.
	userIndexChunkSize  int
	userPubKeyChunkSize int
	extraHeaders        map[string]string

	// Endpoint support.
	enforceJSON              bool
	connectedToDVTMiddleware bool
}

// New creates a new Ethereum 2 client service, connecting with a standard HTTP.
func New(ctx context.Context, params ...Parameter) (eth2client.Service, error) {
	parameters, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	// Set logging.
	log := zerologger.With().Str("service", "client").Str("impl", "http").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}

	if parameters.monitor != nil {
		if err := registerMetrics(ctx, parameters.monitor); err != nil {
			return nil, errors.Wrap(err, "failed to register metrics")
		}
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   parameters.timeout,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:        64,
			MaxConnsPerHost:     64,
			MaxIdleConnsPerHost: 64,
			IdleConnTimeout:     600 * time.Second,
		},
	}

	address := parameters.address
	if !strings.HasPrefix(address, "http") {
		address = fmt.Sprintf("http://%s", address)
	}
	if !strings.HasSuffix(address, "/") {
		address = fmt.Sprintf("%s/", address)
	}
	base, err := url.Parse(address)
	if err != nil {
		return nil, errors.Wrap(err, "invalid URL")
	}

	s := &Service{
		log:                 log,
		base:                base,
		address:             parameters.address,
		client:              client,
		timeout:             parameters.timeout,
		userIndexChunkSize:  parameters.indexChunkSize,
		userPubKeyChunkSize: parameters.pubKeyChunkSize,
		extraHeaders:        parameters.extraHeaders,
		enforceJSON:         parameters.enforceJSON,
	}

	// Fetch static values to confirm the connection is good.
	if err := s.fetchStaticValues(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to confirm node connection")
	}

	// Periodially refetch static values in case of client update.
	s.periodicClearStaticValues(ctx)

	// Handle connection to DVT middleware.
	if err := s.checkDVT(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to check DVT connection")
	}

	// Close the service on context done.
	go func(s *Service) {
		<-ctx.Done()
		log.Trace().Msg("Context done; closing connection")
		s.close()
	}(s)

	return s, nil
}

// fetchStaticValues fetches values that never change.
// This caches the values, avoiding future API calls.
func (s *Service) fetchStaticValues(ctx context.Context) error {
	if _, err := s.Genesis(ctx, &api.GenesisOpts{}); err != nil {
		return errors.Wrap(err, "failed to fetch genesis")
	}
	if _, err := s.Spec(ctx, &api.SpecOpts{}); err != nil {
		return errors.Wrap(err, "failed to fetch spec")
	}
	if _, err := s.DepositContract(ctx, &api.DepositContractOpts{}); err != nil {
		return errors.Wrap(err, "failed to fetch deposit contract")
	}
	if _, err := s.ForkSchedule(ctx, &api.ForkScheduleOpts{}); err != nil {
		return errors.Wrap(err, "failed to fetch fork schedule")
	}
	if _, err := s.NodeVersion(ctx, &api.NodeVersionOpts{}); err != nil {
		return errors.Wrap(err, "failed to fetch node version")
	}

	return nil
}

// periodicClearStaticValues periodically sets static values to nil so they are
// refetched the next time they are required.
func (s *Service) periodicClearStaticValues(ctx context.Context) {
	go func(s *Service, ctx context.Context) {
		// Refreah every 5 minutes.
		refreshTicker := time.NewTicker(5 * time.Minute)
		defer refreshTicker.Stop()
		for {
			select {
			case <-refreshTicker.C:
				s.genesisMutex.Lock()
				s.genesis = nil
				s.genesisMutex.Unlock()
				s.specMutex.Lock()
				s.spec = nil
				s.specMutex.Unlock()
				s.depositContractMutex.Lock()
				s.depositContract = nil
				s.depositContractMutex.Unlock()
				s.forkScheduleMutex.Lock()
				s.forkSchedule = nil
				s.forkScheduleMutex.Unlock()
				s.nodeVersionMutex.Lock()
				s.nodeVersion = ""
				s.nodeVersionMutex.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}(s, ctx)
}

// checkDVT checks if connected to DVT middleware and sets
// internal flags appropriately.
func (s *Service) checkDVT(ctx context.Context) error {
	response, err := s.NodeVersion(ctx, &api.NodeVersionOpts{})
	if err != nil {
		return errors.Wrap(err, "failed to obtain node version for DVT check")
	}

	version := strings.ToLower(response.Data)

	if strings.Contains(version, "charon") {
		s.connectedToDVTMiddleware = true
	}

	return nil
}

// Name provides the name of the service.
func (s *Service) Name() string {
	return "Standard (HTTP)"
}

// Address provides the address for the connection.
func (s *Service) Address() string {
	return s.address
}

// close closes the service, freeing up resources.
func (s *Service) close() {
}
