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

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
)

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
	WithdrawalsRoot  string `json:"withdrawals_root"`
	DataGasUsed      string `json:"data_gas_used"`
	ExcessDataGas    string `json:"excess_data_gas"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayloadHeader) MarshalJSON() ([]byte, error) {
	extraData := "0x"
	if len(e.ExtraData) > 0 {
		extraData = fmt.Sprintf("%#x", e.ExtraData)
	}

	return json.Marshal(&executionPayloadHeaderJSON{
		ParentHash:       e.ParentHash.String(),
		FeeRecipient:     e.FeeRecipient.String(),
		StateRoot:        e.StateRoot.String(),
		ReceiptsRoot:     e.ReceiptsRoot.String(),
		LogsBloom:        fmt.Sprintf("%#x", e.LogsBloom),
		PrevRandao:       fmt.Sprintf("%#x", e.PrevRandao),
		BlockNumber:      fmt.Sprintf("%d", e.BlockNumber),
		GasLimit:         fmt.Sprintf("%d", e.GasLimit),
		GasUsed:          fmt.Sprintf("%d", e.GasUsed),
		Timestamp:        fmt.Sprintf("%d", e.Timestamp),
		ExtraData:        extraData,
		BaseFeePerGas:    e.BaseFeePerGas.Dec(),
		BlockHash:        fmt.Sprintf("%#x", e.BlockHash),
		TransactionsRoot: e.TransactionsRoot.String(),
		WithdrawalsRoot:  e.WithdrawalsRoot.String(),
		DataGasUsed:      fmt.Sprintf("%d", e.DataGasUsed),
		ExcessDataGas:    fmt.Sprintf("%d", e.ExcessDataGas),
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
	if strings.HasPrefix(data.BaseFeePerGas, "0x") {
		e.BaseFeePerGas, err = uint256.FromHex(data.BaseFeePerGas)
	} else {
		e.BaseFeePerGas, err = uint256.FromDecimal(data.BaseFeePerGas)
	}
	if err != nil {
		return errors.Wrap(err, "invalid value for base fee per gas")
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

	if data.WithdrawalsRoot == "" {
		return errors.New("withdrawals root missing")
	}
	withdrawalsRoot, err := hex.DecodeString(strings.TrimPrefix(data.WithdrawalsRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for withdrawals root")
	}
	if len(withdrawalsRoot) != phase0.Hash32Length {
		return errors.New("incorrect length for withdrawals root")
	}
	copy(e.WithdrawalsRoot[:], withdrawalsRoot)

	if data.ExcessDataGas == "" {
		return errors.New("excess data gas missing")
	}
	if data.DataGasUsed == "" {
		return errors.New("data gas used missing")
	}
	dataGasUsed, err := strconv.ParseUint(data.DataGasUsed, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for data gas used")
	}
	e.DataGasUsed = dataGasUsed

	if data.ExcessDataGas == "" {
		return errors.New("excess data gas used missing")
	}
	excessDataGas, err := strconv.ParseUint(data.ExcessDataGas, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for excess data gas")
	}
	e.ExcessDataGas = excessDataGas

	return nil
}
