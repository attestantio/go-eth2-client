// Copyright © 2020 - 2022 Attestant Limited.
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

package mock

import (
	"context"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
)

// Service is a mock Ethereum 2 client service, providing data locally.
type Service struct {
	name    string
	timeout time.Duration

	genesisTime time.Time

	// Various information from the node that does not change during the
	// lifetime of a beacon node.
	// genesis         *api.Genesis
	// spec            map[string]interface{}
	// depositContract *api.DepositContract
	// forkSchedule    []*phase0.Fork
	nodeVersion string

	// Values that can be altered if required.
	HeadSlot     phase0.Slot
	SyncDistance phase0.Slot
}

// log is a service-wide logger.
var log zerolog.Logger

// New creates a new Ethereum 2 client service, mocking connections.
func New(ctx context.Context, params ...Parameter) (*Service, error) {
	parameters, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	// Set logging.
	log = zerologger.With().Str("service", "client").Str("impl", "mock").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}

	s := &Service{
		name:        parameters.name,
		genesisTime: parameters.genesisTime,
		timeout:     parameters.timeout,
		nodeVersion: "mock",

		HeadSlot:     12345,
		SyncDistance: 0,
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
	if _, err := s.Genesis(ctx, &api.GenesisOpts{}); err != nil {
		return errors.Wrap(err, "failed to fetch genesis")
	}
	//	if _, err := s.Spec(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch spec")
	//	}
	//	if _, err := s.DepositContract(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch deposit contract")
	//	}
	//	if _, err := s.ForkSchedule(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch fork schedule")
	//	}

	return nil
}

// Name provides the name of the service.
func (s *Service) Name() string {
	return "Mock"
}

// Address provides the address of the service.
func (s *Service) Address() string {
	return s.name
}

// close closes the service, freeing up resources.
func (s *Service) close() {
}
