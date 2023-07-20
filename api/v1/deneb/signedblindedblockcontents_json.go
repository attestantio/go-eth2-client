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

// signedBlindedBlockContentsJSON is the spec representation of the struct.
type signedBlindedBlockContentsJSON struct {
	SignedBlindedBlock        *SignedBlindedBeaconBlock   `json:"signed_blinded_block"`
	SignedBlindedBlobSidecars []*SignedBlindedBlobSidecar `json:"signed_blinded_blob_sidecars"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedBlindedBlockContents) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedBlindedBlockContentsJSON{
		SignedBlindedBlock:        s.SignedBlindedBlock,
		SignedBlindedBlobSidecars: s.SignedBlindedBlobSidecars,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedBlindedBlockContents) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&signedBlindedBlockContentsJSON{}, input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw["signed_blinded_block"], &s.SignedBlindedBlock); err != nil {
		return errors.Wrap(err, "signed_blinded_block")
	}

	if err := json.Unmarshal(raw["signed_blinded_blob_sidecars"], &s.SignedBlindedBlobSidecars); err != nil {
		return errors.Wrap(err, "signed_blinded_blob_sidecars")
	}

	return nil
}
