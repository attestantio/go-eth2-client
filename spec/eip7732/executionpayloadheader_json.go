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

package eip7732

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// executionPayloadHeaderJSON is the spec representation of the struct.
type executionPayloadHeaderJSON struct {
	ParentBlockHash        string `json:"parent_block_hash"`
	ParentBlockRoot        string `json:"parent_block_root"`
	BlockHash              string `json:"block_hash"`
	GasLimit               string `json:"gas_limit"`
	BuilderIndex           string `json:"builder_index"`
	Slot                   string `json:"slot"`
	Value                  string `json:"value"`
	BlobKZGCommitmentsRoot string `json:"blob_kzg_commitments_root"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayloadHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(&executionPayloadHeaderJSON{
		ParentBlockHash:        fmt.Sprintf("%#x", e.ParentBlockHash),
		ParentBlockRoot:        fmt.Sprintf("%#x", e.ParentBlockRoot),
		BlockHash:              fmt.Sprintf("%#x", e.BlockHash),
		GasLimit:               fmt.Sprintf("%d", e.GasLimit),
		BuilderIndex:           fmt.Sprintf("%d", e.BuilderIndex),
		Slot:                   fmt.Sprintf("%d", e.Slot),
		Value:                  fmt.Sprintf("%d", e.Value),
		BlobKZGCommitmentsRoot: fmt.Sprintf("%#x", e.BlobKZGCommitmentsRoot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionPayloadHeader) UnmarshalJSON(input []byte) error {
	var data executionPayloadHeaderJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	// Parent block hash
	if data.ParentBlockHash == "" {
		return errors.New("parent block hash missing")
	}
	parentBlockHash, err := hex.DecodeString(strings.TrimPrefix(data.ParentBlockHash, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid parent block hash")
	}
	copy(e.ParentBlockHash[:], parentBlockHash)

	// Parent block root
	if data.ParentBlockRoot == "" {
		return errors.New("parent block root missing")
	}
	parentBlockRoot, err := hex.DecodeString(strings.TrimPrefix(data.ParentBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid parent block root")
	}
	copy(e.ParentBlockRoot[:], parentBlockRoot)

	// Block hash
	if data.BlockHash == "" {
		return errors.New("block hash missing")
	}
	blockHash, err := hex.DecodeString(strings.TrimPrefix(data.BlockHash, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid block hash")
	}
	copy(e.BlockHash[:], blockHash)

	// Gas limit
	if data.GasLimit == "" {
		return errors.New("gas limit missing")
	}
	gasLimit, err := strconv.ParseUint(data.GasLimit, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid gas limit")
	}
	e.GasLimit = gasLimit

	// Builder index
	if data.BuilderIndex == "" {
		return errors.New("builder index missing")
	}
	builderIndex, err := strconv.ParseUint(data.BuilderIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid builder index")
	}
	e.BuilderIndex = phase0.ValidatorIndex(builderIndex)

	// Slot
	if data.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid slot")
	}
	e.Slot = phase0.Slot(slot)

	// Value
	if data.Value == "" {
		return errors.New("value missing")
	}
	value, err := strconv.ParseUint(data.Value, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value")
	}
	e.Value = value

	// Blob KZG commitments root
	if data.BlobKZGCommitmentsRoot == "" {
		return errors.New("blob KZG commitments root missing")
	}
	blobKZGCommitmentsRoot, err := hex.DecodeString(strings.TrimPrefix(data.BlobKZGCommitmentsRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid blob KZG commitments root")
	}
	copy(e.BlobKZGCommitmentsRoot[:], blobKZGCommitmentsRoot)

	return nil
}
