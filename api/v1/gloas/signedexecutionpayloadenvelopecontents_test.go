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

package gloas_test

import (
	"testing"

	apiGloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	require "github.com/stretchr/testify/require"
)

func TestSignedExecutionPayloadEnvelopeContentsYAML(t *testing.T) {
	input := &apiGloas.SignedExecutionPayloadEnvelopeContents{
		KZGProofs: []deneb.KZGProof{{0x01}},
		Blobs:     []deneb.Blob{{0x02}},
	}

	marshaled, err := input.MarshalYAML()
	require.NoError(t, err)

	var obtained apiGloas.SignedExecutionPayloadEnvelopeContents
	require.NoError(t, obtained.UnmarshalYAML(marshaled))
	require.Equal(t, input, &obtained)
}
