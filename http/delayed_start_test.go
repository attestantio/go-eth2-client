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
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
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
	_, err := New(ctx, WithAddress("http://localhost:5052/prefix?key=test"), WithHTTPClient(&http.Client{Transport: transport}), WithAllowDelayedStart(true), WithLogLevel(zerolog.Disabled))
	require.NoError(t, err)
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
