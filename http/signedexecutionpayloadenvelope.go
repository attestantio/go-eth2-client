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

	"github.com/attestantio/go-eth2-client/spec/gloas"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	dynssz "github.com/pk910/dynamic-ssz"
)

// SignedExecutionPayloadEnvelope fetches a signed execution payload envelope given a block ID.
func (s *Service) SignedExecutionPayloadEnvelope(ctx context.Context,
	opts *api.SignedExecutionPayloadEnvelopeOpts,
) (
	*api.Response[*gloas.SignedExecutionPayloadEnvelope],
	error,
) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}
	if opts.Block == "" {
		return nil, errors.Join(errors.New("no block specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/eth/v1/beacon/execution_payload/%s", opts.Block)
	httpResponse, err := s.get(ctx, endpoint, "", &opts.Common, true)
	if err != nil {
		return nil, err
	}

	var response *api.Response[*gloas.SignedExecutionPayloadEnvelope]
	switch httpResponse.contentType {
	case ContentTypeSSZ:
		response, err = s.signedExecutionPayloadEnvelopeFromSSZ(ctx, httpResponse)
	case ContentTypeJSON:
		response, err = s.signedExecutionPayloadEnvelopeFromJSON(httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) signedExecutionPayloadEnvelopeFromSSZ(ctx context.Context,
	res *httpResponse,
) (
	*api.Response[*gloas.SignedExecutionPayloadEnvelope],
	error,
) {
	response := &api.Response[*gloas.SignedExecutionPayloadEnvelope]{
		Metadata: metadataFromHeaders(res.headers),
	}

	var dynSSZ *dynssz.DynSsz
	if s.customSpecSupport {
		specs, err := s.Spec(ctx, &api.SpecOpts{})
		if err != nil {
			return nil, errors.Join(errors.New("failed to request specs"), err)
		}

		dynSSZ = dynssz.NewDynSsz(specs.Data)
	}

	if res.consensusVersion != spec.DataVersionGloas {
		return nil, fmt.Errorf("execution payload envelope not available for block version %s", res.consensusVersion)
	}

	var err error
	response.Data = &gloas.SignedExecutionPayloadEnvelope{}
	if s.customSpecSupport {
		err = dynSSZ.UnmarshalSSZ(response.Data, res.body)
	} else {
		err = response.Data.UnmarshalSSZ(res.body)
	}
	if err != nil {
		return nil, errors.Join(errors.New("failed to decode gloas signed execution payload envelope contents"), err)
	}

	return response, nil
}

func (*Service) signedExecutionPayloadEnvelopeFromJSON(res *httpResponse) (
	*api.Response[*gloas.SignedExecutionPayloadEnvelope],
	error,
) {
	response := &api.Response[*gloas.SignedExecutionPayloadEnvelope]{}

	if res.consensusVersion != spec.DataVersionGloas {
		return nil, fmt.Errorf("execution payload envelope not available for block version %s", res.consensusVersion)
	}

	var err error
	response.Data, response.Metadata, err = decodeJSONResponse(bytes.NewReader(res.body),
		&gloas.SignedExecutionPayloadEnvelope{},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}
