// Copyright © 2020 - 2024 Attestant Limited.
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
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestValidatorBalances(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name             string
		opts             *api.ValidatorBalancesOpts
		err              string
		errCodes         []int
		expected         map[phase0.ValidatorIndex]phase0.Gwei
		expectedBalances int
	}{
		{
			name: "NoOpts",
			err:  "no options specified",
		},
		{
			name: "StateInvalid",
			opts: &api.ValidatorBalancesOpts{
				State: "invalid",
			},
			errCodes: []int{400},
		},
		{
			name: "StateUnknown",
			opts: &api.ValidatorBalancesOpts{
				State: "0x0000000000000000000000000000000000000000000000000000000000000000",
			},
			errCodes: []int{404},
		},
		{
			name: "StateFuture",
			opts: &api.ValidatorBalancesOpts{
				State: "9999999999",
			},
			errCodes: []int{404, 500},
		},
		{
			name: "SingleGenesisIndex",
			opts: &api.ValidatorBalancesOpts{
				State:   "0",
				Indices: []phase0.ValidatorIndex{123},
			},
			expected: map[phase0.ValidatorIndex]phase0.Gwei{
				123: 32000000000,
			},
		},
		{
			name: "AllGenesis",
			opts: &api.ValidatorBalancesOpts{
				State: "0",
			},
		},
	}

	service := testService(ctx, t).(client.Service)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.ValidatorBalancesProvider).ValidatorBalances(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			case len(test.errCodes) > 0:
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Contains(t, test.errCodes, apiErr.StatusCode)
				}
			default:
				require.NoError(t, err)
				require.NotNil(t, response)
				if test.expected != nil {
					require.Equal(t, test.expected, response.Data)
				}
				if test.expectedBalances > 0 {
					require.Len(t, response.Data, test.expectedBalances)
				}
			}
		})
	}
}
