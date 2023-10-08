// Copyright © 2020 - 2023 Attestant Limited.
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
	"fmt"
	"os"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestBeaconState(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name     string
		opts     *api.BeaconStateOpts
		expected *phase0.BeaconState
		err      string
		errCode  int
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name: "NilState",
			opts: &api.BeaconStateOpts{},
			err:  "no state specified",
		},
		{
			name: "Genesis",
			opts: &api.BeaconStateOpts{State: "genesis"},
		},
		{
			name: "Altair",
			opts: &api.BeaconStateOpts{State: "2375680"},
		},
		{
			name: "Bellatrix",
			opts: &api.BeaconStateOpts{State: "4636672"},
		},
		{
			name: "Capella",
			opts: &api.BeaconStateOpts{State: "6209536"},
		},
		{
			name: "Head",
			opts: &api.BeaconStateOpts{State: "head"},
		},
	}

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	jsonService, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
		http.WithEnforceJSON(true),
	)
	require.NoError(t, err)

	for _, test := range tests {
		// Run with and without enforced JSON.
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.BeaconStateProvider).BeaconState(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			case test.errCode != 0:
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Equal(t, test.errCode, apiErr.StatusCode)
				}
			default:
				// Possible that the beacon node does not contain the state, so allow a 404.
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Equal(t, 404, apiErr.StatusCode)
				} else {
					require.NoError(t, err)
					require.NotNil(t, response.Data)
				}
			}
		})
		t.Run(fmt.Sprintf("%s (json)", test.name), func(t *testing.T) {
			response, err := jsonService.(client.BeaconStateProvider).BeaconState(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			case test.errCode != 0:
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Equal(t, test.errCode, apiErr.StatusCode)
				}
			default:
				// Possible that the beacon node does not contain the state, so allow a 404.
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Equal(t, 404, apiErr.StatusCode)
				} else {
					require.NoError(t, err)
					require.NotNil(t, response.Data)
				}
			}
		})
	}
}
