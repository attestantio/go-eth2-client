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
	"errors"
	"os"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"
)

func TestBeaconBlockAttestations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name     string
		opts     *api.BeaconBlockAttestationsOpts
		expected *api.Response[[]*spec.VersionedAttestation]
		err      string
		errCode  int
	}{
		{
			name: "NoOpts",
			err:  "no options specified",
		},
		{
			name: "NoBlock",
			opts: &api.BeaconBlockAttestationsOpts{},
			err:  "no block specified",
		},
		{
			name: "InvalidBlock",
			opts: &api.BeaconBlockAttestationsOpts{
				Block: "invalid",
			},
			errCode: 400,
		},
		{
			name: "Genesis",
			opts: &api.BeaconBlockAttestationsOpts{
				Block: "0",
			},
			expected: nil, //TODO
		},
	}

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.BeaconBlockAttestationsProvider).BeaconBlockAttestations(ctx, test.opts)
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

func TestBeaconBlockAttestationsTimeout(t *testing.T) {
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
