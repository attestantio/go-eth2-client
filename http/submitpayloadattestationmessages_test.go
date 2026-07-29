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
	"strings"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestSubmitPayloadAttestationMessages(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	slot := headSlot(ctx, t, service)

	t.Run("NilOpts", func(t *testing.T) {
		err := service.(client.PayloadAttestationMessagesSubmitter).SubmitPayloadAttestationMessages(ctx, nil)
		require.ErrorIs(t, err, client.ErrNoOptions)
	})

	t.Run("NoMessages", func(t *testing.T) {
		err := service.(client.PayloadAttestationMessagesSubmitter).SubmitPayloadAttestationMessages(ctx,
			&api.SubmitPayloadAttestationMessagesOpts{},
		)
		require.ErrorContains(t, err, "no payload attestation messages supplied")
		require.ErrorIs(t, err, client.ErrInvalidOptions)
	})

	// Reaching the server matters more than the submission succeeding: only a
	// PTC member holding the right key can produce a valid signature, so the
	// node is expected to reject this.  What it rejects on is the point.  A
	// missing or wrong Eth-Consensus-Version header is refused before the
	// messages are looked at, so a rejection that names the signature is
	// evidence that the header and the body shape were both accepted.
	t.Run("ReachesServer", func(t *testing.T) {
		err := service.(client.PayloadAttestationMessagesSubmitter).SubmitPayloadAttestationMessages(ctx,
			&api.SubmitPayloadAttestationMessagesOpts{
				Messages: []*spec.VersionedPayloadAttestationMessage{
					{
						Version: spec.DataVersionGloas,
						Gloas: &gloas.PayloadAttestationMessage{
							ValidatorIndex: 8,
							Data: &gloas.PayloadAttestationData{
								BeaconBlockRoot:   phase0.Root{0x01},
								Slot:              slot,
								PayloadPresent:    true,
								BlobDataAvailable: true,
							},
							Signature: phase0.BLSSignature{0xa0},
						},
					},
				},
			},
		)
		require.Error(t, err)

		var apiErr *api.Error
		require.True(t, errors.As(err, &apiErr), "expected a typed api.Error, got %v", err)
		require.Equal(t, 400, apiErr.StatusCode)

		body := strings.ToLower(string(apiErr.Data))
		require.NotContains(t, body, "eth-consensus-version",
			"the node rejected the consensus version header, so the request never reached message validation")
		require.True(t,
			strings.Contains(body, "signature") || strings.Contains(body, "failed validation"),
			"expected a rejection from message validation, got %s", body)
	})
}
