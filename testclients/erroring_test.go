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

package testclients_test

import (
	"context"
	"testing"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestErroringNew(t *testing.T) {
	ctx := context.Background()

	client, err := mock.New(ctx,
		mock.WithLogLevel(zerolog.Disabled),
	)
	require.NoError(t, err)

	tests := []struct {
		name      string
		errorRate float64
		next      consensusclient.Service
		err       string
	}{
		{
			name:      "ClientMissing",
			errorRate: 0.1,
			err:       "no next service supplied",
		},
		{
			name:      "ErrorRateNegative",
			errorRate: -1,
			next:      client,
			err:       "error rate cannot be less than 0",
		},
		{
			name:      "ErrorRateTooHigh",
			errorRate: 1.1,
			next:      client,
			err:       "error rate cannot be more than 1",
		},
		{
			name:      "Good",
			errorRate: 0.5,
			next:      client,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := testclients.NewErroring(ctx, test.errorRate, test.next)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMaybeError(t *testing.T) {
	ctx := context.Background()

	client, err := mock.New(ctx,
		mock.WithLogLevel(zerolog.Disabled),
	)
	require.NoError(t, err)

	errorRate := 0.9
	s, err := testclients.NewErroring(ctx, errorRate, client)
	require.NoError(t, err)

	errors := 0
	for i := 0; i < 100000; i++ {
		_, err := s.(consensusclient.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
		if err != nil {
			errors++
		}
	}
	// Expect approximately 90% of the requests to have errored.
	require.LessOrEqual(t, errors, 90500)
	require.GreaterOrEqual(t, errors, 89500)
}
