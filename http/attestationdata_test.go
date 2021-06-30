// Copyright © 2020 Attestant Limited.
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

	"github.com/attestantio/go-eth2-client/http"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestAttestationData(t *testing.T) {
	tests := []struct {
		name           string
		committeeIndex spec.CommitteeIndex
		slot           int64 // -1 for current
	}{
		// {
		// 	name:           "Future",
		// 	committeeIndex: 1,
		// 	slot:           0x00ffffff,
		// },
		{
			name:           "Good",
			committeeIndex: 1,
			slot:           -1,
		},
	}

	service, err := http.New(context.Background(),
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	// Need to fetch current slot for attestation data.
	genesis, err := service.Genesis(context.Background())
	require.NoError(t, err)
	slotDuration, err := service.SlotDuration(context.Background())
	require.NoError(t, err)

	for _, test := range tests {
		var slot spec.Slot
		if test.slot == -1 {
			slot = spec.Slot(uint64(time.Since(genesis.GenesisTime).Seconds()) / uint64(slotDuration.Seconds()))
		} else {
			slot = spec.Slot(uint64(test.slot))
		}
		t.Run(test.name, func(t *testing.T) {
			attestationData, err := service.AttestationData(context.Background(), slot, test.committeeIndex)
			require.NoError(t, err)
			require.NotNil(t, attestationData)
		})
	}
}
