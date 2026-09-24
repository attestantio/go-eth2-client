// Copyright © 2020 - 2026 Attestant Limited.
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

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
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

// eventTopicData is the set of supported event topics, each with a function returning a new
// value of the type into which the data of an event of that topic decodes.  It is the single
// list of topics in this package: SupportedEventTopics and Event.UnmarshalJSON are both derived
// from it.
var eventTopicData = map[string]func() any{
	"attestation":                 func() any { return &spec.VersionedAttestation{} },
	"attester_slashing":           func() any { return &phase0.AttesterSlashing{} },
	"blob_sidecar":                func() any { return &BlobSidecarEvent{} },
	"block":                       func() any { return &BlockEvent{} },
	"block_gossip":                func() any { return &BlockGossipEvent{} },
	"bls_to_execution_change":     func() any { return &capella.SignedBLSToExecutionChange{} },
	"chain_reorg":                 func() any { return &ChainReorgEvent{} },
	"contribution_and_proof":      func() any { return &altair.SignedContributionAndProof{} },
	"data_column_sidecar":         func() any { return &DataColumnSidecarEvent{} },
	"execution_payload":           func() any { return &ExecutionPayloadEvent{} },
	"execution_payload_available": func() any { return &ExecutionPayloadAvailableEvent{} },
	"execution_payload_bid":       func() any { return &gloas.SignedExecutionPayloadBid{} },
	"execution_payload_gossip":    func() any { return &ExecutionPayloadGossipEvent{} },
	"fast_confirmation":           func() any { return &FastConfirmationEvent{} },
	"finalized_checkpoint":        func() any { return &FinalizedCheckpointEvent{} },
	"head":                        func() any { return &HeadEvent{} },
	"payload_attestation_message": func() any { return &gloas.PayloadAttestationMessage{} },
	"payload_attributes":          func() any { return &PayloadAttributesEvent{} },
	"proposer_preferences":        func() any { return &gloas.SignedProposerPreferences{} },
	"proposer_slashing":           func() any { return &phase0.ProposerSlashing{} },
	"single_attestation":          func() any { return &electra.SingleAttestation{} },
	"voluntary_exit":              func() any { return &phase0.SignedVoluntaryExit{} },
}

// SupportedEventTopics is a map of supported event topics. It is the allow-list
// against which the HTTP client validates Events() subscriptions.
var SupportedEventTopics = supportedEventTopics()

func supportedEventTopics() map[string]bool {
	topics := make(map[string]bool, len(eventTopicData))
	for topic := range eventTopicData {
		topics[topic] = true
	}

	return topics
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

	newData, exists := eventTopicData[eventJSON.Topic]
	if !exists {
		return fmt.Errorf("unsupported event topic %s", eventJSON.Topic)
	}

	e.Data = newData()

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
