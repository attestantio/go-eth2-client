// Copyright © 2023 Attestant Limited.
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

package deneb

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/holiman/uint256"
)

// ExecutionPayloadHeader represents an execution layer payload header.
type ExecutionPayloadHeader struct {
	ParentHash       phase0.Hash32              `ssz-size:"32"`
	FeeRecipient     bellatrix.ExecutionAddress `ssz-size:"20"`
	StateRoot        phase0.Root                `ssz-size:"32"`
	ReceiptsRoot     phase0.Root                `ssz-size:"32"`
	LogsBloom        [256]byte                  `ssz-size:"256"`
	PrevRandao       [32]byte                   `ssz-size:"32"`
	BlockNumber      uint64
	GasLimit         uint64
	GasUsed          uint64
	Timestamp        uint64
	ExtraData        []byte        `ssz-max:"32"`
	BaseFeePerGas    *uint256.Int  `ssz-size:"32"`
	BlockHash        phase0.Hash32 `ssz-size:"32"`
	TransactionsRoot phase0.Root   `ssz-size:"32"`
	WithdrawalsRoot  phase0.Root   `ssz-size:"32"`
	BlobGasUsed      uint64
	ExcessBlobGas    uint64
}

// String returns a string version of the structure.
func (e *ExecutionPayloadHeader) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
