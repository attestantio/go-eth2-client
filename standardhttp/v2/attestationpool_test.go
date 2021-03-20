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
	"github.com/stretchr/testify/require"
)

func TestAttestationPool(t *testing.T) {
	tests := []struct {
		name string
		slot int64 // -1 for current
	}{
		{
			name: "Good",
			slot: -1,
		},
	}

	service, err := standardhttp.New(context.Background(),
		standardhttp.WithTimeout(timeout),
		standardhttp.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	// Need to fetch current slot for attestation pools.
	genesis, err := service.Genesis(context.Background())
	require.NoError(t, err)
	chainSpec, err := service.Spec(context.Background())
	require.NoError(t, err)
	slotDuration := chainSpec["SECONDS_PER_SLOT"].(time.Duration)

	for _, test := range tests {
		var slot spec.Slot
		if test.slot == -1 {
			slot = spec.Slot(uint64(time.Since(genesis.GenesisTime).Seconds()) / uint64(slotDuration.Seconds()))
		} else {
			slot = spec.Slot(uint64(test.slot))
		}
		t.Run(test.name, func(t *testing.T) {
			attestationPool, err := service.AttestationPool(context.Background(), slot)
			require.NoError(t, err)
			require.NotNil(t, attestationPool)
		})
	}
}
