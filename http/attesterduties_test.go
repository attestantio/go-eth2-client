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
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestAttesterDuties(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	// Need to fetch current epoch for duties.
	genesisResponse, err := service.(client.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
	require.NoError(t, err)
	slotDuration, err := service.(client.SlotDurationProvider).SlotDuration(ctx)
	require.NoError(t, err)
	slotsPerEpoch, err := service.(client.SlotsPerEpochProvider).SlotsPerEpoch(ctx)
	require.NoError(t, err)

	tests := []struct {
		name     string
		opts     *api.AttesterDutiesOpts
		expected []*apiv1.AttesterDuty
		err      string
		errCode  int
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name: "NoValidatorIndices",
			opts: &api.AttesterDutiesOpts{
				Epoch:   phase0.Epoch(time.Since(genesisResponse.Data.GenesisTime).Seconds()) / phase0.Epoch(slotDuration.Seconds()) / phase0.Epoch(slotsPerEpoch),
				Indices: []phase0.ValidatorIndex{},
			},
			err: "no validator indices specified",
		},
		{
			name: "Good",
			opts: &api.AttesterDutiesOpts{
				Epoch:   phase0.Epoch(time.Since(genesisResponse.Data.GenesisTime).Seconds()) / phase0.Epoch(slotDuration.Seconds()) / phase0.Epoch(slotsPerEpoch),
				Indices: []phase0.ValidatorIndex{0, 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.AttesterDutiesProvider).AttesterDuties(ctx, test.opts)
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
