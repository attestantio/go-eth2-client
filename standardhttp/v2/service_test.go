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

	standardhttp "github.com/attestantio/go-eth2-client/standardhttp/v2"
	client "github.com/attestantio/go-eth2-client/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		parameters []standardhttp.Parameter
		location   string
		err        string
	}{
		{
			name: "Nil",
			err:  "problem with parameters: no address specified",
		},
		{
			name: "AddressNil",
			parameters: []standardhttp.Parameter{
				standardhttp.WithTimeout(5 * time.Second),
			},
			err: "problem with parameters: no address specified",
		},
		{
			name: "TimeoutZero",
			parameters: []standardhttp.Parameter{
				standardhttp.WithAddress(os.Getenv("HTTP_ADDRESS")),
				standardhttp.WithTimeout(0),
			},
			err: "problem with parameters: no timeout specified",
		},
		{
			name: "AddressInvalid",
			parameters: []standardhttp.Parameter{
				standardhttp.WithAddress(string([]byte{0x01})),
				standardhttp.WithTimeout(5 * time.Second),
			},
			err: "invalid URL: parse \"http://\\x01\": net/url: invalid control character in URL",
		},
		{
			name: "Good",
			parameters: []standardhttp.Parameter{
				standardhttp.WithAddress(os.Getenv("HTTP_ADDRESS")),
				standardhttp.WithTimeout(5 * time.Second),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := standardhttp.New(ctx, test.parameters...)
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
	s, err := standardhttp.New(ctx, standardhttp.WithAddress(os.Getenv("HTTP_ADDRESS")), standardhttp.WithTimeout(5*time.Second))
	require.NoError(t, err)

	// Standard interfacs.
	assert.Implements(t, (*client.AggregateAttestationProvider)(nil), s)
	assert.Implements(t, (*client.AggregateAttestationsSubmitter)(nil), s)
	assert.Implements(t, (*client.AttestationDataProvider)(nil), s)
	assert.Implements(t, (*client.AttestationPoolProvider)(nil), s)
	assert.Implements(t, (*client.AttestationsSubmitter)(nil), s)
	assert.Implements(t, (*client.AttesterDutiesProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockHeadersProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockProposalProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockSubmitter)(nil), s)
	assert.Implements(t, (*client.BeaconCommitteesProvider)(nil), s)
	assert.Implements(t, (*client.BeaconCommitteeSubscriptionsSubmitter)(nil), s)
	assert.Implements(t, (*client.BeaconStateProvider)(nil), s)
	assert.Implements(t, (*client.DepositContractProvider)(nil), s)
	assert.Implements(t, (*client.EventsProvider)(nil), s)
	assert.Implements(t, (*client.ForkProvider)(nil), s)
	assert.Implements(t, (*client.FinalityProvider)(nil), s)
	assert.Implements(t, (*client.ForkProvider)(nil), s)
	assert.Implements(t, (*client.ForkScheduleProvider)(nil), s)
	assert.Implements(t, (*client.GenesisProvider)(nil), s)
	assert.Implements(t, (*client.NodeSyncingProvider)(nil), s)
	assert.Implements(t, (*client.NodeVersionProvider)(nil), s)
	assert.Implements(t, (*client.ProposerDutiesProvider)(nil), s)
	assert.Implements(t, (*client.SignedBeaconBlockProvider)(nil), s)
	assert.Implements(t, (*client.SpecProvider)(nil), s)
	assert.Implements(t, (*client.ValidatorBalancesProvider)(nil), s)
	assert.Implements(t, (*client.ValidatorsProvider)(nil), s)
	assert.Implements(t, (*client.VoluntaryExitSubmitter)(nil), s)

	// Non-standard extensions.
	assert.Implements(t, (*client.DomainProvider)(nil), s)
}
