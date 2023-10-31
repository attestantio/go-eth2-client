// Copyright © 2021, 2023 Attestant Limited.
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
	"fmt"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// AttestationPool obtains the attestation pool for the given options.
func (s *Service) AttestationPool(ctx context.Context,
	opts *api.AttestationPoolOpts,
) (
	*api.Response[[]*phase0.Attestation],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}

	url := fmt.Sprintf("/eth/v1/beacon/pool/attestations?slot=%d", opts.Slot)
	httpResponse, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, err
	}

	switch httpResponse.contentType {
	case ContentTypeJSON:
		return s.attestationPoolFromJSON(ctx, opts, httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}
}

func (s *Service) attestationPoolFromJSON(_ context.Context,
	opts *api.AttestationPoolOpts,
	httpResponse *httpResponse,
) (
	*api.Response[[]*phase0.Attestation],
	error,
) {
	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*phase0.Attestation{})
	if err != nil {
		return nil, err
	}

	if err := verifyAttestationPool(opts, data); err != nil {
		return nil, err
	}

	return &api.Response[[]*phase0.Attestation]{
		Metadata: metadata,
		Data:     data,
	}, nil
}

func verifyAttestationPool(opts *api.AttestationPoolOpts, data []*phase0.Attestation) error {
	for _, datum := range data {
		if datum.Data.Slot != opts.Slot {
			return errors.New("attestation data not for requested slot")
		}
	}

	return nil
}
