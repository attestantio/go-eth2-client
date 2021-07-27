// Copyright © 2020, 2021 Attestant Limited.
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

package prysmgrpc_test

import (
	"context"
	"os"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/prysmgrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	s, err := prysmgrpc.New(context.Background(), prysmgrpc.WithAddress(os.Getenv("PRYSMGRPC_ADDRESS")))
	require.NoError(t, err)
	require.NotNil(t, s)
	require.Equal(t, os.Getenv("PRYSMGRPC_ADDRESS"), s.Address())
}

func TestInterfaces(t *testing.T) {
	var s interface{}
	s, err := prysmgrpc.New(context.Background(), prysmgrpc.WithAddress(os.Getenv("PRYSMGRPC_ADDRESS")))
	require.NoError(t, err)
	require.NotNil(t, s)

	// Standard API.
	assert.Implements(t, (*client.AggregateAttestationsSubmitter)(nil), s)
	assert.Implements(t, (*client.AttestationDataProvider)(nil), s)
	assert.Implements(t, (*client.AttestationsSubmitter)(nil), s)
	assert.Implements(t, (*client.BeaconBlockProposalProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockSubmitter)(nil), s)
	assert.Implements(t, (*client.BeaconCommitteeSubscriptionsSubmitter)(nil), s)
	assert.Implements(t, (*client.ForkProvider)(nil), s)
	assert.Implements(t, (*client.GenesisProvider)(nil), s)
	assert.Implements(t, (*client.SpecProvider)(nil), s)
	assert.Implements(t, (*client.SyncStateProvider)(nil), s)
	assert.Implements(t, (*client.ValidatorsProvider)(nil), s)

	// Non-standard APIs.
	assert.Implements(t, (*client.AttesterDutiesProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockRootProvider)(nil), s)
	assert.Implements(t, (*client.BeaconChainHeadUpdatedSource)(nil), s)
	assert.Implements(t, (*client.GenesisTimeProvider)(nil), s)
	assert.Implements(t, (*client.GenesisValidatorsRootProvider)(nil), s)
	assert.Implements(t, (*client.NodeVersionProvider)(nil), s)
	assert.Implements(t, (*client.ProposerDutiesProvider)(nil), s)
	assert.Implements(t, (*client.SlotDurationProvider)(nil), s)
	assert.Implements(t, (*client.SlotsPerEpochProvider)(nil), s)
	assert.Implements(t, (*client.TargetAggregatorsPerCommitteeProvider)(nil), s)
	assert.Implements(t, (*client.ValidatorsWithoutBalanceProvider)(nil), s)

	// Prysm-specific APIs.
	assert.Implements(t, (*client.PrysmAggregateAttestationProvider)(nil), s)
	assert.Implements(t, (*client.PrysmValidatorBalancesProvider)(nil), s)
}

func TestTLS(t *testing.T) {
	if os.Getenv("PRYSMGRPC_TLS_ADDRESS") == "" {
		t.Skip("PRYSMGRPC_TLS_ADDRESS not specified; not testing secure connection")
	}
	s, err := prysmgrpc.New(context.Background(),
		prysmgrpc.WithAddress(os.Getenv("PRYSMGRPC_TLS_ADDRESS")),
		prysmgrpc.WithTLS(true),
	)
	require.NoError(t, err)
	require.NotNil(t, s)
	require.Equal(t, os.Getenv("PRYSMGRPC_TLS_ADDRESS"), s.Address())
}
