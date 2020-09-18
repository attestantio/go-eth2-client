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

package lighthousehttp_test

import (
	"context"
	"os"
	"testing"

	"github.com/attestantio/go-eth2-client/lighthousehttp"
	"github.com/stretchr/testify/require"
)

func TestSignedBeaconBlockBySlot(t *testing.T) {
	tests := []struct {
		name string
		slot uint64
	}{
		{
			name: "Missing",
			slot: 1,
		},
		{
			name: "Good",
			slot: 12345,
		},
		{
			name: "AttesterSlashing",
			slot: 62390,
		},
		{
			name: "ProposerSlashing",
			slot: 33890,
		},
		{
			name: "Deposits",
			slot: 1005,
		},
	}

	ctx := context.Background()
	service, err := lighthousehttp.New(ctx, lighthousehttp.WithAddress(os.Getenv("LIGHTHOUSEHTTP_ADDRESS")))
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			block, err := service.SignedBeaconBlockBySlot(context.Background(), test.slot)
			require.NoError(t, err)
			require.NotNil(t, block)
		})
	}
}
