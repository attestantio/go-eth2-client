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

package heze

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// ExecutionPayloadBid represents an execution payload bid for EIP-7732.
type ExecutionPayloadBid struct {
	ParentBlockHash    phase0.Hash32              `ssz-size:"32"`
	ParentBlockRoot    phase0.Root                `ssz-size:"32"`
	BlockHash          phase0.Hash32              `ssz-size:"32"`
	PrevRandao         phase0.Root                `ssz-size:"32"`
	FeeRecipient       bellatrix.ExecutionAddress `ssz-size:"20"`
	GasLimit           uint64
	BuilderIndex       gloas.BuilderIndex
	Slot               phase0.Slot
	Value              phase0.Gwei
	ExecutionPayment   phase0.Gwei
	BlobKZGCommitments []deneb.KZGCommitment `dynssz-max:"MAX_BLOB_COMMITMENTS_PER_BLOCK"   ssz-max:"4096" ssz-size:"?,48"`
	InclusionListBits  []byte                `dynssz-size:"INCLUSION_LIST_COMMITTEE_SIZE/8" ssz-size:"2"`
}

// String returns a string version of the structure.
func (e *ExecutionPayloadBid) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
