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

package v1

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/electra"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// PayloadAttributesEvent represents the data of a payload_attributes event.
type PayloadAttributesEvent struct {
	// Version is the fork version of the beacon chain.
	Version spec.DataVersion
	// Data is the data of the event.
	Data *PayloadAttributesData
}

// PayloadAttributesData represents the data of a payload_attributes event.
type PayloadAttributesData struct {
	// ProposerIndex is the index of the proposer.
	ProposerIndex phase0.ValidatorIndex
	// ProposalSlot is the slot of the proposal.
	ProposalSlot phase0.Slot
	// ParentBlockNumber is the number of the parent block.
	ParentBlockNumber uint64
	// ParentBlockRoot is the root of the parent block.
	ParentBlockRoot phase0.Root
	// ParentBlockHash is the hash of the parent block.
	ParentBlockHash phase0.Hash32
	// V1 is the v1 payload attributes.
	V1 *PayloadAttributesV1
	// V2 is the v2 payload attributes.
	V2 *PayloadAttributesV2
	// V3 is the v3 payload attributes.
	V3 *PayloadAttributesV3
	// V4 is the v4 payload attributes.
	V4 *PayloadAttributesV4
}

// PayloadAttributesV1 represents the payload attributes.
type PayloadAttributesV1 struct {
	// Timestamp is the timestamp of the payload.
	Timestamp uint64
	// PrevRandao is the previous randao.
	PrevRandao [32]byte
	// SuggestedFeeRecipient is the suggested fee recipient.
	SuggestedFeeRecipient bellatrix.ExecutionAddress
}

// PayloadAttributesV2 represents the payload attributes v2.
type PayloadAttributesV2 struct {
	// Timestamp is the timestamp of the payload.
	Timestamp uint64
	// PrevRandao is the previous randao.
	PrevRandao [32]byte
	// SuggestedFeeRecipient is the suggested fee recipient.
	SuggestedFeeRecipient bellatrix.ExecutionAddress
	// Withdrawals is the list of withdrawals.
	Withdrawals []*capella.Withdrawal
}

// PayloadAttributesV3 represents the payload attributes v3.
type PayloadAttributesV3 struct {
	// Timestamp is the timestamp of the payload.
	Timestamp uint64
	// PrevRandao is the previous randao.
	PrevRandao [32]byte
	// SuggestedFeeRecipient is the suggested fee recipient.
	SuggestedFeeRecipient bellatrix.ExecutionAddress
	// Withdrawals is the list of withdrawals.
	Withdrawals []*capella.Withdrawal
	// ParentBeaconBlockRoot is the parent beacon block root.
	ParentBeaconBlockRoot phase0.Root
}

// PayloadAttributesV4 represents the payload attributes v4.
type PayloadAttributesV4 struct {
	// Timestamp is the timestamp of the payload.
	Timestamp uint64
	// PrevRandao is the previous randao.
	PrevRandao [32]byte
	// SuggestedFeeRecipient is the suggested fee recipient.
	SuggestedFeeRecipient bellatrix.ExecutionAddress
	// Withdrawals is the list of withdrawals.
	Withdrawals []*capella.Withdrawal
	// ParentBeaconBlockRoot is the parent beacon block root.
	ParentBeaconBlockRoot phase0.Root
	// DepositRequests is the list of deposit receipts.
	DepositRequests []*electra.DepositRequest
	// WithdrawalRequests is the list of withdrawal requests.
	WithdrawalRequests []*electra.WithdrawalRequest
	// ConsolidationRequests is the list of consolidation requests.
	ConsolidationRequests []*electra.ConsolidationRequest
}

// payloadAttributesEventJSON is the spec representation of the event.
type payloadAttributesEventJSON struct {
	Version spec.DataVersion           `json:"version"`
	Data    *payloadAttributesDataJSON `json:"data"`
}

// payloadAttributesDataJSON is the spec representation of the payload attributes data.
type payloadAttributesDataJSON struct {
	ProposerIndex     string          `json:"proposer_index"`
	ProposalSlot      string          `json:"proposal_slot"`
	ParentBlockNumber string          `json:"parent_block_number"`
	ParentBlockRoot   string          `json:"parent_block_root"`
	ParentBlockHash   string          `json:"parent_block_hash"`
	PayloadAttributes json.RawMessage `json:"payload_attributes"`
}

