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

// AttestationData obtains attestation data given the options.
func (s *Service) AttestationData(ctx context.Context,
	opts *api.AttestationDataOpts,
) (
	*api.Response[*phase0.AttestationData],
	error,
) {
	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}

	endpoint := "/eth/v1/validator/attestation_data"
	query := fmt.Sprintf("slot=%d&committee_index=%d", opts.Slot, opts.CommitteeIndex)
	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common, false)
	if err != nil {
		return nil, err
	}

	switch httpResponse.contentType {
	case ContentTypeJSON:
		return s.attestationDataFromJSON(ctx, opts, httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}
}

func (*Service) attestationDataFromJSON(_ context.Context,
	opts *api.AttestationDataOpts,
	httpResponse *httpResponse,
) (
	*api.Response[*phase0.AttestationData],
	error,
) {
	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), phase0.AttestationData{})
	if err != nil {
		return nil, err
	}

	if err := verifyAttestationData(opts, &data); err != nil {
		return nil, err
	}

	return &api.Response[*phase0.AttestationData]{
		Metadata: metadata,
		Data:     &data,
	}, nil
}

func verifyAttestationData(opts *api.AttestationDataOpts, data *phase0.AttestationData) error {
	if data.Slot != opts.Slot {
		return errors.Join(
			fmt.Errorf("attestation data for slot %d; expected %d", data.Slot, opts.Slot),
			client.ErrInconsistentResult,
		)
	}
	if data.Index != opts.CommitteeIndex {
		return errors.Join(
			fmt.Errorf("attestation data for committee index %d; expected %d", data.Index, opts.CommitteeIndex),
			client.ErrInconsistentResult,
		)
	}

	return nil
}
