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
	"os"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestAggregateAttestation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name     string
		opts     *api.AggregateAttestationOpts
		expected *phase0.Attestation
		err      string
		errCode  int
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name:    "NotFound",
			opts:    &api.AggregateAttestationOpts{},
			errCode: 404,
		},
	}

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	// Need to fetch current slot for attestation.
	genesis, err := service.(client.GenesisProvider).Genesis(ctx)
	require.NoError(t, err)
	slotDuration, err := service.(client.SlotDurationProvider).SlotDuration(ctx)
	require.NoError(t, err)

	for _, test := range tests {
		slot := phase0.Slot(time.Since(genesis.GenesisTime).Seconds()) / phase0.Slot(slotDuration.Seconds())
		t.Run(test.name, func(t *testing.T) {
			// Fetch attestation data to generate a root.
			attestationDataResponse, err := service.(client.AttestationDataProvider).AttestationData(ctx, &api.AttestationDataOpts{
				Slot:           slot,
				CommitteeIndex: 0,
			})
			require.NoError(t, err)
			dataRoot, err := attestationDataResponse.Data.HashTreeRoot()
			require.NoError(t, err)

			if test.opts != nil {
				test.opts.AttestationDataRoot = dataRoot
			}

			response, err := service.(client.AggregateAttestationProvider).AggregateAttestation(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			case test.errCode != 0:
				require.Equal(t, test.errCode, err.(api.Error).StatusCode)
			default:
				require.NoError(t, err)
				require.NotNil(t, response)
				if test.expected != nil {
					require.Equal(t, response.Data, test.expected)
				}
			}
		})
	}
}
