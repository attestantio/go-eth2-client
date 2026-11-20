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
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// connectionStateNode is a fake beacon node that answers the endpoints used to check connection
// state.
type connectionStateNode struct {
	// versionFails makes /eth/v1/node/version fail.
	versionFails atomic.Bool
	// versionStallsOnce makes the next /eth/v1/node/version request wait until it is cancelled.
	versionStallsOnce atomic.Bool
	// versionStalled is signalled when a /eth/v1/node/version request starts to stall.
	versionStalled chan struct{}
	// syncingStatus, if set, is the status returned by /eth/v1/node/syncing with no data.
	syncingStatus atomic.Int32
	// syncingStalls makes /eth/v1/node/syncing wait until the request is cancelled.
	syncingStalls atomic.Bool
	// syncingStalled is signalled when a /eth/v1/node/syncing request starts to stall.
	syncingStalled chan struct{}
}

func newConnectionStateNode(t *testing.T) (*connectionStateNode, *httptest.Server) {
	t.Helper()

	node := &connectionStateNode{
		versionStalled: make(chan struct{}, 10),
		syncingStalled: make(chan struct{}, 10),
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/version":
			if node.versionStallsOnce.CompareAndSwap(true, false) {
				node.versionStalled <- struct{}{}
				<-r.Context().Done()

				return
			}
			if node.versionFails.Load() {
				w.WriteHeader(http.StatusInternalServerError)

				return
			}
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/node/syncing":
			if node.syncingStalls.Load() {
				node.syncingStalled <- struct{}{}
				<-r.Context().Done()

				return
			}
			if status := node.syncingStatus.Load(); status != 0 {
				w.WriteHeader(int(status))

				return
			}
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	return node, server
}

// cancelAfterResponse is a transport that, once armed, cancels a context after it has received
// a complete response, so that the caller's context is done by the time the response is handled.
type cancelAfterResponse struct {
	base   http.RoundTripper
	cancel atomic.Pointer[context.CancelFunc]
}

func (c *cancelAfterResponse) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	cancel := c.cancel.Swap(nil)
	if cancel == nil {
		return resp, nil
	}

	// Read the whole response before cancelling, so that the cancellation cannot affect it.
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewReader(body))
	(*cancel)()

	return resp, nil
}

// recordedHooks records each hook call on a channel per hook.
type recordedHooks struct {
	active   chan struct{}
	inactive chan struct{}
	synced   chan struct{}
	desynced chan struct{}
}

func newRecordedHooks() (*recordedHooks, *Hooks) {
	recorded := &recordedHooks{
		active:   make(chan struct{}, 10),
		inactive: make(chan struct{}, 10),
		synced:   make(chan struct{}, 10),
		desynced: make(chan struct{}, 10),
	}

	return recorded, &Hooks{
		OnActive:   func(_ context.Context, _ *Service) { recorded.active <- struct{}{} },
		OnInactive: func(_ context.Context, _ *Service) { recorded.inactive <- struct{}{} },
		OnSynced:   func(_ context.Context, _ *Service) { recorded.synced <- struct{}{} },
		OnDesynced: func(_ context.Context, _ *Service) { recorded.desynced <- struct{}{} },
	}
}

// requireHook fails unless the hook is called within a second.
func requireHook(t *testing.T, ch chan struct{}, name string) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(time.Second):
		require.FailNow(t, "hook not called", name)
	}
}

// requireNoHook fails if any of the hooks is called within the given duration.
func requireNoHook(t *testing.T, recorded *recordedHooks, duration time.Duration) {
	t.Helper()

	select {
	case <-recorded.active:
		require.FailNow(t, "hook called unexpectedly", "OnActive")
	case <-recorded.inactive:
		require.FailNow(t, "hook called unexpectedly", "OnInactive")
	case <-recorded.synced:
		require.FailNow(t, "hook called unexpectedly", "OnSynced")
	case <-recorded.desynced:
		require.FailNow(t, "hook called unexpectedly", "OnDesynced")
	case <-time.After(duration):
	}
}

// requireSignal fails unless the channel is signalled within a second.
func requireSignal(t *testing.T, ch chan struct{}, msg string) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(time.Second):
		require.FailNow(t, msg)
	}
}

func newConnectionStateService(ctx context.Context,
	t *testing.T,
	address string,
	timeout time.Duration,
	transport http.RoundTripper,
	hooks *Hooks,
) *Service {
	t.Helper()

	service, err := New(ctx,
		WithAddress(address),
		WithTimeout(timeout),
		WithHTTPClient(&http.Client{Transport: transport}),
		WithLogLevel(zerolog.Disabled),
		WithAllowDelayedStart(true),
		WithHooks(hooks),
	)
	require.NoError(t, err)
	httpService, isHTTPService := service.(*Service)
	require.True(t, isHTTPService)

	return httpService
}

