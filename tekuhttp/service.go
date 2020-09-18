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

package tekuhttp

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	zerologger "github.com/rs/zerolog/log"
)

// Service is an Ethereum 2 client service.
type Service struct {
	// Hold the initialising context to allow for streams to use it.
	ctx context.Context

	base    *url.URL
	client  *http.Client
	timeout time.Duration

	// Various information from the node that never changes once we have it.
	genesisTime                   *time.Time
	genesisValidatorsRoot         []byte
	slotDuration                  *time.Duration
	slotsPerEpoch                 *uint64
	farFutureEpoch                *uint64
	targetAggregatorsPerCommittee *uint64
	beaconAttesterDomain          []byte
	beaconProposerDomain          []byte
	randaoDomain                  []byte
	depositDomain                 []byte
	voluntaryExitDomain           []byte
	selectionProofDomain          []byte
	aggregateAndProofDomain       []byte

	// Event handlers.
	beaconChainHeadUpdatedMutex    sync.RWMutex
	beaconChainHeadUpdatedHandlers []client.BeaconChainHeadUpdatedHandler
}

// log is a service-wide logger.
var log zerolog.Logger

// New creates a new Ethereum 2 client service, connecting with Teku HTTP.
func New(ctx context.Context, params ...Parameter) (*Service, error) {
	parameters, err := parseAndCheckParameters(params...)
	if err != nil {
		return nil, errors.Wrap(err, "problem with parameters")
	}

	// Set logging.
	log = zerologger.With().Str("service", "client").Str("impl", "tekuhttp").Logger()
	if parameters.logLevel != log.GetLevel() {
		log = log.Level(parameters.logLevel)
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        16,
			MaxIdleConnsPerHost: 16,
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

	return s, nil
}

// fetchStaticValues fetches values that never change.
// This caches the values, avoiding future API calls.
func (s *Service) fetchStaticValues(ctx context.Context) error {
	if _, err := s.GenesisTime(ctx); err != nil {
		return errors.Wrap(err, "failed to fetch genesis time")
	}
	//	if _, err := s.GenesisValidatorsRoot(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch genesis validators root")
	//	}
	//	if _, err := s.SlotDuration(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch slot duration")
	//	}
	//	if _, err := s.SlotsPerEpoch(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch slots per epoch")
	//	}
	//	if _, err := s.FarFutureEpoch(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch far future epoch")
	//	}
	//	if _, err := s.TargetAggregatorsPerCommittee(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch target aggregators per committee")
	//	}
	//	if _, err := s.BeaconAttesterDomain(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch beacon attester domain")
	//	}
	//	if _, err := s.BeaconProposerDomain(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch beacon proposer domain")
	//	}
	//	if _, err := s.RANDAODomain(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch RANDAO domain")
	//	}
	//	if _, err := s.DepositDomain(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch deposit domain")
	//	}
	//	if _, err := s.VoluntaryExitDomain(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch voluntary exit domain")
	//	}
	//	if _, err := s.SelectionProofDomain(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch selection proof domain")
	//	}
	//	if _, err := s.AggregateAndProofDomain(ctx); err != nil {
	//		return errors.Wrap(err, "failed to fetch aggregate and proof domain")
	//	}

	return nil
}

// Name provides the name of the service.
func (s *Service) Name() string {
	return "Teku (HTTP)"
}

// Close the service, freeing up resources.
func (s *Service) Close(ctx context.Context) error {
	return nil
}
