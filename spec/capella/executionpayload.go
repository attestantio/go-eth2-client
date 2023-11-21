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

package capella

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// ExecutionPayload represents an execution layer payload.
type ExecutionPayload struct {
	ParentHash    phase0.Hash32              `ssz-size:"32"`
	FeeRecipient  bellatrix.ExecutionAddress `ssz-size:"20"`
	StateRoot     [32]byte                   `ssz-size:"32"`
	ReceiptsRoot  [32]byte                   `ssz-size:"32"`
	LogsBloom     [256]byte                  `ssz-size:"256"`
	PrevRandao    [32]byte                   `ssz-size:"32"`
	BlockNumber   uint64
	GasLimit      uint64
	GasUsed       uint64
	Timestamp     uint64
	ExtraData     []byte                  `ssz-max:"32"`
	BaseFeePerGas [32]byte                `ssz-size:"32"`
	BlockHash     phase0.Hash32           `ssz-size:"32"`
	Transactions  []bellatrix.Transaction `ssz-max:"1048576,1073741824" ssz-size:"?,?"`
	Withdrawals   []*Withdrawal           `ssz-max:"16"`
}

// executionPayloadJSON is the spec representation of the struct.
type executionPayloadJSON struct {
	ParentHash    string        `json:"parent_hash"`
	FeeRecipient  string        `json:"fee_recipient"`
	StateRoot     string        `json:"state_root"`
	ReceiptsRoot  string        `json:"receipts_root"`
	LogsBloom     string        `json:"logs_bloom"`
	PrevRandao    string        `json:"prev_randao"`
	BlockNumber   string        `json:"block_number"`
	GasLimit      string        `json:"gas_limit"`
	GasUsed       string        `json:"gas_used"`
	Timestamp     string        `json:"timestamp"`
	ExtraData     string        `json:"extra_data"`
	BaseFeePerGas string        `json:"base_fee_per_gas"`
	BlockHash     string        `json:"block_hash"`
	Transactions  []string      `json:"transactions"`
	Withdrawals   []*Withdrawal `json:"withdrawals"`
}

