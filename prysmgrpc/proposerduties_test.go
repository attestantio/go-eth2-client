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
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/prysmgrpc"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func _blsPubKey(input string) spec.BLSPubKey {
	tmp, err := hex.DecodeString(strings.TrimPrefix(input, "0x"))
	if err != nil {
		panic(err)
	}
	var res spec.BLSPubKey
	copy(res[:], tmp)
	return res
}

func TestProposerDuties(t *testing.T) {
	tests := []struct {
		name     string
		epoch    spec.Epoch
		indices  []spec.ValidatorIndex
		expected int
	}{
		{
			name:     "Old",
			epoch:    1,
			expected: 32,
		},
		{
			name:     "Recent",
			epoch:    10990,
			expected: 32,
		},
		{
			name:     "GoodWithValidators",
			epoch:    4092,
			indices:  []spec.ValidatorIndex{33566, 25444},
			expected: 2,
		},
	}

	service, err := prysmgrpc.New(context.Background(),
		prysmgrpc.WithAddress(os.Getenv("PRYSMGRPC_ADDRESS")),
		prysmgrpc.WithTimeout(timeout),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			duties, err := service.ProposerDuties(context.Background(), test.epoch, test.indices)
			require.NoError(t, err)
			require.NotNil(t, duties)
			require.Equal(t, test.expected, len(duties))
		})
	}
}
