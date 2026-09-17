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
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	eth2http "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestSignedBeaconBlockCustomSpecGloas(t *testing.T) {
	ctx := context.Background()
	dynamicSSZ := dynssz.NewDynSsz(map[string]any{
		"SYNC_COMMITTEE_SIZE": uint64(32),
	})
	block := &gloas.SignedBeaconBlock{
		Message: &gloas.BeaconBlock{
			Body: &gloas.BeaconBlockBody{},
		},
	}
	encoded, err := dynamicSSZ.MarshalSSZ(block)
	require.NoError(t, err)
	require.Equal(t, uint32(100), binary.LittleEndian.Uint32(encoded[0:4]))
	require.Equal(t, uint32(84), binary.LittleEndian.Uint32(encoded[180:184]))
	require.Equal(t, uint32(336), binary.LittleEndian.Uint32(encoded[384:388]))

	var staticBlock gloas.SignedBeaconBlock
	require.EqualError(t, staticBlock.UnmarshalSSZ(encoded), "Message.Body.ProposerSlashings:o: incorrect offset: first offset 336 does not match expected 396")

	var mu sync.Mutex
	requests := make([]string, 0, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests = append(requests, r.Method+" "+r.URL.Path)
		mu.Unlock()

		switch r.URL.Path {
		case "/eth/v1/node/syncing":
			require.Equal(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"data":{"head_slot":"0","sync_distance":"0","is_optimistic":false,"is_syncing":false}}`))
			require.NoError(t, err)
		case "/eth/v1/node/version":
			require.Equal(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"data":{"version":"test"}}`))
			require.NoError(t, err)
		case "/eth/v2/beacon/blocks/head":
			require.Equal(t, http.MethodGet, r.Method)
			require.Contains(t, r.Header.Get("Accept"), "application/octet-stream")
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Eth-Consensus-Version", "gloas")
			_, err := w.Write(encoded)
			require.NoError(t, err)
		case "/eth/v1/config/spec":
			require.Equal(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"data":{"SYNC_COMMITTEE_SIZE":"32"}}`))
			require.NoError(t, err)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	service, err := eth2http.New(ctx,
		eth2http.WithAddress(server.URL),
		eth2http.WithCustomSpecSupport(true),
		eth2http.WithLogLevel(zerolog.Disabled),
	)
	require.NoError(t, err)

	response, err := service.(client.SignedBeaconBlockProvider).SignedBeaconBlock(ctx, &api.SignedBeaconBlockOpts{Block: "head"})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Data.Gloas)
	mu.Lock()
	require.Equal(t, []string{"GET /eth/v1/node/syncing", "GET /eth/v1/node/version", "GET /eth/v2/beacon/blocks/head", "GET /eth/v1/config/spec"}, requests)
	mu.Unlock()
}
