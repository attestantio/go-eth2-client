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

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// BlindedBlobSidecar represents a data blob sidecar.
type BlindedBlobSidecar struct {
	BlockRoot       phase0.Root `ssz-size:"32"`
	Index           deneb.BlobIndex
	Slot            phase0.Slot
	BlockParentRoot phase0.Root ` ssz-size:"32"`
	ProposerIndex   phase0.ValidatorIndex
	BlobRoot        phase0.Root         `ssz-size:"32"`
	KzgCommitment   deneb.KzgCommitment `ssz-size:"48"`
	KzgProof        deneb.KzgProof      `ssz-size:"48"`
}

// String returns a string version of the structure.
func (b *BlindedBlobSidecar) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
