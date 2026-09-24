// Copyright © 2021 - 2026 Attestant Limited.
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

package multi_test

import (
	"context"
	"errors"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/multi"
	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestEventsRejectsMissingHandler(t *testing.T) {
	ctx := context.Background()
	client, subscribed := mockCapturingEvents(ctx, t, "mock 1")

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client}),
	)
	require.NoError(t, err)

	err = multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{Topics: []string{"head"}})
	require.ErrorIs(t, err, consensusclient.ErrInvalidOptions)
	require.ErrorContains(t, err, "no handler for head event")
	require.Nil(t, subscribed(), "invalid options were passed to a client")
}

func TestEventsRejectsUnknownTopic(t *testing.T) {
	ctx := context.Background()
	client, subscribed := mockCapturingEvents(ctx, t, "mock 1")

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client}),
	)
	require.NoError(t, err)

	err = multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{
		Topics:  []string{"head", "typo"},
		Handler: func(*apiv1.Event) {},
	})
	require.ErrorIs(t, err, consensusclient.ErrInvalidOptions)
	require.ErrorContains(t, err, "unsupported event topic typo")
	require.Nil(t, subscribed(), "invalid options were passed to a client")
}

func TestEventsRejectsEmptyTopics(t *testing.T) {
	ctx := context.Background()
	client, subscribed := mockCapturingEvents(ctx, t, "mock 1")

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client}),
	)
	require.NoError(t, err)

	err = multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{})
	require.ErrorIs(t, err, consensusclient.ErrInvalidOptions)
	require.Nil(t, subscribed(), "invalid options were passed to a client")
}

func TestEvents(t *testing.T) {
	ctx := context.Background()

	client1, err := mock.New(ctx, mock.WithName("mock 1"))
	require.NoError(t, err)
	erroringClient1, err := testclients.NewErroring(ctx, 0.1, client1)
	require.NoError(t, err)
	client2, err := mock.New(ctx, mock.WithName("mock 2"))
	require.NoError(t, err)
	erroringClient2, err := testclients.NewErroring(ctx, 0.1, client2)
	require.NoError(t, err)
	client3, err := mock.New(ctx, mock.WithName("mock 3"))
	require.NoError(t, err)

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{
			erroringClient1,
			erroringClient2,
			client3,
		}),
	)
	require.NoError(t, err)

	require.NoError(t, multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{
		Topics:       []string{"block"},
		BlockHandler: func(context.Context, *apiv1.BlockEvent) {},
	}))
}

// TestEventsForwardsGenericHandler confirms that an event delivered by an active client
// reaches the caller's own generic handler, exactly once.
func TestEventsForwardsGenericHandler(t *testing.T) {
	ctx := context.Background()

	client, clientOpts := mockCapturingEvents(ctx, t, "mock 1")

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client}),
	)
	require.NoError(t, err)

	received := 0
	require.NoError(t, multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{
		Topics:  []string{"head"},
		Handler: func(*apiv1.Event) { received++ },
	}))

	require.NotNil(t, clientOpts(), "underlying client was never subscribed")
	require.NotNil(t, clientOpts().Handler, "no generic handler passed to the underlying client")

	// Deliver an event as the underlying client would.
	clientOpts().Handler(&apiv1.Event{Topic: "head"})

	require.Equal(t, 1, received)
}

// TestEventsForwardsTopicHandler confirms that an event delivered by an active client reaches
// the caller's own topic-specific handler, exactly once.
func TestEventsForwardsTopicHandler(t *testing.T) {
	ctx := context.Background()

	client, clientOpts := mockCapturingEvents(ctx, t, "mock 1")

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client}),
	)
	require.NoError(t, err)

	received := 0
	require.NoError(t, multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{
		Topics:       []string{"block"},
		BlockHandler: func(context.Context, *apiv1.BlockEvent) { received++ },
	}))

	require.NotNil(t, clientOpts(), "underlying client was never subscribed")
	require.NotNil(t, clientOpts().BlockHandler, "no block handler passed to the underlying client")

	// Deliver an event as the underlying client would.
	clientOpts().BlockHandler(ctx, &apiv1.BlockEvent{})

	require.Equal(t, 1, received)
}

