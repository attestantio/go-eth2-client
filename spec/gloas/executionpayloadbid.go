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

package gloas

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// ExecutionPayloadBid represents an execution payload bid for EIP-7732.
type ExecutionPayloadBid struct {
	ParentBlockHash       phase0.Hash32              `ssz-index:"0"`
	ParentBlockRoot       phase0.Root                `ssz-index:"1"`
	BlockHash             phase0.Hash32              `ssz-index:"2"`
	PrevRandao            phase0.Root                `ssz-index:"3"`
	FeeRecipient          bellatrix.ExecutionAddress `ssz-index:"4"`
	GasLimit              uint64                     `ssz-index:"5"`
	BuilderIndex          BuilderIndex               `ssz-index:"6"`
	Slot                  phase0.Slot                `ssz-index:"7"`
	Value                 phase0.Gwei                `ssz-index:"8"`
	ExecutionPayment      phase0.Gwei                `ssz-index:"9"`
	BlobKZGCommitments    []deneb.KZGCommitment      `ssz-index:"10" ssz-type:"progressive-list"`
	ExecutionRequestsRoot phase0.Root                `ssz-index:"11"`
}

// String returns a string version of the structure.
func (e *ExecutionPayloadBid) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
