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

	spec "github.com/attestantio/go-eth2-client/spec/altair"
	standardhttp "github.com/attestantio/go-eth2-client/standardhttp/v2"
	"github.com/stretchr/testify/require"
)

func TestValidatorBalances(t *testing.T) {
	tests := []struct {
		name       string
		stateID    string
		validators []spec.ValidatorIndex
	}{
		{
			name:       "Single",
			stateID:    "head",
			validators: []spec.ValidatorIndex{1000},
		},
		{
			name:    "All",
			stateID: "head",
		},
	}

	service, err := standardhttp.New(context.Background(),
		standardhttp.WithTimeout(timeout),
		standardhttp.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			balances, err := service.ValidatorBalances(context.Background(), test.stateID, test.validators)
			require.NoError(t, err)
			require.NotNil(t, balances)
		})
	}
}
