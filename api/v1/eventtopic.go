// Copyright © 2026 Attestant Limited.
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
)

// EventTopic describes an event topic: its name, and T, the type into which the data of its
// events decodes.  The topics below are the single list of supported topics: SupportedEventTopics,
// Event.UnmarshalJSON and the handlers in api.EventsOpts are all derived from them.
type EventTopic[T any] struct {
	name string
	// version is the fork of the data of the topic's events, for a topic whose events the
	// beacon-API spec wraps as {"version": "...", "data": {...}}, and DataVersionUnknown
	// otherwise.
	version spec.DataVersion
}

// The supported event topics.
var (
	AttestationEventTopic               = newEventTopic[spec.VersionedAttestation]("attestation")
	AttesterSlashingEventTopic          = newEventTopic[electra.AttesterSlashing]("attester_slashing")
	BlobSidecarEventTopic               = newEventTopic[BlobSidecarEvent]("blob_sidecar")
	BlockEventTopic                     = newEventTopic[BlockEvent]("block")
	BlockGossipEventTopic               = newEventTopic[BlockGossipEvent]("block_gossip")
	BLSToExecutionChangeEventTopic      = newEventTopic[capella.SignedBLSToExecutionChange]("bls_to_execution_change")
	ChainReorgEventTopic                = newEventTopic[ChainReorgEvent]("chain_reorg")
	ContributionAndProofEventTopic      = newEventTopic[altair.SignedContributionAndProof]("contribution_and_proof")
	DataColumnSidecarEventTopic         = newEventTopic[DataColumnSidecarEvent]("data_column_sidecar")
	ExecutionPayloadEventTopic          = newEventTopic[ExecutionPayloadEvent]("execution_payload")
	ExecutionPayloadAvailableEventTopic = newEventTopic[ExecutionPayloadAvailableEvent]("execution_payload_available")
	ExecutionPayloadBidEventTopic       = newVersionedEventTopic[gloas.SignedExecutionPayloadBid](
		"execution_payload_bid", spec.DataVersionGloas)
	ExecutionPayloadGossipEventTopic    = newEventTopic[ExecutionPayloadGossipEvent]("execution_payload_gossip")
	FastConfirmationEventTopic          = newEventTopic[FastConfirmationEvent]("fast_confirmation")
	FinalizedCheckpointEventTopic       = newEventTopic[FinalizedCheckpointEvent]("finalized_checkpoint")
	HeadEventTopic                      = newEventTopic[HeadEvent]("head")
	PayloadAttestationMessageEventTopic = newVersionedEventTopic[gloas.PayloadAttestationMessage](
		"payload_attestation_message", spec.DataVersionGloas)
	PayloadAttributesEventTopic   = newEventTopic[PayloadAttributesEvent]("payload_attributes")
	ProposerPreferencesEventTopic = newVersionedEventTopic[gloas.SignedProposerPreferences](
		"proposer_preferences", spec.DataVersionGloas)
	ProposerSlashingEventTopic  = newEventTopic[phase0.ProposerSlashing]("proposer_slashing")
	SingleAttestationEventTopic = newEventTopic[electra.SingleAttestation]("single_attestation")
	VoluntaryExitEventTopic     = newEventTopic[phase0.SignedVoluntaryExit]("voluntary_exit")
)

// eventTopicData is, by name, a function returning a new T for each topic created with
// newEventTopic or newVersionedEventTopic.  Go initialises it before the topics above, as their
// initialisers refer to it; SupportedEventTopics is derived from it in init, which runs only
// once every topic has been added.
var eventTopicData = map[string]func() any{}

func newEventTopic[T any](name string) EventTopic[T] {
	return newVersionedEventTopic[T](name, spec.DataVersionUnknown)
}

func newVersionedEventTopic[T any](name string, version spec.DataVersion) EventTopic[T] {
	eventTopicData[name] = func() any { return new(T) }

	return EventTopic[T]{
		name:    name,
		version: version,
	}
}

func init() {
	SupportedEventTopics = make(map[string]bool, len(eventTopicData))
	for topic := range eventTopicData {
		SupportedEventTopics[topic] = true
	}
}

// Name returns the name of the topic.
func (t EventTopic[T]) Name() string {
	return t.name
}

// Decode decodes the data of an event of the topic as sent on the events stream.  For a topic
// whose events the beacon-API spec wraps as {"version": "...", "data": {...}}, the version must
// be the fork of T, as data of another fork could decode into T without error while dropping
// fields T does not have.  A bare, unwrapped object is accepted as well, for nodes that do not
// wrap it.
func (t EventTopic[T]) Decode(input []byte) (*T, error) {
	data := new(T)

	if t.version != spec.DataVersionUnknown {
		var wrapper struct {
			Version string          `json:"version"`
			Data    json.RawMessage `json:"data"`
		}

		if err := json.Unmarshal(input, &wrapper); err == nil && len(wrapper.Data) > 0 && wrapper.Version != "" {
			version, err := spec.DataVersionFromString(wrapper.Version)
			if err != nil || version != t.version {
				return nil, fmt.Errorf("unsupported version %q for %s event", wrapper.Version, t.name)
			}

			input = wrapper.Data
		}
	}

	if err := json.Unmarshal(input, data); err != nil {
		return nil, err
	}

	return data, nil
}
