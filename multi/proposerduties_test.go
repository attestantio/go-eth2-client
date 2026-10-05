// Copyright © 2021 Attestant Limited.
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

package multi_test

import (
	"context"
	"testing"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/multi"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type proposerDutiesV1Only struct {
	consensusclient.Service
}

func TestProposerDutiesV2FallsBackFromUnsupportedClient(t *testing.T) {
	ctx := context.Background()

	unsupportedClient, err := mock.New(ctx, mock.WithName("unsupported"))
	require.NoError(t, err)
	supportedClient, err := mock.New(ctx, mock.WithName("supported"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{
			&proposerDutiesV1Only{Service: unsupportedClient},
			supportedClient,
		}),
	)
	require.NoError(t, err)

	response, err := multiClient.(consensusclient.ProposerDutiesV2Provider).ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, "supported", multiClient.Address())
}

func TestProposerDutiesV2(t *testing.T) {
	ctx := context.Background()
	root := phase0.Root{0xaa}
	duties := []*v1.ProposerDuty{{ValidatorIndex: 7, Slot: 65}}

	mockClient, err := mock.New(ctx)
	require.NoError(t, err)
	mockClient.ProposerDutiesV2Func = func(_ context.Context, _ *api.ProposerDutiesOpts) (*api.Response[[]*v1.ProposerDuty], error) {
		return &api.Response[[]*v1.ProposerDuty]{
			Data: duties,
			Metadata: map[string]any{
				"dependent_root":       root,
				"execution_optimistic": true,
			},
		}, nil
	}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{mockClient}),
	)
	require.NoError(t, err)

	provider, supported := multiClient.(consensusclient.ProposerDutiesV2Provider)
	require.True(t, supported)

	response, err := provider.ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{Epoch: 2})
	require.NoError(t, err)
	require.Same(t, duties[0], response.Data[0])
	require.Equal(t, root, response.Metadata["dependent_root"])
	require.Equal(t, true, response.Metadata["execution_optimistic"])
}

func TestProposerDuties(t *testing.T) {
	ctx := context.Background()

	client1, err := mock.New(ctx, mock.WithName("mock 1"))
	require.NoError(t, err)
	erroringClient1, err := testclients.NewErroring(ctx, 0.1, client1)
	require.NoError(t, err)
	client2, err := mock.New(ctx, mock.WithName("mock 2"))
	require.NoError(t, err)
	erroringClient2, err := testclients.NewErroring(ctx, 0.1, client2)
	require.NoError(t, err)
	client3, err := mock.New(ctx, mock.WithName("mock 3"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{
			erroringClient1,
			erroringClient2,
			client3,
		}),
	)
	require.NoError(t, err)

	for i := 0; i < 128; i++ {
		res, err := multiClient.(consensusclient.ProposerDutiesProvider).ProposerDuties(ctx, &api.ProposerDutiesOpts{})
		require.NoError(t, err)
		require.NotNil(t, res)
	}
	// At this point we expect mock 3 to be in active (unless probability hates us).
	require.Equal(t, "mock 3", multiClient.Address())
}
