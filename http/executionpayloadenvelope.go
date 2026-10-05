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
	"net/http"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"go.opentelemetry.io/otel"
)

// ExecutionPayloadEnvelope obtains a cached execution payload envelope.
func (s *Service) ExecutionPayloadEnvelope(ctx context.Context,
	opts *api.ExecutionPayloadEnvelopeOpts,
) (
	*api.Response[*spec.VersionedExecutionPayloadEnvelope],
	error,
) {
	ctx, span := otel.Tracer("attestantio.go-eth2-client.http").Start(ctx, "ExecutionPayloadEnvelope")
	defer span.End()

	httpResponse, err := s.fetchExecutionPayloadEnvelope(ctx, opts)
	if err != nil {
		return nil, err
	}

	response, err := s.executionPayloadEnvelopeFromResponse(ctx, httpResponse)
	if err != nil {
		return nil, err
	}

	if err := assertEnvelopeIsForBlock(response.Data, opts.BeaconBlockRoot); err != nil {
		return nil, err
	}

	return response, nil
}

// assertEnvelopeIsForBlock checks that an envelope belongs to the requested block.
func assertEnvelopeIsForBlock(envelope *spec.VersionedExecutionPayloadEnvelope, expected phase0.Root) error {
	root, err := envelope.BeaconBlockRoot()
	if err != nil {
		return err
	}

	if root != expected {
		return errors.Join(
			fmt.Errorf("execution payload envelope for block %#x; expected %#x", root, expected),
			client.ErrInconsistentResult,
		)
	}

	return nil
}

// fetchExecutionPayloadEnvelope validates the options and performs the GET,
// translating the endpoint's defined 404 into a sentinel.
func (s *Service) fetchExecutionPayloadEnvelope(ctx context.Context,
	opts *api.ExecutionPayloadEnvelopeOpts,
) (*httpResponse, error) {
	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	if opts.Slot == 0 {
		return nil, errors.Join(errors.New("no slot specified"), client.ErrInvalidOptions)
	}

	if opts.BeaconBlockRoot.IsZero() {
		return nil, errors.Join(errors.New("no beacon block root specified"), client.ErrInvalidOptions)
	}

	endpoint := fmt.Sprintf("/eth/v1/validator/execution_payload_envelopes/%d/%#x", opts.Slot, opts.BeaconBlockRoot)

	httpResponse, err := s.getWithResponseLimit(ctx, endpoint, "", &opts.Common, true, maxEPBSResponseSize)
	if err != nil {
		return nil, notFoundToSentinel(err)
	}

	return httpResponse, nil
}

// notFoundToSentinel preserves a 404 API error while adding the cache-miss sentinel.
func notFoundToSentinel(err error) error {
	var apiErr *api.Error
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
		return errors.Join(client.ErrNoExecutionPayloadEnvelope, apiErr)
	}

	return err
}

// executionPayloadEnvelopeFromResponse decodes a fetched envelope response.  It
// is separate from the fetch so that both encodings stay testable without a
// node holding a cached envelope, which it will only do for one slot at a time.
func (s *Service) executionPayloadEnvelopeFromResponse(ctx context.Context,
	res *httpResponse,
) (
	*api.Response[*spec.VersionedExecutionPayloadEnvelope],
	error,
) {
	if res.consensusVersion != spec.DataVersionGloas {
		return nil, fmt.Errorf("execution payload envelope not available for version %s", res.consensusVersion)
	}

	// The envelope is a gloas-onwards container, so the wire bytes always parse
	// into *gloas.ExecutionPayloadEnvelope.
	envelope := &gloas.ExecutionPayloadEnvelope{}
	metadata := metadataFromHeaders(res.headers)

	switch res.contentType {
	case ContentTypeSSZ:
		// Decoded through the request-scoped codec, as the sibling beacon-side
		// envelope endpoint does.  Measured against a minimal-preset node this
		// container encodes byte-identically to mainnet, since nothing in its
		// fixed part is preset-derived, so the generated codec would read it too.
		// The dynamic one is kept because it cannot be wrong for a preset this
		// code has not been run against, and the endpoint is called at most once
		// per proposal, so the codec build does not sit in a hot path.
		ds, err := s.dynSSZForRequest(ctx)
		if err != nil {
			return nil, err
		}

		if err := ds.UnmarshalSSZ(envelope, res.body); err != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to decode %s SSZ execution payload envelope", res.consensusVersion),
				err,
			)
		}
	case ContentTypeJSON:
		// Seeded with a typed nil pointer so that a body with no data key in it
		// is distinguishable from one carrying a zero-valued envelope.
		decoded, jsonMetadata, err := decodeJSONResponse(bytes.NewReader(res.body), (*gloas.ExecutionPayloadEnvelope)(nil))
		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("failed to decode %s JSON execution payload envelope", res.consensusVersion),
				err,
			)
		}

		if decoded == nil {
			return nil, fmt.Errorf("no %s execution payload envelope in response", res.consensusVersion)
		}

		envelope = decoded

		maps.Copy(metadata, jsonMetadata)
	default:
		return nil, fmt.Errorf("unhandled content type %v", res.contentType)
	}

	return &api.Response[*spec.VersionedExecutionPayloadEnvelope]{
		Data: &spec.VersionedExecutionPayloadEnvelope{
			Version: res.consensusVersion,
			Gloas:   envelope,
		},
		Metadata: metadata,
	}, nil
}
