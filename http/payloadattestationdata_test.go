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

package http_test

import (
	"context"
	"errors"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/stretchr/testify/require"
)

// TestPayloadAttestationData exercises the endpoint against a live node.
//
// A beacon node only produces payload attestation data for its current slot,
// so the success path is not reachable on demand; it is covered in-process by
// TestPayloadAttestationDataFromResponse, which also covers the SSZ body and
// the 204 "no block seen" signal that a node will not emit to order.  What is
// verified here is that a structurally-valid request reaches the server and
// that its rejection surfaces as a typed api.Error rather than a decode
// failure.
func TestPayloadAttestationData(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	t.Run("NilOpts", func(t *testing.T) {
		_, err := service.(client.PayloadAttestationDataProvider).PayloadAttestationData(ctx, nil)
		require.ErrorIs(t, err, client.ErrNoOptions)
	})

	t.Run("PastSlot", func(t *testing.T) {
		// Slot 1 is long past, so the node must refuse to produce data for it.
		_, err := service.(client.PayloadAttestationDataProvider).PayloadAttestationData(ctx,
			&api.PayloadAttestationDataOpts{Slot: 1},
		)
		require.Error(t, err)

		var apiErr *api.Error
		require.True(t, errors.As(err, &apiErr), "expected a typed api.Error, got %v", err)
		require.Equal(t, 400, apiErr.StatusCode)
	})
}