// newActiveConnectionStateService returns a service that is active and synced, with the hooks
// for that already received.
func newActiveConnectionStateService(ctx context.Context,
	t *testing.T,
	address string,
	timeout time.Duration,
	transport http.RoundTripper,
	hooks *Hooks,
	recorded *recordedHooks,
) *Service {
	t.Helper()

	service := newConnectionStateService(ctx, t, address, timeout, transport, hooks)
	requireHook(t, recorded.active, "OnActive")
	requireHook(t, recorded.synced, "OnSynced")
	require.True(t, service.IsActive())
	require.True(t, service.IsSynced())

	return service
}

// TestCheckConnectionStateCallerContextDone confirms that a probe that fails because the
// caller's context is done leaves the connection state as it was.  The probe can run with the
// context of a request that has since completed, which says nothing about the node.
func TestCheckConnectionStateCallerContextDone(t *testing.T) {
	tests := []struct {
		name      string
		callerCtx func(ctx context.Context) (context.Context, context.CancelFunc)
		// cancelCaller is true if the test cancels the caller's context once the probe is in
		// flight, rather than leaving it to its deadline.
		cancelCaller bool
	}{
		{
			name:         "Cancelled",
			callerCtx:    context.WithCancel,
			cancelCaller: true,
		},
		{
			name: "Deadline",
			callerCtx: func(ctx context.Context) (context.Context, context.CancelFunc) {
				return context.WithTimeout(ctx, 200*time.Millisecond)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			node, server := newConnectionStateNode(t)
			recorded, hooks := newRecordedHooks()
			service := newActiveConnectionStateService(ctx, t, server.URL, 2*time.Second, http.DefaultTransport, hooks, recorded)

			node.syncingStalls.Store(true)
			callerCtx, cancelCaller := test.callerCtx(ctx)
			defer cancelCaller()
			checked := make(chan struct{})
			go func() {
				service.CheckConnectionState(callerCtx)
				close(checked)
			}()
			requireSignal(t, node.syncingStalled, "probe did not reach the node")
			if test.cancelCaller {
				cancelCaller()
			}
			requireSignal(t, checked, "probe did not return")

			require.True(t, service.IsActive())
			require.True(t, service.IsSynced())
			requireNoHook(t, recorded, 200*time.Millisecond)
		})
	}
}

// TestCheckConnectionStateNodeFailure confirms that a probe that fails for a reason other than
// the caller's context marks the client inactive, even if the caller's context is done by the
// time the failure is handled.
func TestCheckConnectionStateNodeFailure(t *testing.T) {
	tests := []struct {
		name string
		// fail makes the node fail the next probe.  It returns the context to probe with.
		fail func(ctx context.Context, node *connectionStateNode, server *httptest.Server, transport *cancelAfterResponse) context.Context
	}{
		{
			name: "Unreachable",
			fail: func(ctx context.Context, _ *connectionStateNode, server *httptest.Server, _ *cancelAfterResponse) context.Context {
				server.Close()

				return ctx
			},
		},
		{
			name: "ClientTimeout",
			fail: func(ctx context.Context, node *connectionStateNode, _ *httptest.Server, _ *cancelAfterResponse) context.Context {
				node.syncingStalls.Store(true)

				return ctx
			},
		},
		{
			name: "ErrorStatusThenCallerCancelled",
			fail: func(ctx context.Context, node *connectionStateNode, _ *httptest.Server, transport *cancelAfterResponse) context.Context {
				node.syncingStatus.Store(http.StatusServiceUnavailable)
				callerCtx, cancelCaller := context.WithCancel(ctx)
				transport.cancel.Store(&cancelCaller)

				return callerCtx
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			node, server := newConnectionStateNode(t)
			transport := &cancelAfterResponse{base: http.DefaultTransport}
			recorded, hooks := newRecordedHooks()
			service := newActiveConnectionStateService(ctx, t, server.URL, 200*time.Millisecond, transport, hooks, recorded)

			probeCtx := test.fail(ctx, node, server, transport)
			checked := make(chan struct{})
			go func() {
				service.CheckConnectionState(probeCtx)
				close(checked)
			}()
			requireSignal(t, checked, "probe did not return")

			require.False(t, service.IsActive())
			require.False(t, service.IsSynced())
			requireHook(t, recorded.inactive, "OnInactive")
			requireHook(t, recorded.desynced, "OnDesynced")
		})
	}
}

