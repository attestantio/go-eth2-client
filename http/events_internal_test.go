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

package http

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/OffchainLabs/go-bitfield"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/r3labs/sse/v2"
	"github.com/stretchr/testify/require"
)

func TestHandleEventDispatchesEverySupportedTopic(t *testing.T) {
	tests := []struct {
		topic     string
		data      any
		expected  any
		versioned bool
	}{
		{
			topic: "attestation",
			data: &phase0.Attestation{
				AggregationBits: bitfield.NewBitlist(1),
				Data:            testAttestationData(),
			},
			expected: &spec.VersionedAttestation{},
		},
		{
			topic: "attester_slashing",
			data: &electra.AttesterSlashing{
				Attestation1: testIndexedAttestation(),
				Attestation2: testIndexedAttestation(),
			},
			expected: &electra.AttesterSlashing{},
		},
		{topic: "blob_sidecar", data: &apiv1.BlobSidecarEvent{}, expected: &apiv1.BlobSidecarEvent{}},
		{topic: "block", data: &apiv1.BlockEvent{}, expected: &apiv1.BlockEvent{}},
		{topic: "block_gossip", data: &apiv1.BlockGossipEvent{}, expected: &apiv1.BlockGossipEvent{}},
		{
			topic:    "bls_to_execution_change",
			data:     &capella.SignedBLSToExecutionChange{Message: &capella.BLSToExecutionChange{}},
			expected: &capella.SignedBLSToExecutionChange{},
		},
		{topic: "chain_reorg", data: &apiv1.ChainReorgEvent{}, expected: &apiv1.ChainReorgEvent{}},
		{
			topic: "contribution_and_proof",
			data: &altair.SignedContributionAndProof{Message: &altair.ContributionAndProof{
				AggregatorIndex: 1,
				Contribution: &altair.SyncCommitteeContribution{
					AggregationBits: bitfield.NewBitvector128(),
				},
			}},
			expected: &altair.SignedContributionAndProof{},
		},
		{
			topic: "data_column_sidecar",
			data: &apiv1.DataColumnSidecarEvent{
				KZGCommitments: []deneb.KZGCommitment{{}},
			},
			expected: &apiv1.DataColumnSidecarEvent{},
		},
		{topic: "execution_payload", data: &apiv1.ExecutionPayloadEvent{}, expected: &apiv1.ExecutionPayloadEvent{}, versioned: true},
		{topic: "execution_payload_available", data: &apiv1.ExecutionPayloadAvailableEvent{}, expected: &apiv1.ExecutionPayloadAvailableEvent{}},
		{
			topic:     "execution_payload_bid",
			data:      &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{}},
			expected:  &gloas.SignedExecutionPayloadBid{},
			versioned: true,
		},
		{topic: "execution_payload_gossip", data: &apiv1.ExecutionPayloadEvent{}, expected: &apiv1.ExecutionPayloadEvent{}, versioned: true},
		{topic: "fast_confirmation", data: &apiv1.FastConfirmationEvent{}, expected: &apiv1.FastConfirmationEvent{}},
		{topic: "finalized_checkpoint", data: &apiv1.FinalizedCheckpointEvent{}, expected: &apiv1.FinalizedCheckpointEvent{}},
		{
			topic:    "head",
			data:     json.RawMessage(`{"slot":"231192","block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","state":"0x61099b2c1dee0104c93ce0e14e5f5fc4b6faceff4cb863278d055bdfb73b7dc7","epoch_transition":false,"previous_duty_dependent_root":"0xa692c095bbca3eeaf99eeabada78874c028c02b176ccf691f3e8fa075d67f5c6","current_duty_dependent_root":"0x92c6b763f610d5941d2041906007bf9449d37772aacf0483a76275ac27c096b4"}`),
			expected: &apiv1.HeadEvent{},
		},
		{
			topic:     "payload_attestation_message",
			data:      &gloas.PayloadAttestationMessage{Data: &gloas.PayloadAttestationData{}},
			expected:  &gloas.PayloadAttestationMessage{},
			versioned: true,
		},
		{
			topic:    "payload_attributes",
			data:     json.RawMessage(`{"version":"bellatrix","data":{"proposer_index":"123","proposal_slot":"10","parent_block_number":"9","parent_block_root":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","parent_block_hash":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_attributes":{"timestamp":"123456","prev_randao":"0xcf8e0d4e9587369b2301d0790347320302cc0943d5a1884560367e8208d920f2","suggested_fee_recipient":"0x0000000000000000000000000000000000000000"}}}`),
			expected: &apiv1.PayloadAttributesEvent{},
		},
		{
			topic:     "proposer_preferences",
			data:      &gloas.SignedProposerPreferences{Message: &gloas.ProposerPreferences{}},
			expected:  &gloas.SignedProposerPreferences{},
			versioned: true,
		},
		{
			topic: "proposer_slashing",
			data: &phase0.ProposerSlashing{
				SignedHeader1: &phase0.SignedBeaconBlockHeader{Message: &phase0.BeaconBlockHeader{}},
				SignedHeader2: &phase0.SignedBeaconBlockHeader{Message: &phase0.BeaconBlockHeader{}},
			},
			expected: &phase0.ProposerSlashing{},
		},
		{
			topic: "single_attestation",
			data: &electra.SingleAttestation{
				CommitteeIndex: 1,
				AttesterIndex:  1,
				Data:           testAttestationData(),
			},
			expected: &electra.SingleAttestation{},
		},
		{
			topic:    "voluntary_exit",
			data:     &phase0.SignedVoluntaryExit{Message: &phase0.VoluntaryExit{}},
			expected: &phase0.SignedVoluntaryExit{},
		},
	}

	topics := make(map[string]bool, len(tests))
	for _, test := range tests {
		require.NotContains(t, topics, test.topic)
		topics[test.topic] = true
	}
	require.Equal(t, apiv1.SupportedEventTopics, topics)

	for _, test := range tests {
		t.Run(test.topic, func(t *testing.T) {
			require.True(t, apiv1.SupportedEventTopics[test.topic])

			var event *apiv1.Event
			service := &Service{}
			service.handleEvent(context.Background(), &sse.Event{
				Event: []byte(test.topic),
				Data:  eventData(t, test.data, test.versioned),
			}, &api.EventsOpts{
				Handler: func(received *apiv1.Event) {
					event = received
				},
			})

			require.NotNil(t, event)
			require.Equal(t, test.topic, event.Topic)
			require.IsType(t, test.expected, event.Data)
		})
	}
}

