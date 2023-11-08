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
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/pkg/errors"
)

// signedBlockContentsJSON is the spec representation of the struct.
type signedBlockContentsJSON struct {
	SignedBlock *deneb.SignedBeaconBlock `json:"signed_block"`
	KZGProofs   []deneb.KZGProof         `json:"kzg_proofs"`
	Blobs       []deneb.Blob             `json:"blobs"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedBlockContents) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedBlockContentsJSON{
		SignedBlock: s.SignedBlock,
		KZGProofs:   s.KZGProofs,
		Blobs:       s.Blobs,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedBlockContents) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&signedBlockContentsJSON{}, input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw["signed_block"], &s.SignedBlock); err != nil {
		return errors.Wrap(err, "signed_block")
	}

	if err := json.Unmarshal(raw["kzg_proofs"], &s.KZGProofs); err != nil {
		return errors.Wrap(err, "kzg_proofs")
	}

	if err := json.Unmarshal(raw["blobs"], &s.Blobs); err != nil {
		return errors.Wrap(err, "blobs")
	}

	return nil
}
