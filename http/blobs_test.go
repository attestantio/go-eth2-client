// Copyright © 2025 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"math"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"
)

func TestBlobs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	headBlock, err := service.(client.SignedBeaconBlockProvider).SignedBeaconBlock(ctx, &api.SignedBeaconBlockOpts{
		Block: "head",
	})
	require.NoError(t, err)
	require.NotNil(t, headBlock)

	if headBlock.Data.Version < spec.DataVersionFulu {
		t.Logf("Client does not support Fulu, skipping tests")
		return
	}

	headBlockSlot, err := headBlock.Data.Slot()
	require.NoError(t, err)
	headBlockRoot, err := headBlock.Data.Root()
	require.NoError(t, err)

	assertBlobs := func(t *testing.T, response *api.Response[v1.Blobs], err error) {
		require.NoError(t, err)
		require.NotNil(t, response)
		require.NotNil(t, response.Data)
		require.Greater(t, len(response.Data), 0)

		// Check a Blob for correctness
		// GetTree panics with the error: size of tree should be a power of 2
		_, err = response.Data.GetTree()
		require.NoError(t, err)
		_, err = response.Data.HashTreeRoot()
		require.NoError(t, err)
		size := response.Data.SizeSSZ()
		require.Equal(t, 131072*len(response.Data), size)
	}

	tests := []struct {
		name   string
		opts   *api.BlobsOpts
		assert func(t *testing.T, response *api.Response[v1.Blobs], err error)
	}{
		{
			name: "Head Block Root",
			opts: &api.BlobsOpts{
				Block: headBlockRoot.String(),
			},
			assert: assertBlobs,
		},
		{
			name: "Head Block Slot",
			opts: &api.BlobsOpts{
				Block: fmt.Sprintf("%d", headBlockSlot),
			},
			assert: assertBlobs,
		},
		{
			name: "Invalid Block",
			opts: &api.BlobsOpts{
				Block: "invalid",
			},
			assert: func(t *testing.T, response *api.Response[v1.Blobs], err error) {
				var apiError *api.Error
				if errors.As(err, &apiError) {
					require.Equal(t, 400, apiError.StatusCode)
				}
			},
		},
		{
			name: "Slot to far in the future",
			opts: &api.BlobsOpts{
				Block: fmt.Sprintf("%d", math.MaxInt64),
			},
			assert: func(t *testing.T, response *api.Response[v1.Blobs], err error) {
				var apiError *api.Error
				if errors.As(err, &apiError) {
					require.Equal(t, 404, apiError.StatusCode)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.BlobsProvider).Blobs(ctx, test.opts)
			test.assert(t, response, err)
		})
	}
}
