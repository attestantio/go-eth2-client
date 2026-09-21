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

// SupportedEventTopics is the event-topic catalog from beacon-APIs commit
// ef98d512c03c8ca6b9d7cbdc45b9293ec2b24722:
// https://github.com/ethereum/beacon-APIs/blob/ef98d512c03c8ca6b9d7cbdc45b9293ec2b24722/apis/eventstream/index.yaml.
// The HTTP client uses it to validate Events() subscriptions.
var SupportedEventTopics = map[string]bool{
	"attestation":                    true,
	"attester_slashing":              true,
	"block":                          true,
	"block_gossip":                   true,
	"bls_to_execution_change":        true,
	"chain_reorg":                    true,
	"contribution_and_proof":         true,
	"data_column_sidecar":            true,
	"execution_payload":              true,
	"execution_payload_available":    true,
	"execution_payload_bid":          true,
	"execution_payload_gossip":       true,
	"fast_confirmation":              true,
	"finalized_checkpoint":           true,
	"head":                           true,
	"head_v2":                        true,
	"light_client_finality_update":   true,
	"light_client_optimistic_update": true,
	"payload_attestation_message":    true,
	"payload_attributes":             true,
	"proposer_preferences":           true,
	"proposer_slashing":              true,
	"single_attestation":             true,
	"voluntary_exit":                 true,
}

// eventJSON is the spec representation of the struct.
type eventJSON struct {
	Topic string         `json:"topic"`
	Data  map[string]any `json:"data"`
}

// MarshalJSON implements json.Marshaler.
func (e *Event) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(e.Data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal data")
	}

	var data map[string]any
	if err := json.Unmarshal(marshalled, &data); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal data")
	}
	if isVersionedGloasEventTopic(e.Topic) {
		data = map[string]any{
			"version": spec.DataVersionGloas.String(),
			"data":    data,
		}
	}

	return json.Marshal(&eventJSON{
		Topic: e.Topic,
		Data:  data,
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
		e.Data = &spec.VersionedAttestation{}
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
	case "data_column_sidecar":
		e.Data = &DataColumnSidecarEvent{}
	case "execution_payload", "execution_payload_gossip":
		e.Data = &ExecutionPayloadEvent{}
	case "execution_payload_available":
		e.Data = &ExecutionPayloadAvailableEvent{}
	case "execution_payload_bid":
		e.Data = &gloas.SignedExecutionPayloadBid{}
	case "fast_confirmation":
		e.Data = &FastConfirmationEvent{}
	case "finalized_checkpoint":
		e.Data = &FinalizedCheckpointEvent{}
	case "head":
		e.Data = &HeadEvent{}
	case "head_v2":
		e.Data = &HeadEventV2{}
	case "light_client_finality_update":
		e.Data = &LightClientFinalityUpdateEvent{}
	case "light_client_optimistic_update":
		e.Data = &LightClientOptimisticUpdateEvent{}
	case "payload_attestation_message":
		e.Data = &gloas.PayloadAttestationMessage{}
	case "payload_attributes":
		e.Data = &PayloadAttributesEvent{}
	case "proposer_preferences":
		e.Data = &gloas.SignedProposerPreferences{}
	case "proposer_slashing":
		e.Data = &phase0.ProposerSlashing{}
	case "single_attestation":
		e.Data = &electra.SingleAttestation{}
	case "voluntary_exit":
		e.Data = &phase0.SignedVoluntaryExit{}
	default:
		return fmt.Errorf("unsupported event topic %s", eventJSON.Topic)
	}

	data, err := json.Marshal(eventJSON.Data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}
	if isVersionedGloasEventTopic(eventJSON.Topic) {
		var versionedData struct {
			Version spec.DataVersion `json:"version"`
			Data    json.RawMessage  `json:"data"`
		}
		if err := json.Unmarshal(data, &versionedData); err != nil {
			return errors.Wrap(err, "invalid versioned event data")
		}
		if versionedData.Version != spec.DataVersionGloas {
			return fmt.Errorf("unsupported event data version %s", versionedData.Version)
		}
		if len(versionedData.Data) == 0 {
			return errors.New("event data missing")
		}
		data = versionedData.Data
	}

	if err := json.Unmarshal(data, e.Data); err != nil {
		return errors.Wrap(err, "invalid event data")
	}

	return nil
}

func isVersionedGloasEventTopic(topic string) bool {
	switch topic {
	case "execution_payload_bid", "payload_attestation_message", "proposer_preferences":
		return true
	default:
		return false
	}
}

// String returns a string version of the structure.
func (e *Event) String() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
