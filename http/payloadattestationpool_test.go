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
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"
)

func TestPayloadAttestationPool(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	slot := headSlot(ctx, t, service)

	t.Run("NilOpts", func(t *testing.T) {
		_, err := service.(client.PayloadAttestationPoolProvider).PayloadAttestationPool(ctx, nil)
		require.ErrorIs(t, err, client.ErrNoOptions)
	})

	t.Run("Good", func(t *testing.T) {
		response, err := service.(client.PayloadAttestationPoolProvider).PayloadAttestationPool(ctx,
			&api.PayloadAttestationPoolOpts{Slot: &slot},
		)
		require.NoError(t, err)
		require.NotNil(t, response)

		// The pool for a historical slot is usually empty; whatever it holds
		// must be a gloas attestation for the slot that was asked for.
		for _, attestation := range response.Data {
			require.Equal(t, spec.DataVersionGloas, attestation.Version)

			data, err := attestation.Data()
			require.NoError(t, err)
			require.Equal(t, slot, data.Slot)
		}
	})
}
