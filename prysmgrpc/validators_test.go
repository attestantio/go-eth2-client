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
	"bytes"
	"context"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/prysmgrpc"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func _bytes(input string) []byte {
	res, err := hex.DecodeString(strings.TrimPrefix(input, "0x"))
	if err != nil {
		panic(err)
	}
	return res
}

func TestValidators(t *testing.T) {
	tests := []struct {
		name        string
		genesisFork []byte
		stateID     string
		validators  []spec.BLSPubKey
	}{
		{
			name:    "Genesis",
			stateID: "genesis",
		},
		{
			name:    "Old",
			stateID: "32",
		},
		{
			name:        "Single",
			genesisFork: _bytes("0x00002009"),
			stateID:     "head",
			validators: []spec.BLSPubKey{
				_blsPubKey("0xb1aa1fbe5851d7477ba12042f05bf406771471a118252bb1d455a184af23a4f317d854668f683aba629dcd3f698ba7b7"),
			},
		},
		{
			name:        "Unknown",
			genesisFork: _bytes("0x00002009"),
			stateID:     "head",
			validators: []spec.BLSPubKey{
				_blsPubKey("0xb1aa1fbe5851d7477ba12042f05bf406771471a118252bb1d455a184af23a4f317d854668f683aba629dcd3f698ba7b7"),
				_blsPubKey("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			},
		},
		{
			name:        "Issue4",
			genesisFork: _bytes("0x00002009"),
			stateID:     "head",
			validators: []spec.BLSPubKey{
				_blsPubKey("0xb1aa1fbe5851d7477ba12042f05bf406771471a118252bb1d455a184af23a4f317d854668f683aba629dcd3f698ba7b7"),
				_blsPubKey("0xa7e6da76277e0ab2fbfc69ce532fa71cf7ea977102042ac5ed9b4ea6adf1346a56bce68d8f731d682d6eeaa8cab06cc8"),
				_blsPubKey("0x90c39607ff913c77b1d5565143d2d75a16e07ef32c943ea055aef73af66c2d430c0f59980041bd99b55bac67dd597cf4"),
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

	config, err := service.Spec(context.Background())
	require.NoError(t, err)
	genesisFork, exists := config["GENESIS_FORK_VERSION"]
	require.True(t, exists)
	fork := genesisFork.(spec.Version)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if len(test.genesisFork) > 0 && !bytes.Equal(fork[:], test.genesisFork) {
				t.Skip("test not for service fork version")
			}
			validators, err := service.ValidatorsByPubKey(context.Background(), test.stateID, test.validators)
			require.NoError(t, err)
			require.NotNil(t, validators)
			require.NotEqual(t, 0, len(validators))
			// Not all validators will have a balance, but at least one must.
			hasBalance := false
			for i := range validators {
				if validators[i].Balance != 0 {
					hasBalance = true
					break
				}
			}
			require.True(t, hasBalance)
		})
	}
}

func TestValidatorsWithoutBalance(t *testing.T) {
	tests := []struct {
		name        string
		genesisFork []byte
		stateID     string
		validators  []spec.BLSPubKey
	}{
		{
			name:    "Genesis",
			stateID: "genesis",
		},
		{
			name:    "Old",
			stateID: "32",
		},
		{
			name:        "Single",
			genesisFork: _bytes("0x00002009"),
			stateID:     "head",
			validators: []spec.BLSPubKey{
				_blsPubKey("0xb1aa1fbe5851d7477ba12042f05bf406771471a118252bb1d455a184af23a4f317d854668f683aba629dcd3f698ba7b7"),
			},
		},
		{
			name:        "Unknown",
			genesisFork: _bytes("0x00002009"),
			stateID:     "head",
			validators: []spec.BLSPubKey{
				_blsPubKey("0xb1aa1fbe5851d7477ba12042f05bf406771471a118252bb1d455a184af23a4f317d854668f683aba629dcd3f698ba7b7"),
				_blsPubKey("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
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

	config, err := service.Spec(context.Background())
	require.NoError(t, err)
	genesisFork, exists := config["GENESIS_FORK_VERSION"]
	require.True(t, exists)
	fork := genesisFork.(spec.Version)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if len(test.genesisFork) > 0 && !bytes.Equal(fork[:], test.genesisFork) {
				t.Skip("test not for service fork version")
			}
			validators, err := service.ValidatorsWithoutBalanceByPubKey(context.Background(), test.stateID, test.validators)
			require.NoError(t, err)
			require.NotNil(t, validators)
			require.NotEqual(t, 0, len(validators))
			for i := range validators {
				require.Equal(t, spec.Gwei(0), validators[i].Balance)
			}
		})
	}
}
