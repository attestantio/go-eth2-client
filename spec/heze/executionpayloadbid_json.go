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

package heze

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// executionPayloadBidJSON is the spec representation of the struct.
type executionPayloadBidJSON struct {
	ParentBlockHash    string   `json:"parent_block_hash"`
	ParentBlockRoot    string   `json:"parent_block_root"`
	BlockHash          string   `json:"block_hash"`
	PrevRandao         string   `json:"prev_randao"`
	FeeRecipient       string   `json:"fee_recipient"`
	GasLimit           string   `json:"gas_limit"`
	BuilderIndex       string   `json:"builder_index"`
	Slot               string   `json:"slot"`
	Value              string   `json:"value"`
	ExecutionPayment   string   `json:"execution_payment"`
	BlobKZGCommitments []string `json:"blob_kzg_commitments"`
	InclusionListBits  string   `json:"inclusion_list_bits"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayloadBid) MarshalJSON() ([]byte, error) {
	blobKZGCommitments := make([]string, len(e.BlobKZGCommitments))
	for i := range e.BlobKZGCommitments {
		blobKZGCommitments[i] = fmt.Sprintf("%#x", e.BlobKZGCommitments[i])
	}

	return json.Marshal(&executionPayloadBidJSON{
		ParentBlockHash:    fmt.Sprintf("%#x", e.ParentBlockHash),
		ParentBlockRoot:    fmt.Sprintf("%#x", e.ParentBlockRoot),
		BlockHash:          fmt.Sprintf("%#x", e.BlockHash),
		PrevRandao:         fmt.Sprintf("%#x", e.PrevRandao),
		FeeRecipient:       fmt.Sprintf("%#x", e.FeeRecipient),
		GasLimit:           fmt.Sprintf("%d", e.GasLimit),
		BuilderIndex:       fmt.Sprintf("%d", e.BuilderIndex),
		Slot:               fmt.Sprintf("%d", e.Slot),
		Value:              fmt.Sprintf("%d", e.Value),
		ExecutionPayment:   fmt.Sprintf("%d", e.ExecutionPayment),
		BlobKZGCommitments: blobKZGCommitments,
		InclusionListBits:  fmt.Sprintf("%#x", e.InclusionListBits),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionPayloadBid) UnmarshalJSON(input []byte) error {
	var data executionPayloadBidJSON
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

	// Prev randao
	if data.PrevRandao == "" {
		return errors.New("prev randao missing")
	}
	prevRandao, err := hex.DecodeString(strings.TrimPrefix(data.PrevRandao, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid prev randao")
	}
	copy(e.PrevRandao[:], prevRandao)

	// Fee recipient
	if data.FeeRecipient == "" {
		return errors.New("fee recipient missing")
	}
	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.FeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid fee recipient")
	}
	copy(e.FeeRecipient[:], feeRecipient)

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
	e.BuilderIndex = gloas.BuilderIndex(builderIndex)

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
	e.Value = phase0.Gwei(value)

	// Execution payment
	if data.ExecutionPayment == "" {
		return errors.New("execution payment missing")
	}
	executionPayment, err := strconv.ParseUint(data.ExecutionPayment, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid execution payment")
	}
	e.ExecutionPayment = phase0.Gwei(executionPayment)

	// Blob KZG commitments
	if data.BlobKZGCommitments == nil {
		data.BlobKZGCommitments = []string{}
	}
	e.BlobKZGCommitments = make([]deneb.KZGCommitment, len(data.BlobKZGCommitments))
	for i, commitment := range data.BlobKZGCommitments {
		commitmentBytes, err := hex.DecodeString(strings.TrimPrefix(commitment, "0x"))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid blob KZG commitment %d", i))
		}
		if len(commitmentBytes) != deneb.KZGCommitmentLength {
			return fmt.Errorf("blob KZG commitment %d has incorrect length %d", i, len(commitmentBytes))
		}
		copy(e.BlobKZGCommitments[i][:], commitmentBytes)
	}

	// Inclusion list bits
	if data.InclusionListBits == "" {
		return errors.New("inclusion list bits missing")
	}
	inclusionListBits, err := hex.DecodeString(strings.TrimPrefix(data.InclusionListBits, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid inclusion list bits")
	}
	copy(e.InclusionListBits, inclusionListBits)

	return nil
}
