// Copyright © 2020 Attestant Limited.
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

package tekuhttp_test

import (
	"context"
	"os"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/tekuhttp"
	"github.com/stretchr/testify/require"
)

func TestSubmitBeaconCommitteeSubscriptions(t *testing.T) {
	tests := []struct {
		name          string
		epoch         uint64
		subscriptions []*client.BeaconCommitteeSubscription
	}{
		{
			name:  "Good",
			epoch: 1,
			subscriptions: []*client.BeaconCommitteeSubscription{
				{
					Slot:           1,
					CommitteeIndex: 1,
					CommitteeSize:  128,
					ValidatorIndex: 2,
					ValidatorPubKey: []byte{
						0x8c, 0x2f, 0x53, 0x5d, 0x3b, 0xec, 0x65, 0xf9, 0x5c, 0xb4, 0xba, 0x45, 0x55, 0x66, 0xe4, 0xec,
						0x3d, 0xe8, 0xda, 0x5c, 0x13, 0xa6, 0x81, 0x69, 0x9e, 0x0f, 0x80, 0xd7, 0x94, 0x2d, 0x6f, 0xdc,
						0xbc, 0xef, 0x18, 0xc8, 0xcf, 0x18, 0xf9, 0xda, 0x14, 0xaa, 0x37, 0x9b, 0xdd, 0x6d, 0x29, 0xc5,
					},
					Aggregate: true,
				},
			},
		},
	}

	service, err := tekuhttp.New(context.Background(),
		tekuhttp.WithAddress(os.Getenv("TEKUHTTP_ADDRESS")),
		tekuhttp.WithTimeout(timeout),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.SubmitBeaconCommitteeSubscriptions(context.Background(), test.subscriptions)
			require.NoError(t, err)
		})
	}
}
