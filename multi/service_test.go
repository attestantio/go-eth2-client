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

package multi_test

import (
	"context"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/multi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	ctx := context.Background()

	consensusclient1, err := mock.New(ctx)
	require.NoError(t, err)
	consensusclient2, err := mock.New(ctx)
	require.NoError(t, err)
	consensusclient3, err := mock.New(ctx)
	require.NoError(t, err)
	inactiveconsensusclient1, err := mock.New(ctx)
	require.NoError(t, err)
	inactiveconsensusclient1.SyncDistance = 10
	tests := []struct {
		name   string
		params []multi.Parameter
		err    string
	}{
		{
			name: "ClientsMissing",
			params: []multi.Parameter{
				multi.WithLogLevel(zerolog.Disabled),
			},
			err: "problem with parameters: no Ethereum 2 clients specified",
		},
		{
			name: "AllClientsInactive",
			params: []multi.Parameter{
				multi.WithLogLevel(zerolog.Disabled),
				multi.WithClients([]client.Service{
					inactiveconsensusclient1,
				}),
			},
			err: "No providers active, cannot proceed",
		},
		{
			name: "Good",
			params: []multi.Parameter{
				multi.WithLogLevel(zerolog.Disabled),
				multi.WithClients([]client.Service{
					consensusclient1,
					consensusclient2,
					consensusclient3,
				}),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := multi.New(context.Background(), test.params...)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestInterfaces(t *testing.T) {
	ctx := context.Background()

	consensusclient1, err := mock.New(ctx)
	require.NoError(t, err)
	consensusclient2, err := mock.New(ctx)
	require.NoError(t, err)
	s, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]client.Service{
			consensusclient1,
			consensusclient2,
		}),
	)
	require.NoError(t, err)

	// Standard interfacs.
	assert.Implements(t, (*client.AggregateAttestationProvider)(nil), s)
	assert.Implements(t, (*client.AggregateAttestationsSubmitter)(nil), s)
	assert.Implements(t, (*client.AttestationDataProvider)(nil), s)
	assert.Implements(t, (*client.AttestationPoolProvider)(nil), s)
	assert.Implements(t, (*client.AttestationsSubmitter)(nil), s)
	assert.Implements(t, (*client.AttesterDutiesProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockHeadersProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockRootProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockSubmitter)(nil), s)
	assert.Implements(t, (*client.BeaconCommitteeSubscriptionsSubmitter)(nil), s)
	assert.Implements(t, (*client.BeaconStateProvider)(nil), s)
	assert.Implements(t, (*client.BlindedBeaconBlockSubmitter)(nil), s)
	assert.Implements(t, (*client.ValidatorRegistrationsSubmitter)(nil), s)
	assert.Implements(t, (*client.DepositContractProvider)(nil), s)
	assert.Implements(t, (*client.EventsProvider)(nil), s)
	assert.Implements(t, (*client.FinalityProvider)(nil), s)
	assert.Implements(t, (*client.ForkProvider)(nil), s)
	assert.Implements(t, (*client.ForkScheduleProvider)(nil), s)
	assert.Implements(t, (*client.GenesisProvider)(nil), s)
	assert.Implements(t, (*client.NodeSyncingProvider)(nil), s)
	assert.Implements(t, (*client.ProposalPreparationsSubmitter)(nil), s)
	assert.Implements(t, (*client.ProposalProvider)(nil), s)
	assert.Implements(t, (*client.ProposerDutiesProvider)(nil), s)
	assert.Implements(t, (*client.SpecProvider)(nil), s)
	assert.Implements(t, (*client.SyncCommitteeContributionProvider)(nil), s)
	assert.Implements(t, (*client.SyncCommitteeContributionsSubmitter)(nil), s)
	assert.Implements(t, (*client.SyncCommitteeDutiesProvider)(nil), s)
	assert.Implements(t, (*client.SyncCommitteeMessagesSubmitter)(nil), s)
	assert.Implements(t, (*client.SyncCommitteesProvider)(nil), s)
	assert.Implements(t, (*client.SyncCommitteeSubscriptionsSubmitter)(nil), s)
	assert.Implements(t, (*client.ValidatorBalancesProvider)(nil), s)
	assert.Implements(t, (*client.ValidatorsProvider)(nil), s)
	assert.Implements(t, (*client.VoluntaryExitSubmitter)(nil), s)
	assert.Implements(t, (*client.VoluntaryExitPoolProvider)(nil), s)

	// Non-standard extensions.
	assert.Implements(t, (*client.DomainProvider)(nil), s)
	assert.Implements(t, (*client.GenesisTimeProvider)(nil), s)
}