// TestCheckConnectionStateCallerContextDoneDuringDVTCheck confirms that a probe whose DVT check
// is ended by the caller's context does not overwrite the state set by a probe that ran
// meanwhile.
//
// NodeVersion holds its lock while it fetches, so a second probe's DVT check waits for this
// one's request to return, and can then finish before this probe publishes its state.  That
// ordering cannot be forced with two real probes, so the test sets the state that the second
// probe publishes.
func TestCheckConnectionStateCallerContextDoneDuringDVTCheck(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node, server := newConnectionStateNode(t)
	node.versionFails.Store(true)
	recorded, hooks := newRecordedHooks()
	service := newConnectionStateService(ctx, t, server.URL, 2*time.Second, http.DefaultTransport, hooks)
	require.False(t, service.IsActive())
	require.False(t, service.IsSynced())

	// The probe stalls in its DVT check.
	node.versionFails.Store(false)
	node.versionStallsOnce.Store(true)
	callerCtx, cancelCaller := context.WithCancel(ctx)
	defer cancelCaller()
	checked := make(chan struct{})
	go func() {
		service.CheckConnectionState(callerCtx)
		close(checked)
	}()
	requireSignal(t, node.versionStalled, "DVT check did not reach the node")

	// Another probe activates the client meanwhile.
	service.connectionMu.Lock()
	service.connectionActive = true
	service.connectionSynced = true
	service.connectionMu.Unlock()

	// The stalled probe's caller goes away.
	cancelCaller()
	requireSignal(t, checked, "probe did not return")

	require.True(t, service.IsActive())
	require.True(t, service.IsSynced())
	requireNoHook(t, recorded, 200*time.Millisecond)
}

// TestCheckConnectionStateDVTCheckFailure confirms that a client that fails the DVT check when
// it becomes active is neither active nor synced, so that both hooks fire once the check passes.
func TestCheckConnectionStateDVTCheckFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node, server := newConnectionStateNode(t)
	node.versionFails.Store(true)
	recorded, hooks := newRecordedHooks()
	service := newConnectionStateService(ctx, t, server.URL, 2*time.Second, http.DefaultTransport, hooks)

	require.False(t, service.IsActive())
	require.False(t, service.IsSynced())
	requireNoHook(t, recorded, 200*time.Millisecond)

	node.versionFails.Store(false)
	service.CheckConnectionState(ctx)
	require.True(t, service.IsActive())
	require.True(t, service.IsSynced())
	requireHook(t, recorded.active, "OnActive")
	requireHook(t, recorded.synced, "OnSynced")
}

// TestFailedForCaller confirms that a failure is put down to the caller only when the caller's
// context ended the request, and not when the client's own timeout came first.
func TestFailedForCaller(t *testing.T) {
	service := &Service{timeout: 2 * time.Second}

	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	expired, cancelExpired := context.WithDeadline(context.Background(), time.Now().Add(-time.Millisecond))
	defer cancelExpired()
	deadline, _ := expired.Deadline()

	deadlineErr := errors.Join(errors.New("failed to call GET endpoint"), context.DeadlineExceeded)
	tests := []struct {
		name     string
		ctx      context.Context
		started  time.Time
		err      error
		expected bool
	}{
		{
			name:    "NoError",
			ctx:     cancelled,
			started: time.Now(),
		},
		{
			name:    "ClientTimeoutCallerLive",
			ctx:     context.Background(),
			started: time.Now().Add(-3 * time.Second),
			err:     deadlineErr,
		},
		{
			name:     "CallerCancelled",
			ctx:      cancelled,
			started:  time.Now(),
			err:      errors.Join(errors.New("failed to call GET endpoint"), context.Canceled),
			expected: true,
		},
		{
			name:    "ErrorStatusCallerCancelled",
			ctx:     cancelled,
			started: time.Now(),
			err:     &api.Error{StatusCode: http.StatusServiceUnavailable},
		},
		{
			name:     "CallerDeadlineBeforeClientTimeout",
			ctx:      expired,
			started:  deadline.Add(-time.Second),
			err:      deadlineErr,
			expected: true,
		},
		{
			name:     "CallerDeadlineAtClientTimeout",
			ctx:      expired,
			started:  deadline.Add(-2 * time.Second),
			err:      deadlineErr,
			expected: true,
		},
		{
			name:    "ClientTimeoutBeforeCallerDeadline",
			ctx:     expired,
			started: deadline.Add(-3 * time.Second),
			err:     deadlineErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, service.failedForCaller(test.ctx, test.started, test.err))
		})
	}
}
