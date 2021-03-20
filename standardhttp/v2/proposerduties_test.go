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

package v2_test

import (
	"context"
	"os"
	"testing"
	"time"

	spec "github.com/attestantio/go-eth2-client/spec/altair"
	standardhttp "github.com/attestantio/go-eth2-client/standardhttp/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestProposerDuties(t *testing.T) {
	tests := []struct {
		name             string
		epoch            int64 // -1 for current
		validatorIndices []spec.ValidatorIndex
		expected         int
	}{
		{
			name:     "Epoch",
			epoch:    0,
			expected: 31,
		},
		{
			name:     "Old",
			epoch:    1,
			expected: 32,
		},
		{
			name:     "Current",
			epoch:    -1,
			expected: 32,
		},
	}

	service, err := standardhttp.New(context.Background(),
		standardhttp.WithLogLevel(zerolog.TraceLevel),
		standardhttp.WithTimeout(timeout),
		standardhttp.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	// Needed to fetch current epoch.
	genesis, err := service.Genesis(context.Background())
	require.NoError(t, err)
	chainSpec, err := service.Spec(context.Background())
	require.NoError(t, err)
	slotDuration := chainSpec["SECONDS_PER_SLOT"].(time.Duration)
	slotsPerEpoch := chainSpec["SLOTS_PER_EPOCH"].(uint64)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var epoch spec.Epoch
			if test.epoch == -1 {
				epoch = spec.Epoch(uint64(time.Since(genesis.GenesisTime).Seconds()) / (uint64(slotDuration.Seconds()) * slotsPerEpoch))
			} else {
				epoch = spec.Epoch(test.epoch)
			}
			duties, err := service.ProposerDuties(context.Background(), epoch, test.validatorIndices)
			require.NoError(t, err)
			require.NotNil(t, duties)
			require.Equal(t, test.expected, len(duties))
		})
	}
}
