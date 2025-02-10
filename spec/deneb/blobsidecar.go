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

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
)

// BlobSidecar represents a data blob sidecar.
type BlobSidecar struct {
	Index                       BlobIndex
	Blob                        Blob          `ssz-size:"131072"`
	KZGCommitment               KZGCommitment `ssz-size:"48"`
	KZGProof                    KZGProof      `ssz-size:"48"`
	SignedBlockHeader           *phase0.SignedBeaconBlockHeader
	KZGCommitmentInclusionProof KZGCommitmentInclusionProof `ssz-size:"544"`
}

// String returns a string version of the structure.
func (b *BlobSidecar) String() string {
	data, err := yaml.Marshal(b)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
