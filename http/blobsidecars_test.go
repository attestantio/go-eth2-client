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
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/stretchr/testify/require"
)

func TestBlobsSidecars(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	headBlock, err := service.(client.SignedBeaconBlockProvider).SignedBeaconBlock(ctx, &api.SignedBeaconBlockOpts{
		Block: "head",
	})
	require.NoError(t, err)
	require.NotNil(t, headBlock)

	headBlockSlot, err := headBlock.Data.Slot()
	require.NoError(t, err)
	headBlockRoot, err := headBlock.Data.Root()
	require.NoError(t, err)

	assertBlobs := func(t *testing.T, response *api.Response[[]*deneb.BlobSidecar], err error) {
		require.NoError(t, err)
		require.NotNil(t, response)
		require.NotNil(t, response.Data)

		if len(response.Data) == 0 {
			t.Skip("No blobs found or endpoint is already deprecated")
		}

		// Check a Blob for correctness
		_, err = response.Data[0].GetTree()
		require.NoError(t, err)
		_, err = response.Data[0].HashTreeRoot()
		require.NoError(t, err)
		size := response.Data[0].SizeSSZ()
		require.Greater(t, size, 0)
	}

	tests := []struct {
		name   string
		opts   *api.BlobSidecarsOpts
		assert func(t *testing.T, response *api.Response[[]*deneb.BlobSidecar], err error)
	}{
		{
			name: "Head Block Root",
			opts: &api.BlobSidecarsOpts{
				Block: headBlockRoot.String(),
			},
			assert: assertBlobs,
		},
		{
			name: "Head Block Slot",
			opts: &api.BlobSidecarsOpts{
				Block: fmt.Sprintf("%d", headBlockSlot),
			},
			assert: assertBlobs,
		},
		{
			name: "Invalid Block",
			opts: &api.BlobSidecarsOpts{
				Block: "invalid",
			},
			assert: func(t *testing.T, response *api.Response[[]*deneb.BlobSidecar], err error) {
				var apiError *api.Error
				if errors.As(err, &apiError) {
					require.Equal(t, 400, apiError.StatusCode)
				}
			},
		},
		{
			name: "Slot to far in the future",
			opts: &api.BlobSidecarsOpts{
				Block: fmt.Sprintf("%d", math.MaxInt64),
			},
			assert: func(t *testing.T, response *api.Response[[]*deneb.BlobSidecar], err error) {
				var apiError *api.Error
				if errors.As(err, &apiError) {
					require.Equal(t, 404, apiError.StatusCode)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.BlobSidecarsProvider).BlobSidecars(ctx, test.opts)
			test.assert(t, response, err)
		})
	}
}