// executionPayloadYAML is the spec representation of the struct.
type executionPayloadYAML struct {
	ParentHash    string        `yaml:"parent_hash"`
	FeeRecipient  string        `yaml:"fee_recipient"`
	StateRoot     string        `yaml:"state_root"`
	ReceiptsRoot  string        `yaml:"receipts_root"`
	LogsBloom     string        `yaml:"logs_bloom"`
	PrevRandao    string        `yaml:"prev_randao"`
	BlockNumber   uint64        `yaml:"block_number"`
	GasLimit      uint64        `yaml:"gas_limit"`
	GasUsed       uint64        `yaml:"gas_used"`
	Timestamp     uint64        `yaml:"timestamp"`
	ExtraData     string        `yaml:"extra_data"`
	BaseFeePerGas string        `yaml:"base_fee_per_gas"`
	BlockHash     string        `yaml:"block_hash"`
	Transactions  []string      `yaml:"transactions"`
	Withdrawals   []*Withdrawal `yaml:"withdrawals"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayload) MarshalJSON() ([]byte, error) {
	transactions := make([]string, len(e.Transactions))
	for i := range e.Transactions {
		transactions[i] = fmt.Sprintf("%#x", e.Transactions[i])
	}

	extraData := "0x"
	if len(e.ExtraData) > 0 {
		extraData = fmt.Sprintf("%#x", e.ExtraData)
	}

	// base fee per gas is stored little-endian but we need it
	// big-endian for big.Int.
	var baseFeePerGasBEBytes [32]byte
	for i := 0; i < 32; i++ {
		baseFeePerGasBEBytes[i] = e.BaseFeePerGas[32-1-i]
	}
	baseFeePerGas := new(big.Int).SetBytes(baseFeePerGasBEBytes[:])

	return json.Marshal(&executionPayloadJSON{
		ParentHash:    fmt.Sprintf("%#x", e.ParentHash),
		FeeRecipient:  e.FeeRecipient.String(),
		StateRoot:     fmt.Sprintf("%#x", e.StateRoot),
		ReceiptsRoot:  fmt.Sprintf("%#x", e.ReceiptsRoot),
		LogsBloom:     fmt.Sprintf("%#x", e.LogsBloom),
		PrevRandao:    fmt.Sprintf("%#x", e.PrevRandao),
		BlockNumber:   strconv.FormatUint(e.BlockNumber, 10),
		GasLimit:      strconv.FormatUint(e.GasLimit, 10),
		GasUsed:       strconv.FormatUint(e.GasUsed, 10),
		Timestamp:     strconv.FormatUint(e.Timestamp, 10),
		ExtraData:     extraData,
		BaseFeePerGas: baseFeePerGas.String(),
		BlockHash:     fmt.Sprintf("%#x", e.BlockHash),
		Transactions:  transactions,
		Withdrawals:   e.Withdrawals,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionPayload) UnmarshalJSON(input []byte) error {
	var data executionPayloadJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return e.unpack(&data)
}

//nolint:gocyclo
func (e *ExecutionPayload) unpack(data *executionPayloadJSON) error {
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
	if len(feeRecipient) != bellatrix.FeeRecipientLength {
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
	e.Timestamp, err = strconv.ParseUint(data.Timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for timestamp")
	}

	if data.ExtraData == "" {
		return errors.New("extra data missing")
	}
	switch {
	case data.ExtraData == "0x", data.ExtraData == "0":
		e.ExtraData = []byte{}
	default:
		data.ExtraData = strings.TrimPrefix(data.ExtraData, "0x")
		if len(data.ExtraData)%2 == 1 {
			data.ExtraData = fmt.Sprintf("0%s", data.ExtraData)
		}
		extraData, err := hex.DecodeString(data.ExtraData)
		if err != nil {
			return errors.Wrap(err, "invalid value for extra data")
		}
		if len(extraData) > 32 {
			return errors.New("incorrect length for extra data")
		}
		e.ExtraData = extraData
	}

	if data.BaseFeePerGas == "" {
		return errors.New("base fee per gas missing")
	}
	baseFeePerGas := new(big.Int)
	var ok bool
	input := data.BaseFeePerGas
	if strings.HasPrefix(data.BaseFeePerGas, "0x") {
		input = strings.TrimPrefix(input, "0x")
		if len(data.BaseFeePerGas)%2 == 1 {
			input = fmt.Sprintf("0%s", input)
		}
		baseFeePerGas, ok = baseFeePerGas.SetString(input, 16)
	} else {
		baseFeePerGas, ok = baseFeePerGas.SetString(input, 10)
	}
	if !ok {
		return errors.New("invalid value for base fee per gas")
	}
	if baseFeePerGas.Cmp(maxBaseFeePerGas) > 0 {
		return errors.New("overflow for base fee per gas")
	}
	// We need to store internally as little-endian, but big.Int uses
	// big-endian so do it manually.
	baseFeePerGasBEBytes := baseFeePerGas.Bytes()
	var baseFeePerGasLEBytes [32]byte
	baseFeeLen := len(baseFeePerGasBEBytes)
	for i := 0; i < baseFeeLen; i++ {
		baseFeePerGasLEBytes[i] = baseFeePerGasBEBytes[baseFeeLen-1-i]
	}
	copy(e.BaseFeePerGas[:], baseFeePerGasLEBytes[:])

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

	if data.Transactions == nil {
		return errors.New("transactions missing")
	}
	transactions := make([]bellatrix.Transaction, len(data.Transactions))
	for i := range data.Transactions {
		if data.Transactions[i] == "" {
			return errors.New("transaction missing")
		}
		tmp, err := hex.DecodeString(strings.TrimPrefix(data.Transactions[i], "0x"))
		if err != nil {
			return errors.Wrap(err, "invalid value for transaction")
		}
		transactions[i] = bellatrix.Transaction(tmp)
	}
	e.Transactions = transactions

	if data.Withdrawals == nil {
		return errors.New("withdrawals missing")
	}
	e.Withdrawals = data.Withdrawals

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (e *ExecutionPayload) MarshalYAML() ([]byte, error) {
	transactions := make([]string, len(e.Transactions))
	for i := range e.Transactions {
		transactions[i] = fmt.Sprintf("%#x", e.Transactions[i])
	}

	extraData := "0x"
	if len(e.ExtraData) > 0 {
		extraData = fmt.Sprintf("%#x", e.ExtraData)
	}

	// base fee per gas is stored little-endian but we need it
	// big-endian for big.Int.
	var baseFeePerGasBEBytes [32]byte
	for i := 0; i < 32; i++ {
		baseFeePerGasBEBytes[i] = e.BaseFeePerGas[32-1-i]
	}
	baseFeePerGas := new(big.Int).SetBytes(baseFeePerGasBEBytes[:])

	yamlBytes, err := yaml.MarshalWithOptions(&executionPayloadYAML{
		ParentHash:    fmt.Sprintf("%#x", e.ParentHash),
		FeeRecipient:  e.FeeRecipient.String(),
		StateRoot:     fmt.Sprintf("%#x", e.StateRoot),
		ReceiptsRoot:  fmt.Sprintf("%#x", e.ReceiptsRoot),
		LogsBloom:     fmt.Sprintf("%#x", e.LogsBloom),
		PrevRandao:    fmt.Sprintf("%#x", e.PrevRandao),
		BlockNumber:   e.BlockNumber,
		GasLimit:      e.GasLimit,
		GasUsed:       e.GasUsed,
		Timestamp:     e.Timestamp,
		ExtraData:     extraData,
		BaseFeePerGas: baseFeePerGas.String(),
		BlockHash:     fmt.Sprintf("%#x", e.BlockHash),
		Transactions:  transactions,
		Withdrawals:   e.Withdrawals,
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (e *ExecutionPayload) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data executionPayloadJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return e.unpack(&data)
}

// String returns a string version of the structure.
func (e *ExecutionPayload) String() string {
	data, err := yaml.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
