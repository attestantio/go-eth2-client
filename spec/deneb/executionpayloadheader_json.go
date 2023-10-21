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
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
)

// executionPayloadHeaderJSON is the spec representation of the struct.
type executionPayloadHeaderJSON struct {
	ParentHash       phase0.Hash32              `json:"parent_hash"`
	FeeRecipient     bellatrix.ExecutionAddress `json:"fee_recipient"`
	StateRoot        phase0.Root                `json:"state_root"`
	ReceiptsRoot     phase0.Root                `json:"receipts_root"`
	LogsBloom        string                     `json:"logs_bloom"`
	PrevRandao       string                     `json:"prev_randao"`
	BlockNumber      string                     `json:"block_number"`
	GasLimit         string                     `json:"gas_limit"`
	GasUsed          string                     `json:"gas_used"`
	Timestamp        string                     `json:"timestamp"`
	ExtraData        string                     `json:"extra_data"`
	BaseFeePerGas    string                     `json:"base_fee_per_gas"`
	BlockHash        phase0.Hash32              `json:"block_hash"`
	TransactionsRoot phase0.Root                `json:"transactions_root"`
	WithdrawalsRoot  phase0.Root                `json:"withdrawals_root"`
	BlobGasUsed      string                     `json:"blob_gas_used"`
	ExcessBlobGas    string                     `json:"excess_blob_gas"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayloadHeader) MarshalJSON() ([]byte, error) {
	extraData := "0x"
	if len(e.ExtraData) > 0 {
		extraData = fmt.Sprintf("%#x", e.ExtraData)
	}

	return json.Marshal(&executionPayloadHeaderJSON{
		ParentHash:       e.ParentHash,
		FeeRecipient:     e.FeeRecipient,
		StateRoot:        e.StateRoot,
		ReceiptsRoot:     e.ReceiptsRoot,
		LogsBloom:        fmt.Sprintf("%#x", e.LogsBloom),
		PrevRandao:       fmt.Sprintf("%#x", e.PrevRandao),
		BlockNumber:      strconv.FormatUint(e.BlockNumber, 10),
		GasLimit:         strconv.FormatUint(e.GasLimit, 10),
		GasUsed:          strconv.FormatUint(e.GasUsed, 10),
		Timestamp:        strconv.FormatUint(e.Timestamp, 10),
		ExtraData:        extraData,
		BaseFeePerGas:    e.BaseFeePerGas.Dec(),
		BlockHash:        e.BlockHash,
		TransactionsRoot: e.TransactionsRoot,
		WithdrawalsRoot:  e.WithdrawalsRoot,
		BlobGasUsed:      strconv.FormatUint(e.BlobGasUsed, 10),
		ExcessBlobGas:    strconv.FormatUint(e.ExcessBlobGas, 10),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
//
//nolint:gocyclo
func (e *ExecutionPayloadHeader) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&executionPayloadHeaderJSON{}, input)
	if err != nil {
		return err
	}

	if err := e.ParentHash.UnmarshalJSON(raw["parent_hash"]); err != nil {
		return errors.Wrap(err, "parent_hash")
	}

	if err := e.FeeRecipient.UnmarshalJSON(raw["fee_recipient"]); err != nil {
		return errors.Wrap(err, "fee_recipient")
	}

	if err := e.StateRoot.UnmarshalJSON(raw["state_root"]); err != nil {
		return errors.Wrap(err, "state_root")
	}

	if err := e.ReceiptsRoot.UnmarshalJSON(raw["receipts_root"]); err != nil {
		return errors.Wrap(err, "receipts_root")
	}

	logsBloom := raw["logs_bloom"]
	if !bytes.HasPrefix(logsBloom, []byte{'"', '0', 'x'}) {
		return errors.New("logs_bloom: invalid prefix")
	}
	if !bytes.HasSuffix(logsBloom, []byte{'"'}) {
		return errors.New("logs_bloom: invalid suffix")
	}
	if len(logsBloom) != 1+2+256*2+1 {
		return errors.New("logs_bloom: incorrect length")
	}
	length, err := hex.Decode(e.LogsBloom[:], logsBloom[3:3+256*2])
	if err != nil {
		return errors.Wrap(err, "logs_bloom")
	}
	if length != 256 {
		return errors.New("logs_bloom: incorrect length")
	}

	prevRandao := raw["prev_randao"]
	if !bytes.HasPrefix(prevRandao, []byte{'"', '0', 'x'}) {
		return errors.New("prev_randao: invalid prefix")
	}
	if !bytes.HasSuffix(prevRandao, []byte{'"'}) {
		return errors.New("prev_randao: invalid suffix")
	}
	if len(prevRandao) != 1+2+32*2+1 {
		return errors.New("prev_randao: incorrect length")
	}
	length, err = hex.Decode(e.PrevRandao[:], prevRandao[3:3+32*2])
	if err != nil {
		return errors.Wrap(err, "prev_randao")
	}
	if length != 32 {
		return errors.New("prev_randao: incorrect length")
	}

	tmpUint, err := strconv.ParseUint(string(bytes.Trim(raw["block_number"], `"`)), 10, 64)
	if err != nil {
		return errors.Wrap(err, "block_number")
	}
	e.BlockNumber = tmpUint

	tmpUint, err = strconv.ParseUint(string(bytes.Trim(raw["gas_limit"], `"`)), 10, 64)
	if err != nil {
		return errors.Wrap(err, "gas_limit")
	}
	e.GasLimit = tmpUint

	tmpUint, err = strconv.ParseUint(string(bytes.Trim(raw["gas_used"], `"`)), 10, 64)
	if err != nil {
		return errors.Wrap(err, "gas_used")
	}
	e.GasUsed = tmpUint

	tmpUint, err = strconv.ParseUint(string(bytes.Trim(raw["timestamp"], `"`)), 10, 64)
	if err != nil {
		return errors.Wrap(err, "timestamp")
	}
	e.Timestamp = tmpUint

	var tmpBytes []byte
	switch {
	case bytes.Equal(raw["extra_data"], []byte{'0', 'x'}), bytes.Equal(raw["extra_data"], []byte{'0'}):
		// Empty.
	default:
		tmpBytes = bytes.TrimPrefix(bytes.Trim(raw["extra_data"], `"`), []byte{'0', 'x'})
		if len(tmpBytes)%2 == 1 {
			tmpBytes = []byte(fmt.Sprintf("0%s", string(tmpBytes)))
		}
		tmp, err := hex.DecodeString(string(tmpBytes))
		if err != nil {
			return errors.Wrap(err, "extra_data")
		}
		if len(tmp) > 32 {
			return errors.New("extra_data: incorrect length")
		}
		e.ExtraData = tmp
	}

	tmpBytes = bytes.Trim(raw["base_fee_per_gas"], `"`)
	if bytes.HasPrefix(tmpBytes, []byte{'0', 'x'}) {
		e.BaseFeePerGas, err = uint256.FromHex(string(tmpBytes))
	} else {
		e.BaseFeePerGas, err = uint256.FromDecimal(string(tmpBytes))
	}
	if err != nil {
		return errors.Wrap(err, "base_fee_per_gas")
	}

	if err := e.BlockHash.UnmarshalJSON(raw["block_hash"]); err != nil {
		return errors.Wrap(err, "block_hash")
	}

	if err := e.TransactionsRoot.UnmarshalJSON(raw["transactions_root"]); err != nil {
		return errors.Wrap(err, "transactions_root")
	}

	if err := e.WithdrawalsRoot.UnmarshalJSON(raw["withdrawals_root"]); err != nil {
		return errors.Wrap(err, "withdrawals_root")
	}

	tmpUint, err = strconv.ParseUint(string(bytes.Trim(raw["blob_gas_used"], `"`)), 10, 64)
	if err != nil {
		return errors.Wrap(err, "blob_gas_used")
	}
	e.BlobGasUsed = tmpUint

	tmpUint, err = strconv.ParseUint(string(bytes.Trim(raw["excess_blob_gas"], `"`)), 10, 64)
	if err != nil {
		return errors.Wrap(err, "excess_blob_gas")
	}
	e.ExcessBlobGas = tmpUint

	return nil
}
