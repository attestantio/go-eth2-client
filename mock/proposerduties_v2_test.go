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

package mock_test

import (
	"context"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestProposerDutiesV2Callback(t *testing.T) {
	ctx := context.Background()
	expectedOpts := &api.ProposerDutiesOpts{Epoch: 4}
	expectedResponse := &api.Response[[]*apiv1.ProposerDuty]{
		Data: []*apiv1.ProposerDuty{{ValidatorIndex: 9, Slot: 128}},
		Metadata: map[string]any{
			"dependent_root": phase0.Root{0xaa},
		},
	}

	service, err := mock.New(ctx)
	require.NoError(t, err)
	service.ProposerDutiesV2Func = func(_ context.Context, opts *api.ProposerDutiesOpts) (*api.Response[[]*apiv1.ProposerDuty], error) {
		require.Same(t, expectedOpts, opts)
		return expectedResponse, nil
	}

	response, err := service.ProposerDutiesV2(ctx, expectedOpts)
	require.NoError(t, err)
	require.Same(t, expectedResponse, response)
}

func TestProposerDutiesV2Default(t *testing.T) {
	ctx := context.Background()
	service, err := mock.New(ctx)
	require.NoError(t, err)

	provider, supported := any(service).(client.ProposerDutiesV2Provider)
	require.True(t, supported)

	response, err := provider.ProposerDutiesV2(ctx, &api.ProposerDutiesOpts{
		Indices: []phase0.ValidatorIndex{3, 5},
	})
	require.NoError(t, err)
	require.Equal(t, []*apiv1.ProposerDuty{
		{ValidatorIndex: 3},
		{ValidatorIndex: 5},
	}, response.Data)
	require.Empty(t, response.Metadata)
}