// TestEventsFiltersNonPrimaryClient confirms that events are forwarded only from the client
// that is currently the primary active one.  Were events from every active client forwarded we
// could end up with inconsistent results, a `head` event arriving from one client while a
// subsequent fetch is answered by another that is behind it.
func TestEventsFiltersNonPrimaryClient(t *testing.T) {
	ctx := context.Background()

	primary, primaryOpts := mockCapturingEvents(ctx, t, "mock 1")
	secondary, secondaryOpts := mockCapturingEvents(ctx, t, "mock 2")

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{primary, secondary}),
	)
	require.NoError(t, err)
	require.Equal(t, "mock 1", multiClient.Address(), "the first client must be the primary")

	received := 0
	require.NoError(t, multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{
		Topics:  []string{"head"},
		Handler: func(*apiv1.Event) { received++ },
	}))

	require.NotNil(t, secondaryOpts(), "the secondary client was never subscribed")
	secondaryOpts().Handler(&apiv1.Event{Topic: "head"})
	require.Equal(t, 0, received, "an event from a non-primary client must not be forwarded")

	require.NotNil(t, primaryOpts(), "the primary client was never subscribed")
	primaryOpts().Handler(&apiv1.Event{Topic: "head"})
	require.Equal(t, 1, received, "an event from the primary client must be forwarded")
}

// TestEventsFiltersClientThatBecomesActive confirms that a client which was not synced when the
// subscription was made is subject to the same active-address filtering as one that was.  Such a
// client is subscribed later, from a goroutine, once it reports itself synced; that subscription
// must not bypass the filter, or events from a client that is not the primary one are forwarded.
func TestEventsFiltersClientThatBecomesActive(t *testing.T) {
	ctx := context.Background()

	primary, _ := mockCapturingEvents(ctx, t, "mock 1")

	late, lateOpts := mockCapturingEvents(ctx, t, "mock 2")
	// Not synced, so it starts out inactive, but reports itself synced when next asked.
	late.SyncDistance = 1
	late.NodeSyncingFunc = func(context.Context, *api.NodeSyncingOpts) (*api.Response[*apiv1.SyncState], error) {
		return &api.Response[*apiv1.SyncState]{Data: &apiv1.SyncState{IsSyncing: false}}, nil
	}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{primary, late}),
	)
	require.NoError(t, err)
	require.Equal(t, "mock 1", multiClient.Address(), "the synced client must be the primary")

	received := 0
	require.NoError(t, multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{
		Topics:  []string{"head"},
		Handler: func(*apiv1.Event) { received++ },
	}))

	require.Eventually(t, func() bool { return lateOpts() != nil }, 5*time.Second, 10*time.Millisecond,
		"the client that became active was never subscribed")

	lateOpts().Handler(&apiv1.Event{Topic: "head"})
	require.Equal(t, 0, received, "an event from a non-primary client must not be forwarded")
}

// TestEventsPassesOptionsToClient confirms the caller's topics and common options reach the
// underlying client.  A client rejects a subscription naming no topics, so an active client
// handed options stripped of them does not subscribe at all: it reports a failure, and is
// dropped from the active list.
func TestEventsPassesOptionsToClient(t *testing.T) {
	ctx := context.Background()

	client, clientOpts := mockCapturingEvents(ctx, t, "mock 1")

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client}),
	)
	require.NoError(t, err)

	require.NoError(t, multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{
		Common:  api.CommonOpts{Timeout: 17 * time.Second},
		Topics:  []string{"head", "block"},
		Handler: func(*apiv1.Event) {},
	}))

	require.NotNil(t, clientOpts(), "underlying client was never subscribed")
	require.Equal(t, []string{"head", "block"}, clientOpts().Topics)
	require.Equal(t, 17*time.Second, clientOpts().Common.Timeout, "common options did not reach the client")
}

// TestEventsOwnsTopics confirms that a subscription retains the requested topics even if the
// caller later reuses its options slice.
func TestEventsOwnsTopics(t *testing.T) {
	ctx := context.Background()

	client, clientOpts := mockCapturingEvents(ctx, t, "mock 1")

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{client}),
	)
	require.NoError(t, err)

	topics := []string{"head", "block"}
	require.NoError(t, multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{
		Topics:  topics,
		Handler: func(*apiv1.Event) {},
	}))

	topics[0] = "attestation"

	require.NotNil(t, clientOpts(), "underlying client was never subscribed")
	require.Equal(t, []string{"head", "block"}, clientOpts().Topics)
}

