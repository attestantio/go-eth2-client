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

package eip7732

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// ExecutionPayloadHeader represents an execution payload header for EIP-7732.
type ExecutionPayloadHeader struct {
	ParentBlockHash        phase0.Hash32 `ssz-size:"32"`
	ParentBlockRoot        phase0.Root   `ssz-size:"32"`
	BlockHash              phase0.Hash32 `ssz-size:"32"`
	GasLimit               uint64
	BuilderIndex           phase0.ValidatorIndex
	Slot                   phase0.Slot
	Value                  uint64      // Gwei
	BlobKZGCommitmentsRoot phase0.Root `ssz-size:"32"`
}

// String returns a string version of the structure.
func (e *ExecutionPayloadHeader) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
