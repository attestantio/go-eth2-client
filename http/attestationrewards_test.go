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

func TestAttestationRewards(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name              string
		opts              *api.AttestationRewardsOpts
		expectedErrorCode int
		expectedResponse  string
	}{
		{
			name:              "EpochFarFuture",
			opts:              &api.AttestationRewardsOpts{Epoch: 0xffffffffffffffff},
			expectedErrorCode: 404,
		},
		{
			name: "MixedIndicesAndPubKeys",
			opts: &api.AttestationRewardsOpts{
				Epoch: 335909,
				Indices: []phase0.ValidatorIndex{
					0, 1,
				},
				PubKeys: []phase0.BLSPubKey{
					*mustParsePubKey("0xb2ff4716ed345b05dd1dfc6a5a9fa70856d8c75dcc9e881dd2f766d5f891326f0d10e96f3a444ce6c912b69c22c6754d"),
					*mustParsePubKey("0x8e323fd501233cd4d1b9d63d74076a38de50f2f584b001a5ac2412e4e46adb26d2fb2a6041e7e8c57cd4df0916729219"),
				},
			},
			expectedResponse: `{"ideal_rewards":[{"effective_balance":"1000000000","head":"75","target":"140","source":"75","inactivity":"0"},{"effective_balance":"2000000000","head":"151","target":"281","source":"151","inactivity":"0"},{"effective_balance":"3000000000","head":"226","target":"422","source":"227","inactivity":"0"},{"effective_balance":"4000000000","head":"302","target":"562","source":"302","inactivity":"0"},{"effective_balance":"5000000000","head":"377","target":"703","source":"378","inactivity":"0"},{"effective_balance":"6000000000","head":"453","target":"844","source":"454","inactivity":"0"},{"effective_balance":"7000000000","head":"528","target":"985","source":"530","inactivity":"0"},{"effective_balance":"8000000000","head":"604","target":"1125","source":"605","inactivity":"0"},{"effective_balance":"9000000000","head":"679","target":"1266","source":"681","inactivity":"0"},{"effective_balance":"10000000000","head":"755","target":"1407","source":"757","inactivity":"0"},{"effective_balance":"11000000000","head":"830","target":"1547","source":"833","inactivity":"0"},{"effective_balance":"12000000000","head":"906","target":"1688","source":"908","inactivity":"0"},{"effective_balance":"13000000000","head":"981","target":"1829","source":"984","inactivity":"0"},{"effective_balance":"14000000000","head":"1057","target":"1970","source":"1060","inactivity":"0"},{"effective_balance":"15000000000","head":"1133","target":"2110","source":"1136","inactivity":"0"},{"effective_balance":"16000000000","head":"1208","target":"2251","source":"1211","inactivity":"0"},{"effective_balance":"17000000000","head":"1284","target":"2392","source":"1287","inactivity":"0"},{"effective_balance":"18000000000","head":"1359","target":"2532","source":"1363","inactivity":"0"},{"effective_balance":"19000000000","head":"1435","target":"2673","source":"1439","inactivity":"0"},{"effective_balance":"20000000000","head":"1510","target":"2814","source":"1514","inactivity":"0"},{"effective_balance":"21000000000","head":"1586","target":"2955","source":"1590","inactivity":"0"},{"effective_balance":"22000000000","head":"1661","target":"3095","source":"1666","inactivity":"0"},{"effective_balance":"23000000000","head":"1737","target":"3236","source":"1742","inactivity":"0"},{"effective_balance":"24000000000","head":"1812","target":"3377","source":"1817","inactivity":"0"},{"effective_balance":"25000000000","head":"1888","target":"3517","source":"1893","inactivity":"0"},{"effective_balance":"26000000000","head":"1963","target":"3658","source":"1969","inactivity":"0"},{"effective_balance":"27000000000","head":"2039","target":"3799","source":"2045","inactivity":"0"},{"effective_balance":"28000000000","head":"2115","target":"3940","source":"2120","inactivity":"0"},{"effective_balance":"29000000000","head":"2190","target":"4080","source":"2196","inactivity":"0"},{"effective_balance":"30000000000","head":"2266","target":"4221","source":"2272","inactivity":"0"},{"effective_balance":"31000000000","head":"2341","target":"4362","source":"2348","inactivity":"0"},{"effective_balance":"32000000000","head":"2417","target":"4502","source":"2423","inactivity":"0"}],"total_rewards":[{"validator_index":"0","head":"2417","target":"4502","source":"2423","inactivity":"0"},{"validator_index":"1","head":"2417","target":"4502","source":"2423","inactivity":"0"},{"validator_index":"2","head":"2417","target":"4502","source":"2423","inactivity":"0"},{"validator_index":"3","head":"2417","target":"4502","source":"2423","inactivity":"0"}]}`,
		},
		{
			name: "NegativeRewards",
			opts: &api.AttestationRewardsOpts{
				Epoch: 335909,
				Indices: []phase0.ValidatorIndex{
					63,
				},
			},
			expectedResponse: `{"ideal_rewards":[{"effective_balance":"1000000000","head":"75","target":"140","source":"75","inactivity":"0"},{"effective_balance":"2000000000","head":"151","target":"281","source":"151","inactivity":"0"},{"effective_balance":"3000000000","head":"226","target":"422","source":"227","inactivity":"0"},{"effective_balance":"4000000000","head":"302","target":"562","source":"302","inactivity":"0"},{"effective_balance":"5000000000","head":"377","target":"703","source":"378","inactivity":"0"},{"effective_balance":"6000000000","head":"453","target":"844","source":"454","inactivity":"0"},{"effective_balance":"7000000000","head":"528","target":"985","source":"530","inactivity":"0"},{"effective_balance":"8000000000","head":"604","target":"1125","source":"605","inactivity":"0"},{"effective_balance":"9000000000","head":"679","target":"1266","source":"681","inactivity":"0"},{"effective_balance":"10000000000","head":"755","target":"1407","source":"757","inactivity":"0"},{"effective_balance":"11000000000","head":"830","target":"1547","source":"833","inactivity":"0"},{"effective_balance":"12000000000","head":"906","target":"1688","source":"908","inactivity":"0"},{"effective_balance":"13000000000","head":"981","target":"1829","source":"984","inactivity":"0"},{"effective_balance":"14000000000","head":"1057","target":"1970","source":"1060","inactivity":"0"},{"effective_balance":"15000000000","head":"1133","target":"2110","source":"1136","inactivity":"0"},{"effective_balance":"16000000000","head":"1208","target":"2251","source":"1211","inactivity":"0"},{"effective_balance":"17000000000","head":"1284","target":"2392","source":"1287","inactivity":"0"},{"effective_balance":"18000000000","head":"1359","target":"2532","source":"1363","inactivity":"0"},{"effective_balance":"19000000000","head":"1435","target":"2673","source":"1439","inactivity":"0"},{"effective_balance":"20000000000","head":"1510","target":"2814","source":"1514","inactivity":"0"},{"effective_balance":"21000000000","head":"1586","target":"2955","source":"1590","inactivity":"0"},{"effective_balance":"22000000000","head":"1661","target":"3095","source":"1666","inactivity":"0"},{"effective_balance":"23000000000","head":"1737","target":"3236","source":"1742","inactivity":"0"},{"effective_balance":"24000000000","head":"1812","target":"3377","source":"1817","inactivity":"0"},{"effective_balance":"25000000000","head":"1888","target":"3517","source":"1893","inactivity":"0"},{"effective_balance":"26000000000","head":"1963","target":"3658","source":"1969","inactivity":"0"},{"effective_balance":"27000000000","head":"2039","target":"3799","source":"2045","inactivity":"0"},{"effective_balance":"28000000000","head":"2115","target":"3940","source":"2120","inactivity":"0"},{"effective_balance":"29000000000","head":"2190","target":"4080","source":"2196","inactivity":"0"},{"effective_balance":"30000000000","head":"2266","target":"4221","source":"2272","inactivity":"0"},{"effective_balance":"31000000000","head":"2341","target":"4362","source":"2348","inactivity":"0"},{"effective_balance":"32000000000","head":"2417","target":"4502","source":"2423","inactivity":"0"}],"total_rewards":[{"validator_index":"63","head":"0","target":"-4511","source":"-2429","inactivity":"0"}]}`,
		},
	}

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.AttestationRewardsProvider).AttestationRewards(ctx, test.opts)
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
