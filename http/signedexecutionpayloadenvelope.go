// Copyright © 2026 Attestant Limited.
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
	"maps"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	dynssz "github.com/pk910/dynamic-ssz"
)

// SignedExecutionPayloadEnvelope fetches a signed execution payload envelope
// given a block ID and decodes it directly into the per-fork view stored on a
// *spec.VersionedSignedExecutionPayloadEnvelope. No intermediate copy.
func (s *Service) SignedExecutionPayloadEnvelope(ctx context.Context,
	opts *api.SignedExecutionPayloadEnvelopeOpts,
) (
	*api.Response[*spec.VersionedSignedExecutionPayloadEnvelope],
	error,
) {
	httpResponse, err := s.fetchSignedExecutionPayloadEnvelope(ctx, opts)
	if err != nil {
		return nil, err
	}

	// The envelope is a gloas-onwards container, so the wire bytes always
	// parse into *gloas.SignedExecutionPayloadEnvelope.
	envelope := &gloas.SignedExecutionPayloadEnvelope{}
	metadata := metadataFromHeaders(httpResponse.headers)

	switch httpResponse.contentType {
	case ContentTypeSSZ:
		ds, err := s.dynSSZForRequest(ctx)
		if err != nil {
			return nil, err
		}

		if err := ds.UnmarshalSSZ(envelope, httpResponse.body); err != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to decode %s signed execution payload envelope", httpResponse.consensusVersion),
				err,
			)
		}
	case ContentTypeJSON:
		decoded, jsonMetadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), envelope)
		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to decode %s signed execution payload envelope", httpResponse.consensusVersion),
				err,
			)
		}

		envelope = decoded

		maps.Copy(metadata, jsonMetadata)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}

	return &api.Response[*spec.VersionedSignedExecutionPayloadEnvelope]{
		Data: &spec.VersionedSignedExecutionPayloadEnvelope{
			Version: httpResponse.consensusVersion,
			Gloas:   envelope,
		},
		Metadata: metadata,
	}, nil
}

// fetchSignedExecutionPayloadEnvelope performs the GET request for
// SignedExecutionPayloadEnvelope: validates opts, hits the endpoint, and
// rejects responses for forks where the envelope doesn't apply.
func (s *Service) fetchSignedExecutionPayloadEnvelope(ctx context.Context,
	opts *api.SignedExecutionPayloadEnvelopeOpts,
) (*httpResponse, error) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	if opts.Block == "" {
		return nil, errors.Join(errors.New("no block specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/eth/v1/beacon/execution_payload_envelopes/%s", opts.Block)

	httpResponse, err := s.get(ctx, endpoint, "", &opts.Common, true)
	if err != nil {
		return nil, err
	}

	if httpResponse.consensusVersion != spec.DataVersionGloas {
		return nil, fmt.Errorf("execution payload envelope not available for block version %s", httpResponse.consensusVersion)
	}

	return httpResponse, nil
}

// dynSSZForRequest returns the dynamic-SSZ codec to use for marshaling and
// unmarshaling ePBS payloads. With custom spec support it builds a codec from
// the node's fetched spec (needed for non-mainnet presets); otherwise it
// returns the shared global codec built from the default (mainnet) preset.
func (s *Service) dynSSZForRequest(ctx context.Context) (*dynssz.DynSsz, error) {
	if !s.customSpecSupport {
		return dynssz.GetGlobalDynSsz(), nil
	}

	specs, err := s.Spec(ctx, &api.SpecOpts{})
	if err != nil {
		return nil, errors.Join(errors.New("failed to request specs"), err)
	}

	return dynssz.NewDynSsz(specs.Data), nil
}
