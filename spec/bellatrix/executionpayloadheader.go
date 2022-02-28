// Copyright © 2022 Attestant Limited.
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

package bellatrix

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
)

// ExecutionPayloadHeader represents an execution layer payload header.
type ExecutionPayloadHeader struct {
	ParentHash       phase0.Hash32
	FeeRecipient     ExecutionAddress
	StateRoot        [32]byte  `ssz-size:"32"`
	ReceiptsRoot     [32]byte  `ssz-size:"32"`
	LogsBloom        [256]byte // BYTES_PER_LOGS_BLOOM from spec
	PrevRandao       [32]byte  `ssz-size:"32"`
	BlockNumber      uint64
	GasLimit         uint64
	GasUsed          uint64
	Timestamp        time.Time
	ExtraData        []byte
	BaseFeePerGas    *uint256.Int
	BlockHash        phase0.Hash32
	TransactionsRoot phase0.Root
}

// executionPayloadHeaderJSON is the spec representation of the struct.
type executionPayloadHeaderJSON struct {
	ParentHash       string `json:"parent_hash"`
	FeeRecipient     string `json:"fee_recipient"`
	StateRoot        string `json:"state_root"`
	ReceiptsRoot     string `json:"receipts_root"`
	LogsBloom        string `json:"logs_bloom"`
	PrevRandao       string `json:"prev_randao"`
	BlockNumber      string `json:"block_number"`
	GasLimit         string `json:"gas_limit"`
	GasUsed          string `json:"gas_used"`
	Timestamp        string `json:"timestamp"`
	ExtraData        string `json:"extra_data"`
	BaseFeePerGas    string `json:"base_fee_per_gas"`
	BlockHash        string `json:"block_hash"`
	TransactionsRoot string `json:"transactions_root"`
}

