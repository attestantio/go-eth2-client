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

package http_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/http"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBeaconBlockProposal(t *testing.T) {
	tests := []struct {
		name         string
		randaoReveal spec.BLSSignature
		graffiti     []byte
	}{
		{
			name: "Good",
			randaoReveal: spec.BLSSignature([96]byte{
				0x8d, 0x7b, 0x2a, 0x32, 0xb0, 0x26, 0xe9, 0xc7, 0x9a, 0xae, 0x6e, 0xc6, 0xb8, 0x3e, 0xab, 0xae,
				0x89, 0xd6, 0x0c, 0xac, 0xd6, 0x5a, 0xc4, 0x1e, 0xd7, 0xd2, 0xf4, 0xbe, 0x9d, 0xd8, 0xc8, 0x9c,
				0x1b, 0xf7, 0xcd, 0x3d, 0x70, 0x03, 0x74, 0xe1, 0x8d, 0x03, 0xd1, 0x2f, 0x6a, 0x05, 0x4c, 0x23,
				0x00, 0x6f, 0x64, 0xf0, 0xe4, 0xe8, 0xb7, 0xcf, 0x37, 0xd6, 0xac, 0x9a, 0x4c, 0x7d, 0x81, 0x5c,
				0x85, 0x81, 0x20, 0xc5, 0x46, 0x73, 0xb7, 0xd3, 0xcb, 0x2b, 0xb1, 0x55, 0x0a, 0x4d, 0x65, 0x9e,
				0xaf, 0x46, 0xe3, 0x45, 0x15, 0x67, 0x7c, 0x67, 0x8b, 0x70, 0xd6, 0xf6, 0x2d, 0xbf, 0x89, 0xf0,
			}),
			graffiti: []byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			},
		},
	}

	service, err := http.New(context.Background(),
		http.WithTimeout(timeout),
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	// Need to fetch current slot for proposal.
	genesis, err := service.Genesis(context.Background())
	require.NoError(t, err)
	slotDuration, err := service.SlotDuration(context.Background())
	require.NoError(t, err)

	for _, test := range tests {
		nextSlot := spec.Slot(uint64(time.Since(genesis.GenesisTime).Seconds())/uint64(slotDuration.Seconds())) + 1
		t.Run(test.name, func(t *testing.T) {
			resp, err := service.BeaconBlockProposal(context.Background(), nextSlot, test.randaoReveal, test.graffiti)
			require.NoError(t, err)
			require.NotNil(t, resp)
			if resp.Phase0 != nil {
				assert.Equal(t, test.graffiti, resp.Phase0.Body.Graffiti)
				assert.Equal(t, test.randaoReveal, resp.Phase0.Body.RANDAOReveal)
			}
			if resp.Altair != nil {
				assert.Equal(t, test.graffiti, resp.Altair.Body.Graffiti)
				assert.Equal(t, test.randaoReveal, resp.Altair.Body.RANDAOReveal)
			}
		})
	}
}
