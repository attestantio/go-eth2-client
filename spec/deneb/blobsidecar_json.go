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
	BlockRoot       phase0.Root   `json:"block_root"`
	Index           string        `json:"index"`
	Slot            string        `json:"slot"`
	BlockParentRoot phase0.Root   `json:"block_parent_root"`
	ProposerIndex   string        `json:"proposer_index"`
	Blob            Blob          `json:"blob"`
	KzgCommitment   KzgCommitment `json:"kzg_commitment"`
	KzgProof        KzgProof      `json:"kzg_proof"`
}

func (b *BlobSidecar) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blobSidecarJSON{
		BlockRoot:       b.BlockRoot,
		Index:           fmt.Sprintf("%d", b.Index),
		Slot:            fmt.Sprintf("%d", b.Slot),
		BlockParentRoot: b.BlockParentRoot,
		ProposerIndex:   fmt.Sprintf("%d", b.ProposerIndex),
		Blob:            b.Blob,
		KzgCommitment:   b.KzgCommitment,
		KzgProof:        b.KzgProof,
	})
}

func (b *BlobSidecar) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&blobSidecarJSON{}, input)
	if err != nil {
		return err
	}

	if err := b.BlockRoot.UnmarshalJSON(raw["block_root"]); err != nil {
		return errors.Wrap(err, "block_root")
	}

	if err := b.Index.UnmarshalJSON(raw["index"]); err != nil {
		return errors.Wrap(err, "index")
	}

	if err := b.Slot.UnmarshalJSON(raw["slot"]); err != nil {
		return errors.Wrap(err, "slot")
	}

	if err := b.BlockParentRoot.UnmarshalJSON(raw["block_parent_root"]); err != nil {
		return errors.Wrap(err, "block_parent_root")
	}

	if err := b.ProposerIndex.UnmarshalJSON(raw["proposer_index"]); err != nil {
		return errors.Wrap(err, "proposer_index")
	}

	if err := b.Blob.UnmarshalJSON(raw["blob"]); err != nil {
		return errors.Wrap(err, "blob")
	}

	if err := b.KzgCommitment.UnmarshalJSON(raw["kzg_commitment"]); err != nil {
		return errors.Wrap(err, "kzg_commitment")
	}

	if err := b.KzgProof.UnmarshalJSON(raw["kzg_proof"]); err != nil {
		return errors.Wrap(err, "kzg_proof")
	}

	return nil
}
