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
	"fmt"

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// blobSidecarJSON is the spec representation of the struct.
type blobSidecarJSON struct {
	Index                       string                          `json:"index"`
	Blob                        Blob                            `json:"blob"`
	KZGCommitment               KZGCommitment                   `json:"kzg_commitment"`
	KZGProof                    KZGProof                        `json:"kzg_proof"`
	SignedBlockHeader           *phase0.SignedBeaconBlockHeader `json:"signed_block_header"`
	KZGCommitmentInclusionProof KZGCommitmentInclusionProof     `json:"kzg_commitment_inclusion_proof"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlobSidecar) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blobSidecarJSON{
		Index:                       fmt.Sprintf("%d", b.Index),
		Blob:                        b.Blob,
		KZGCommitment:               b.KZGCommitment,
		KZGProof:                    b.KZGProof,
		SignedBlockHeader:           b.SignedBlockHeader,
		KZGCommitmentInclusionProof: b.KZGCommitmentInclusionProof,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlobSidecar) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&blobSidecarJSON{}, input)
	if err != nil {
		return err
	}

	if err := b.Index.UnmarshalJSON(raw["index"]); err != nil {
		return errors.Wrap(err, "index")
	}

	if err := b.Blob.UnmarshalJSON(raw["blob"]); err != nil {
		return errors.Wrap(err, "blob")
	}

	if err := b.KZGCommitment.UnmarshalJSON(raw["kzg_commitment"]); err != nil {
		return errors.Wrap(err, "kzg_commitment")
	}

	if err := b.KZGProof.UnmarshalJSON(raw["kzg_proof"]); err != nil {
		return errors.Wrap(err, "kzg_proof")
	}

	b.SignedBlockHeader = &phase0.SignedBeaconBlockHeader{}
	if err := b.SignedBlockHeader.UnmarshalJSON(raw["signed_block_header"]); err != nil {
		return errors.Wrap(err, "signed_block_header")
	}

	if err := b.KZGCommitmentInclusionProof.UnmarshalJSON(raw["kzg_commitment_inclusion_proof"]); err != nil {
		return errors.Wrap(err, "kzg_commitment_inclusion_proof")
	}

	return nil
}
