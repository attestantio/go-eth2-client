// Copyright © 2020 Attestant Limited.
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

package v1

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
)

// Service is an Ethereum 2 client service.
type Service struct {
	// Hold the initialising context to use for streams.
	ctx context.Context

	base    *url.URL
	client  *http.Client
	timeout time.Duration

	// Various information from the node that does not change during the
	// lifetime of a beacon node.
	// TODO needs to be refreshed on a reconnect.
	genesis         *api.Genesis
	spec            map[string]interface{}
	depositContract *api.DepositContract
	forkSchedule    []*spec.Fork
	nodeVersion     string

	// Event handlers.
	beaconChainHeadUpdatedMutex    sync.RWMutex
	beaconChainHeadUpdatedHandlers []client.BeaconChainHeadUpdatedHandler
}

// log is a service-wide logger.
var log zerolog.Logger

// New creates a new Ethereum 2 client service, connecting with a standard HTTP.
func New(ctx context.Context, params ...Parameter) (*Service, error) {
	parameters, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	// Set logging.
	log = zerologger.With().Str("service", "client").Str("impl", "standardv1").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:        64,
			MaxIdleConnsPerHost: 64,
			IdleConnTimeout:     384 * time.Second,
		},
	}

	base, err := url.Parse(parameters.address)
	if err != nil {
		return nil, errors.Wrap(err, "invalid URL")
	}

	s := &Service{
		ctx:     ctx,
		base:    base,
		client:  client,
		timeout: parameters.timeout,
	}

	// Fetch static values to confirm the connection is good.
	if err := s.fetchStaticValues(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to confirm node connection")
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

	return nil
}

// Name provides the name of the service.
func (s *Service) Name() string {
	return "Standard (HTTP)"
}

// close closes the service, freeing up resources.
func (s *Service) close() {
}
