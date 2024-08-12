// Copyright © 2020, 2021 Attestant Limited.
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
	"encoding/json"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// Event is the container for events sent from the API.
type Event struct {
	// Topic is the topic of the event.
	Topic string
	// Data is the data of the event.
	Data any
}

// SupportedEventTopics is a map of supported event topics.
var SupportedEventTopics = map[string]bool{
	"attestation":             true,
	"attester_slashing":       true,
	"blob_sidecar":            true,
	"block":                   true,
	"block_gossip":            true,
	"bls_to_execution_change": true,
	"chain_reorg":             true,
	"contribution_and_proof":  true,
	"finalized_checkpoint":    true,
	"head":                    true,
	"payload_attributes":      true,
	"proposer_slashing":       true,
	"voluntary_exit":          true,
}

// eventJSON is the spec representation of the struct.
type eventJSON struct {
	Topic string         `json:"topic"`
	Data  map[string]any `json:"data"`
}

// MarshalJSON implements json.Marshaler.
func (e *Event) MarshalJSON() ([]byte, error) {
	// Need to turn event data in to a generic map.
	data, err := json.Marshal(e.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal data")
	}
	var unmarshalled map[string]any
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal data")
	}

	return json.Marshal(&eventJSON{
		Topic: e.Topic,
		Data:  unmarshalled,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (e *Event) UnmarshalJSON(input []byte) error {
	var err error

	var eventJSON eventJSON
	if err = json.Unmarshal(input, &eventJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if eventJSON.Topic == "" {
		return errors.New("topic missing")
	}
	e.Topic = eventJSON.Topic

	if eventJSON.Data == nil {
		return errors.New("data missing")
	}
	switch eventJSON.Topic {
	case "attestation":
		e.Data = &phase0.Attestation{}
	case "attester_slashing":
		e.Data = &phase0.AttesterSlashing{}
	case "blob_sidecar":
		e.Data = &BlobSidecarEvent{}
	case "block":
		e.Data = &BlockEvent{}
	case "block_gossip":
		e.Data = &BlockGossipEvent{}
	case "bls_to_execution_change":
		e.Data = &capella.SignedBLSToExecutionChange{}
	case "chain_reorg":
		e.Data = &ChainReorgEvent{}
	case "contribution_and_proof":
		e.Data = &altair.SignedContributionAndProof{}
	case "finalized_checkpoint":
		e.Data = &FinalizedCheckpointEvent{}
	case "head":
		e.Data = &HeadEvent{}
	case "payload_attributes":
		e.Data = &PayloadAttributesEvent{}
	case "proposer_slashing":
		e.Data = &phase0.ProposerSlashing{}
	case "voluntary_exit":
		e.Data = &phase0.SignedVoluntaryExit{}
	default:
		return fmt.Errorf("unsupported event topic %s", eventJSON.Topic)
	}
	data, err := json.Marshal(eventJSON.Data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}
	if err := json.Unmarshal(data, &e.Data); err != nil {
		return errors.New("data missing")
	}
	e.Data = eventJSON.Data

	return nil
}

// String returns a string version of the structure.
func (e *Event) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
