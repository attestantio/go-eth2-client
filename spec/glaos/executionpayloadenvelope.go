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

package glaos

import (
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// ExecutionPayloadEnvelope represents an execution payload envelope.
type ExecutionPayloadEnvelope struct {
	Payload            *deneb.ExecutionPayload
	ExecutionRequests  *electra.ExecutionRequests
	BuilderIndex       phase0.ValidatorIndex
	BeaconBlockRoot    phase0.Root `ssz-size:"32"`
	Slot               phase0.Slot
	BlobKZGCommitments []deneb.KZGCommitment `dynssz-max:"MAX_BLOB_COMMITMENTS_PER_BLOCK" ssz-max:"4096" ssz-size:"?,48"`
	StateRoot          phase0.Root           `ssz-size:"32"`
}

// String returns a string version of the structure.
func (e *ExecutionPayloadEnvelope) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
