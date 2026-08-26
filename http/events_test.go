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

package http_test

import (
	"context"
	"sync"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

func TestEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name   string
		topics []string
	}{
		{
			name:   "Good",
			topics: []string{"head", "chain_reorg"},
		},
	}

	service := testService(ctx, t).(client.Service)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			eventsMu := sync.Mutex{}
			events := 0
			err := service.(client.EventsProvider).Events(ctx, &api.EventsOpts{
				Topics: test.topics,
				Handler: func(*apiv1.Event) {
					eventsMu.Lock()
					events++
					eventsMu.Unlock()
				},
			})
			require.NoError(t, err)
			time.Sleep(30 * time.Second)
			eventsMu.Lock()
			defer eventsMu.Unlock()
			require.NotEqual(t, 0, events)
			cancel()
		})
	}
}

// gloasEventTopics are the SSE topics introduced by the Gloas (ePBS) fork.
var gloasEventTopics = []string{
	"execution_payload",
	"execution_payload_available",
	"execution_payload_bid",
	"execution_payload_gossip",
	"fast_confirmation",
	"payload_attestation_message",
	"proposer_preferences",
}

// TestEventsGloasTopics confirms the public Events() entry point accepts a
// subscription naming any of the Gloas event topics.
//
// Events() gates every requested topic on apiv1.SupportedEventTopics, via
// checkEventsOpts, before it builds the subscription, so a topic absent from
// that allow-list is refused here and never reaches the wire, however
// completely the dispatcher downstream handles it. Any test that calls
// handleEvent directly sits downstream of this check and so cannot exercise
// it, which leaves this test the only coverage of the gate itself.
//
// Every topic is subscribed to with its own topic-specific handler, which must
// be accepted outright — that also pins each topic to the EventsOpts field
// checkEventSpecificHandler consults for it, which a topic-agnostic generic
// handler would not. One topic is additionally subscribed to with no handler at
// all: "no handler for <topic> event" is only reachable past the allow-list
// check, so it is direct evidence of clearing the gate rather than of merely
// getting a nil error afterwards.
func TestEventsGloasTopics(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.EventsProvider)

	tests := []struct {
		name string
		opts *api.EventsOpts
		err  string
	}{
		{
			// The gate's two rejections are near-identical in wording, so this
			// asserts the exact message: "unsupported event topic X" is the
			// allow-list check's own return, while "no handler for X event"
			// can only be reached once that check has passed.
			name: "ExecutionPayloadNoHandler",
			opts: &api.EventsOpts{Topics: []string{"execution_payload"}},
			err:  "no handler for execution_payload event",
		},
		{
			name: "ExecutionPayloadSpecificHandler",
			opts: &api.EventsOpts{
				Topics:                  []string{"execution_payload"},
				ExecutionPayloadHandler: func(context.Context, *apiv1.ExecutionPayloadEvent) {},
			},
		},
		{
			name: "ExecutionPayloadAvailableSpecificHandler",
			opts: &api.EventsOpts{
				Topics:                           []string{"execution_payload_available"},
				ExecutionPayloadAvailableHandler: func(context.Context, *apiv1.ExecutionPayloadAvailableEvent) {},
			},
		},
		{
			name: "ExecutionPayloadBidSpecificHandler",
			opts: &api.EventsOpts{
				Topics:                     []string{"execution_payload_bid"},
				ExecutionPayloadBidHandler: func(context.Context, *gloas.SignedExecutionPayloadBid) {},
			},
		},
		{
			name: "ExecutionPayloadGossipSpecificHandler",
			opts: &api.EventsOpts{
				Topics:                        []string{"execution_payload_gossip"},
				ExecutionPayloadGossipHandler: func(context.Context, *apiv1.ExecutionPayloadEvent) {},
			},
		},
		{
			name: "FastConfirmationSpecificHandler",
			opts: &api.EventsOpts{
				Topics:                  []string{"fast_confirmation"},
				FastConfirmationHandler: func(context.Context, *apiv1.FastConfirmationEvent) {},
			},
		},
		{
			name: "PayloadAttestationMessageSpecificHandler",
			opts: &api.EventsOpts{
				Topics:                           []string{"payload_attestation_message"},
				PayloadAttestationMessageHandler: func(context.Context, *gloas.PayloadAttestationMessage) {},
			},
		},
		{
			name: "ProposerPreferencesSpecificHandler",
			opts: &api.EventsOpts{
				Topics:                     []string{"proposer_preferences"},
				ProposerPreferencesHandler: func(context.Context, *gloas.SignedProposerPreferences) {},
			},
		},
		{
			name: "AllGloasTopicsGenericHandler",
			opts: &api.EventsOpts{
				Topics:  gloasEventTopics,
				Handler: func(*apiv1.Event) {},
			},
		},
		{
			// Control: shows what a topic the allow-list does not carry looks
			// like, so the assertions above cannot pass vacuously.
			name: "UnknownTopic",
			opts: &api.EventsOpts{
				Topics:  []string{"not_an_event_topic"},
				Handler: func(*apiv1.Event) {},
			},
			err: "unsupported event topic not_an_event_topic",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// A fresh context per subtest, so that any subscription Events()
			// does accept is torn down as soon as the subtest ends.
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			err := service.Events(ctx, test.opts)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
