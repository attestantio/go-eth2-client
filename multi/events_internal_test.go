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

package multi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestEventsDoesNotWriteServiceClientLists confirms that an active client failing to subscribe
// is queued for retry on a list of Events' own.  activateClient and deactivateClient leave the
// service's lists with spare capacity, so appending to a copy of the slice header would write
// into the service's backing array outside the lock, racing with those functions' own appends.
func TestEventsDoesNotWriteServiceClientLists(t *testing.T) {
	// The failing client is retried until the context is done, so end it with the test.  The
	// mocks are given a context of their own, as ending theirs has them race on mock's logger.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	failing, err := mock.New(context.Background(), mock.WithName("failing"))
	require.NoError(t, err)
	failing.EventsFunc = func(context.Context, *api.EventsOpts) error {
		return errors.New("failed to subscribe")
	}

	inactive, err := mock.New(context.Background(), mock.WithName("inactive"))
	require.NoError(t, err)

	inactiveClients := make([]consensusclient.Service, 1, 2)
	inactiveClients[0] = inactive

	s := &Service{
		log:                 zerolog.Nop(),
		activeClients:       []consensusclient.Service{failing},
		inactiveClients:     inactiveClients,
		eventsRetryInterval: time.Millisecond,
	}

	require.NoError(t, s.Events(ctx, &api.EventsOpts{
		Topics:  []string{"head"},
		Handler: func(*apiv1.Event) {},
	}))

	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()
	require.Nil(t, s.inactiveClients[:cap(s.inactiveClients)][1], "Events wrote into the service's inactive client list")
}

// TestEventsRetriesWithDefaultIntervalWhenUnset confirms that a Service whose retry interval is
// unset, as one not built by New, waits between attempts rather than retrying without pause.
func TestEventsRetriesWithDefaultIntervalWhenUnset(t *testing.T) {
	// As in TestEventsDoesNotWriteServiceClientLists, the mock has a context of its own.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var attempts atomic.Int32
	failing, err := mock.New(context.Background(), mock.WithName("failing"))
	require.NoError(t, err)
	failing.EventsFunc = func(context.Context, *api.EventsOpts) error {
		attempts.Add(1)

		return errors.New("failed to subscribe")
	}

	s := &Service{
		log:             zerolog.Nop(),
		inactiveClients: []consensusclient.Service{failing},
	}

	require.NoError(t, s.Events(ctx, &api.EventsOpts{
		Topics:  []string{"head"},
		Handler: func(*apiv1.Event) {},
	}))

	time.Sleep(50 * time.Millisecond)
	require.LessOrEqual(t, attempts.Load(), int32(1), "retried without waiting for an interval")
}

// syncingOnlyClient is a client that reports its sync state but does not provide events.
type syncingOnlyClient struct{}

func (syncingOnlyClient) Name() string    { return "syncing only" }
func (syncingOnlyClient) Address() string { return "syncing only" }
func (syncingOnlyClient) IsActive() bool  { return true }
func (syncingOnlyClient) IsSynced() bool  { return true }

func (syncingOnlyClient) NodeSyncing(context.Context, *api.NodeSyncingOpts) (*api.Response[*apiv1.SyncState], error) {
	return &api.Response[*apiv1.SyncState]{Data: &apiv1.SyncState{}}, nil
}

// TestEventsSkipsClientsWithoutEvents confirms that a client that does not provide events, be it
// active or awaiting retry, is skipped rather than panicking Events, and that Events fails when
// that leaves no client to provide them.
func TestEventsSkipsClientsWithoutEvents(t *testing.T) {
	provider, err := mock.New(context.Background(), mock.WithName("provider"))
	require.NoError(t, err)
	provider.EventsFunc = func(context.Context, *api.EventsOpts) error { return nil }

	tests := []struct {
		name            string
		activeClients   []consensusclient.Service
		inactiveClients []consensusclient.Service
		err             string
	}{
		{
			name:            "NoneProvide",
			activeClients:   []consensusclient.Service{syncingOnlyClient{}},
			inactiveClients: []consensusclient.Service{syncingOnlyClient{}},
			err:             "no client can provide events",
		},
		{
			name:            "OtherActiveProvides",
			activeClients:   []consensusclient.Service{syncingOnlyClient{}, provider},
			inactiveClients: []consensusclient.Service{syncingOnlyClient{}},
		},
		{
			name:            "OtherDeferredProvides",
			activeClients:   []consensusclient.Service{syncingOnlyClient{}},
			inactiveClients: []consensusclient.Service{syncingOnlyClient{}, provider},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := &Service{
				log:                 zerolog.Nop(),
				activeClients:       test.activeClients,
				inactiveClients:     test.inactiveClients,
				eventsRetryInterval: time.Millisecond,
			}

			var err error
			require.NotPanics(t, func() {
				err = s.Events(ctx, &api.EventsOpts{
					Topics:  []string{"head"},
					Handler: func(*apiv1.Event) {},
				})
			})
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRetryInterval(t *testing.T) {
	tests := []struct {
		name     string
		interval time.Duration
		expected time.Duration
	}{
		{name: "Set", interval: time.Millisecond, expected: time.Millisecond},
		{name: "Zero", interval: 0, expected: defaultEventsRetryInterval},
		{name: "Negative", interval: -time.Second, expected: defaultEventsRetryInterval},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, (&Service{eventsRetryInterval: test.interval}).retryInterval())
		})
	}
}

