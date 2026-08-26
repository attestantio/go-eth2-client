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

package http_test

import (
	"context"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

func TestSubmitExecutionPayloadBid(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	tests := []struct {
		name string
		opts *api.SubmitExecutionPayloadBidOpts
		err  string
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name: "NilBid",
			opts: &api.SubmitExecutionPayloadBidOpts{},
			err:  "no bid supplied",
		},
		{
			name: "UnsupportedVersion",
			opts: &api.SubmitExecutionPayloadBidOpts{
				SignedExecutionPayloadBid: &spec.VersionedSignedExecutionPayloadBid{
					Version: spec.DataVersionPhase0,
				},
			},
			err: "unsupported bid version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.(client.ExecutionPayloadBidSubmitter).SubmitExecutionPayloadBid(ctx, test.opts)
			require.ErrorContains(t, err, test.err)
		})
	}

	t.Run("PostsBid", func(t *testing.T) {
		// The submission is recorded here and asserted on the test's own
		// goroutine below.  Asserting inside the handler would call
		// runtime.Goexit from a goroutine that is not the test's, abandoning the
		// request and reporting a wrong header as a transport error instead.
		var (
			received bool
			method   string
			headers  nethttp.Header
			body     []byte
			readErr  error
		)
		server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			switch r.URL.Path {
			case "/eth/v1/node/version":
				_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
			case "/eth/v1/node/syncing":
				_, _ = w.Write([]byte(`{"data":{"is_syncing":false,"is_optimistic":false,"el_offline":false,"head_slot":"1","sync_distance":"0"}}`))
			case "/eth/v1/beacon/execution_payload_bids":
				received = true
				method = r.Method
				headers = r.Header.Clone()
				body, readErr = io.ReadAll(r.Body)
				w.WriteHeader(nethttp.StatusNoContent)
			default:
				w.WriteHeader(nethttp.StatusNotFound)
			}
		}))
		defer server.Close()

		localService, err := http.New(ctx, http.WithAddress(server.URL))
		require.NoError(t, err)
		err = localService.(client.ExecutionPayloadBidSubmitter).SubmitExecutionPayloadBid(ctx, &api.SubmitExecutionPayloadBidOpts{
			SignedExecutionPayloadBid: &spec.VersionedSignedExecutionPayloadBid{
				Version: spec.DataVersionGloas,
				Gloas: &gloas.SignedExecutionPayloadBid{
					Message: &gloas.ExecutionPayloadBid{},
				},
			},
		})
		require.NoError(t, err)
		require.True(t, received)
		require.Equal(t, nethttp.MethodPost, method)
		require.Equal(t, "gloas", headers.Get("Eth-Consensus-Version"))
		require.Equal(t, "application/octet-stream", headers.Get("Content-Type"))
		require.NoError(t, readErr)
		require.NotEmpty(t, body)
	})
}
