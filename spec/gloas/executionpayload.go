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

package gloas

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/holiman/uint256"
)

// ExecutionPayload represents an execution layer payload.
//
//nolint:revive
type ExecutionPayload struct {
	ParentHash      phase0.Hash32              `ssz-size:"32"`
	FeeRecipient    bellatrix.ExecutionAddress `ssz-size:"20"`
	StateRoot       phase0.Root                `ssz-size:"32"`
	ReceiptsRoot    phase0.Root                `ssz-size:"32"`
	LogsBloom       [256]byte                  `ssz-size:"256"`
	PrevRandao      [32]byte                   `ssz-size:"32"`
	BlockNumber     uint64
	GasLimit        uint64
	GasUsed         uint64
	Timestamp       uint64
	ExtraData       []byte                  `dynssz-max:"MAX_EXTRA_DATA_BYTES"                                   ssz-max:"32"`
	BaseFeePerGas   *uint256.Int            `ssz-size:"32"`
	BlockHash       phase0.Hash32           `ssz-size:"32"`
	Transactions    []bellatrix.Transaction `dynssz-max:"MAX_TRANSACTIONS_PER_PAYLOAD,MAX_BYTES_PER_TRANSACTION" ssz-max:"1048576,1073741824" ssz-size:"?,?"`
	Withdrawals     []*capella.Withdrawal   `dynssz-max:"MAX_WITHDRAWALS_PER_PAYLOAD"                            ssz-max:"16"`
	BlobGasUsed     uint64
	ExcessBlobGas   uint64
	BlockAccessList BlockAccessList `dynssz-max:"MAX_BYTES_PER_TRANSACTION" ssz-max:"1073741824"`
	SlotNumber      uint64
}

// String returns a string version of the structure.
func (e *ExecutionPayload) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
