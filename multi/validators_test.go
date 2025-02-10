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
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/multi"
	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestValidators(t *testing.T) {
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
		res, err := multiClient.(consensusclient.ValidatorsProvider).Validators(ctx, &api.ValidatorsOpts{})
		require.NoError(t, err)
		require.NotNil(t, res)
	}
	// At this point we expect mock 3 to be in active (unless probability hates us).
	require.Equal(t, "mock 3", multiClient.Address())
}
