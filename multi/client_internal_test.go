// Copyright © 2022 Attestant Limited.
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

package multi

import (
	"context"
	"sync"
	"testing"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestDeactivateMulti ensures that multiple concurrent calls to deactivateClient
// do not result in a bad list of active and inactive clients.
func TestDeactivateMulti(t *testing.T) {
	ctx := context.Background()

	consensusClient, err := mock.New(ctx)
	require.NoError(t, err)
	erroringClient1, err := testclients.NewErroring(ctx, 0.00000000001, consensusClient)
	require.NoError(t, err)
	erroringClient2, err := testclients.NewErroring(ctx, 0.00000000002, consensusClient)
	require.NoError(t, err)
	erroringClient3, err := testclients.NewErroring(ctx, 0.00000000003, consensusClient)
	require.NoError(t, err)

	s, err := New(ctx,
		WithLogLevel(zerolog.Disabled),
		WithClients([]consensusclient.Service{
			erroringClient1,
			erroringClient2,
			erroringClient3,
			consensusClient,
		}),
	)
	require.NoError(t, err)
	multi := s.(*Service)

	var wg sync.WaitGroup
	starter := make(chan interface{})
	for i := 0; i < 256; i++ {
		wg.Add(1)
		go func() {
			<-starter
			defer wg.Done()
			multi.deactivateClient(ctx, erroringClient1)
		}()
	}
	close(starter)
	wg.Wait()

	require.Len(t, multi.activeClients, 3)
	require.Len(t, multi.inactiveClients, 1)
	require.Equal(t, erroringClient1, multi.inactiveClients[0])
}

// TestActivateMulti ensures that multiple concurrent calls to activateClient
// do not result in a bad list of active and inactive clients.
func TestActivateMulti(t *testing.T) {
	ctx := context.Background()

	consensusClient, err := mock.New(ctx)
	require.NoError(t, err)
	erroringClient1, err := testclients.NewErroring(ctx, 0.00000000001, consensusClient)
	require.NoError(t, err)
	erroringClient2, err := testclients.NewErroring(ctx, 0.00000000002, consensusClient)
	require.NoError(t, err)
	erroringClient3, err := testclients.NewErroring(ctx, 0.00000000003, consensusClient)
	require.NoError(t, err)

	s, err := New(ctx,
		WithLogLevel(zerolog.Disabled),
		WithClients([]consensusclient.Service{
			erroringClient1,
			erroringClient2,
			erroringClient3,
			consensusClient,
		}),
	)
	require.NoError(t, err)
	multi := s.(*Service)

	multi.deactivateClient(ctx, erroringClient1)
	multi.deactivateClient(ctx, erroringClient2)

	var wg sync.WaitGroup
	starter := make(chan interface{})
	for i := 0; i < 256; i++ {
		wg.Add(1)
		go func() {
			<-starter
			defer wg.Done()
			multi.activateClient(ctx, erroringClient2)
		}()
	}
	close(starter)
	wg.Wait()

	require.Len(t, multi.activeClients, 3)
	require.Len(t, multi.inactiveClients, 1)
	require.Equal(t, erroringClient1, multi.inactiveClients[0])
}

// TestRecheck tests the recheck functionality when no nodes are available.
func TestRecheck(t *testing.T) {
	ctx := context.Background()

	consensusClient, err := mock.New(ctx)
	require.NoError(t, err)

	s, err := New(ctx,
		WithLogLevel(zerolog.Disabled),
		WithClients([]consensusclient.Service{
			consensusClient,
		}),
	)
	require.NoError(t, err)
	multi := s.(*Service)

	_, err = s.(consensusclient.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
	require.NoError(t, err)

	multi.deactivateClient(ctx, consensusClient)

	_, err = s.(consensusclient.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
	// Should re-activate in recheck so not return an error.
	require.NoError(t, err)
}