type gloasHandlerTest struct {
	topic      string
	data       any
	setHandler func(*api.EventsOpts, func())
}

func gloasHandlerTests() []gloasHandlerTest {
	return []gloasHandlerTest{
		{
			topic: "execution_payload",
			data:  &apiv1.ExecutionPayloadEvent{},
			setHandler: func(opts *api.EventsOpts, handled func()) {
				opts.ExecutionPayloadHandler = func(context.Context, *apiv1.ExecutionPayloadEvent) { handled() }
			},
		},
		{
			topic: "execution_payload_available",
			data:  &apiv1.ExecutionPayloadAvailableEvent{},
			setHandler: func(opts *api.EventsOpts, handled func()) {
				opts.ExecutionPayloadAvailableHandler = func(context.Context, *apiv1.ExecutionPayloadAvailableEvent) { handled() }
			},
		},
		{
			topic: "execution_payload_bid",
			data:  &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{}},
			setHandler: func(opts *api.EventsOpts, handled func()) {
				opts.ExecutionPayloadBidHandler = func(context.Context, *gloas.SignedExecutionPayloadBid) { handled() }
			},
		},
		{
			topic: "execution_payload_gossip",
			data:  &apiv1.ExecutionPayloadEvent{},
			setHandler: func(opts *api.EventsOpts, handled func()) {
				opts.ExecutionPayloadGossipHandler = func(context.Context, *apiv1.ExecutionPayloadEvent) { handled() }
			},
		},
		{
			topic: "fast_confirmation",
			data:  &apiv1.FastConfirmationEvent{},
			setHandler: func(opts *api.EventsOpts, handled func()) {
				opts.FastConfirmationHandler = func(context.Context, *apiv1.FastConfirmationEvent) { handled() }
			},
		},
		{
			topic: "payload_attestation_message",
			data:  &gloas.PayloadAttestationMessage{Data: &gloas.PayloadAttestationData{}},
			setHandler: func(opts *api.EventsOpts, handled func()) {
				opts.PayloadAttestationMessageHandler = func(context.Context, *gloas.PayloadAttestationMessage) { handled() }
			},
		},
		{
			topic: "proposer_preferences",
			data:  &gloas.SignedProposerPreferences{Message: &gloas.ProposerPreferences{}},
			setHandler: func(opts *api.EventsOpts, handled func()) {
				opts.ProposerPreferencesHandler = func(context.Context, *gloas.SignedProposerPreferences) { handled() }
			},
		},
	}

}

