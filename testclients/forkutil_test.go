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

package testclients_test

import (
	"context"
	nethttp "net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/stretchr/testify/require"
)

// TestGloasPredicatesAgainstLiveNode is the only place the predicates meet a real
// node's wire format, which is the one thing no stub can establish: that
// GLOAS_FORK_EPOCH survives the spec parser as a key of that name, and that the
// head block's version really is reported as gloas.
//
// It skips without HTTP_ADDRESS, as TestNetworkName does — the same seam, and the
// reason the predicates take `any` rather than a concrete client type.
func TestGloasPredicatesAgainstLiveNode(t *testing.T) {
	if os.Getenv("HTTP_ADDRESS") == "" {
		t.Skip("HTTP_ADDRESS not set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Custom spec support is not incidental here: OnGloas reads the head block, so
	// the service has to be able to decode one, and against a non-mainnet preset
	// the compiled-in codec cannot — it fails on the offsets. The http suite builds
	// its shared service the same way for the same reason.
	service, err := http.New(ctx,
		http.WithAddress(os.Getenv("HTTP_ADDRESS")),
		http.WithCustomSpecSupport(true),
		http.WithAllowDelayedStart(true),
	)
	require.NoError(t, err)

	require.True(t, testclients.KnowsGloas(ctx, service),
		"GLOAS_FORK_EPOCH absent from the node's spec")
	require.True(t, testclients.OnGloas(ctx, service),
		"the node's head is not a gloas block")
	require.Equal(t, "gloas", testclients.HeadVersion(ctx, service))
}

// TestKnowsGloasReadsTheSpec covers the case the validation devnet cannot: a node
// that answers, and whose answer does not mention Gloas. Only a node without the
// fork tells "absent" apart from "present", so the spec comes from a stub rather
// than from hunting down a pre-Gloas endpoint to point at.
//
// That direction is the one that matters. A predicate stuck at false keeps its
// tests dark forever with nothing to surface it; one stuck at true fails loudly on
// the first run.
//
// The present row is not padding — it is what makes the absent row mean anything.
// KnowsGloas answers false both for a spec without the key and for a spec it never
// managed to read, and the stub is reached through a service whose start-up ping
// this stub does not serve. Without a row that reaches the same stub and comes back
// true, an absent-row pass is equally consistent with the spec fetch failing
// outright, which would assert nothing about the key at all.
func TestKnowsGloasReadsTheSpec(t *testing.T) {
	tests := []struct {
		name       string
		forkEpochs string
		expected   bool
	}{
		{
			name:       "GloasAbsent",
			forkEpochs: `"ELECTRA_FORK_EPOCH":"0","FULU_FORK_EPOCH":"1"`,
			expected:   false,
		},
		{
			// Far future on purpose: a client that knows the fork publishes an
			// epoch long before reaching it, and its endpoints exist to refuse
			// pre-fork requests from that moment, so this must read as knowing.
			name:       "GloasFarFuture",
			forkEpochs: `"FULU_FORK_EPOCH":"1","GLOAS_FORK_EPOCH":"99999999"`,
			expected:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			service := stubSpecService(ctx, t, test.forkEpochs)

			require.Equal(t, test.expected, testclients.KnowsGloas(ctx, service))
		})
	}
}

// TestGloasPredicatesDegrade pins what the predicates do when there is nothing to
// ask.  Both answer "no", and HeadVersion names the fork "unknown" rather than
// returning an empty string, because a caller puts it straight into the message it
// prints when declining to run and "head is " reads as a bug in the harness.
//
// Answering rather than failing is the deliberate part: these run on the gating
// endpoint, where a node that cannot be reached must leave the fork-dependent tests
// out rather than take the suite down with it.
func TestGloasPredicatesDegrade(t *testing.T) {
	tests := []struct {
		name    string
		service func(context.Context, *testing.T) any
	}{
		{
			// Not a beacon client at all: the predicates take `any`, so this is
			// a reachable mistake rather than a hypothetical one.
			name: "NotAProvider",
			service: func(_ context.Context, _ *testing.T) any {
				return "http://a-string-is-not-a-service"
			},
		},
		{
			// A real service whose node is not listening, which is what a stale
			// address or a reassigned devnet port looks like from here.
			name: "UnreachableNode",
			service: func(ctx context.Context, t *testing.T) any {
				t.Helper()

				service, err := http.New(ctx,
					http.WithAddress("http://127.0.0.1:1"),
					http.WithAllowDelayedStart(true),
				)
				require.NoError(t, err)

				return service
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			service := test.service(ctx, t)

			require.False(t, testclients.KnowsGloas(ctx, service))
			require.False(t, testclients.OnGloas(ctx, service))
			require.Equal(t, "unknown", testclients.HeadVersion(ctx, service))
		})
	}
}

// stubSpecService returns a service whose spec is the given fork epochs and nothing
// else, which is all the predicate reads.  SLOTS_PER_EPOCH is included because the
// client's spec parser expects it to be there.
//
// The sync endpoint is served as well, and has to be: it is what New pings on the
// way up, and a service that fails that ping answers every later request with
// "client is not active" instead of calling out at all.  A stub serving only the
// spec makes KnowsGloas return false whatever the spec says, which is precisely
// what the present row above exists to catch.
func stubSpecService(ctx context.Context, t *testing.T, forkEpochs string) *http.Service {
	t.Helper()

	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		// Both of these are what New asks for on the way up, and both have to
		// answer for the service to come up active.
		case "/eth/v1/node/syncing":
			_, _ = w.Write([]byte(`{"data":{"head_slot":"32","sync_distance":"0",` +
				`"is_syncing":false,"is_optimistic":false,"el_offline":false}}`))
		case "/eth/v1/node/version":
			_, _ = w.Write([]byte(`{"data":{"version":"stub/v0.0.0"}}`))
		case "/eth/v1/config/spec":
			_, _ = w.Write([]byte(`{"data":{"SLOTS_PER_EPOCH":"32",` + forkEpochs + `}}`))
		default:
			w.WriteHeader(nethttp.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)

	service, err := http.New(ctx,
		http.WithAddress(server.URL),
		http.WithAllowDelayedStart(true),
	)
	require.NoError(t, err)

	return service.(*http.Service)
}