// TestEventsOwnsTopicsForDeferredClients confirms that a client subscribed later, once it
// reports itself synced, retains the requested topics when the caller reuses the options slice
// after Events returns.
func TestEventsOwnsTopicsForDeferredClients(t *testing.T) {
	ctx := context.Background()

	primary, _ := mockCapturingEvents(ctx, t, "mock 1")
	deferred, deferredOpts := mockCapturingEvents(ctx, t, "mock 2")
	deferred.SyncDistance = 1
	syncStarted := make(chan struct{})
	continueSync := make(chan struct{})
	deferred.NodeSyncingFunc = func(context.Context, *api.NodeSyncingOpts) (*api.Response[*apiv1.SyncState], error) {
		close(syncStarted)
		<-continueSync

		return &api.Response[*apiv1.SyncState]{Data: &apiv1.SyncState{IsSyncing: false}}, nil
	}

	multiClient, err := multi.New(ctx,
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{primary, deferred}),
	)
	require.NoError(t, err)

	topics := []string{"head", "block"}
	require.NoError(t, multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{
		Topics:  topics,
		Handler: func(*apiv1.Event) {},
	}))

	<-syncStarted
	topics[0] = "attestation"
	close(continueSync)

	require.Eventually(t, func() bool { return deferredOpts() != nil }, 5*time.Second, 10*time.Millisecond,
		"the deferred client was never subscribed")
	require.Equal(t, []string{"head", "block"}, deferredOpts().Topics)
}

// TestEventsRetriesDeferredClient confirms that a client which was not synced when the
// subscription was made is retried after a failure to check its sync state or to subscribe,
// rather than abandoned.  An abandoned client can still become the primary one, and as events
// are forwarded only from the primary the caller would then receive none at all.
func TestEventsRetriesDeferredClient(t *testing.T) {
	tests := []struct {
		name        string
		syncErrors  int32
		eventErrors int32
	}{
		{
			name:       "SyncStateFailure",
			syncErrors: 1,
		},
		{
			name:        "SubscribeFailure",
			eventErrors: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()

			primary, _ := mockCapturingEvents(ctx, t, "mock 1")
			deferred, deferredOpts := mockCapturingEvents(ctx, t, "mock 2")
			deferred.SyncDistance = 1

			var syncCalls atomic.Int32
			deferred.NodeSyncingFunc = func(context.Context, *api.NodeSyncingOpts) (*api.Response[*apiv1.SyncState], error) {
				if syncCalls.Add(1) <= test.syncErrors {
					return nil, errors.New("failed to obtain sync state")
				}

				return &api.Response[*apiv1.SyncState]{Data: &apiv1.SyncState{IsSyncing: false}}, nil
			}

			subscribe := deferred.EventsFunc
			var eventCalls atomic.Int32
			deferred.EventsFunc = func(ctx context.Context, opts *api.EventsOpts) error {
				if eventCalls.Add(1) <= test.eventErrors {
					return errors.New("failed to subscribe")
				}

				return subscribe(ctx, opts)
			}

			multiClient, err := multi.New(ctx,
				multi.WithLogLevel(zerolog.Disabled),
				multi.WithClients([]consensusclient.Service{primary, deferred}),
			)
			require.NoError(t, err)
			multi.SetEventsRetryInterval(multiClient, time.Millisecond)

			require.NoError(t, multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{
				Topics:  []string{"head"},
				Handler: func(*apiv1.Event) {},
			}))

			require.Eventually(t, func() bool { return deferredOpts() != nil }, 5*time.Second, time.Millisecond,
				"the deferred client was abandoned after a failure")
		})
	}
}

