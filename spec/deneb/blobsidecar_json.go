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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// blobSidecarJSON is the spec representation of the struct.
type blobSidecarJSON struct {
	BlockRoot       string `json:"block_root"`
	Index           string `json:"index"`
	Slot            string `json:"slot"`
	BlockParentRoot string `json:"block_parent_root"`
	ProposerIndex   string `json:"proposer_index"`
	Blob            string `json:"blob"`
	KzgCommitment   string `json:"kzg_commitment"`
	KzgProof        string `json:"kzg_proof"`
}

// MarshalJSON implements json.Marshaler.
func (b *BlobSidecar) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blobSidecarJSON{
		BlockRoot:       b.BlockRoot.String(),
		Index:           fmt.Sprintf("%d", b.Index),
		Slot:            fmt.Sprintf("%d", b.Slot),
		BlockParentRoot: b.BlockParentRoot.String(),
		ProposerIndex:   fmt.Sprintf("%d", b.ProposerIndex),
		Blob:            fmt.Sprintf("%#x", b.Blob),
		KzgCommitment:   b.KzgCommitment.String(),
		KzgProof:        b.KzgProof.String(),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *BlobSidecar) UnmarshalJSON(input []byte) error {
	var data blobSidecarJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	return b.unpack(&data)
}

func (b *BlobSidecar) unpack(data *blobSidecarJSON) error {
	if data.BlockRoot == "" {
		return errors.New("block root missing")
	}
	blockRoot, err := hex.DecodeString(strings.TrimPrefix(data.BlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block root")
	}
	if len(blockRoot) != phase0.RootLength {
		return errors.New("incorrect length for block root")
	}
	copy(b.BlockRoot[:], blockRoot)

	if data.Index == "" {
		return errors.New("index missing")
	}
	index, err := strconv.ParseUint(data.Index, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for index")
	}
	b.Index = BlobIndex(index)

	if data.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for slot")
	}
	b.Slot = phase0.Slot(slot)

	if data.BlockParentRoot == "" {
		return errors.New("block parent root missing")
	}
	blockParentRoot, err := hex.DecodeString(strings.TrimPrefix(data.BlockParentRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block parent root")
	}
	if len(blockParentRoot) != phase0.RootLength {
		return errors.New("incorrect length for block parent root")
	}
	copy(b.BlockParentRoot[:], blockParentRoot)

	if data.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}
	proposerIndex, err := strconv.ParseUint(data.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}
	b.ProposerIndex = phase0.ValidatorIndex(proposerIndex)

	if data.Blob == "" {
		return errors.New("blob missing")
	}
	blob, err := hex.DecodeString(strings.TrimPrefix(data.Blob, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for blob")
	}
	if len(blob) != BlobLength {
		return errors.New("incorrect length for blob")
	}
	copy(b.Blob[:], blob)

	if data.KzgCommitment == "" {
		return errors.New("kzg commitment missing")
	}
	KzgCommitment, err := hex.DecodeString(strings.TrimPrefix(data.KzgCommitment, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for kzg commitment")
	}
	if len(KzgCommitment) != KzgCommitmentLength {
		return errors.New("incorrect length for kzg commitment")
	}
	copy(b.KzgCommitment[:], KzgCommitment)

	if data.KzgProof == "" {
		return errors.New("kzg proof missing")
	}
	KzgProof, err := hex.DecodeString(strings.TrimPrefix(data.KzgProof, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for kzg proof")
	}
	if len(KzgProof) != KzgProofLength {
		return errors.New("incorrect length for kzg proof")
	}
	copy(b.KzgProof[:], KzgProof)

	return nil
}
