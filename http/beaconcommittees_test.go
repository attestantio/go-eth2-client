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

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestBeaconCommittees(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	epoch := phase0.Epoch(0)

	tests := []struct {
		name string
		opts *api.BeaconCommitteesOpts
	}{
		{
			name: "Good",
			opts: &api.BeaconCommitteesOpts{
				State: "head",
			},
		},
		{
			name: "AtGenesisEpoch",
			opts: &api.BeaconCommitteesOpts{
				State: "genesis",
				Epoch: &epoch,
			},
		},
	}

	service, err := http.New(ctx,
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			beaconCommittees, err := service.(client.BeaconCommitteesProvider).BeaconCommittees(ctx, test.opts)
			require.NoError(t, err)
			require.NotNil(t, beaconCommittees)
		})
	}
}
