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
	"errors"
	"fmt"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// AggregateAttestation fetches the aggregate attestation for the given options.
func (s *Service) AggregateAttestation(ctx context.Context,
	opts *api.AggregateAttestationOpts,
) (
	*api.Response[*spec.VersionedAttestation],
	error,
) {
	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}
	if opts.AttestationDataRoot.IsZero() {
		return nil, errors.Join(errors.New("no attestation data root specified"), client.ErrInvalidOptions)
	}

	endpoint := "/eth/v2/validator/aggregate_attestation"
	query := fmt.Sprintf("slot=%d&attestation_data_root=%#x&committee_index=%d",
		opts.Slot, opts.AttestationDataRoot, opts.CommitteeIndex)
	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common, false)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeAggregateAttestation(httpResponse)
	if err != nil {
		return nil, err
	}

	// Confirm the attestation is for the requested slot.
	attestationData, err := data.Data()
	if err != nil {
		return nil,
			errors.Join(
				errors.New("failed to extract attestation data from response"),
				err,
			)
	}
	if attestationData.Slot != opts.Slot {
		return nil,
			errors.Join(
				fmt.Errorf("aggregate attestation for slot %d; expected %d", attestationData.Slot, opts.Slot),
				client.ErrInconsistentResult,
			)
	}

	// Confirm the attestation data is correct.
	dataRoot, err := attestationData.HashTreeRoot()
	if err != nil {
		return nil, errors.Join(errors.New("failed to obtain hash tree root of aggregate attestation data"), err)
	}
	if !bytes.Equal(dataRoot[:], opts.AttestationDataRoot[:]) {
		return nil, errors.Join(
			fmt.Errorf("aggregate attestation has data root %#x; expected %#x", dataRoot[:], opts.AttestationDataRoot[:]),
			client.ErrInconsistentResult,
		)
	}

	return &api.Response[*spec.VersionedAttestation]{
		Metadata: metadata,
		Data:     data,
	}, nil
}

func decodeAggregateAttestation(httpResponse *httpResponse) (*spec.VersionedAttestation, map[string]any, error) {
	var metadata map[string]any
	data := &spec.VersionedAttestation{
		Version: httpResponse.consensusVersion,
	}
	switch httpResponse.consensusVersion {
	case spec.DataVersionPhase0:
		phase0Data, phase0Metadata, decodeErr := decodeJSONResponse(bytes.NewReader(httpResponse.body), &phase0.Attestation{})
		metadata = phase0Metadata
		data.Phase0 = phase0Data
		if decodeErr != nil {
			return &spec.VersionedAttestation{}, nil, decodeErr
		}

		return data, metadata, nil
	case spec.DataVersionAltair:
		phase0Data, phase0Metadata, decodeErr := decodeJSONResponse(bytes.NewReader(httpResponse.body), &phase0.Attestation{})
		metadata = phase0Metadata
		data.Altair = phase0Data
		if decodeErr != nil {
			return &spec.VersionedAttestation{}, nil, decodeErr
		}

		return data, metadata, nil
	case spec.DataVersionBellatrix:
		phase0Data, phase0Metadata, decodeErr := decodeJSONResponse(bytes.NewReader(httpResponse.body), &phase0.Attestation{})
		metadata = phase0Metadata
		data.Bellatrix = phase0Data
		if decodeErr != nil {
			return &spec.VersionedAttestation{}, nil, decodeErr
		}

		return data, metadata, nil
	case spec.DataVersionCapella:
		phase0Data, phase0Metadata, decodeErr := decodeJSONResponse(bytes.NewReader(httpResponse.body), &phase0.Attestation{})
		metadata = phase0Metadata
		data.Capella = phase0Data
		if decodeErr != nil {
			return &spec.VersionedAttestation{}, nil, decodeErr
		}

		return data, metadata, nil
	case spec.DataVersionDeneb:
		phase0Data, phase0Metadata, decodeErr := decodeJSONResponse(bytes.NewReader(httpResponse.body), &phase0.Attestation{})
		metadata = phase0Metadata
		data.Deneb = phase0Data
		if decodeErr != nil {
			return &spec.VersionedAttestation{}, nil, decodeErr
		}

		return data, metadata, nil
	case spec.DataVersionElectra:
		electraData, electraMetadata, decodeErr := decodeJSONResponse(bytes.NewReader(httpResponse.body), &electra.Attestation{})
		metadata = electraMetadata
		data.Electra = electraData
		if decodeErr != nil {
			return &spec.VersionedAttestation{}, nil, decodeErr
		}

		return data, metadata, nil
	default:
		return &spec.VersionedAttestation{}, nil, errors.New("unknown consensus version")
	}
}
