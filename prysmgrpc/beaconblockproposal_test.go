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
	"time"

	"github.com/attestantio/go-eth2-client/prysmgrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBeaconBlockProposal(t *testing.T) {
	tests := []struct {
		name         string
		randaoReveal []byte
		graffiti     []byte
	}{
		{
			name: "Good",
			randaoReveal: []byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
				0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
				0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f,
				0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
			},
			graffiti: []byte("TestGood"),
		},
	}

	service, err := prysmgrpc.New(context.Background(),
		prysmgrpc.WithAddress(os.Getenv("PRYSMGRPC_ADDRESS")),
		prysmgrpc.WithTimeout(timeout),
	)
	require.NoError(t, err)

	// Need to fetch current slot for proposal.
	genesisTime, err := service.GenesisTime(context.Background())
	require.NoError(t, err)
	slotDuration, err := service.SlotDuration(context.Background())
	require.NoError(t, err)

	for _, test := range tests {
		nextSlot := uint64(time.Since(genesisTime).Seconds()/slotDuration.Seconds()) + 1
		t.Run(test.name, func(t *testing.T) {
			block, err := service.BeaconBlockProposal(context.Background(), nextSlot, test.randaoReveal, test.graffiti)
			require.NoError(t, err)
			require.NotNil(t, block)

			// Graffiti is always fixed 32 bytes.
			fixedGraffiti := make([]byte, 32)
			copy(fixedGraffiti, test.graffiti)
			assert.Equal(t, fixedGraffiti, block.Body.Graffiti)
			assert.Equal(t, test.randaoReveal, block.Body.RANDAOReveal)
		})
	}
}
