// Copyright © 2020 - 2024 Attestant Limited.
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
	"encoding/json"
	"errors"
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
)

// SubmitAggregateAttestations submits aggregate attestations.
func (s *Service) SubmitAggregateAttestations(ctx context.Context, opts *api.SubmitAggregateAttestationsOpts) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}
	if opts == nil {
		return client.ErrNoOptions
	}
	if len(opts.SignedAggregateAndProofs) == 0 {
		return errors.Join(errors.New("no aggregate and proofs supplied"), client.ErrInvalidOptions)
	}
	aggregateAndProofs := opts.SignedAggregateAndProofs
	unversionedAggregates, err := createUnversionedAggregates(aggregateAndProofs)
	if err != nil {
		return err
	}

	specJSON, err := json.Marshal(unversionedAggregates)
	if err != nil {
		return errors.Join(errors.New("failed to marshal JSON"), err)
	}

	endpoint := "/eth/v2/validator/aggregate_and_proofs"
	query := ""

	headers := make(map[string]string)
	headers["Eth-Consensus-Version"] = strings.ToLower(aggregateAndProofs[0].Version.String())
	if _, err = s.post(ctx,
		endpoint,
		query,
		&opts.Common,
		bytes.NewReader(specJSON),
		ContentTypeJSON,
		headers,
	); err != nil {
		return errors.Join(errors.New("failed to submit versioned aggregate and proofs"), err)
	}

	return nil
}

func createUnversionedAggregates(aggregateAndProofs []*spec.VersionedSignedAggregateAndProof) ([]any, error) {
	var version spec.DataVersion
	var unversionedAggregates []any

	for i := range aggregateAndProofs {
		if aggregateAndProofs[i] == nil {
			return nil, errors.Join(errors.New("nil aggregate and proof version supplied"), client.ErrInvalidOptions)
		}

		// Ensure consistent versioning.
		if version == spec.DataVersionUnknown {
			version = aggregateAndProofs[i].Version
		} else if version != aggregateAndProofs[i].Version {
			return nil, errors.Join(errors.New("aggregate and proofs must all be of the same version"), client.ErrInvalidOptions)
		}

		// Append to unversionedAggregates.
		switch aggregateAndProofs[i].Version {
		case spec.DataVersionPhase0:
			unversionedAggregates = append(unversionedAggregates, aggregateAndProofs[i].Phase0)
		case spec.DataVersionAltair:
			unversionedAggregates = append(unversionedAggregates, aggregateAndProofs[i].Altair)
		case spec.DataVersionBellatrix:
			unversionedAggregates = append(unversionedAggregates, aggregateAndProofs[i].Bellatrix)
		case spec.DataVersionCapella:
			unversionedAggregates = append(unversionedAggregates, aggregateAndProofs[i].Capella)
		case spec.DataVersionDeneb:
			unversionedAggregates = append(unversionedAggregates, aggregateAndProofs[i].Deneb)
		case spec.DataVersionElectra:
			unversionedAggregates = append(unversionedAggregates, aggregateAndProofs[i].Electra)
		default:
			return nil, errors.Join(errors.New("unknown aggregate and proof version"), client.ErrInvalidOptions)
		}
	}

	return unversionedAggregates, nil
}