// TestEventsStopsRetryingOnContextDone confirms that a deferred client is retried only for as
// long as the subscription's context lives, so that one which never syncs does not leave a
// goroutine polling it for the life of the process.
func TestEventsStopsRetryingOnContextDone(t *testing.T) {
	primary, _ := mockCapturingEvents(context.Background(), t, "mock 1")
	deferred, _ := mockCapturingEvents(context.Background(), t, "mock 2")
	deferred.SyncDistance = 1

	var syncCalls atomic.Int32
	deferred.NodeSyncingFunc = func(context.Context, *api.NodeSyncingOpts) (*api.Response[*apiv1.SyncState], error) {
		syncCalls.Add(1)

		return nil, errors.New("failed to obtain sync state")
	}

	multiClient, err := multi.New(context.Background(),
		multi.WithLogLevel(zerolog.Disabled),
		multi.WithClients([]consensusclient.Service{primary, deferred}),
	)
	require.NoError(t, err)
	multi.SetEventsRetryInterval(multiClient, time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	require.NoError(t, multiClient.(consensusclient.EventsProvider).Events(ctx, &api.EventsOpts{
		Topics:  []string{"head"},
		Handler: func(*apiv1.Event) {},
	}))

	require.Eventually(t, func() bool { return syncCalls.Load() > 1 }, 5*time.Second, time.Millisecond,
		"the deferred client was not retried")
	cancel()

	// Allow a retry already under way to finish, then confirm no further ones start.
	time.Sleep(50 * time.Millisecond)
	calls := syncCalls.Load()
	time.Sleep(50 * time.Millisecond)
	require.Equal(t, calls, syncCalls.Load(), "the deferred client was retried after its context was done")
}

// TestEventsForwardsEveryHandler confirms that every handler in api.EventsOpts is substituted
// and forwarded back to the caller.  The handler fields are walked reflectively rather than
// named, so that a handler added to api.EventsOpts in the future cannot be silently left
// unwired here: it joins this test automatically, and fails it until it is forwarded.
func TestEventsForwardsEveryHandler(t *testing.T) {
	ctx := context.Background()

	handlers := handlerFields()
	require.NotEmpty(t, handlers, "no handler fields found in api.EventsOpts")

	for _, handler := range handlers {
		t.Run(handler.Name, func(t *testing.T) {
			client, clientOpts := mockCapturingEvents(ctx, t, "mock 1")

			multiClient, err := multi.New(ctx,
				multi.WithLogLevel(zerolog.Disabled),
				multi.WithClients([]consensusclient.Service{client}),
			)
			require.NoError(t, err)

			// Supply this specific handler and a generic handler for the head subscription.
			received := 0
			opts := &api.EventsOpts{
				Topics:  []string{"head"},
				Handler: func(*apiv1.Event) {},
			}
			reflect.ValueOf(opts).Elem().FieldByName(handler.Name).Set(
				reflect.MakeFunc(handler.Type, func([]reflect.Value) []reflect.Value {
					received++

					return nil
				}),
			)

			require.NoError(t, multiClient.(consensusclient.EventsProvider).Events(ctx, opts))

			require.NotNil(t, clientOpts(), "underlying client was never subscribed")
			substituted := reflect.ValueOf(clientOpts()).Elem().FieldByName(handler.Name)
			require.False(t, substituted.IsNil(), "handler not passed to the underlying client")

			// No other specific handler is substituted.  The generic handler makes the head
			// subscription valid even when testing another topic's handler.
			for _, other := range handlers {
				if other.Name == handler.Name || other.Name == "Handler" {
					continue
				}

				unsupplied := reflect.ValueOf(clientOpts()).Elem().FieldByName(other.Name)
				require.True(t, unsupplied.IsNil(), "unsupplied handler %s was substituted", other.Name)
			}

			// Deliver an event as the underlying client would.
			substituted.Call(handlerArgs(ctx, handler.Type))

			require.Equal(t, 1, received, "event did not reach the caller's handler exactly once")
		})
	}
}

// handlerFields returns the handler fields of api.EventsOpts, being those that are functions.
func handlerFields() []reflect.StructField {
	optsType := reflect.TypeFor[api.EventsOpts]()

	fields := make([]reflect.StructField, 0, optsType.NumField())

	for i := range optsType.NumField() {
		if field := optsType.Field(i); field.Type.Kind() == reflect.Func {
			fields = append(fields, field)
		}
	}

	return fields
}

// handlerArgs builds an argument list with which a handler of the given type can be called: the
// context where one is wanted, and a pointer to a zero value elsewhere, given that handlers can
// be expected to dereference the event they are passed.
func handlerArgs(ctx context.Context, handler reflect.Type) []reflect.Value {
	ctxType := reflect.TypeFor[context.Context]()

	args := make([]reflect.Value, 0, handler.NumIn())

	for i := range handler.NumIn() {
		switch arg := handler.In(i); {
		case arg == ctxType:
			args = append(args, reflect.ValueOf(ctx))
		case arg.Kind() == reflect.Pointer:
			args = append(args, reflect.New(arg.Elem()))
		default:
			args = append(args, reflect.Zero(arg))
		}
	}

	return args
}

// mockCapturingEvents returns a synced mock client that captures the options it is handed by
// Events(), along with an accessor returning those options, or nil if it was never subscribed.
// The accessor is safe to call from another goroutine, as a client that starts out inactive is
// subscribed from one.
func mockCapturingEvents(ctx context.Context, t *testing.T, name string) (*mock.Service, func() *api.EventsOpts) {
	t.Helper()

	client, err := mock.New(ctx, mock.WithName(name))
	require.NoError(t, err)

	var captured atomic.Pointer[api.EventsOpts]

	client.EventsFunc = func(_ context.Context, opts *api.EventsOpts) error {
		captured.Store(opts)

		return nil
	}

	return client, captured.Load
}
