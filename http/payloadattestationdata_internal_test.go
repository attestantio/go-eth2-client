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
	"net/http"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestPayloadAttestationDataFromResponse covers the decode half of the
// endpoint, which takes an already-fetched response and so needs no beacon
// node.  Two of the paths here are unreachable against a live node — a 204 is
// only emitted when no block has been seen for the requested slot, and neither
// client serves this endpoint as SSZ — so this is the only place they get
// exercised at all.
func TestPayloadAttestationDataFromResponse(t *testing.T) {
	datum := &gloas.PayloadAttestationData{
		BeaconBlockRoot:   phase0.Root{0x01, 0x02, 0x03},
		Slot:              60300,
		PayloadPresent:    true,
		BlobDataAvailable: false,
	}

	sszBody, err := datum.MarshalSSZ()
	require.NoError(t, err)

	jsonBody := []byte(`{"version":"gloas","data":{"beacon_block_root":"0x0102030000000000000000000000000000000000000000000000000000000000","slot":"60300","payload_present":true,"blob_data_available":false}}`)

	t.Run("JSON", func(t *testing.T) {
		response, err := payloadAttestationDataFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             jsonBody,
			headers:          map[string]string{},
		})
		require.NoError(t, err)
		require.Equal(t, spec.DataVersionGloas, response.Data.Version)
		require.Equal(t, datum, response.Data.Gloas)
	})

	t.Run("SSZ", func(t *testing.T) {
		response, err := payloadAttestationDataFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeSSZ,
			consensusVersion: spec.DataVersionGloas,
			body:             sszBody,
			headers:          map[string]string{},
		})
		require.NoError(t, err)
		require.Equal(t, spec.DataVersionGloas, response.Data.Version)
		require.Equal(t, datum, response.Data.Gloas)
	})

	// A 204 means the node has seen no block for the slot and the validator
	// must not attest.  Returning an empty-but-valid response here would let a
	// caller that forgets to check sign a zero-valued datum, so it fails
	// closed with a sentinel the caller can match on.
	t.Run("NoContent", func(t *testing.T) {
		_, err := payloadAttestationDataFromResponse(&httpResponse{
			statusCode: http.StatusNoContent,
			headers:    map[string]string{},
		})
		require.ErrorIs(t, err, client.ErrNoPayloadAttestationData)
	})

	t.Run("PreGloasVersion", func(t *testing.T) {
		_, err := payloadAttestationDataFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionElectra,
			body:             jsonBody,
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "payload attestation data not available for version electra")
	})

	t.Run("UnknownContentType", func(t *testing.T) {
		_, err := payloadAttestationDataFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeUnknown,
			consensusVersion: spec.DataVersionGloas,
			body:             jsonBody,
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "unhandled content type")
	})

	// The JSON decoder fills in only the keys it finds, so a body carrying no
	// data key at all leaves the seeded datum untouched.  That must be an
	// error rather than a success wrapping a zero-valued datum, which a caller
	// cannot tell apart from a genuine vote for the zero root.
	t.Run("MissingDataKey", func(t *testing.T) {
		_, err := payloadAttestationDataFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             []byte(`{"version":"gloas"}`),
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "no gloas payload attestation data in response")
	})

	// An explicit null decodes to the same nil pointer as an absent key, so
	// the same check catches it.  It is pinned separately because it used to
	// take a different route: the arm came back nil and only the endpoint's
	// own Slot() call downstream refused it, which left this function
	// reporting success on a response carrying nothing to attest to.
	t.Run("NullData", func(t *testing.T) {
		_, err := payloadAttestationDataFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             []byte(`{"version":"gloas","data":null}`),
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "no gloas payload attestation data in response")
	})

	t.Run("CorruptJSON", func(t *testing.T) {
		_, err := payloadAttestationDataFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             jsonBody[:len(jsonBody)-1],
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "failed to decode gloas payload attestation data")
	})

	t.Run("CorruptSSZ", func(t *testing.T) {
		_, err := payloadAttestationDataFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeSSZ,
			consensusVersion: spec.DataVersionGloas,
			body:             sszBody[:len(sszBody)-1],
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "failed to decode gloas payload attestation data")
	})
}