// executionPayloadHeaderYAML is the spec representation of the struct.
type executionPayloadHeaderYAML struct {
	ParentHash       string `yaml:"parent_hash"`
	FeeRecipient     string `yaml:"fee_recipient"`
	StateRoot        string `yaml:"state_root"`
	ReceiptsRoot     string `yaml:"receipts_root"`
	LogsBloom        string `yaml:"logs_bloom"`
	PrevRandao       string `yaml:"prev_randao"`
	BlockNumber      uint64 `yaml:"block_number"`
	GasLimit         uint64 `yaml:"gas_limit"`
	GasUsed          uint64 `yaml:"gas_used"`
	Timestamp        uint64 `yaml:"timestamp"`
	ExtraData        string `yaml:"extra_data"`
	BaseFeePerGas    string `yaml:"base_fee_per_gas"`
	BlockHash        string `yaml:"block_hash"`
	TransactionsRoot string `yaml:"transactions_root"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayloadHeader) MarshalJSON() ([]byte, error) {
	return json.Marshal(&executionPayloadHeaderJSON{
		ParentHash:       fmt.Sprintf("%#x", e.ParentHash),
		FeeRecipient:     fmt.Sprintf("%#x", e.FeeRecipient),
		StateRoot:        fmt.Sprintf("%#x", e.StateRoot),
		ReceiptsRoot:     fmt.Sprintf("%#x", e.ReceiptsRoot),
		LogsBloom:        fmt.Sprintf("%#x", e.LogsBloom),
		PrevRandao:       fmt.Sprintf("%#x", e.PrevRandao),
		BlockNumber:      fmt.Sprintf("%d", e.BlockNumber),
		GasLimit:         fmt.Sprintf("%d", e.GasLimit),
		GasUsed:          fmt.Sprintf("%d", e.GasUsed),
		Timestamp:        fmt.Sprintf("%d", e.Timestamp.Unix()),
		ExtraData:        fmt.Sprintf("%#x", e.ExtraData),
		BaseFeePerGas:    fmt.Sprintf("%d", e.BaseFeePerGas),
		BlockHash:        fmt.Sprintf("%#x", e.BlockHash),
		TransactionsRoot: fmt.Sprintf("%#x", e.TransactionsRoot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionPayloadHeader) UnmarshalJSON(input []byte) error {
	var data executionPayloadHeaderJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	return e.unpack(&data)
}

// nolint:gocyclo
func (e *ExecutionPayloadHeader) unpack(data *executionPayloadHeaderJSON) error {
	if data.ParentHash == "" {
		return errors.New("parent hash missing")
	}
	parentHash, err := hex.DecodeString(strings.TrimPrefix(data.ParentHash, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for parent hash")
	}
	if len(parentHash) != phase0.Hash32Length {
		return errors.New("incorrect length for parent hash")
	}
	copy(e.ParentHash[:], parentHash)

	if data.FeeRecipient == "" {
		return errors.New("fee recipient missing")
	}
	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.FeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for fee recipient")
	}
	if len(feeRecipient) != FeeRecipientLength {
		return errors.New("incorrect length for fee recipient")
	}
	copy(e.FeeRecipient[:], feeRecipient)

	if data.StateRoot == "" {
		return errors.New("state root missing")
	}
	stateRoot, err := hex.DecodeString(strings.TrimPrefix(data.StateRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for state root")
	}
	if len(stateRoot) != 32 {
		return errors.New("incorrect length for state root")
	}
	copy(e.StateRoot[:], stateRoot)

	if data.ReceiptsRoot == "" {
		return errors.New("receipts root missing")
	}
	receiptsRoot, err := hex.DecodeString(strings.TrimPrefix(data.ReceiptsRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for receipts root")
	}
	if len(receiptsRoot) != 32 {
		return errors.New("incorrect length for receipts root")
	}
	copy(e.ReceiptsRoot[:], receiptsRoot)

	if data.LogsBloom == "" {
		return errors.New("logs bloom missing")
	}
	logsBloom, err := hex.DecodeString(strings.TrimPrefix(data.LogsBloom, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for logs bloom")
	}
	if len(logsBloom) != 256 {
		return errors.New("incorrect length for logs bloom")
	}
	copy(e.LogsBloom[:], logsBloom)

	if data.PrevRandao == "" {
		return errors.New("prev randao missing")
	}
	prevRandao, err := hex.DecodeString(strings.TrimPrefix(data.PrevRandao, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for prev randao")
	}
	if len(prevRandao) != 32 {
		return errors.New("incorrect length for prev randao")
	}
	copy(e.PrevRandao[:], prevRandao)

	if data.BlockNumber == "" {
		return errors.New("block number missing")
	}
	blockNumber, err := strconv.ParseUint(data.BlockNumber, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for block number")
	}
	e.BlockNumber = blockNumber

	if data.GasLimit == "" {
		return errors.New("gas limit missing")
	}
	gasLimit, err := strconv.ParseUint(data.GasLimit, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for gas limit")
	}
	e.GasLimit = gasLimit

	if data.GasUsed == "" {
		return errors.New("gas used missing")
	}
	gasUsed, err := strconv.ParseUint(data.GasUsed, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for gas used")
	}
	e.GasUsed = gasUsed

	if data.Timestamp == "" {
		return errors.New("timestamp missing")
	}
	tmp, err := strconv.ParseInt(data.Timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for timestamp")
	}
	e.Timestamp = time.Unix(tmp, 0)

	if data.ExtraData == "" {
		return errors.New("extra data missing")
	}
	extraData, err := hex.DecodeString(strings.TrimPrefix(data.ExtraData, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for extra data")
	}
	if len(extraData) > 32 {
		return errors.New("incorrect length for extra data")
	}
	e.ExtraData = extraData

	if data.BaseFeePerGas == "" {
		return errors.New("base fee per gas missing")
	}
	baseFeePerGas := new(big.Int)
	baseFeePerGas, ok := baseFeePerGas.SetString(data.BaseFeePerGas, 10)
	if !ok {
		return errors.New("invalid value for base fee per gas")
	}
	var overflow bool
	e.BaseFeePerGas, overflow = uint256.FromBig(baseFeePerGas)
	if overflow {
		return errors.New("overflow for base fee per gas")
	}

	if data.BlockHash == "" {
		return errors.New("block hash missing")
	}
	blockHash, err := hex.DecodeString(strings.TrimPrefix(data.BlockHash, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for block hash")
	}
	if len(blockHash) != phase0.Hash32Length {
		return errors.New("incorrect length for block hash")
	}
	copy(e.BlockHash[:], blockHash)

	if data.TransactionsRoot == "" {
		return errors.New("transactions root missing")
	}
	transactionsRoot, err := hex.DecodeString(strings.TrimPrefix(data.TransactionsRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for transactions root")
	}
	if len(transactionsRoot) != phase0.Hash32Length {
		return errors.New("incorrect length for transactions root")
	}
	copy(e.TransactionsRoot[:], transactionsRoot)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (e *ExecutionPayloadHeader) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&executionPayloadHeaderYAML{
		ParentHash:       fmt.Sprintf("%#x", e.ParentHash),
		FeeRecipient:     fmt.Sprintf("%#x", e.FeeRecipient),
		StateRoot:        fmt.Sprintf("%#x", e.StateRoot),
		ReceiptsRoot:     fmt.Sprintf("%#x", e.ReceiptsRoot),
		LogsBloom:        fmt.Sprintf("%#x", e.LogsBloom),
		PrevRandao:       fmt.Sprintf("%#x", e.PrevRandao),
		BlockNumber:      e.BlockNumber,
		GasLimit:         e.GasLimit,
		GasUsed:          e.GasUsed,
		Timestamp:        uint64(e.Timestamp.Unix()),
		ExtraData:        fmt.Sprintf("%#x", e.ExtraData),
		BaseFeePerGas:    fmt.Sprintf("%s", e.BaseFeePerGas),
		BlockHash:        fmt.Sprintf("%#x", e.BlockHash),
		TransactionsRoot: fmt.Sprintf("%#x", e.TransactionsRoot),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}
	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (e *ExecutionPayloadHeader) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data executionPayloadHeaderJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}
	return e.unpack(&data)
}

// String returns a string version of the structure.
func (e *ExecutionPayloadHeader) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