// payloadAttributesV1JSON is the spec representation of the payload attributes.
type payloadAttributesV1JSON struct {
	Timestamp             string `json:"timestamp"`
	PrevRandao            string `json:"prev_randao"`
	SuggestedFeeRecipient string `json:"suggested_fee_recipient"`
}

// payloadAttributesV2JSON is the spec representation of the payload attributes v2.
type payloadAttributesV2JSON struct {
	Timestamp             string                `json:"timestamp"`
	PrevRandao            string                `json:"prev_randao"`
	SuggestedFeeRecipient string                `json:"suggested_fee_recipient"`
	Withdrawals           []*capella.Withdrawal `json:"withdrawals"`
}

// payloadAttributesV3JSON is the spec representation of the payload attributes v3.
type payloadAttributesV3JSON struct {
	Timestamp             string                `json:"timestamp"`
	PrevRandao            string                `json:"prev_randao"`
	SuggestedFeeRecipient string                `json:"suggested_fee_recipient"`
	Withdrawals           []*capella.Withdrawal `json:"withdrawals"`
	ParentBeaconBlockRoot string                `json:"parent_beacon_block_root"`
}

// payloadAttributesV4JSON is the spec representation of the payload attributes v4.
type payloadAttributesV4JSON struct {
	Timestamp             string                          `json:"timestamp"`
	PrevRandao            string                          `json:"prev_randao"`
	SuggestedFeeRecipient string                          `json:"suggested_fee_recipient"`
	Withdrawals           []*capella.Withdrawal           `json:"withdrawals"`
	ParentBeaconBlockRoot string                          `json:"parent_beacon_block_root"`
	DepositRequests       []*electra.DepositRequest       `json:"deposit_requests"`
	WithdrawalRequests    []*electra.WithdrawalRequest    `json:"withdrawal_requests"`
	ConsolidationRequests []*electra.ConsolidationRequest `json:"consolidation_requests"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributesV1) UnmarshalJSON(input []byte) error {
	var payloadAttributes payloadAttributesV1JSON
	if err := json.Unmarshal(input, &payloadAttributes); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&payloadAttributes)
}

func (p *PayloadAttributesV1) unpack(data *payloadAttributesV1JSON) error {
	var err error

	if data.Timestamp == "" {
		return errors.New("payload attributes timestamp missing")
	}
	p.Timestamp, err = strconv.ParseUint(data.Timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes timestamp")
	}

	if data.PrevRandao == "" {
		return errors.New("payload attributes prev randao missing")
	}
	prevRandao, err := hex.DecodeString(strings.TrimPrefix(data.PrevRandao, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes prev randao")
	}
	if len(prevRandao) != 32 {
		return errors.New("incorrect length for payload attributes prev randao")
	}
	copy(p.PrevRandao[:], prevRandao)

	if data.SuggestedFeeRecipient == "" {
		return errors.New("payload attributes suggested fee recipient missing")
	}
	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.SuggestedFeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes suggested fee recipient")
	}
	if len(feeRecipient) != bellatrix.FeeRecipientLength {
		return errors.New("incorrect length for payload attributes suggested fee recipient")
	}
	copy(p.SuggestedFeeRecipient[:], feeRecipient)

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributesV2) UnmarshalJSON(input []byte) error {
	var payloadAttributes payloadAttributesV2JSON
	if err := json.Unmarshal(input, &payloadAttributes); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&payloadAttributes)
}

func (p *PayloadAttributesV2) unpack(data *payloadAttributesV2JSON) error {
	var err error

	if data.Timestamp == "" {
		return errors.New("payload attributes timestamp missing")
	}
	p.Timestamp, err = strconv.ParseUint(data.Timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes timestamp")
	}

	if data.PrevRandao == "" {
		return errors.New("payload attributes prev randao missing")
	}
	prevRandao, err := hex.DecodeString(strings.TrimPrefix(data.PrevRandao, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes prev randao")
	}
	if len(prevRandao) != 32 {
		return errors.New("incorrect length for payload attributes prev randao")
	}
	copy(p.PrevRandao[:], prevRandao)

	if data.SuggestedFeeRecipient == "" {
		return errors.New("payload attributes suggested fee recipient missing")
	}
	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.SuggestedFeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes suggested fee recipient")
	}
	if len(feeRecipient) != bellatrix.FeeRecipientLength {
		return errors.New("incorrect length for payload attributes suggested fee recipient")
	}
	copy(p.SuggestedFeeRecipient[:], feeRecipient)

	if data.Withdrawals == nil {
		return errors.New("payload attributes withdrawals missing")
	}
	for i := range data.Withdrawals {
		if data.Withdrawals[i] == nil {
			return fmt.Errorf("withdrawals entry %d missing", i)
		}
	}
	p.Withdrawals = data.Withdrawals

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributesV3) UnmarshalJSON(input []byte) error {
	var payloadAttributes payloadAttributesV3JSON
	if err := json.Unmarshal(input, &payloadAttributes); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&payloadAttributes)
}

func (p *PayloadAttributesV3) unpack(data *payloadAttributesV3JSON) error {
	var err error

	if data.Timestamp == "" {
		return errors.New("payload attributes timestamp missing")
	}
	p.Timestamp, err = strconv.ParseUint(data.Timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes timestamp")
	}

	if data.PrevRandao == "" {
		return errors.New("payload attributes prev randao missing")
	}
	prevRandao, err := hex.DecodeString(strings.TrimPrefix(data.PrevRandao, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes prev randao")
	}
	if len(prevRandao) != 32 {
		return errors.New("incorrect length for payload attributes prev randao")
	}
	copy(p.PrevRandao[:], prevRandao)

	if data.SuggestedFeeRecipient == "" {
		return errors.New("payload attributes suggested fee recipient missing")
	}
	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.SuggestedFeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes suggested fee recipient")
	}
	if len(feeRecipient) != bellatrix.FeeRecipientLength {
		return errors.New("incorrect length for payload attributes suggested fee recipient")
	}
	copy(p.SuggestedFeeRecipient[:], feeRecipient)

	if data.Withdrawals == nil {
		return errors.New("payload attributes withdrawals missing")
	}
	for i := range data.Withdrawals {
		if data.Withdrawals[i] == nil {
			return fmt.Errorf("withdrawals entry %d missing", i)
		}
	}
	p.Withdrawals = data.Withdrawals

	if data.ParentBeaconBlockRoot == "" {
		return errors.New("payload attributes parent beacon block root missing")
	}
	parentBeaconBlockRoot, err := hex.DecodeString(strings.TrimPrefix(data.ParentBeaconBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes parent beacon block root")
	}
	if len(parentBeaconBlockRoot) != phase0.RootLength {
		return errors.New("incorrect length for payload attributes parent beacon block root")
	}
	copy(p.ParentBeaconBlockRoot[:], parentBeaconBlockRoot)

	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PayloadAttributesV4) UnmarshalJSON(input []byte) error {
	var payloadAttributes payloadAttributesV4JSON
	if err := json.Unmarshal(input, &payloadAttributes); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return p.unpack(&payloadAttributes)
}

func (p *PayloadAttributesV4) unpack(data *payloadAttributesV4JSON) error {
	var err error

	if data.Timestamp == "" {
		return errors.New("payload attributes timestamp missing")
	}
	p.Timestamp, err = strconv.ParseUint(data.Timestamp, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes timestamp")
	}

	if data.PrevRandao == "" {
		return errors.New("payload attributes prev randao missing")
	}
	prevRandao, err := hex.DecodeString(strings.TrimPrefix(data.PrevRandao, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes prev randao")
	}
	if len(prevRandao) != 32 {
		return errors.New("incorrect length for payload attributes prev randao")
	}
	copy(p.PrevRandao[:], prevRandao)

	if data.SuggestedFeeRecipient == "" {
		return errors.New("payload attributes suggested fee recipient missing")
	}
	feeRecipient, err := hex.DecodeString(strings.TrimPrefix(data.SuggestedFeeRecipient, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes suggested fee recipient")
	}
	if len(feeRecipient) != bellatrix.FeeRecipientLength {
		return errors.New("incorrect length for payload attributes suggested fee recipient")
	}
	copy(p.SuggestedFeeRecipient[:], feeRecipient)

	if data.Withdrawals == nil {
		return errors.New("payload attributes withdrawals missing")
	}
	for i := range data.Withdrawals {
		if data.Withdrawals[i] == nil {
			return fmt.Errorf("withdrawals entry %d missing", i)
		}
	}
	p.Withdrawals = data.Withdrawals

	if data.ParentBeaconBlockRoot == "" {
		return errors.New("payload attributes parent beacon block root missing")
	}
	parentBeaconBlockRoot, err := hex.DecodeString(strings.TrimPrefix(data.ParentBeaconBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for payload attributes parent beacon block root")
	}
	if len(parentBeaconBlockRoot) != phase0.RootLength {
		return errors.New("incorrect length for payload attributes parent beacon block root")
	}
	copy(p.ParentBeaconBlockRoot[:], parentBeaconBlockRoot)

	if data.DepositRequests == nil {
		return errors.New("payload attributes deposit requests missing")
	}
	for i := range data.DepositRequests {
		if data.DepositRequests[i] == nil {
			return fmt.Errorf("deposit requests entry %d missing", i)
		}
	}
	p.DepositRequests = data.DepositRequests

	if data.WithdrawalRequests == nil {
		return errors.New("payload attributes withdraw requests missing")
	}
	for i := range data.WithdrawalRequests {
		if data.WithdrawalRequests[i] == nil {
			return fmt.Errorf("withdraw requests entry %d missing", i)
		}
	}
	p.WithdrawalRequests = data.WithdrawalRequests

	if data.ConsolidationRequests == nil {
		return errors.New("payload attributes consolidation requests missing")
	}
	for i := range data.ConsolidationRequests {
		if data.ConsolidationRequests[i] == nil {
			return fmt.Errorf("consolidation requests entry %d missing", i)
		}
	}
	p.ConsolidationRequests = data.ConsolidationRequests

	return nil
}

// MarshalJSON implements json.Marshaler.
func (e *PayloadAttributesEvent) MarshalJSON() ([]byte, error) {
	var payloadAttributes []byte
	var err error

	switch e.Version {
	case spec.DataVersionBellatrix:
		if e.Data.V1 == nil {
			return nil, errors.New("no payload attributes v1 data")
		}
		payloadAttributes, err = json.Marshal(&payloadAttributesV1JSON{
			Timestamp:             strconv.FormatUint(e.Data.V1.Timestamp, 10),
			PrevRandao:            fmt.Sprintf("%#x", e.Data.V1.PrevRandao),
			SuggestedFeeRecipient: e.Data.V1.SuggestedFeeRecipient.String(),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal payload attributes v1")
		}
	case spec.DataVersionCapella:
		if e.Data.V2 == nil {
			return nil, errors.New("no payload attributes v2 data")
		}
		payloadAttributes, err = json.Marshal(&payloadAttributesV2JSON{
			Timestamp:             strconv.FormatUint(e.Data.V2.Timestamp, 10),
			PrevRandao:            fmt.Sprintf("%#x", e.Data.V2.PrevRandao),
			SuggestedFeeRecipient: e.Data.V2.SuggestedFeeRecipient.String(),
			Withdrawals:           e.Data.V2.Withdrawals,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal payload attributes v2")
		}
	case spec.DataVersionDeneb:
		if e.Data.V3 == nil {
			return nil, errors.New("no payload attributes v3 data")
		}
		payloadAttributes, err = json.Marshal(&payloadAttributesV3JSON{
			Timestamp:             strconv.FormatUint(e.Data.V3.Timestamp, 10),
			PrevRandao:            fmt.Sprintf("%#x", e.Data.V3.PrevRandao),
			SuggestedFeeRecipient: e.Data.V3.SuggestedFeeRecipient.String(),
			Withdrawals:           e.Data.V3.Withdrawals,
			ParentBeaconBlockRoot: fmt.Sprintf("%#x", e.Data.V3.ParentBeaconBlockRoot),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal payload attributes v3")
		}
	case spec.DataVersionElectra:
		if e.Data.V4 == nil {
			return nil, errors.New("no payload attributes v4 data")
		}
		payloadAttributes, err = json.Marshal(&payloadAttributesV4JSON{
			Timestamp:             strconv.FormatUint(e.Data.V4.Timestamp, 10),
			PrevRandao:            fmt.Sprintf("%#x", e.Data.V4.PrevRandao),
			SuggestedFeeRecipient: e.Data.V4.SuggestedFeeRecipient.String(),
			Withdrawals:           e.Data.V4.Withdrawals,
			ParentBeaconBlockRoot: fmt.Sprintf("%#x", e.Data.V4.ParentBeaconBlockRoot),
			DepositRequests:       e.Data.V4.DepositRequests,
			WithdrawalRequests:    e.Data.V4.WithdrawalRequests,
			ConsolidationRequests: e.Data.V4.ConsolidationRequests,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal payload attributes v4")
		}
	default:
		return nil, fmt.Errorf("unsupported payload attributes version: %s", e.Version)
	}

	data := payloadAttributesDataJSON{
		ProposerIndex:     fmt.Sprintf("%d", e.Data.ProposerIndex),
		ProposalSlot:      fmt.Sprintf("%d", e.Data.ProposalSlot),
		ParentBlockNumber: strconv.FormatUint(e.Data.ParentBlockNumber, 10),
		ParentBlockRoot:   fmt.Sprintf("%#x", e.Data.ParentBlockRoot),
		ParentBlockHash:   fmt.Sprintf("%#x", e.Data.ParentBlockHash),
		PayloadAttributes: payloadAttributes,
	}

	return json.Marshal(&payloadAttributesEventJSON{
		Version: e.Version,
		Data:    &data,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *PayloadAttributesEvent) UnmarshalJSON(input []byte) error {
	var event payloadAttributesEventJSON
	if err := json.Unmarshal(input, &event); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return e.unpack(&event)
}

func (e *PayloadAttributesEvent) unpack(data *payloadAttributesEventJSON) error {
	var err error

	if data.Data == nil {
		return errors.New("payload attributes data missing")
	}
	e.Data = &PayloadAttributesData{}

	if data.Data.ProposerIndex == "" {
		return errors.New("proposer index missing")
	}
	proposerIndex, err := strconv.ParseUint(data.Data.ProposerIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposer index")
	}
	e.Data.ProposerIndex = phase0.ValidatorIndex(proposerIndex)

	if data.Data.ProposalSlot == "" {
		return errors.New("proposal slot missing")
	}
	proposalSlot, err := strconv.ParseUint(data.Data.ProposalSlot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for proposal slot")
	}
	e.Data.ProposalSlot = phase0.Slot(proposalSlot)

	if data.Data.ParentBlockNumber == "" {
		return errors.New("parent block number missing")
	}
	parentBlockNumber, err := strconv.ParseUint(data.Data.ParentBlockNumber, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid value for parent block number")
	}
	e.Data.ParentBlockNumber = parentBlockNumber

	if data.Data.ParentBlockRoot == "" {
		return errors.New("parent block root missing")
	}
	parentBlockRoot, err := hex.DecodeString(strings.TrimPrefix(data.Data.ParentBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for parent block root")
	}
	if len(parentBlockRoot) != phase0.RootLength {
		return errors.New("incorrect length for parent block root")
	}
	copy(e.Data.ParentBlockRoot[:], parentBlockRoot)

	if data.Data.ParentBlockHash == "" {
		return errors.New("parent block hash missing")
	}
	parentBlockHash, err := hex.DecodeString(strings.TrimPrefix(data.Data.ParentBlockHash, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for parent block hash")
	}
	if len(parentBlockHash) != phase0.Hash32Length {
		return errors.New("incorrect length for parent block hash")
	}
	copy(e.Data.ParentBlockHash[:], parentBlockHash)

	if data.Data.PayloadAttributes == nil {
		return errors.New("payload attributes missing")
	}

	switch data.Version {
	case spec.DataVersionBellatrix:
		var payloadAttributes PayloadAttributesV1
		err = json.Unmarshal(data.Data.PayloadAttributes, &payloadAttributes)
		if err != nil {
			return err
		}
		e.Data.V1 = &payloadAttributes
	case spec.DataVersionCapella:
		var payloadAttributes PayloadAttributesV2
		err = json.Unmarshal(data.Data.PayloadAttributes, &payloadAttributes)
		if err != nil {
			return err
		}
		e.Data.V2 = &payloadAttributes
	case spec.DataVersionDeneb:
		var payloadAttributes PayloadAttributesV3
		err = json.Unmarshal(data.Data.PayloadAttributes, &payloadAttributes)
		if err != nil {
			return err
		}
		e.Data.V3 = &payloadAttributes
	case spec.DataVersionElectra:
		var payloadAttributes PayloadAttributesV4
		err = json.Unmarshal(data.Data.PayloadAttributes, &payloadAttributes)
		if err != nil {
			return err
		}
		e.Data.V4 = &payloadAttributes
	default:
		return errors.New("unsupported data version")
	}
	e.Version = data.Version

	return nil
}

// String returns a string version of the structure.
func (e *PayloadAttributesEvent) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
