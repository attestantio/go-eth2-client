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
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// newDataOf returns a new value of the type into which the data of the topic's events decodes.
func newDataOf[T any](EventTopic[T]) any {
	return new(T)
}

// TestEventTopicDecodeTypes pins, for each topic, its name and the type into which its data
// decodes, both through the typed topic and through the lookup by name that
// Event.UnmarshalJSON uses, so that a topic given the wrong name or type fails here rather than
// handing handlers a different type than documented.
func TestEventTopicDecodeTypes(t *testing.T) {
	tests := []struct {
		name     string
		topic    string
		data     any
		expected any
	}{
		{name: "attestation", topic: AttestationEventTopic.Name(), data: newDataOf(AttestationEventTopic), expected: &spec.VersionedAttestation{}},
		{name: "attester_slashing", topic: AttesterSlashingEventTopic.Name(), data: newDataOf(AttesterSlashingEventTopic), expected: &electra.AttesterSlashing{}},
		{name: "blob_sidecar", topic: BlobSidecarEventTopic.Name(), data: newDataOf(BlobSidecarEventTopic), expected: &BlobSidecarEvent{}},
		{name: "block", topic: BlockEventTopic.Name(), data: newDataOf(BlockEventTopic), expected: &BlockEvent{}},
		{name: "block_gossip", topic: BlockGossipEventTopic.Name(), data: newDataOf(BlockGossipEventTopic), expected: &BlockGossipEvent{}},
		{name: "bls_to_execution_change", topic: BLSToExecutionChangeEventTopic.Name(), data: newDataOf(BLSToExecutionChangeEventTopic), expected: &capella.SignedBLSToExecutionChange{}},
		{name: "chain_reorg", topic: ChainReorgEventTopic.Name(), data: newDataOf(ChainReorgEventTopic), expected: &ChainReorgEvent{}},
		{name: "contribution_and_proof", topic: ContributionAndProofEventTopic.Name(), data: newDataOf(ContributionAndProofEventTopic), expected: &altair.SignedContributionAndProof{}},
		{name: "data_column_sidecar", topic: DataColumnSidecarEventTopic.Name(), data: newDataOf(DataColumnSidecarEventTopic), expected: &DataColumnSidecarEvent{}},
		{name: "execution_payload", topic: ExecutionPayloadEventTopic.Name(), data: newDataOf(ExecutionPayloadEventTopic), expected: &ExecutionPayloadEvent{}},
		{name: "execution_payload_available", topic: ExecutionPayloadAvailableEventTopic.Name(), data: newDataOf(ExecutionPayloadAvailableEventTopic), expected: &ExecutionPayloadAvailableEvent{}},
		{name: "execution_payload_bid", topic: ExecutionPayloadBidEventTopic.Name(), data: newDataOf(ExecutionPayloadBidEventTopic), expected: &gloas.SignedExecutionPayloadBid{}},
		{name: "execution_payload_gossip", topic: ExecutionPayloadGossipEventTopic.Name(), data: newDataOf(ExecutionPayloadGossipEventTopic), expected: &ExecutionPayloadGossipEvent{}},
		{name: "fast_confirmation", topic: FastConfirmationEventTopic.Name(), data: newDataOf(FastConfirmationEventTopic), expected: &FastConfirmationEvent{}},
		{name: "finalized_checkpoint", topic: FinalizedCheckpointEventTopic.Name(), data: newDataOf(FinalizedCheckpointEventTopic), expected: &FinalizedCheckpointEvent{}},
		{name: "head", topic: HeadEventTopic.Name(), data: newDataOf(HeadEventTopic), expected: &HeadEvent{}},
		{name: "payload_attestation_message", topic: PayloadAttestationMessageEventTopic.Name(), data: newDataOf(PayloadAttestationMessageEventTopic), expected: &gloas.PayloadAttestationMessage{}},
		{name: "payload_attributes", topic: PayloadAttributesEventTopic.Name(), data: newDataOf(PayloadAttributesEventTopic), expected: &PayloadAttributesEvent{}},
		{name: "proposer_preferences", topic: ProposerPreferencesEventTopic.Name(), data: newDataOf(ProposerPreferencesEventTopic), expected: &gloas.SignedProposerPreferences{}},
		{name: "proposer_slashing", topic: ProposerSlashingEventTopic.Name(), data: newDataOf(ProposerSlashingEventTopic), expected: &phase0.ProposerSlashing{}},
		{name: "single_attestation", topic: SingleAttestationEventTopic.Name(), data: newDataOf(SingleAttestationEventTopic), expected: &electra.SingleAttestation{}},
		{name: "voluntary_exit", topic: VoluntaryExitEventTopic.Name(), data: newDataOf(VoluntaryExitEventTopic), expected: &phase0.SignedVoluntaryExit{}},
	}

	names := make(map[string]bool, len(tests))
	for _, test := range tests {
		require.NotContains(t, names, test.name)
		names[test.name] = true
	}
	require.Equal(t, SupportedEventTopics, names)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.name, test.topic)
			require.IsType(t, test.expected, test.data)
			require.Contains(t, eventTopicData, test.name)
			require.IsType(t, test.expected, eventTopicData[test.name]())
		})
	}
}

func TestEventTopicDecode(t *testing.T) {
	message := &gloas.PayloadAttestationMessage{
		ValidatorIndex: 123,
		Data:           &gloas.PayloadAttestationData{Slot: 10},
	}
	bare, err := json.Marshal(message)
	require.NoError(t, err)

	wrapped := func(version string) []byte {
		t.Helper()

		data, err := json.Marshal(map[string]any{"version": version, "data": json.RawMessage(bare)})
		require.NoError(t, err)

		return data
	}

	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{name: "Wrapped", input: wrapped("gloas")},
		{name: "Bare", input: bare},
		{name: "OtherFork", input: wrapped("fulu"), err: `unsupported version "fulu" for payload_attestation_message event`},
		{name: "UnknownFork", input: wrapped("unknown"), err: `unsupported version "unknown" for payload_attestation_message event`},
		{name: "Malformed", input: []byte(`invalid`), err: "invalid character 'i' looking for beginning of value"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := PayloadAttestationMessageEventTopic.Decode(test.input)
			if test.err != "" {
				require.EqualError(t, err, test.err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, phase0.ValidatorIndex(123), data.ValidatorIndex)
			require.Equal(t, phase0.Slot(10), data.Data.Slot)
		})
	}
}
