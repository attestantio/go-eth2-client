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
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// AggregateAttestation fetches the aggregate attestation for the given options.
func (s *Service) AggregateAttestation(ctx context.Context,
	opts *api.AggregateAttestationOpts,
) (
	*api.Response[*phase0.Attestation],
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

	endpoint := "/eth/v1/validator/aggregate_attestation"
	query := fmt.Sprintf("slot=%d&attestation_data_root=%#x", opts.Slot, opts.AttestationDataRoot)
	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), phase0.Attestation{})
	if err != nil {
		return nil, err
	}

	// Confirm the attestation is for the requested slot.
	if data.Data.Slot != opts.Slot {
		return nil, errors.Join(fmt.Errorf("aggregate attestation for slot %d; expected %d", data.Data.Slot, opts.Slot), client.ErrInconsistentResult)
	}

	// Confirm the attestation data is correct.
	dataRoot, err := data.Data.HashTreeRoot()
	if err != nil {
		return nil, errors.Join(errors.New("failed to obtain hash tree root of aggregate attestation data"), err)
	}
	if !bytes.Equal(dataRoot[:], opts.AttestationDataRoot[:]) {
		return nil, errors.Join(fmt.Errorf("aggregate attestation has data root %#x; expected %#x", dataRoot[:], opts.AttestationDataRoot[:]), client.ErrInconsistentResult)
	}

	return &api.Response[*phase0.Attestation]{
		Metadata: metadata,
		Data:     &data,
	}, nil
}
