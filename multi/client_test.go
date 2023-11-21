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

func TestClient(t *testing.T) {
	ctx := context.Background()

	consensusClient, err := mock.New(ctx)
	require.NoError(t, err)
	eth1ClientErroring1, err := testclients.NewErroring(ctx, 0.1, consensusClient)
	require.NoError(t, err)
	eth1ClientErroring2, err := testclients.NewErroring(ctx, 0.1, consensusClient)
	require.NoError(t, err)
	eth1ClientErroring3, err := testclients.NewErroring(ctx, 0.1, consensusClient)
	require.NoError(t, err)

	s, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{
			eth1ClientErroring1,
			eth1ClientErroring2,
			eth1ClientErroring3,
			consensusClient,
		}),
	)
	require.NoError(t, err)

	for i := 0; i < 1024; i++ {
		_, err := s.(consensusclient.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
		require.NoError(t, err)
	}
}
