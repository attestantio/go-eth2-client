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

package prysmgrpc_test

import (
	"context"
	"os"
	"testing"

	"github.com/attestantio/go-eth2-client/prysmgrpc"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestValidatorBalances(t *testing.T) {
	tests := []struct {
		name       string
		stateID    string
		validators []spec.BLSPubKey
	}{
		{
			name:    "Old",
			stateID: "32",
		},
		{
			name:    "Single",
			stateID: "head",
			validators: []spec.BLSPubKey{
				_blsPubKey("0xb2007d1354db791b924fd35a6b0a8525266a021765b54641f4d415daa50c511204d6acc213a23468f2173e60cc950e26"),
			},
		},
		{
			name:    "All",
			stateID: "head",
		},
	}

	service, err := prysmgrpc.New(context.Background(),
		prysmgrpc.WithAddress(os.Getenv("PRYSMGRPC_ADDRESS")),
		prysmgrpc.WithTimeout(timeout),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			balances, err := service.PrysmValidatorBalances(context.Background(), test.stateID, test.validators)
			require.NoError(t, err)
			require.NotNil(t, balances)
		})
	}
}
