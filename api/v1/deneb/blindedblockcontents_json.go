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

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/pkg/errors"
)

// blindedBlockContentsJSON is the spec representation of the struct.
type blindedBlockContentsJSON struct {
	BlindedBlock        *BlindedBeaconBlock   `json:"blinded_block"`
	BlindedBlobSidecars []*BlindedBlobSidecar `json:"blinded_blob_sidecars"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlindedBlockContents) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blindedBlockContentsJSON{
		BlindedBlock:        b.BlindedBlock,
		BlindedBlobSidecars: b.BlindedBlobSidecars,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlindedBlockContents) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&blindedBlockContentsJSON{}, input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw["blinded_block"], &b.BlindedBlock); err != nil {
		return errors.Wrap(err, "blinded_block")
	}

	if err := json.Unmarshal(raw["blinded_blob_sidecars"], &b.BlindedBlobSidecars); err != nil {
		return errors.Wrap(err, "blinded_blob_sidecars")
	}

	return nil
}
