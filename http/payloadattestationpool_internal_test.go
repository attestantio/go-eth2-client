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
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// poolAttestationJSON builds one wire-format payload attestation for the given
// slot.  Aggregation bits are a 64-byte bitvector and the signature 96 bytes.
func poolAttestationJSON(slot phase0.Slot) string {
	return fmt.Sprintf(
		`{"aggregation_bits":"0x08%s","data":{"beacon_block_root":"0x0102030000000000000000000000000000000000000000000000000000000000","slot":"%d","payload_present":true,"blob_data_available":true},"signature":"0xaa%s"}`,
		strings.Repeat("00", 63),
		slot,
		strings.Repeat("00", 95),
	)
}

func TestPayloadAttestationPoolFromResponse(t *testing.T) {
	body := func(attestations ...string) []byte {
		return []byte(`{"version":"gloas","data":[` + strings.Join(attestations, ",") + `]}`)
	}

	t.Run("JSON", func(t *testing.T) {
		response, err := payloadAttestationPoolFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body(poolAttestationJSON(60300), poolAttestationJSON(60301)),
			headers:          map[string]string{},
		})
		require.NoError(t, err)
		require.Len(t, response.Data, 2)

		for i, expectedSlot := range []phase0.Slot{60300, 60301} {
			require.Equal(t, spec.DataVersionGloas, response.Data[i].Version)

			data, err := response.Data[i].Data()
			require.NoError(t, err)
			require.Equal(t, expectedSlot, data.Slot)

			// Every element must carry its own bits, not just the first: a loop that
			// set them once outside the body would still pass a length check.  Read
			// through the container's accessor, not the bitvector's mainnet-only BitAt,
			// which this mainnet-width fixture would hide.
			require.NotNil(t, response.Data[i].Gloas)
			attested, err := response.Data[i].Gloas.AttestedAt(3)
			require.NoError(t, err)
			require.True(t, attested)
		}
	})

	t.Run("EmptyPool", func(t *testing.T) {
		response, err := payloadAttestationPoolFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body(),
			headers:          map[string]string{},
		})
		require.NoError(t, err)
		require.Empty(t, response.Data)
	})

	t.Run("PreGloasVersion", func(t *testing.T) {
		_, err := payloadAttestationPoolFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionFulu,
			body:             body(),
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "payload attestations not available for version fulu")
	})

	t.Run("UnknownContentType", func(t *testing.T) {
		_, err := payloadAttestationPoolFromResponse(&httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeUnknown,
			consensusVersion: spec.DataVersionGloas,
			body:             body(),
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "unhandled content type")
	})
}

// TestVerifyPayloadAttestationPool covers the filter check applied when a slot
// was requested.  It is pure, so it needs no beacon node.
func TestVerifyPayloadAttestationPool(t *testing.T) {
	attestation := func(slot phase0.Slot) *spec.VersionedPayloadAttestation {
		return &spec.VersionedPayloadAttestation{
			Version: spec.DataVersionGloas,
			Gloas: &gloas.PayloadAttestation{
				Data: &gloas.PayloadAttestationData{Slot: slot},
			},
		}
	}
	slot := phase0.Slot(60300)
	other := phase0.Slot(60301)

	tests := []struct {
		name         string
		slot         *phase0.Slot
		attestations []*spec.VersionedPayloadAttestation
		err          string
	}{
		{
			name:         "NoSlotRequested",
			attestations: []*spec.VersionedPayloadAttestation{attestation(60300), attestation(60301)},
		},
		{
			name:         "SlotMatches",
			slot:         &slot,
			attestations: []*spec.VersionedPayloadAttestation{attestation(60300)},
		},
		{
			name:         "SlotDiffers",
			slot:         &other,
			attestations: []*spec.VersionedPayloadAttestation{attestation(60300)},
			err:          "payload attestation for slot 60300; expected 60301",
		},
		{
			// A datum-less attestation must not read as "matches whatever was
			// asked for"; the filter has nothing to compare against.
			name:         "NilData",
			slot:         &slot,
			attestations: []*spec.VersionedPayloadAttestation{{Version: spec.DataVersionGloas, Gloas: &gloas.PayloadAttestation{}}},
			err:          "no gloas payload attestation data",
		},
		{
			// Asking for the whole pool waives the slot filter, not the check
			// that there is an attestation there at all: a nil arm reaches the
			// caller as an element whose every accessor errors.
			name:         "NilArmNoSlotRequested",
			attestations: []*spec.VersionedPayloadAttestation{{Version: spec.DataVersionGloas}},
			err:          "no gloas payload attestation",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := verifyPayloadAttestationPool(test.slot, test.attestations)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
