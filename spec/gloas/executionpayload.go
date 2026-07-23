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

// BlockAccessList is the RLP-encoded block access list (EIP-7928).
type BlockAccessList []byte

// ExecutionPayload represents an execution layer payload.
type ExecutionPayload struct {
	ParentHash      phase0.Hash32              `ssz-index:"0"`
	FeeRecipient    bellatrix.ExecutionAddress `ssz-index:"1"`
	StateRoot       phase0.Root                `ssz-index:"2"`
	ReceiptsRoot    phase0.Root                `ssz-index:"3"`
	LogsBloom       [256]byte                  `ssz-index:"4"`
	PrevRandao      [32]byte                   `ssz-index:"5"`
	BlockNumber     uint64                     `ssz-index:"6"`
	GasLimit        uint64                     `ssz-index:"7"`
	GasUsed         uint64                     `ssz-index:"8"`
	Timestamp       uint64                     `ssz-index:"9"`
	ExtraData       []byte                     `dynssz-max:"MAX_EXTRA_DATA_BYTES" ssz-index:"10"                               ssz-max:"32"`
	BaseFeePerGas   *uint256.Int               `ssz-index:"11"                    ssz-type:"uint256"`
	BlockHash       phase0.Hash32              `ssz-index:"12"`
	Transactions    []bellatrix.Transaction    `ssz-index:"13"                    ssz-type:"progressive-list,progressive-list"`
	Withdrawals     []*capella.Withdrawal      `ssz-index:"14"                    ssz-type:"progressive-list"`
	BlobGasUsed     uint64                     `ssz-index:"15"`
	ExcessBlobGas   uint64                     `ssz-index:"16"`
	BlockAccessList BlockAccessList            `ssz-index:"17"                    ssz-type:"progressive-list"`
	SlotNumber      uint64                     `ssz-index:"18"`
}

// String returns a string version of the structure.
func (e *ExecutionPayload) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
