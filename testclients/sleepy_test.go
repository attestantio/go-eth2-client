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
	"time"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestSleepyNew(t *testing.T) {
	ctx := context.Background()

	client, err := mock.New(ctx,
		mock.WithLogLevel(zerolog.Disabled),
	)
	require.NoError(t, err)

	tests := []struct {
		name     string
		minSleep time.Duration
		maxSleep time.Duration
		next     consensusclient.Service
		err      string
	}{
		{
			name:     "ClientMissing",
			minSleep: time.Second,
			maxSleep: 2 * time.Second,
			err:      "no next service supplied",
		},
		{
			name:     "MaxSleepLowerThanMinSleep",
			minSleep: 2 * time.Second,
			maxSleep: time.Second,
			next:     client,
			err:      "max sleep less than min sleep",
		},
		{
			name:     "Good",
			minSleep: time.Second,
			maxSleep: 2 * time.Second,
			next:     client,
		},
		{
			name:     "EqualSleep",
			minSleep: time.Second,
			maxSleep: time.Second,
			next:     client,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := testclients.NewSleepy(ctx, test.minSleep, test.maxSleep, test.next)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSleep(t *testing.T) {
	ctx := context.Background()

	client, err := mock.New(ctx,
		mock.WithLogLevel(zerolog.Disabled),
	)
	require.NoError(t, err)

	minSleep := 100 * time.Millisecond
	maxSleep := 200 * time.Millisecond
	s, err := testclients.NewSleepy(ctx, minSleep, maxSleep, client)
	require.NoError(t, err)

	for i := 0; i < 16; i++ {
		started := time.Now()
		_, err := s.(consensusclient.GenesisProvider).Genesis(ctx, &api.GenesisOpts{})
		require.NoError(t, err)
		duration := time.Since(started)
		require.GreaterOrEqual(t, duration.Milliseconds(), minSleep.Milliseconds())
		// We can be a bit longer than max sleep due to the call itself taking time.
		require.LessOrEqual(t, duration.Milliseconds(), (maxSleep + 50*time.Millisecond).Milliseconds())
	}
}
