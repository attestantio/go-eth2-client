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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestAggregateAttestation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	// Need to fetch current slot for attestation.
	genesisResponse, err := service.(client.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
	require.NoError(t, err)
	slotDuration, err := service.(client.SlotDurationProvider).SlotDuration(ctx)
	require.NoError(t, err)

	// Fetch attestation data to generate a root.
	attestationDataResponse, err := service.(client.AttestationDataProvider).AttestationData(ctx, &api.AttestationDataOpts{
		Slot:           phase0.Slot(time.Since(genesisResponse.Data.GenesisTime).Seconds()) / phase0.Slot(slotDuration.Seconds()),
		CommitteeIndex: 0,
	})
	require.NoError(t, err)
	dataRoot, err := attestationDataResponse.Data.HashTreeRoot()
	require.NoError(t, err)

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
			name: "NilRoot",
			opts: &api.AggregateAttestationOpts{
				Slot: phase0.Slot(time.Since(genesisResponse.Data.GenesisTime).Seconds()) / phase0.Slot(slotDuration.Seconds()),
			},
			err: "no attestation data root specified",
		},
		{
			// Will generally get "not found" because the beacon node is unlikely to be an aggregator.
			name: "NotFound",
			opts: &api.AggregateAttestationOpts{
				Slot:                phase0.Slot(time.Since(genesisResponse.Data.GenesisTime).Seconds()) / phase0.Slot(slotDuration.Seconds()),
				AttestationDataRoot: dataRoot,
			},
			errCode: 404,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.AggregateAttestationProvider).AggregateAttestation(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			case test.errCode != 0:
				var apiErr *api.Error
				if errors.As(err, &apiErr) {
					require.Equal(t, test.errCode, apiErr.StatusCode)
				}
			default:
				require.NoError(t, err)
				require.NotNil(t, response)
				if test.expected != nil {
					require.Equal(t, test.expected, response.Data)
				}
			}
		})
	}
}
