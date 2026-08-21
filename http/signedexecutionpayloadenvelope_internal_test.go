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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// validSignedExecutionPayloadEnvelope wraps the package's shared valid envelope
// (defined alongside the unsigned endpoint's test) in a signature so its own
// marshalers run.
func validSignedExecutionPayloadEnvelope() *gloas.SignedExecutionPayloadEnvelope {
	return &gloas.SignedExecutionPayloadEnvelope{
		Message:   validExecutionPayloadEnvelope(),
		Signature: phase0.BLSSignature{0x01, 0x02},
	}
}

// TestSignedExecutionPayloadEnvelopeFromResponse covers the decode half of the
// endpoint.  Splitting decode from the fetch is what lets the SSZ body and the
// error paths be exercised without a node serving the route.
func TestSignedExecutionPayloadEnvelopeFromResponse(t *testing.T) {
	ctx := context.Background()
	s := &Service{}

	envelope := validSignedExecutionPayloadEnvelope()

	envelopeJSON, err := json.Marshal(envelope)
	require.NoError(t, err)

	jsonBody := fmt.Appendf(nil, `{"version":"gloas","data":%s}`, envelopeJSON)

	t.Run("JSON", func(t *testing.T) {
		response, err := s.signedExecutionPayloadEnvelopeFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             jsonBody,
			headers:          map[string]string{},
		})
		require.NoError(t, err)
		require.Equal(t, spec.DataVersionGloas, response.Data.Version)
		require.NotNil(t, response.Data.Gloas)
		require.False(t, response.Data.IsEmpty())
	})

	t.Run("SSZ", func(t *testing.T) {
		sszBody, err := envelope.MarshalSSZ()
		require.NoError(t, err)

		response, err := s.signedExecutionPayloadEnvelopeFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeSSZ,
			consensusVersion: spec.DataVersionGloas,
			body:             sszBody,
			headers:          map[string]string{},
		})
		require.NoError(t, err)

		// SSZ carries no field names, so a wrongly ordered codec shows up as a
		// differing hash tree root rather than a decode error.
		want, err := envelope.HashTreeRoot()
		require.NoError(t, err)
		got, err := response.Data.Gloas.HashTreeRoot()
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	// The regression this file exists for.  decodeJSONResponse only unmarshals
	// when a data key is present and does not error when it is absent, so a body
	// with no data key (or an explicit null) must be rejected rather than
	// returned as a zero-valued envelope whose every accessor reads garbage.
	// Before the typed-nil seed and this guard, the non-nil seed made both
	// bodies decode as success.
	t.Run("NoDatum", func(t *testing.T) {
		for _, body := range []string{
			`{"version":"gloas"}`,
			`{"version":"gloas","data":null}`,
		} {
			_, err := s.signedExecutionPayloadEnvelopeFromResponse(ctx, &httpResponse{
				statusCode:       http.StatusOK,
				contentType:      ContentTypeJSON,
				consensusVersion: spec.DataVersionGloas,
				body:             []byte(body),
				headers:          map[string]string{},
			})
			require.ErrorContains(t, err, "no gloas signed execution payload envelope in response")
		}
	})

	t.Run("CorruptSSZ", func(t *testing.T) {
		sszBody, err := envelope.MarshalSSZ()
		require.NoError(t, err)

		_, err = s.signedExecutionPayloadEnvelopeFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeSSZ,
			consensusVersion: spec.DataVersionGloas,
			body:             sszBody[:len(sszBody)-1],
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "failed to decode gloas signed execution payload envelope")
	})

	t.Run("UnknownContentType", func(t *testing.T) {
		_, err := s.signedExecutionPayloadEnvelopeFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeUnknown,
			consensusVersion: spec.DataVersionGloas,
			body:             jsonBody,
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "unhandled content type")
	})
}
