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

package eventdispatch_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/internal/eventdispatch"
	"github.com/stretchr/testify/require"
)

// TestBindsEveryTopicHandler confirms that each specific handler field of
// api.EventsOpts is bound to exactly one supported topic, and each supported topic to exactly
// one field.  The fields are walked reflectively, so a handler added to api.EventsOpts but not
// bound fails here.
func TestBindsEveryTopicHandler(t *testing.T) {
	optsType := reflect.TypeFor[api.EventsOpts]()

	boundTopics := make(map[string]bool)
	for i := range optsType.NumField() {
		field := optsType.Field(i)
		if field.Type.Kind() != reflect.Func || field.Name == "Handler" {
			continue
		}

		t.Run(field.Name, func(t *testing.T) {
			opts := &api.EventsOpts{}
			reflect.ValueOf(opts).Elem().Field(i).Set(reflect.MakeFunc(field.Type,
				func([]reflect.Value) []reflect.Value { return nil }))

			var topics []string
			for topic := range apiv1.SupportedEventTopics {
				if eventdispatch.HasTopicHandler(opts, topic) {
					topics = append(topics, topic)
				}
			}
			require.Len(t, topics, 1, "handler is not bound to exactly one topic")
			require.NotContains(t, boundTopics, topics[0], "topic bound to more than one handler")
			boundTopics[topics[0]] = true
		})
	}

	require.Equal(t, apiv1.SupportedEventTopics, boundTopics)
}

func TestHandle(t *testing.T) {
	ctx := context.Background()
	data := []byte(`{"slot":"10","block_root":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf"}`)

	t.Run("Specific", func(t *testing.T) {
		var received *apiv1.ExecutionPayloadAvailableEvent
		opts := &api.EventsOpts{
			Handler: func(*apiv1.Event) { require.Fail(t, "generic handler called") },
			ExecutionPayloadAvailableHandler: func(_ context.Context, event *apiv1.ExecutionPayloadAvailableEvent) {
				received = event
			},
		}
		require.NoError(t, eventdispatch.Handle(ctx, opts, "execution_payload_available", data))
		require.NotNil(t, received)
	})

	t.Run("Generic", func(t *testing.T) {
		var received *apiv1.Event
		opts := &api.EventsOpts{Handler: func(event *apiv1.Event) { received = event }}
		require.NoError(t, eventdispatch.Handle(ctx, opts, "execution_payload_available", data))
		require.Equal(t, "execution_payload_available", received.Topic)
		require.IsType(t, &apiv1.ExecutionPayloadAvailableEvent{}, received.Data)
	})

	t.Run("NoHandler", func(t *testing.T) {
		require.NoError(t, eventdispatch.Handle(ctx, &api.EventsOpts{}, "execution_payload_available", data))
	})

	t.Run("UnsupportedTopic", func(t *testing.T) {
		require.EqualError(t, eventdispatch.Handle(ctx, &api.EventsOpts{}, "unknown", data), "unsupported event topic unknown")
	})

	t.Run("Malformed", func(t *testing.T) {
		opts := &api.EventsOpts{Handler: func(*apiv1.Event) { require.Fail(t, "handler called") }}
		require.Error(t, eventdispatch.Handle(ctx, opts, "execution_payload_available", []byte(`invalid`)))
	})
}

func TestFiltered(t *testing.T) {
	ctx := context.Background()
	data := []byte(`{"slot":"10","block_root":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf"}`)

	specific := 0
	generic := 0
	opts := &api.EventsOpts{
		Topics:  []string{"execution_payload_available", "head"},
		Handler: func(*apiv1.Event) { generic++ },
		ExecutionPayloadAvailableHandler: func(context.Context, *apiv1.ExecutionPayloadAvailableEvent) {
			specific++
		},
	}

	forward := false
	var forwardedTopics []string
	filtered := eventdispatch.Filtered(opts, func(topic string) bool {
		forwardedTopics = append(forwardedTopics, topic)

		return forward
	})

	require.NotSame(t, opts, filtered)
	require.Equal(t, opts.Topics, filtered.Topics)
	require.Nil(t, filtered.HeadHandler, "unsupplied handler was filled in")

	require.NoError(t, eventdispatch.Handle(ctx, filtered, "execution_payload_available", data))
	require.Zero(t, specific, "event forwarded although forward refused it")

	forward = true
	require.NoError(t, eventdispatch.Handle(ctx, filtered, "execution_payload_available", data))
	require.Equal(t, 1, specific)
	filtered.Handler(&apiv1.Event{Topic: "head"})
	require.Equal(t, 1, generic)
	require.Equal(t, []string{"execution_payload_available", "execution_payload_available", "head"}, forwardedTopics)

	filtered.Topics[0] = "changed"
	require.Equal(t, "execution_payload_available", opts.Topics[0], "topics shared with the original options")
}