// eventsOnlyClient is a client that provides events, failing to subscribe, but does not report its
// sync state.
type eventsOnlyClient struct{}

func (eventsOnlyClient) Name() string    { return "events only" }
func (eventsOnlyClient) Address() string { return "events only" }
func (eventsOnlyClient) IsActive() bool  { return true }
func (eventsOnlyClient) IsSynced() bool  { return true }

func (eventsOnlyClient) Events(context.Context, *api.EventsOpts) error {
	return errors.New("subscription refused")
}

// lockedBuffer is a buffer that can be written by one goroutine while read by another.
type lockedBuffer struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buffer.Write(p)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buffer.String()
}

// TestEventsLogsWhyClientIsDropped confirms that each client Events does not subscribe is logged
// once, saying whether it will be retried and, if not, what it lacks.
func TestEventsLogsWhyClientIsDropped(t *testing.T) {
	failing, err := mock.New(context.Background(), mock.WithName("failing"))
	require.NoError(t, err)
	failing.EventsFunc = func(context.Context, *api.EventsOpts) error {
		return errors.New("subscription refused")
	}

	provider, err := mock.New(context.Background(), mock.WithName("provider"))
	require.NoError(t, err)
	provider.EventsFunc = func(context.Context, *api.EventsOpts) error { return nil }

	tests := []struct {
		name            string
		activeClients   []consensusclient.Service
		inactiveClients []consensusclient.Service
		err             string
		expected        string
	}{
		{
			name:          "ActiveFailsDeferrable",
			activeClients: []consensusclient.Service{failing},
			expected:      `{"level":"warn","address":"failing","topics":["head"],"error":"subscription refused","message":"Failed to set up events handler; will retry"}`,
		},
		{
			name:          "ActiveFailsNotDeferrable",
			activeClients: []consensusclient.Service{provider, eventsOnlyClient{}},
			expected:      `{"level":"error","address":"events only","topics":["head"],"error":"subscription refused","message":"Failed to set up events handler; not a node syncing provider, so will not retry"}`,
		},
		{
			name:          "OnlyClientFailsNotDeferrable",
			activeClients: []consensusclient.Service{eventsOnlyClient{}},
			err:           "no client can provide events",
			expected:      `{"level":"error","address":"events only","topics":["head"],"error":"subscription refused","message":"Failed to set up events handler; not a node syncing provider, so will not retry"}`,
		},
		{
			name:            "InactiveNotEventsProvider",
			activeClients:   []consensusclient.Service{provider},
			inactiveClients: []consensusclient.Service{syncingOnlyClient{}},
			expected:        `{"level":"error","address":"syncing only","topics":["head"],"message":"Not an events provider"}`,
		},
		{
			name:            "InactiveNotSyncingProvider",
			activeClients:   []consensusclient.Service{provider},
			inactiveClients: []consensusclient.Service{eventsOnlyClient{}},
			expected:        `{"level":"error","address":"events only","topics":["head"],"message":"Not a node syncing provider; not subscribing to events"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// The deferred client is retried until the context is done, so end it with the test.
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			output := &lockedBuffer{}
			s := &Service{
				log:                 zerolog.New(output).Level(zerolog.WarnLevel),
				activeClients:       test.activeClients,
				inactiveClients:     test.inactiveClients,
				eventsRetryInterval: time.Hour,
			}

			err := s.Events(ctx, &api.EventsOpts{
				Topics:  []string{"head"},
				Handler: func(*apiv1.Event) {},
			})
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
			}

			// Events writes the line before it returns.  A client that failed to subscribe waits
			// an interval before its retry goroutine tries again, so no further line follows.
			require.Never(t, func() bool {
				return strings.Count(output.String(), "\n") > 1
			}, 50*time.Millisecond, 5*time.Millisecond, "expected exactly one log line")
			lines := strings.Split(strings.TrimSpace(output.String()), "\n")
			require.Len(t, lines, 1, "expected exactly one log line")

			// Each line carries a random call id, which is dropped before comparing.
			var logged map[string]any
			require.NoError(t, json.Unmarshal([]byte(lines[0]), &logged))
			delete(logged, "id")
			actual, err := json.Marshal(logged)
			require.NoError(t, err)
			require.JSONEq(t, test.expected, string(actual))
		})
	}
}
