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
	"reflect"
	"testing"

	"github.com/attestantio/go-eth2-client/internal/eventtopic"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestEventTopicDecodeTypes pins, for each topic, the type into which its data decodes, both by
// name, as Event.UnmarshalJSON finds it, and by type, as the HTTP and multi clients find it
// through the eventtopic registry, so that a topic given the wrong name or type fails here
// rather than handing handlers a different type than documented.
func TestEventTopicDecodeTypes(t *testing.T) {
	tests := []struct {
		name     string
		expected any
	}{
		{name: "attestation", expected: &spec.VersionedAttestation{}},
		{name: "attester_slashing", expected: &electra.AttesterSlashing{}},
		{name: "blob_sidecar", expected: &BlobSidecarEvent{}},
		{name: "block", expected: &BlockEvent{}},
		{name: "block_gossip", expected: &BlockGossipEvent{}},
		{name: "bls_to_execution_change", expected: &capella.SignedBLSToExecutionChange{}},
		{name: "chain_reorg", expected: &ChainReorgEvent{}},
		{name: "contribution_and_proof", expected: &altair.SignedContributionAndProof{}},
		{name: "data_column_sidecar", expected: &DataColumnSidecarEvent{}},
		{name: "execution_payload", expected: &ExecutionPayloadEvent{}},
		{name: "execution_payload_available", expected: &ExecutionPayloadAvailableEvent{}},
		{name: "execution_payload_bid", expected: &gloas.SignedExecutionPayloadBid{}},
		{name: "execution_payload_gossip", expected: &ExecutionPayloadGossipEvent{}},
		{name: "fast_confirmation", expected: &FastConfirmationEvent{}},
		{name: "finalized_checkpoint", expected: &FinalizedCheckpointEvent{}},
		{name: "head", expected: &HeadEvent{}},
		{name: "payload_attestation_message", expected: &gloas.PayloadAttestationMessage{}},
		{name: "payload_attributes", expected: &PayloadAttributesEvent{}},
		{name: "proposer_preferences", expected: &gloas.SignedProposerPreferences{}},
		{name: "proposer_slashing", expected: &phase0.ProposerSlashing{}},
		{name: "single_attestation", expected: &electra.SingleAttestation{}},
		{name: "voluntary_exit", expected: &phase0.SignedVoluntaryExit{}},
	}

	names := make(map[string]bool, len(tests))
	for _, test := range tests {
		require.NotContains(t, names, test.name)
		names[test.name] = true
	}
	require.Equal(t, SupportedEventTopics, names)
	require.Len(t, eventTopics, len(tests), "a topic is listed more than once")

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Contains(t, eventTopicsByName, test.name)
			require.IsType(t, test.expected, eventTopicsByName[test.name].NewData())

			registered, exists := eventtopic.LookupType(reflect.TypeOf(test.expected).Elem())
			require.True(t, exists, "no topic registered for the type")
			require.Equal(t, test.name, registered.Name())
		})
	}
}

// TestEventTopicsRejectNullData confirms that the data type of every topic rejects null.  Decode
// relies on this: a type without an UnmarshalJSON of its own would decode null into a zero value
// without error, handing handlers a value with none of its fields set.
func TestEventTopicsRejectNullData(t *testing.T) {
	for _, topic := range eventTopics {
		t.Run(topic.Name(), func(t *testing.T) {
			require.Error(t, json.Unmarshal([]byte(`null`), topic.NewData()))
		})
	}
}
