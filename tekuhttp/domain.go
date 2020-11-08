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
	"fmt"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// Domain provides a domain for a given domain typ eat a given epoch.
func (s *Service) Domain(ctx context.Context, domainType spec.DomainType, epoch spec.Epoch) ([]byte, error) {
	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain slots per epoch")
	}

	// Obtain the fork.
	fork, err := s.Fork(ctx, fmt.Sprintf("%d", uint64(epoch)*slotsPerEpoch))
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain fork")
	}

	// Obtain the genesis validators root.
	genesisValidatorsRoot, err := s.GenesisValidatorsRoot(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain genesis validators root")
	}

	// Calculate the domain.
	var forkVersion spec.Version
	if epoch < fork.Epoch {
		forkVersion = fork.PreviousVersion
	} else {
		forkVersion = fork.CurrentVersion
	}
	if len(forkVersion) != 4 {
		return nil, errors.New("fork version is invalid")
	}

	forkData := &spec.ForkData{
		CurrentVersion: forkVersion,
	}
	copy(forkData.GenesisValidatorsRoot[:], genesisValidatorsRoot)
	root, err := forkData.HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate signature domain")
	}

	domain := make([]byte, 32)
	copy(domain, domainType[:])
	copy(domain[4:], root[:])
	return domain, nil
}
