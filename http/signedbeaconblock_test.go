// Copyright © 2020, 2021 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestSignedBeaconBlockGloas verifies that a block read back from a gloas chain
// lands in the gloas arm and is usable.
//
// TestSignedBeaconBlock below tolerates any api.Error, since the node may not
// hold the block asked for, which also makes it tolerate a version the client
// cannot decode.  This asserts the arm directly instead.
func TestSignedBeaconBlockGloas(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	requireOnGloas(ctx, t, service)

	response, err := service.(client.SignedBeaconBlockProvider).SignedBeaconBlock(ctx,
		&api.SignedBeaconBlockOpts{Block: "head"},
	)
	require.NoError(t, err)
	require.Equal(t, spec.DataVersionGloas, response.Data.Version)
	require.NotNil(t, response.Data.Gloas)
	require.NotNil(t, response.Data.Gloas.Message)

	// Post-Gloas the block commits to an execution payload bid rather than
	// carrying a payload, so the bid is what must have survived decoding.
	require.NotNil(t, response.Data.Gloas.Message.Body.SignedExecutionPayloadBid)

	slot, err := response.Data.Slot()
	require.NoError(t, err)
	require.NotZero(t, slot)

	root, err := response.Data.Root()
	require.NoError(t, err)
	require.NotEqual(t, phase0.Root{}, root)
}

func TestSignedBeaconBlock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name string
		opts *api.SignedBeaconBlockOpts
	}{
		{
			name: "Good",
			opts: &api.SignedBeaconBlockOpts{Block: "head"},
		},
		{
			name: "WithProposerSlashing",
			opts: &api.SignedBeaconBlockOpts{Block: "139"},
		},
	}

	service := testService(ctx, t).(client.Service)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.SignedBeaconBlockProvider).SignedBeaconBlock(ctx, test.opts)
			if err != nil {
				// The beacon node we are talking to may not have the block.
				var apiError *api.Error
				require.True(t, errors.As(err, &apiError))
				require.Equal(t, 404, apiError.StatusCode)
			} else {
				require.NoError(t, err)
				require.NotNil(t, response)
			}
		})
	}
}
