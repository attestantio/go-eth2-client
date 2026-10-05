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
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

// TestMarshalRequestBody pins the request-body content-type negotiation shared
// by the ePBS submitters: SSZ by default, JSON when enforced, and the error
// wrapping of each arm.
//
// This needs no beacon node. The negotiation is pure logic given a Service:
// with customSpecSupport false, dynSSZForRequest returns the process-global
// codec and makes no network call, and neither arm touches any other field.
func TestMarshalRequestBody(t *testing.T) {
	ctx := context.Background()

	bid := &gloas.SignedExecutionPayloadBid{
		Message: &gloas.ExecutionPayloadBid{},
	}

	// Expected SSZ from the type's own generated static codec: fastssz, a
	// different implementation from the reflective dynssz codec under test.
	wantSSZ, err := bid.MarshalSSZ()
	require.NoError(t, err)

	t.Run("SSZByDefault", func(t *testing.T) {
		s := &Service{}

		body, contentType, err := s.marshalRequestBody(ctx, bid)
		require.NoError(t, err)
		require.Equal(t, ContentTypeSSZ, contentType)
		require.Equal(t, wantSSZ, body)
	})

	t.Run("JSONWhenEnforced", func(t *testing.T) {
		s := &Service{enforceJSON: true}

		body, contentType, err := s.marshalRequestBody(ctx, bid)
		require.NoError(t, err)
		require.Equal(t, ContentTypeJSON, contentType)

		// JSONEq rather than a decode-and-deep-equal round trip: this type's
		// JSON codec normalises a nil BlobKZGCommitments to an empty slice, so
		// a round trip would fail on the codec, not on the negotiation.
		wantJSON, err := bid.MarshalJSON()
		require.NoError(t, err)
		require.JSONEq(t, string(wantJSON), string(body))
	})

	t.Run("JSONMarshalFailure", func(t *testing.T) {
		s := &Service{enforceJSON: true}

		body, contentType, err := s.marshalRequestBody(ctx, make(chan int))
		require.EqualError(t, err, "failed to marshal JSON\njson: unsupported type: chan int")
		require.Nil(t, body)
		require.Equal(t, ContentTypeUnknown, contentType)
	})

	t.Run("SSZMarshalFailure", func(t *testing.T) {
		s := &Service{}

		body, contentType, err := s.marshalRequestBody(ctx, make(chan int))
		// ErrorContains, not EqualError: the text past our own wrapper comes
		// from dynssz, whose wording this package does not own.
		require.ErrorContains(t, err, "failed to marshal SSZ")
		require.Nil(t, body)
		require.Equal(t, ContentTypeUnknown, contentType)
	})
}
