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

// SubmitAttestations submits versioned attestations.
func (s *Service) SubmitAttestations(ctx context.Context, opts *api.SubmitAttestationsOpts) error {
	if err := s.assertIsSynced(ctx); err != nil {
		return err
	}
	if opts == nil {
		return client.ErrNoOptions
	}
	if len(opts.Attestations) == 0 {
		return errors.Join(errors.New("no attestations supplied"), client.ErrInvalidOptions)
	}
	attestations := opts.Attestations
	unversionedAttestations, err := s.createUnversionedAttestations(attestations)
	if err != nil {
		return err
	}

	specJSON, err := json.Marshal(unversionedAttestations)
	if err != nil {
		return errors.Join(errors.New("failed to marshal JSON"), err)
	}

	endpoint := "/eth/v2/beacon/pool/attestations"
	query := ""

	headers := make(map[string]string)
	headers["Eth-Consensus-Version"] = strings.ToLower(attestations[0].Version.String())
	if _, err = s.post(ctx,
		endpoint,
		query,
		&opts.Common,
		bytes.NewReader(specJSON),
		ContentTypeJSON,
		headers,
	); err != nil {
		return errors.Join(errors.New("failed to submit versioned beacon attestations"), err)
	}

	return nil
}

func (s *Service) createUnversionedAttestations(attestations []*spec.VersionedAttestation) ([]any, error) {
	var version spec.DataVersion
	var unversionedAttestations []any

	for i := range attestations {
		if attestations[i] == nil {
			return nil, errors.Join(errors.New("nil attestation version supplied"), client.ErrInvalidOptions)
		}

		// Ensure consistent versioning.
		if version == spec.DataVersionUnknown {
			version = attestations[i].Version
		} else if version != attestations[i].Version {
			return nil, errors.Join(errors.New("attestations must all be of the same version"), client.ErrInvalidOptions)
		}

		// Append to unversionedAttestations.
		switch attestations[i].Version {
		case spec.DataVersionPhase0:
			unversionedAttestations = append(unversionedAttestations, attestations[i].Phase0)
		case spec.DataVersionAltair:
			unversionedAttestations = append(unversionedAttestations, attestations[i].Altair)
		case spec.DataVersionBellatrix:
			unversionedAttestations = append(unversionedAttestations, attestations[i].Bellatrix)
		case spec.DataVersionCapella:
			unversionedAttestations = append(unversionedAttestations, attestations[i].Capella)
		case spec.DataVersionDeneb:
			unversionedAttestations = append(unversionedAttestations, attestations[i].Deneb)
		case spec.DataVersionElectra:
			singleAttestation, err := attestations[i].Electra.ToSingleAttestation(attestations[i].ValidatorIndex)
			if err != nil {
				s.log.Warn().Err(err).Msg("Failed to convert attestation to single attestation")

				continue
			}
			unversionedAttestations = append(unversionedAttestations, singleAttestation)
		default:
			return nil, errors.Join(errors.New("unknown attestation version"), client.ErrInvalidOptions)
		}
	}

	return unversionedAttestations, nil
}
