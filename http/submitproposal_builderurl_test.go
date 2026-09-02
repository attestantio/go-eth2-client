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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	eth2http "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestSubmitProposalEchoesBuilderURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/syncing":
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"data":{"head_slot":"0","sync_distance":"0","is_optimistic":false,"is_syncing":false}}`))
			require.NoError(t, err)
		case "/eth/v2/beacon/blocks":
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "https://builder.example", r.Header.Get("Eth-Builder-Url"))
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	service, err := eth2http.New(context.Background(),
		eth2http.WithAddress(server.URL),
		eth2http.WithAllowDelayedStart(true),
		eth2http.WithEnforceJSON(true),
		eth2http.WithLogLevel(zerolog.Disabled),
	)
	require.NoError(t, err)

	err = service.(interface {
		SubmitProposal(context.Context, *api.SubmitProposalOpts) error
	}).SubmitProposal(context.Background(), &api.SubmitProposalOpts{
		Proposal: &api.VersionedSignedProposal{
			Version: spec.DataVersionGloas,
			Gloas:   &gloas.SignedBeaconBlock{Message: &gloas.BeaconBlock{}},
		},
		BuilderURL: "https://builder.example",
	})
	require.NoError(t, err)
}
