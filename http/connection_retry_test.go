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

type connectionRetryTransport func(*http.Request) (*http.Response, error)

func (f connectionRetryTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestSyncingFailureWithQueryDoesNotRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var requests atomic.Int64
	transport := connectionRetryTransport(func(r *http.Request) (*http.Response, error) {
		n := requests.Add(1)
		require.Equal(t, "/prefix/eth/v1/node/syncing", r.URL.Path)
		require.Equal(t, "test", r.URL.Query().Get("key"))
		if n > 50 {
			return nil, context.Canceled // Bound an unfixed retry loop.
		}
		return nil, errors.New("offline endpoint")
	})
	service, err := New(ctx, WithAddress("http://localhost:5052/prefix?key=test"), WithHTTPClient(&http.Client{Transport: transport}), WithAllowDelayedStart(true), WithLogLevel(zerolog.Disabled))
	require.NoError(t, err)
	_, err = service.(client.NodeSyncingProvider).NodeSyncing(ctx, &api.NodeSyncingOpts{})
	require.Error(t, err)

	time.Sleep(50 * time.Millisecond)
	require.EqualValues(t, 2, requests.Load(), "constructor and explicit syncing requests must not schedule retries")
}

func TestActivationVersionFailureDoesNotRecursivelyRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	var requests atomic.Int64
	transport := connectionRetryTransport(func(r *http.Request) (*http.Response, error) {
		n := requests.Add(1)
		if n > 50 {
			return nil, context.Canceled // Bound an unfixed retry loop.
		}
		if r.URL.Path == "/eth/v1/node/syncing" {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": {"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"data":{"head_slot":"1","sync_distance":"0","is_syncing":false,"is_optimistic":false}}`)),
				Request:    r,
			}, nil
		}
		return nil, errors.New("version endpoint reset the connection")
	})
	service, err := New(ctx, WithAddress("http://localhost:5052"), WithHTTPClient(&http.Client{Transport: transport}), WithAllowDelayedStart(true), WithLogLevel(zerolog.Disabled))
	require.NoError(t, err)
	require.False(t, service.IsActive())

	time.Sleep(50 * time.Millisecond)
	require.EqualValues(t, 2, requests.Load(), "activation must not schedule itself after its version request fails")
}
