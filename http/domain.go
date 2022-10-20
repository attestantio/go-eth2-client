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
	"bytes"
	"context"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// Domain provides a domain for a given domain type at a given epoch.
func (s *Service) Domain(ctx context.Context, domainType phase0.DomainType, epoch phase0.Epoch) (phase0.Domain, error) {
	// Obtain the fork for the epoch.
	fork, err := s.forkAtEpoch(ctx, epoch)
	if err != nil {
		return phase0.Domain{}, errors.Wrap(err, "failed to obtain fork")
	}

	// Calculate the domain.
	var forkVersion phase0.Version
	if epoch < fork.Epoch {
		forkVersion = fork.PreviousVersion
	} else {
		forkVersion = fork.CurrentVersion
	}
	if len(forkVersion) != 4 {
		return phase0.Domain{}, errors.New("fork version is invalid")
	}

	forkData := &phase0.ForkData{
		CurrentVersion: forkVersion,
	}

	if !bytes.Equal(domainType[:], []byte{0x00, 0x00, 0x00, 0x01}) {
		// Use the chain's genesis validators root for non-application domain types.
		genesis, err := s.Genesis(ctx)
		if err != nil {
			return phase0.Domain{}, errors.Wrap(err, "failed to obtain genesis")
		}

		forkData.GenesisValidatorsRoot = genesis.GenesisValidatorsRoot
	}

	root, err := forkData.HashTreeRoot()
	if err != nil {
		return phase0.Domain{}, errors.Wrap(err, "failed to calculate signature domain")
	}

	var domain phase0.Domain
	copy(domain[:], domainType[:])
	copy(domain[4:], root[:])
	return domain, nil
}

// forkAtEpoch works through the fork schedule to obtain the current fork.
func (s *Service) forkAtEpoch(ctx context.Context, epoch phase0.Epoch) (*phase0.Fork, error) {
	forkSchedule, err := s.ForkSchedule(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain fork schedule")
	}

	if len(forkSchedule) == 0 {
		return nil, errors.New("no fork schedule returned")
	}

	currentFork := forkSchedule[0]
	for i := range forkSchedule {
		if forkSchedule[i].Epoch > epoch {
			break
		}
		currentFork = forkSchedule[i]
	}
	return currentFork, nil
}
