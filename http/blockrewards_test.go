// Copyright © 2025 Attestant Limited.
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
	"encoding/json"
	"strconv"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/stretchr/testify/require"
)

func TestBlockRewards(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name              string
		opts              *api.BlockRewardsOpts
		expectedErrorCode int
		expectedResponse  string
		network           string
	}{
		{
			name:              "BlockInvalid",
			opts:              &api.BlockRewardsOpts{Block: "current"},
			expectedErrorCode: 400,
		},
		{
			name: "Good",
			opts: &api.BlockRewardsOpts{
				Block: "10760040",
			},
			expectedResponse: `{"proposer_index":"1515282","total":"42294581","attestations":"40696997","sync_aggregate":"1597584","proposer_slashings":"0","attester_slashings":"0"}`,
			network:          "mainnet",
		},
		{
			name: "GoodHoodi",
			opts: &api.BlockRewardsOpts{
				Block: "1714850",
			},
			expectedResponse: `{"proposer_index":"196644","total":"48373688","attestations":"46647929","sync_aggregate":"1725759","proposer_slashings":"0","attester_slashings":"0"}`,
			network:          "hoodi",
		},
	}

	service := testService(ctx, t).(client.Service)
	network := testclients.NetworkName(ctx, service)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.network != "" && test.network != network {
				t.Skipf("Skipping test %s on network %s", test.name, network)
			}
			response, err := service.(client.BlockRewardsProvider).BlockRewards(ctx, test.opts)
			if test.expectedErrorCode != 0 {
				require.Contains(t, err.Error(), strconv.Itoa(test.expectedErrorCode))
			} else {
				require.NoError(t, err)
				require.NotNil(t, response)
				require.NotNil(t, response.Data)
				require.NotNil(t, response.Metadata)
				if test.expectedResponse != "" {
					responseJSON, err := json.Marshal(response.Data)
					require.NoError(t, err)
					require.Equal(t, test.expectedResponse, string(responseJSON))
				}
			}
		})
	}
}
