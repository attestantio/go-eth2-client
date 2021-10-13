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

package http_test

import (
	"context"
	"os"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	v1 "github.com/attestantio/go-eth2-client/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		parameters []v1.Parameter
		location   string
		err        string
	}{
		{
			name: "Nil",
			err:  "problem with parameters: no address specified",
		},
		{
			name: "AddressNil",
			parameters: []v1.Parameter{
				v1.WithTimeout(5 * time.Second),
			},
			err: "problem with parameters: no address specified",
		},
		{
			name: "TimeoutZero",
			parameters: []v1.Parameter{
				v1.WithAddress(os.Getenv("HTTP_ADDRESS")),
				v1.WithTimeout(0),
			},
			err: "problem with parameters: no timeout specified",
		},
		{
			name: "AddressInvalid",
			parameters: []v1.Parameter{
				v1.WithAddress(string([]byte{0x01})),
				v1.WithTimeout(5 * time.Second),
			},
			err: "invalid URL: parse \"http://\\x01\": net/url: invalid control character in URL",
		},
		{
			name: "IndexChunkSizeZero",
			parameters: []v1.Parameter{
				v1.WithAddress(os.Getenv("HTTP_ADDRESS")),
				v1.WithTimeout(5 * time.Second),
				v1.WithIndexChunkSize(0),
			},
			err: "problem with parameters: no index chunk size specified",
		},
		{
			name: "PubKeyChunkSizeZero",
			parameters: []v1.Parameter{
				v1.WithAddress(os.Getenv("HTTP_ADDRESS")),
				v1.WithTimeout(5 * time.Second),
				v1.WithPubKeyChunkSize(0),
			},
			err: "problem with parameters: no public key chunk size specified",
		},
		{
			name: "Good",
			parameters: []v1.Parameter{
				v1.WithAddress(os.Getenv("HTTP_ADDRESS")),
				v1.WithTimeout(5 * time.Second),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := v1.New(ctx, test.parameters...)
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
	s, err := v1.New(ctx, v1.WithAddress(os.Getenv("HTTP_ADDRESS")), v1.WithTimeout(5*time.Second))
	require.NoError(t, err)

	// Standard interfacs.
	assert.Implements(t, (*client.AggregateAttestationProvider)(nil), s)
	assert.Implements(t, (*client.AggregateAttestationsSubmitter)(nil), s)
	assert.Implements(t, (*client.AttestationDataProvider)(nil), s)
	assert.Implements(t, (*client.AttestationsSubmitter)(nil), s)
	assert.Implements(t, (*client.AttesterDutiesProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockHeadersProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockProposalProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockRootProvider)(nil), s)
	assert.Implements(t, (*client.BeaconBlockSubmitter)(nil), s)
	assert.Implements(t, (*client.BeaconCommitteeSubscriptionsSubmitter)(nil), s)
	assert.Implements(t, (*client.BeaconStateProvider)(nil), s)
	assert.Implements(t, (*client.DepositContractProvider)(nil), s)
	assert.Implements(t, (*client.EventsProvider)(nil), s)
	assert.Implements(t, (*client.ForkProvider)(nil), s)
	assert.Implements(t, (*client.ForkScheduleProvider)(nil), s)
	assert.Implements(t, (*client.GenesisProvider)(nil), s)
	assert.Implements(t, (*client.NodeSyncingProvider)(nil), s)
	assert.Implements(t, (*client.ProposerDutiesProvider)(nil), s)
	assert.Implements(t, (*client.SpecProvider)(nil), s)
	assert.Implements(t, (*client.SyncCommitteeDutiesProvider)(nil), s)
	assert.Implements(t, (*client.SyncCommitteeMessagesSubmitter)(nil), s)
	assert.Implements(t, (*client.SyncCommitteeSubscriptionsSubmitter)(nil), s)
	assert.Implements(t, (*client.ValidatorBalancesProvider)(nil), s)
	assert.Implements(t, (*client.ValidatorsProvider)(nil), s)
	assert.Implements(t, (*client.VoluntaryExitSubmitter)(nil), s)

	// Non-standard extensions.
	assert.Implements(t, (*client.DomainProvider)(nil), s)
	assert.Implements(t, (*client.GenesisTimeProvider)(nil), s)
}
