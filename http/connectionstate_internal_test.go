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
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// connectionStateNode is a fake beacon node that answers the endpoints used to check connection
// state.
type connectionStateNode struct {
	// versionFails makes /eth/v1/node/version fail.
	versionFails atomic.Bool
	// syncingStalls makes /eth/v1/node/syncing wait until the request is cancelled.
	syncingStalls atomic.Bool
	// syncingStalled is signalled when a /eth/v1/node/syncing request starts to stall.
	syncingStalled chan struct{}
}

func newConnectionStateNode(t *testing.T) (*connectionStateNode, *httptest.Server) {
	t.Helper()

	node := &connectionStateNode{
		syncingStalled: make(chan struct{}, 10),
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/version":
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
			_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	return node, server
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

func newConnectionStateService(ctx context.Context, t *testing.T, address string, hooks *Hooks) *Service {
	t.Helper()

	service, err := New(ctx,
		WithAddress(address),
		WithTimeout(2*time.Second),
		WithLogLevel(zerolog.Disabled),
		WithAllowDelayedStart(true),
		WithHooks(hooks),
	)
	require.NoError(t, err)
	httpService, isHTTPService := service.(*Service)
	require.True(t, isHTTPService)

	return httpService
}

// TestCheckConnectionStateCallerContextDone confirms that a probe that fails because the
// caller's context is done leaves the connection state as it was.  The probe can run with the
// context of a request that has since completed, which says nothing about the node.
func TestCheckConnectionStateCallerContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node, server := newConnectionStateNode(t)
	recorded, hooks := newRecordedHooks()
	service := newConnectionStateService(ctx, t, server.URL, hooks)
	requireHook(t, recorded.active, "OnActive")
	requireHook(t, recorded.synced, "OnSynced")
	require.True(t, service.IsActive())
	require.True(t, service.IsSynced())

	// The caller's context is cancelled while the probe is in flight.
	node.syncingStalls.Store(true)
	callerCtx, cancelCaller := context.WithCancel(ctx)
	checked := make(chan struct{})
	go func() {
		service.CheckConnectionState(callerCtx)
		close(checked)
	}()
	select {
	case <-node.syncingStalled:
	case <-time.After(time.Second):
		require.FailNow(t, "probe did not reach the node")
	}
	cancelCaller()
	select {
	case <-checked:
	case <-time.After(time.Second):
		require.FailNow(t, "probe did not return")
	}

	require.True(t, service.IsActive())
	require.True(t, service.IsSynced())
	requireNoHook(t, recorded, 200*time.Millisecond)

	// A probe with a live context that fails still marks the node inactive.
	node.syncingStalls.Store(false)
	server.Close()
	service.CheckConnectionState(ctx)
	require.False(t, service.IsActive())
	require.False(t, service.IsSynced())
	requireHook(t, recorded.inactive, "OnInactive")
	requireHook(t, recorded.desynced, "OnDesynced")
}

// TestCheckConnectionStateDVTCheckFailure confirms that a client that fails the DVT check when
// it becomes active is neither active nor synced, so that both hooks fire once the check passes.
func TestCheckConnectionStateDVTCheckFailure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node, server := newConnectionStateNode(t)
	node.versionFails.Store(true)
	recorded, hooks := newRecordedHooks()
	service := newConnectionStateService(ctx, t, server.URL, hooks)

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
