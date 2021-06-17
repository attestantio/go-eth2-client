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

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/multi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	ctx := context.Background()

	eth2Client1, err := mock.New(ctx)
	require.NoError(t, err)
	eth2Client2, err := mock.New(ctx)
	require.NoError(t, err)
	eth2Client3, err := mock.New(ctx)
	require.NoError(t, err)
	inactiveEth2Client1, err := mock.New(ctx)
	require.NoError(t, err)
	inactiveEth2Client1.SyncDistance = 10
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
				multi.WithClients([]eth2client.Service{
					inactiveEth2Client1,
				}),
			},
			err: "No  clients active, cannot proceed",
		},
		{
			name: "Good",
			params: []multi.Parameter{
				multi.WithLogLevel(zerolog.Disabled),
				multi.WithClients([]eth2client.Service{
					eth2Client1,
					eth2Client2,
					eth2Client3,
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
