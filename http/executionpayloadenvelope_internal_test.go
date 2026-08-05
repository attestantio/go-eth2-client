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

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

// validExecutionPayloadEnvelope returns an envelope populated far enough for its
// own marshalers to run: BaseFeePerGas is dereferenced by the payload's.
func validExecutionPayloadEnvelope() *gloas.ExecutionPayloadEnvelope {
	return &gloas.ExecutionPayloadEnvelope{
		Payload:               &gloas.ExecutionPayload{BaseFeePerGas: uint256.NewInt(7)},
		ExecutionRequests:     &gloas.ExecutionRequests{},
		BuilderIndex:          12,
		BeaconBlockRoot:       phase0.Root{0x0a, 0x0b},
		ParentBeaconBlockRoot: phase0.Root{0x0c, 0x0d},
	}
}

// TestExecutionPayloadEnvelopeFromResponse covers the decode half of the
// endpoint.  A node caches an envelope for one slot at a time, so the success
// path cannot be provoked on demand against a live node; this is where it, the
// SSZ body and the error paths are exercised.
func TestExecutionPayloadEnvelopeFromResponse(t *testing.T) {
	ctx := context.Background()
	s := &Service{}

	envelope := validExecutionPayloadEnvelope()

	envelopeJSON, err := json.Marshal(envelope)
	require.NoError(t, err)

	jsonBody := fmt.Appendf(nil, `{"version":"gloas","data":%s}`, envelopeJSON)

	t.Run("JSON", func(t *testing.T) {
		response, err := s.executionPayloadEnvelopeFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             jsonBody,
			headers:          map[string]string{},
		})
		require.NoError(t, err)
		require.Equal(t, spec.DataVersionGloas, response.Data.Version)
		require.False(t, response.Data.IsEmpty())

		builderIndex, err := response.Data.BuilderIndex()
		require.NoError(t, err)
		require.Equal(t, gloas.BuilderIndex(12), builderIndex)

		root, err := response.Data.BeaconBlockRoot()
		require.NoError(t, err)
		require.Equal(t, phase0.Root{0x0a, 0x0b}, root)

		payload, err := response.Data.Payload()
		require.NoError(t, err)
		require.Equal(t, uint256.NewInt(7), payload.BaseFeePerGas)
	})

	t.Run("SSZ", func(t *testing.T) {
		sszBody, err := envelope.MarshalSSZ()
		require.NoError(t, err)

		response, err := s.executionPayloadEnvelopeFromResponse(ctx, &httpResponse{
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

	// An envelope without its payload is not a usable envelope: the whole point
	// of retrieving one is to publish the payload it carries.  The spec marks the
	// field required, and the container's own decoder enforces it, which matters
	// because a node still on the removed blinded-envelope schema sends a
	// payload_root in its place and encoding/json would otherwise ignore the
	// unrecognised key and hand back an envelope with no payload at all.
	t.Run("JSONWithoutPayload", func(t *testing.T) {
		blinded := []byte(`{"version":"gloas","data":{` +
			`"payload_root":"0xfd2b34f63745b836c3ce3153439315284de96536d78ceac0f62408c9b1eebfd1",` +
			`"execution_requests":{"deposits":[],"withdrawals":[],"consolidations":[],` +
			`"builder_deposits":[],"builder_exits":[]},` +
			`"builder_index":"18446744073709551615",` +
			`"beacon_block_root":"0xea5d1f1b8bd38f4c8f133689dcfc0d333bf5fcf96d791beffef6023cd5146860",` +
			`"parent_beacon_block_root":"0x60ee9ee6812eaa4b111e3d37f777eda3b3baf5b530afe7e4c1a34dcf31f156fc"}}`)

		_, err := s.executionPayloadEnvelopeFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             blinded,
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "failed to decode gloas JSON execution payload envelope")
		require.ErrorContains(t, err, "payload missing")
	})

	t.Run("PreGloasVersion", func(t *testing.T) {
		_, err := s.executionPayloadEnvelopeFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionFulu,
			body:             jsonBody,
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "execution payload envelope not available for version fulu")
	})

	t.Run("UnknownContentType", func(t *testing.T) {
		_, err := s.executionPayloadEnvelopeFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeUnknown,
			consensusVersion: spec.DataVersionGloas,
			body:             jsonBody,
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "unhandled content type")
	})

	// The decoder writes only the keys it finds, so a body with no data key at
	// all leaves the seeded nil pointer untouched.  Returning success there would
	// hand back a wrapper whose every accessor errors.
	t.Run("NoDatum", func(t *testing.T) {
		for _, body := range []string{
			`{"version":"gloas"}`,
			`{"version":"gloas","data":null}`,
		} {
			_, err := s.executionPayloadEnvelopeFromResponse(ctx, &httpResponse{
				statusCode:       http.StatusOK,
				contentType:      ContentTypeJSON,
				consensusVersion: spec.DataVersionGloas,
				body:             []byte(body),
				headers:          map[string]string{},
			})
			require.ErrorContains(t, err, "no gloas execution payload envelope in response")
		}
	})

	// The block root is the only field the request and the response have in
	// common, and it is what binds the envelope to the block the caller proposed.
	// A node answering with a different one hands back an envelope that cannot be
	// published, so it is refused rather than passed on.
	t.Run("BlockRootMismatch", func(t *testing.T) {
		wrapped := &spec.VersionedExecutionPayloadEnvelope{
			Version: spec.DataVersionGloas,
			Gloas:   validExecutionPayloadEnvelope(),
		}

		// The fixture's own root, so the matching case must be accepted.
		require.NoError(t, assertEnvelopeIsForBlock(wrapped, phase0.Root{0x0a, 0x0b}))

		err := assertEnvelopeIsForBlock(wrapped, phase0.Root{0xff, 0xee})
		require.ErrorIs(t, err, client.ErrInconsistentResult)
		require.ErrorContains(t, err, "execution payload envelope for block")
	})

	t.Run("CorruptSSZ", func(t *testing.T) {
		sszBody, err := envelope.MarshalSSZ()
		require.NoError(t, err)

		_, err = s.executionPayloadEnvelopeFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeSSZ,
			consensusVersion: spec.DataVersionGloas,
			body:             sszBody[:len(sszBody)-1],
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "failed to decode gloas SSZ execution payload envelope")
	})
}