func TestHandleEventGloasSpecificHandlers(t *testing.T) {
	tests := gloasHandlerTests()

	for _, test := range tests {
		t.Run(test.topic, func(t *testing.T) {
			var genericHandled bool
			var specificHandled bool
			opts := &api.EventsOpts{
				Handler: func(*apiv1.Event) {
					genericHandled = true
				},
			}
			test.setHandler(opts, func() { specificHandled = true })

			(&Service{}).handleEvent(context.Background(), &sse.Event{
				Event: []byte(test.topic),
				Data:  eventData(t, test.data, false),
			}, opts)

			require.True(t, specificHandled)
			require.False(t, genericHandled)
		})
	}
}

func TestHandleEventGloasMalformedPayloads(t *testing.T) {
	tests := gloasHandlerTests()

	for _, test := range tests {
		t.Run(test.topic, func(t *testing.T) {
			var genericHandled bool
			var specificHandled bool
			opts := &api.EventsOpts{
				Handler: func(*apiv1.Event) {
					genericHandled = true
				},
			}
			test.setHandler(opts, func() { specificHandled = true })

			(&Service{}).handleEvent(context.Background(), &sse.Event{
				Event: []byte(test.topic),
				Data:  []byte(`invalid`),
			}, opts)

			require.False(t, specificHandled)
			require.False(t, genericHandled)
		})
	}
}

func TestHandleEventControls(t *testing.T) {
	tests := []struct {
		name    string
		message *sse.Event
	}{
		{name: "Nil"},
		{name: "Keepalive", message: &sse.Event{}},
		{name: "UnknownTopic", message: &sse.Event{Event: []byte("unknown")}},
		{
			name: "NoHandlers",
			message: &sse.Event{
				Event: []byte("head"),
				Data:  eventData(t, json.RawMessage(`{"slot":"231192","block":"0xbe36e714a6114cf718e35dafc4ac530ce8f01e4a9a360e78098eb129772dcc39","state":"0x61099b2c1dee0104c93ce0e14e5f5fc4b6faceff4cb863278d055bdfb73b7dc7","epoch_transition":false,"previous_duty_dependent_root":"0xa692c095bbca3eeaf99eeabada78874c028c02b176ccf691f3e8fa075d67f5c6","current_duty_dependent_root":"0x92c6b763f610d5941d2041906007bf9449d37772aacf0483a76275ac27c096b4"}`), false),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				(&Service{}).handleEvent(context.Background(), test.message, &api.EventsOpts{})
			})
		})
	}
}

func TestUnmarshalVersionedEventData(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		expected  *apiv1.ExecutionPayloadEvent
		wantError bool
	}{
		{
			name:     "Wrapped",
			input:    eventData(t, &apiv1.ExecutionPayloadEvent{Slot: 1, BuilderIndex: 2}, true),
			expected: &apiv1.ExecutionPayloadEvent{Slot: 1, BuilderIndex: 2},
		},
		{
			name:     "Bare",
			input:    eventData(t, &apiv1.ExecutionPayloadEvent{Slot: 3, BuilderIndex: 4}, false),
			expected: &apiv1.ExecutionPayloadEvent{Slot: 3, BuilderIndex: 4},
		},
		{
			name:      "Malformed",
			input:     []byte(`invalid`),
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var event apiv1.ExecutionPayloadEvent
			err := unmarshalVersionedEventData(test.input, &event)
			if test.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, &event)
			}
		})
	}
}

func testAttestationData() *phase0.AttestationData {
	return &phase0.AttestationData{
		Source: &phase0.Checkpoint{},
		Target: &phase0.Checkpoint{},
	}
}

func testIndexedAttestation() *electra.IndexedAttestation {
	return &electra.IndexedAttestation{
		AttestingIndices: []uint64{1},
		Data:             testAttestationData(),
	}
}

func eventData(t *testing.T, data any, versioned bool) []byte {
	t.Helper()

	payload, err := json.Marshal(data)
	require.NoError(t, err)
	if !versioned {
		return payload
	}

	wrapped, err := json.Marshal(struct {
		Version string          `json:"version"`
		Data    json.RawMessage `json:"data"`
	}{
		Version: "gloas",
		Data:    payload,
	})
	require.NoError(t, err)

	return wrapped
}
