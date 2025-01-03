// Copyright © 2025 Attestant Limited.
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

// SubmitVersionedAttestations submits versioned attestations.
func (s *Service) SubmitVersionedAttestations(ctx context.Context,
	opts *api.SubmitAttestationsOpts,
) error {
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
	var version *spec.DataVersion
	var unversionedAttestations []any

	for i := range attestations {
		if attestations[i] == nil {
			return errors.Join(errors.New("nil attestation version supplied"), client.ErrInvalidOptions)
		}

		// Ensure consistent versioning.
		if version == nil {
			version = &attestations[i].Version
		} else if *version != attestations[i].Version {
			return errors.Join(errors.New("attestations must all be of the same version"), client.ErrInvalidOptions)
		}

		// Append to unversionedAttestations.
		switch attestations[i].Version {
		case spec.DataVersionElectra:
			unversionedAttestations = append(unversionedAttestations, attestations[i].Electra)
		default:
			return errors.Join(errors.New("unknown attestation version"), client.ErrInvalidOptions)
		}
	}

	specJSON, err := json.Marshal(unversionedAttestations)
	if err != nil {
		return errors.Join(errors.New("failed to marshal JSON"), err)
	}

	// TODO: It looks like this endpoint accepts both electra.Attestation and SingleAttestation containers.
	// Reference: https://github.com/ethereum/beacon-APIs/blob/cee75f936fb1c1d8b1daf68f9be8c4d463f9fde9/apis/beacon/pool/attestations.v2.yaml#L55-L85.
	// Should we consider introducing a SubmitSingleAttestations interface or even transform the versioned Attestation to SingleAttestation on the fly in this method?
	// I'm not sure what benefits we get from submitting the SingleAttestation, but I'm wary if we have a SubmitVersionedAttestation and a SubmitSingleAttestation, then
	// this could lead to a SubmitVersionedSingleAttestation etc and start to become quite a lot of overhead.
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
