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
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	eth2http "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

const buildersTestResponse = `{"execution_optimistic":false,"finalized":false,"data":[{"index":"3","status":"active","builder":{"pubkey":"0x863917c482adda67fba8305a8ef2cce39cfcd0d13311c86c0cb5b648102b6a481d0b25931d7cfbb5fa3e21c114ec0951","version":"0","execution_address":"0x67565aa29bcf60c3fdcd0b79a36cd85cc13c7a91","balance":"39550000000","deposit_epoch":"1434","withdrawable_epoch":"18446744073709551615"}}]}`

// buildersTestServer serves the builders endpoint from a local handler,
// recording the request body and answering with the given status and body.
func buildersTestServer(t *testing.T, status int, response string, body *string) client.BuildersProvider {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/syncing":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"head_slot":"0","sync_distance":"0","is_optimistic":false,"is_syncing":false}}`))
		case "/eth/v1/node/version":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"data":{"version":"test"}}`))
		case "/eth/v1/beacon/states/head/builders":
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			data, err := io.ReadAll(r.Body)
			if err != nil || r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			*body = string(data)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(response))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	service, err := eth2http.New(context.Background(),
		eth2http.WithAddress(server.URL),
		eth2http.WithAllowDelayedStart(true),
		eth2http.WithLogLevel(zerolog.Disabled),
	)
	require.NoError(t, err)

	return service.(client.BuildersProvider)
}

// TestBuildersRequest checks the request always carries a JSON object body:
// Prysm rejects an array body with 400.
func TestBuildersRequest(t *testing.T) {
	tests := []struct {
		name string
		opts *api.BuildersOpts
		body string
	}{
		{
			name: "NoFilter",
			opts: &api.BuildersOpts{State: "head"},
			body: `{}`,
		},
		{
			name: "Filtered",
			opts: &api.BuildersOpts{
				State:           "head",
				Indices:         []gloas.BuilderIndex{3},
				PubKeys:         []phase0.BLSPubKey{{0x86}},
				BuilderStatuses: []apiv1.BuilderStatus{apiv1.BuilderStatusActive, apiv1.BuilderStatusExited},
			},
			body: `{"ids":["3","0x860000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"],"statuses":["active","exited"]}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var body string
			provider := buildersTestServer(t, http.StatusOK, buildersTestResponse, &body)

			response, err := provider.Builders(context.Background(), test.opts)
			require.NoError(t, err)
			require.JSONEq(t, test.body, body)
			require.Len(t, response.Data, 1)
			require.Equal(t, gloas.BuilderIndex(3), response.Data[3].Index)
			require.Equal(t, apiv1.BuilderStatusActive, response.Data[3].Status)
			require.Equal(t, phase0.Gwei(39550000000), response.Data[3].Builder.Balance)
			require.Equal(t, false, response.Metadata["finalized"])
		})
	}
}

// TestBuildersErrors checks failures reach the caller as an api.Error with the
// node's status code: 404 when the node has no route (Nimbus), as distinct from
// an empty result, and 400 for a pre-Gloas state.
func TestBuildersErrors(t *testing.T) {
	for _, status := range []int{http.StatusNotFound, http.StatusBadRequest} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			var body string
			provider := buildersTestServer(t, status, `{"code":0,"message":"test"}`, &body)

			response, err := provider.Builders(context.Background(), &api.BuildersOpts{State: "head"})
			require.Nil(t, response)
			var apiErr *api.Error
			require.True(t, errors.As(err, &apiErr))
			require.Equal(t, status, apiErr.StatusCode)
			require.Equal(t, `{}`, body)
		})
	}

	t.Run("Empty", func(t *testing.T) {
		var body string
		provider := buildersTestServer(t, http.StatusOK, `{"execution_optimistic":false,"finalized":false,"data":[]}`, &body)

		response, err := provider.Builders(context.Background(), &api.BuildersOpts{State: "head"})
		require.NoError(t, err)
		require.Empty(t, response.Data)
	})
}

func TestBuildersOpts(t *testing.T) {
	var body string
	provider := buildersTestServer(t, http.StatusOK, buildersTestResponse, &body)

	_, err := provider.Builders(context.Background(), nil)
	require.ErrorIs(t, err, client.ErrNoOptions)

	_, err = provider.Builders(context.Background(), &api.BuildersOpts{})
	require.ErrorIs(t, err, client.ErrInvalidOptions)
	require.Empty(t, body)
}

// TestBuildersGloas runs against a live Gloas node.
func TestBuildersGloas(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)
	requireOnGloas(ctx, t, service)
	provider := service.(client.BuildersProvider)

	all, err := provider.Builders(ctx, &api.BuildersOpts{State: "head"})
	var apiErr *api.Error
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
		// Nimbus 26.10 has no route for this endpoint.
		gate(t, "node does not serve the builders endpoint")
	}
	require.NoError(t, err)
	require.NotEmpty(t, all.Data, "devnet has no builders")

	var first *apiv1.Builder
	for _, builder := range all.Data {
		first = builder
		break
	}

	t.Run("ByIndex", func(t *testing.T) {
		res, err := provider.Builders(ctx, &api.BuildersOpts{State: "head", Indices: []gloas.BuilderIndex{first.Index}})
		require.NoError(t, err)
		require.Len(t, res.Data, 1)
		require.Equal(t, first.Builder.PublicKey, res.Data[first.Index].Builder.PublicKey)
	})

	t.Run("ByPubKey", func(t *testing.T) {
		res, err := provider.Builders(ctx, &api.BuildersOpts{State: "head", PubKeys: []phase0.BLSPubKey{first.Builder.PublicKey}})
		require.NoError(t, err)
		require.Len(t, res.Data, 1)
		require.Contains(t, res.Data, first.Index)
	})

	t.Run("ByStatus", func(t *testing.T) {
		res, err := provider.Builders(ctx, &api.BuildersOpts{State: "head", BuilderStatuses: []apiv1.BuilderStatus{first.Status}})
		require.NoError(t, err)
		require.Contains(t, res.Data, first.Index)
		for _, builder := range res.Data {
			require.Equal(t, first.Status, builder.Status)
		}
	})

	t.Run("PreGloasState", func(t *testing.T) {
		_, err := provider.Builders(ctx, &api.BuildersOpts{State: "genesis"})
		var apiErr *api.Error
		require.True(t, errors.As(err, &apiErr))
		// Teku answers 404 here, as it does not serve the genesis state.
		require.Contains(t, []int{http.StatusBadRequest, http.StatusNotFound}, apiErr.StatusCode)
	})
}
