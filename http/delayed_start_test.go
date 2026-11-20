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
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

type startupTransport func(*http.Request) (*http.Response, error)

func (f startupTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestDeferredStartDoesNotWaitForConnectionCheck(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	started := make(chan struct{}, 3)
	transport := startupTransport(func(r *http.Request) (*http.Response, error) {
		started <- struct{}{}
		<-r.Context().Done()
		return nil, r.Context().Err()
	})
	services := make(chan []*Service, 1)
	done := make(chan error, 1)
	go func() {
		clients := make([]*Service, 0, 3)
		for range 3 {
			service, err := New(ctx, WithAddress("http://localhost:5052"), WithHTTPClient(&http.Client{Transport: transport}), WithTimeout(time.Minute), WithLogLevel(zerolog.Disabled), WithDeferredStart(true))
			if err != nil {
				done <- err
				return
			}
			clients = append(clients, service.(*Service))
		}
		services <- clients
		done <- nil
	}()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("deferred start waited for a hanging endpoint")
	}
	require.Empty(t, started, "construction must not issue any readiness requests")
	for _, service := range <-services {
		go service.CheckConnectionState(ctx)
	}
	for range 3 {
		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatal("on-demand connection check did not start")
		}
	}
}

func TestSyncingFailureWithQueryDoesNotRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var requests atomic.Int64
	first := make(chan *http.Request, 1)
	transport := startupTransport(func(r *http.Request) (*http.Response, error) {
		n := requests.Add(1)
		if n == 1 {
			first <- r
		}
		if n > 50 {
			return nil, context.Canceled // Bound an unfixed retry loop.
		}
		return nil, errors.New("offline endpoint")
	})
	service, err := New(ctx, WithAddress("http://localhost:5052/prefix?key=test"), WithHTTPClient(&http.Client{Transport: transport}), WithDeferredStart(true), WithLogLevel(zerolog.Disabled))
	require.NoError(t, err)
	_, err = service.(client.NodeSyncingProvider).NodeSyncing(ctx, &api.NodeSyncingOpts{})
	require.Error(t, err)
	select {
	case request := <-first:
		require.Equal(t, "/prefix/eth/v1/node/syncing", request.URL.Path)
		require.Equal(t, "test", request.URL.Query().Get("key"))
	case <-time.After(time.Second):
		t.Fatal("initial connection check did not start")
	}
	time.Sleep(50 * time.Millisecond)
	require.EqualValues(t, 1, requests.Load(), "syncing transport failures must not trigger recursive connection checks")
}

func TestNonDelayedStartStillRequiresActiveClient(t *testing.T) {
	transport := startupTransport(func(_ *http.Request) (*http.Response, error) {
		return nil, errors.New("offline endpoint")
	})
	_, err := New(t.Context(), WithAddress("http://localhost:5052"), WithHTTPClient(&http.Client{Transport: transport}), WithLogLevel(zerolog.Disabled))
	require.ErrorIs(t, err, client.ErrNotActive)
}

func startupJSONResponse(r *http.Request, body string) *http.Response {
	return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}
}

func TestConcurrentFirstUseWaitsForCompleteActivation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	syncStarted, versionStarted := make(chan struct{}, 1), make(chan struct{}, 1)
	releaseSync, releaseVersion := make(chan struct{}), make(chan struct{})
	transport := startupTransport(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/eth/v1/node/syncing":
			syncStarted <- struct{}{}
			select {
			case <-releaseSync:
			case <-r.Context().Done():
				return nil, r.Context().Err()
			}
			return startupJSONResponse(r, `{"data":{"head_slot":"1","sync_distance":"0","is_syncing":false,"is_optimistic":false}}`), nil
		case "/eth/v1/node/version":
			versionStarted <- struct{}{}
			select {
			case <-releaseVersion:
			case <-r.Context().Done():
				return nil, r.Context().Err()
			}
			return startupJSONResponse(r, `{"data":{"version":"Lighthouse/v8.0.0"}}`), nil
		case "/eth/v1/beacon/genesis":
			return startupJSONResponse(r, `{"data":{"genesis_time":"1742213400","genesis_validators_root":"0x210000000000000000000000000000000000000000000000000000000000002f","genesis_fork_version":"0x10000910"}}`), nil
		default:
			return nil, fmt.Errorf("unexpected request: %s", r.URL.Path)
		}
	})
	service, err := New(ctx, WithAddress("http://localhost:5052"), WithHTTPClient(&http.Client{Transport: transport}), WithDeferredStart(true), WithLogLevel(zerolog.Disabled))
	require.NoError(t, err)
	provider := service.(client.GenesisProvider)
	first, second := make(chan error, 1), make(chan error, 1)
	go func() { _, err := provider.Genesis(ctx, &api.GenesisOpts{}); first <- err }()
	select {
	case <-syncStarted:
	case <-time.After(time.Second):
		t.Fatal("syncing check did not start")
	}
	go func() { _, err := provider.Genesis(ctx, &api.GenesisOpts{}); second <- err }()
	select {
	case err := <-second:
		t.Fatalf("concurrent reader did not wait for syncing: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	close(releaseSync)
	select {
	case <-versionStarted:
	case <-time.After(time.Second):
		t.Fatal("version activation check did not start")
	}
	select {
	case err := <-second:
		t.Fatalf("concurrent reader did not wait for version activation: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	waitCtx, stopWaiting := context.WithCancel(ctx)
	canceled := make(chan error, 1)
	go func() { _, err := provider.Genesis(waitCtx, &api.GenesisOpts{}); canceled <- err }()
	select {
	case err := <-canceled:
		t.Fatalf("waiter returned before cancellation or activation: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	stopWaiting()
	select {
	case err := <-canceled:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("activation waiter ignored cancellation")
	}
	close(releaseVersion)
	for _, result := range []chan error{first, second} {
		select {
		case err := <-result:
			require.NoError(t, err, "both readers must see completed activation")
		case <-time.After(time.Second):
			t.Fatal("reader remained blocked after activation")
		}
	}
}

func TestActivationVersionFailureDoesNotRecursivelyRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var requests atomic.Int64
	transport := startupTransport(func(r *http.Request) (*http.Response, error) {
		n := requests.Add(1)
		if n > 50 {
			return nil, context.Canceled // Bound an unfixed retry chain.
		}
		if r.URL.Path == "/eth/v1/node/syncing" {
			return startupJSONResponse(r, `{"data":{"head_slot":"1","sync_distance":"0","is_syncing":false,"is_optimistic":false}}`), nil
		}
		return nil, errors.New("version endpoint reset the connection")
	})
	service, err := New(ctx, WithAddress("http://localhost:5052"), WithHTTPClient(&http.Client{Transport: transport}), WithDeferredStart(true), WithLogLevel(zerolog.Disabled))
	require.NoError(t, err)
	service.(*Service).CheckConnectionState(ctx)
	require.False(t, service.IsActive())
	require.False(t, service.IsSynced())
	time.Sleep(50 * time.Millisecond)
	require.EqualValues(t, 2, requests.Load(), "activation must not schedule itself after its version request fails")
}
