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
	"github.com/attestantio/go-eth2-client/internal/eventtopic"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// eventTopics is the single list of supported event topics, each with the type into which the
// data of its events decodes.  SupportedEventTopics and Event.UnmarshalJSON are derived from it,
// and the HTTP and multi clients reach it through the eventtopic package, with which it is
// registered.
var eventTopics = []eventtopic.Descriptor{
	eventtopic.New[spec.VersionedAttestation]("attestation"),
	eventtopic.New[electra.AttesterSlashing]("attester_slashing"),
	eventtopic.New[BlobSidecarEvent]("blob_sidecar"),
	eventtopic.New[BlockEvent]("block"),
	eventtopic.New[BlockGossipEvent]("block_gossip"),
	eventtopic.New[capella.SignedBLSToExecutionChange]("bls_to_execution_change"),
	eventtopic.New[ChainReorgEvent]("chain_reorg"),
	eventtopic.New[altair.SignedContributionAndProof]("contribution_and_proof"),
	eventtopic.New[DataColumnSidecarEvent]("data_column_sidecar"),
	eventtopic.New[ExecutionPayloadEvent]("execution_payload"),
	eventtopic.New[ExecutionPayloadAvailableEvent]("execution_payload_available"),
	eventtopic.NewVersioned[gloas.SignedExecutionPayloadBid]("execution_payload_bid", spec.DataVersionGloas),
	eventtopic.New[ExecutionPayloadGossipEvent]("execution_payload_gossip"),
	eventtopic.New[FastConfirmationEvent]("fast_confirmation"),
	eventtopic.New[FinalizedCheckpointEvent]("finalized_checkpoint"),
	eventtopic.New[HeadEvent]("head"),
	eventtopic.NewVersioned[gloas.PayloadAttestationMessage]("payload_attestation_message", spec.DataVersionGloas),
	eventtopic.New[PayloadAttributesEvent]("payload_attributes"),
	eventtopic.NewVersioned[gloas.SignedProposerPreferences]("proposer_preferences", spec.DataVersionGloas),
	eventtopic.New[phase0.ProposerSlashing]("proposer_slashing"),
	eventtopic.New[electra.SingleAttestation]("single_attestation"),
	eventtopic.New[phase0.SignedVoluntaryExit]("voluntary_exit"),
}

// eventTopicsByName is eventTopics by name.
var eventTopicsByName = func() map[string]eventtopic.Descriptor {
	byName := make(map[string]eventtopic.Descriptor, len(eventTopics))
	for _, topic := range eventTopics {
		byName[topic.Name()] = topic
	}

	return byName
}()

// SupportedEventTopics is a map of supported event topics.  It is informational: the clients
// validate Events() subscriptions against the topics themselves, so changing it has no effect.
var SupportedEventTopics = func() map[string]bool {
	supported := make(map[string]bool, len(eventTopics))
	for _, topic := range eventTopics {
		supported[topic.Name()] = true
	}

	return supported
}()

func init() {
	eventtopic.Register(eventTopics...)
}
