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

package multi_test

import (
	"context"
	"testing"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/multi"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestBuilders(t *testing.T) {
	ctx := context.Background()

	client1, err := mock.New(ctx, mock.WithName("mock 1"))
	require.NoError(t, err)
	builders := map[gloas.BuilderIndex]*apiv1.Builder{3: {Index: 3, Status: apiv1.BuilderStatusActive, Builder: &gloas.Builder{}}}
	var received *api.BuildersOpts
	client1.BuildersFunc = func(_ context.Context, opts *api.BuildersOpts) (*api.Response[map[gloas.BuilderIndex]*apiv1.Builder], error) {
		received = opts
		return &api.Response[map[gloas.BuilderIndex]*apiv1.Builder]{Data: builders, Metadata: map[string]any{}}, nil
	}
	client2, err := mock.New(ctx, mock.WithName("mock 2"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client1, client2}),
	)
	require.NoError(t, err)

	opts := &api.BuildersOpts{State: "head", Indices: []gloas.BuilderIndex{3}}
	res, err := multiClient.(consensusclient.BuildersProvider).Builders(ctx, opts)
	require.NoError(t, err)
	require.Same(t, opts, received)
	require.Equal(t, builders, res.Data)
}
