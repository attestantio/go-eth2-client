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
	"fmt"
	"strconv"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/pkg/errors"
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
		network           string
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
			network:          "mainnet",
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
			network:          "mainnet",
		},
		{
			name: "MixedIndicesAndPubKeysHoodi",
			opts: &api.SyncCommitteeRewardsOpts{
				Block: "1714544",
				Indices: []phase0.ValidatorIndex{
					290742,
				},
				PubKeys: []phase0.BLSPubKey{
					*mustParsePubKey("0x89c3d75c9fa8daa39cf721fd3caf441de9b43dd59ae1275dd482fa48dbd8463737038b3b8cb2f53c1a8635c8733fb7ca"),
				},
			},
			expectedResponse: `[{"validator_index":"290742","reward":"25016"}, {"validator_index":"525735","reward":"-25016"}]`,
			network:          "hoodi",
		},
		{
			name: "NegativeRewardsHoodi",
			opts: &api.SyncCommitteeRewardsOpts{
				Block: "1714544",
				Indices: []phase0.ValidatorIndex{
					525735,
				},
			},
			expectedResponse: `[{"validator_index":"525735","reward":"-25016"}]`,
			network:          "hoodi",
		},
	}

	service := testService(ctx, t).(client.Service)
	network := testclients.NetworkName(ctx, service)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.network != "" && test.network != network {
				t.Skipf("Skipping test %s on network %s. Client network: %s", test.name, test.network, network)
			}
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
					err = jsonEqualCommitteeRewards(test.expectedResponse, string(responseJSON))
					require.NoError(t, err)
				}
			}
		})
	}
}

type committeeReward struct {
	ValidatorIndex string `json:"validator_index"`
	Reward         string `json:"reward"`
}

// require.JSONEq fails for these tests because the order of the elements is not guaranteed.
func jsonEqualCommitteeRewards(expectedJson, actualJson string) error {
	var expectedData []committeeReward
	var actualData []committeeReward

	if err := json.Unmarshal([]byte(expectedJson), &expectedData); err != nil {
		return errors.Wrap(err, "could not unmarshal json1")
	}

	if err := json.Unmarshal([]byte(actualJson), &actualData); err != nil {
		return errors.Wrap(err, "could not unmarshal json2")
	}

	if len(expectedData) != len(actualData) {
		return errors.New("number of rewards is different")
	}

	for i := range expectedData {
		found := false
		for j := range actualData {
			if expectedData[i].ValidatorIndex == actualData[j].ValidatorIndex {
				if expectedData[i].Reward != actualData[j].Reward {
					return errors.New(fmt.Sprintf("response does not contain the expected reward for validator index %s: expected %s, got %s", expectedData[i].ValidatorIndex, expectedData[i].Reward, actualData[j].Reward))
				}
				found = true
				break
			}
		}
		if !found {
			return errors.New(fmt.Sprintf("response does not contain the expected validator index: expected %s, got %s", expectedData[i].ValidatorIndex, actualData[i].ValidatorIndex))
		}
	}
	return nil
}
