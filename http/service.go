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
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
)

// Service is an Ethereum 2 client service.
type Service struct {
	// log is a service-wide logger.
	log zerolog.Logger

	// Hold the initialising context to use for streams.
	ctx context.Context

	base    *url.URL
	address string
	client  *http.Client
	timeout time.Duration

	// Various information from the node that does not change during the
	// lifetime of a beacon node.
	genesis              *api.Genesis
	genesisMutex         sync.RWMutex
	spec                 map[string]interface{}
	specMutex            sync.RWMutex
	depositContract      *api.DepositContract
	depositContractMutex sync.RWMutex
	forkSchedule         []*phase0.Fork
	forkScheduleMutex    sync.RWMutex
	nodeVersion          string
	nodeVersionMutex     sync.RWMutex

	// API support.
	supportsV2BeaconBlocks    bool
	supportsV2BeaconState     bool
	supportsV2ValidatorBlocks bool

	// User-specified chunk sizes.
	userIndexChunkSize  int
	userPubKeyChunkSize int
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

	client := &http.Client{
		Timeout: parameters.timeout,
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
		ctx:                 ctx,
		base:                base,
		address:             parameters.address,
		client:              client,
		timeout:             parameters.timeout,
		userIndexChunkSize:  parameters.indexChunkSize,
		userPubKeyChunkSize: parameters.pubKeyChunkSize,
	}

	// Fetch static values to confirm the connection is good.
	if err := s.fetchStaticValues(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to confirm node connection")
	}

	// Periodially refetch static values in case of client update.
	if err := s.periodicClearStaticValues(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to set update ticker")
	}

	// Handle flags for API versioning.
	if err := s.checkAPIVersioning(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to check API versioning")
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
	if _, err := s.Genesis(ctx); err != nil {
		return errors.Wrap(err, "failed to fetch genesis")
	}
	if _, err := s.Spec(ctx); err != nil {
		return errors.Wrap(err, "failed to fetch spec")
	}
	if _, err := s.DepositContract(ctx); err != nil {
		return errors.Wrap(err, "failed to fetch deposit contract")
	}
	if _, err := s.ForkSchedule(ctx); err != nil {
		return errors.Wrap(err, "failed to fetch fork schedule")
	}
	if _, err := s.NodeVersion(ctx); err != nil {
		return errors.Wrap(err, "failed to fetch node version")
	}

	return nil
}

// periodicClearStaticValues periodically sets static values to nil so they are
// refetched the next time they are required.
func (s *Service) periodicClearStaticValues(ctx context.Context) error {
	go func(s *Service, ctx context.Context) {
		// Refreah every 5 minutes.
		refreshTicker := time.NewTicker(5 * time.Minute)
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
	return nil
}

// checkAPIVersioning checks the versions of some APIs and sets
// internal flags appropriately.
func (s *Service) checkAPIVersioning(ctx context.Context) error {
	// Start by setting the API v2 flag for blocks and fetching block 0.
	s.supportsV2BeaconBlocks = true
	_, err := s.SignedBeaconBlock(ctx, "0")
	if err == nil {
		// It's good.  Assume that other V2 APIs introduced with Altair
		// are present.
		s.supportsV2BeaconState = true
		s.supportsV2ValidatorBlocks = true
	} else {
		// Assume this is down to the V2 endpoint missing rather than
		// some other failure.
		s.supportsV2BeaconBlocks = false
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
