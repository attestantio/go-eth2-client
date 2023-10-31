// Copyright © 2023 Attestant Limited.
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
	"os"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestBeaconBlockRoot(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name     string
		opts     *api.BeaconBlockRootOpts
		expected *phase0.Root
		err      string
		errCode  int
	}{
		{
			name: "NoOpts",
			err:  "no options specified",
		},
		{
			name: "NoBlock",
			opts: &api.BeaconBlockRootOpts{},
			err:  "no block specified",
		},
		{
			name: "InvalidBlock",
			opts: &api.BeaconBlockRootOpts{
				Block: "invalid",
			},
			errCode: 400,
		},
		{
			name: "Genesis",
			opts: &api.BeaconBlockRootOpts{
				Block: "0",
			},
			expected: mustParseRoot("0x4d611d5b93fdab69013a7f0a2f961caca0c853f87cfe9595fe50038163079360"),
		},
	}

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.BeaconBlockRootProvider).BeaconBlockRoot(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			case test.errCode != 0:
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Equal(t, test.errCode, apiErr.StatusCode)
				}
			default:
				require.NoError(t, err)
				require.NotNil(t, response)
				require.Equal(t, test.expected, response.Data)
			}
		})
	}
}

func TestBeaconBlockRootTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	_, err = service.(client.BeaconBlockRootProvider).BeaconBlockRoot(ctx, &api.BeaconBlockRootOpts{
		Common: api.CommonOpts{
			Timeout: time.Millisecond,
		},
		Block: "0",
	})
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}
