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

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// executionPayloadEnvelopeJSON is the spec representation of the struct.
type executionPayloadEnvelopeJSON struct {
	Payload            *deneb.ExecutionPayload    `json:"payload"`
	ExecutionRequests  *electra.ExecutionRequests `json:"execution_requests"`
	BuilderIndex       string                     `json:"builder_index"`
	BeaconBlockRoot    string                     `json:"beacon_block_root"`
	Slot               string                     `json:"slot"`
	BlobKZGCommitments []string                   `json:"blob_kzg_commitments"`
	PayloadWithheld    bool                       `json:"payload_withheld"`
	StateRoot          string                     `json:"state_root"`
}

// MarshalJSON implements json.Marshaler.
func (e *ExecutionPayloadEnvelope) MarshalJSON() ([]byte, error) {
	blobCommitments := make([]string, len(e.BlobKZGCommitments))
	for i := range e.BlobKZGCommitments {
		blobCommitments[i] = fmt.Sprintf("%#x", e.BlobKZGCommitments[i])
	}

	return json.Marshal(&executionPayloadEnvelopeJSON{
		Payload:            e.Payload,
		ExecutionRequests:  e.ExecutionRequests,
		BuilderIndex:       fmt.Sprintf("%d", e.BuilderIndex),
		BeaconBlockRoot:    fmt.Sprintf("%#x", e.BeaconBlockRoot),
		Slot:               fmt.Sprintf("%d", e.Slot),
		BlobKZGCommitments: blobCommitments,
		StateRoot:          fmt.Sprintf("%#x", e.StateRoot),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *ExecutionPayloadEnvelope) UnmarshalJSON(input []byte) error {
	var data executionPayloadEnvelopeJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	if data.Payload == nil {
		return errors.New("payload missing")
	}
	e.Payload = data.Payload

	if data.ExecutionRequests == nil {
		return errors.New("execution requests missing")
	}
	e.ExecutionRequests = data.ExecutionRequests

	if data.BuilderIndex == "" {
		return errors.New("builder index missing")
	}
	builderIndex, err := strconv.ParseUint(data.BuilderIndex, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid builder index")
	}
	e.BuilderIndex = phase0.ValidatorIndex(builderIndex)

	if data.BeaconBlockRoot == "" {
		return errors.New("beacon block root missing")
	}
	root, err := hex.DecodeString(strings.TrimPrefix(data.BeaconBlockRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid beacon block root")
	}
	copy(e.BeaconBlockRoot[:], root)

	if data.Slot == "" {
		return errors.New("slot missing")
	}
	slot, err := strconv.ParseUint(data.Slot, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid slot")
	}
	e.Slot = phase0.Slot(slot)

	e.BlobKZGCommitments = make([]deneb.KZGCommitment, len(data.BlobKZGCommitments))
	for i := range data.BlobKZGCommitments {
		commitment, err := hex.DecodeString(strings.TrimPrefix(data.BlobKZGCommitments[i], "0x"))
		if err != nil {
			return errors.Wrap(err, "invalid blob KZG commitment")
		}
		copy(e.BlobKZGCommitments[i][:], commitment)
	}

	if data.StateRoot == "" {
		return errors.New("state root missing")
	}
	stateRoot, err := hex.DecodeString(strings.TrimPrefix(data.StateRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid state root")
	}
	copy(e.StateRoot[:], stateRoot)

	return nil
}
