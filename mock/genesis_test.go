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

package mock_test

import (
	"context"
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Good",
		},
	}

	service, err := mock.New(context.Background())
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.Genesis(context.Background(), &api.GenesisOpts{})
			require.NoError(t, err)
			require.NotNil(t, response)
			require.NotNil(t, response.Data)
			require.NotNil(t, response.Data.GenesisTime)
			require.NotNil(t, response.Data.GenesisValidatorsRoot)
			require.NotNil(t, response.Data.GenesisForkVersion)
		})
	}
}
