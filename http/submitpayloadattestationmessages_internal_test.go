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
	"encoding/json"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func gloasMessage(validatorIndex phase0.ValidatorIndex) *spec.VersionedPayloadAttestationMessage {
	return &spec.VersionedPayloadAttestationMessage{
		Version: spec.DataVersionGloas,
		Gloas: &gloas.PayloadAttestationMessage{
			ValidatorIndex: validatorIndex,
			Data: &gloas.PayloadAttestationData{
				BeaconBlockRoot:   phase0.Root{0x01},
				Slot:              60300,
				PayloadPresent:    true,
				BlobDataAvailable: true,
			},
			Signature: phase0.BLSSignature{0xaa},
		},
	}
}

// TestCreateUnversionedPayloadAttestationMessages covers the step that strips
// the version wrapper before the messages go on the wire.  It is pure, so it
// needs no beacon node.
func TestCreateUnversionedPayloadAttestationMessages(t *testing.T) {
	tests := []struct {
		name     string
		messages []*spec.VersionedPayloadAttestationMessage
		expected int
		err      string
	}{
		{
			name:     "Single",
			messages: []*spec.VersionedPayloadAttestationMessage{gloasMessage(8)},
			expected: 1,
		},
		{
			name:     "Multiple",
			messages: []*spec.VersionedPayloadAttestationMessage{gloasMessage(8), gloasMessage(18)},
			expected: 2,
		},
		{
			name:     "NilMessage",
			messages: []*spec.VersionedPayloadAttestationMessage{nil},
			err:      "nil payload attestation message supplied",
		},
		{
			// Naming a version whose arm is nil must be rejected rather than
			// appended: an unchecked arm marshals to a body of [null], which
			// the node accepts as well-formed JSON and then mis-handles.
			name: "NilArm",
			messages: []*spec.VersionedPayloadAttestationMessage{
				{Version: spec.DataVersionGloas},
			},
			err: "no gloas payload attestation message supplied",
		},
		{
			name: "MixedVersions",
			messages: []*spec.VersionedPayloadAttestationMessage{
				gloasMessage(8),
				{Version: spec.DataVersionElectra},
			},
			err: "payload attestation messages must all be of the same version",
		},
		{
			name: "UnsupportedVersion",
			messages: []*spec.VersionedPayloadAttestationMessage{
				{Version: spec.DataVersionElectra},
			},
			err: "unsupported payload attestation message version electra",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			unversioned, version, err := createUnversionedPayloadAttestationMessages(test.messages)
			if test.err != "" {
				require.ErrorContains(t, err, test.err)
				require.ErrorIs(t, err, client.ErrInvalidOptions)

				return
			}

			require.NoError(t, err)
			require.Len(t, unversioned, test.expected)
			require.Equal(t, spec.DataVersionGloas, version)

			// The body must be an array of bare messages, with no version
			// wrapper and no nulls in it.
			encoded, err := json.Marshal(unversioned)
			require.NoError(t, err)
			require.NotContains(t, string(encoded), "null")
			require.NotContains(t, string(encoded), "\"version\"")
			require.Contains(t, string(encoded), "\"validator_index\":\"8\"")
		})
	}
}
