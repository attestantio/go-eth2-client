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
	"os"
	"strconv"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestSyncCommitteeRewards(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name              string
		opts              *api.SyncCommitteeRewardsOpts
		expectedErrorCode int
		expectedResponse  string
	}{
		{
			name:              "BlockInvalid",
			opts:              &api.SyncCommitteeRewardsOpts{Block: "current"},
			expectedErrorCode: 400,
		},
		{
			name: "MixedIndicesAndPubKeys",
			opts: &api.SyncCommitteeRewardsOpts{
				Block: "10760058",
				Indices: []phase0.ValidatorIndex{
					286437,
				},
				PubKeys: []phase0.BLSPubKey{
					*mustParsePubKey("0xb7dd1c63cfe60163ffcb889d502b0af3b8ab41cb0dc95edb46eccfeb79e984886fe54f800e813ae09d48e98087010a10"),
				},
			},
			expectedResponse: `[{"validator_index":"286437","reward":"22456"},{"validator_index":"1674334","reward":"22456"}]`,
		},
		{
			name: "NegativeRewards",
			opts: &api.SyncCommitteeRewardsOpts{
				Block: "10760058",
				Indices: []phase0.ValidatorIndex{
					1055307,
				},
			},
			expectedResponse: `[{"validator_index":"1055307","reward":"-22456"}]`,
		},
	}

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.SyncCommitteeRewardsProvider).SyncCommitteeRewards(ctx, test.opts)
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
