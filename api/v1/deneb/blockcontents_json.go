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
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/pkg/errors"
)

// blockContentsJSON is the spec representation of the struct.
type blockContentsJSON struct {
	Block        *deneb.BeaconBlock   `json:"block"`
	BlobSidecars []*deneb.BlobSidecar `json:"blob_sidecars"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlockContents) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blockContentsJSON{
		Block:        b.Block,
		BlobSidecars: b.BlobSidecars,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlockContents) UnmarshalJSON(input []byte) error {
	var data blockContentsJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	return b.unpack(&data)
}

func (b *BlockContents) unpack(data *blockContentsJSON) error {
	if data.Block == nil {
		return errors.New("block missing")
	}
	b.Block = data.Block

	if data.BlobSidecars == nil {
		return errors.New("blob sidecars missing")
	}
	b.BlobSidecars = data.BlobSidecars

	return nil
}
