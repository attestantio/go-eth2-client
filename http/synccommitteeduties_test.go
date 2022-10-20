// Copyright © 2021 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestSyncCommitteeDuties(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	// Needed to fetch current epoch.
	genesis, err := service.(client.GenesisProvider).Genesis(ctx)
	require.NoError(t, err)
	slotDuration, err := service.(client.SlotDurationProvider).SlotDuration(ctx)
	require.NoError(t, err)
	slotsPerEpoch, err := service.(client.SlotsPerEpochProvider).SlotsPerEpoch(ctx)
	require.NoError(t, err)

	tests := []struct {
		name             string
		epoch            int64 // -1 for current
		validatorIndices []phase0.ValidatorIndex
	}{
		{
			name:  "Current",
			epoch: -1,
			validatorIndices: []phase0.ValidatorIndex{
				0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
				16, 17, 18, 29, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
				32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47,
				48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63,
				64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
				80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95,
				96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111,
				112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var epoch phase0.Epoch
			if test.epoch == -1 {
				epoch = phase0.Epoch(uint64(time.Since(genesis.GenesisTime).Seconds()) / (uint64(slotDuration.Seconds()) * slotsPerEpoch))
			} else {
				epoch = phase0.Epoch(test.epoch)
			}
			duties, err := service.(client.SyncCommitteeDutiesProvider).SyncCommitteeDuties(ctx, epoch, test.validatorIndices)
			require.NoError(t, err)
			require.NotNil(t, duties)
			// No guaratee that any included indices will be have a sync duty.
		})
	}
}
