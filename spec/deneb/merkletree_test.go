// Copyright © 2026 Attestant Limited.
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

package deneb_test

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/pk910/dynamic-ssz/sszutils"
	require "github.com/stretchr/testify/require"
)

// TestMerkleTreeMatchesHashTreeRoot checks that the merkle tree of the execution payload containers hashes to their hash
// tree root.  BaseFeePerGas is hashed through a buffered append followed by fields added directly, which dynamic-ssz
// before v1.3.3 placed out of order in the tree, so proofs generated from it did not verify (issue #321).
func TestMerkleTreeMatchesHashTreeRoot(t *testing.T) {
	payload := &deneb.ExecutionPayload{
		ParentHash:    phase0.Hash32{0x01},
		FeeRecipient:  bellatrix.ExecutionAddress{0x02},
		StateRoot:     phase0.Root{0x03},
		ReceiptsRoot:  phase0.Root{0x04},
		BlockNumber:   5,
		GasLimit:      6,
		GasUsed:       7,
		Timestamp:     8,
		ExtraData:     []byte{0x09},
		BaseFeePerGas: uint256.NewInt(10),
		BlockHash:     phase0.Hash32{0x0b},
		Transactions:  []bellatrix.Transaction{{0x0c}},
		Withdrawals:   []*capella.Withdrawal{{Index: 13, ValidatorIndex: 14, Address: bellatrix.ExecutionAddress{0x0f}, Amount: 16}},
		BlobGasUsed:   17,
		ExcessBlobGas: 18,
	}
	header := &deneb.ExecutionPayloadHeader{
		ParentHash:       payload.ParentHash,
		FeeRecipient:     payload.FeeRecipient,
		StateRoot:        payload.StateRoot,
		ReceiptsRoot:     payload.ReceiptsRoot,
		BlockNumber:      payload.BlockNumber,
		GasLimit:         payload.GasLimit,
		GasUsed:          payload.GasUsed,
		Timestamp:        payload.Timestamp,
		ExtraData:        payload.ExtraData,
		BaseFeePerGas:    payload.BaseFeePerGas,
		BlockHash:        payload.BlockHash,
		TransactionsRoot: phase0.Root{0x13},
		WithdrawalsRoot:  phase0.Root{0x14},
		BlobGasUsed:      payload.BlobGasUsed,
		ExcessBlobGas:    payload.ExcessBlobGas,
	}

	tests := []struct {
		name  string
		input sszutils.FastsszHashRoot
	}{
		{
			name:  "ExecutionPayload",
			input: payload,
		},
		{
			name:  "ExecutionPayloadHeader",
			input: header,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root, err := test.input.HashTreeRoot()
			require.NoError(t, err)
			tree, err := dynssz.GetGlobalDynSsz().GetTree(test.input)
			require.NoError(t, err)
			require.Equal(t, root[:], tree.Hash())
		})
	}
}
