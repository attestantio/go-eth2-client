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
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// timeout for tests.
var timeout = 60 * time.Second

// newTestService creates a service against HTTP_ADDRESS for this package's internal
// tests, authenticating with HTTP_BEARER_TOKEN when one is set. It is the package http
// counterpart of newTestService in main_test.go, which package http cannot reach: the
// two are separate packages, so every internal test that wants a live service has to
// build one for itself, and this is the one place that does it.
//
// It returns *Service rather than client.Service because an internal test invariably
// wants the concrete type, and asserting for it at each call site established nothing.
// customSpecSupport picks which branch the SSZ codecs take.
//
// params are applied after the defaults, so a caller can override any of them — with
// one trap worth naming, since it is this helper's own defect in miniature:
// WithExtraHeaders replaces the header map rather than merging into it, so a caller
// passing its own would drop the Authorization header again. Nothing does today.
//
// The call sites this replaced were every one package http had, from gopls references
// on New: epbsguardwiring_internal_test.go (disagreeingService),
// epbsproposal_internal_test.go (customSpecService) and events_internal_test.go
// (twice, one arm per bearer-token branch). Two of the four omitted the token. The
// same query now answers with a single hit in package http, the New call below.
func newTestService(ctx context.Context,
	t *testing.T,
	customSpecSupport bool,
	params ...Parameter,
) *Service {
	t.Helper()

	address := os.Getenv("HTTP_ADDRESS")
	if address == "" {
		// Belt only, and worth knowing it is: the one TestMain governing both
		// packages lives in main_test.go and returns without running anything at
		// all when this is unset, so nothing currently reaches this skip. It
		// starts mattering if that is ever fixed.
		t.Skip("HTTP_ADDRESS not set")
	}

	parameters := []Parameter{
		WithTimeout(timeout),
		WithAddress(address),
		WithCustomSpecSupport(customSpecSupport),
	}

	if token := os.Getenv("HTTP_BEARER_TOKEN"); token != "" {
		parameters = append(parameters, WithExtraHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", token),
		}))
	}

	service, err := New(ctx, append(parameters, params...)...)
	require.NoError(t, err)

	return service.(*Service)
}

// TestInternalTestServicesAuthenticate pins the one thing every service this
// package's tests build against the live node must do: send the bearer token.
//
// It is asserted per constructor rather than once on the shared builder because
// what broke was a call site, not the builder — newTestService above lists which
// ones and why. A missing token fails only against a node that requires auth, and
// then only as "client is not active", which names nothing.
//
// The rows are therefore redundant by design while all three constructors route
// through newTestService, since they exercise one header assembly between them.
// They stop being redundant the moment one drifts back off it, which is the whole
// regression here: re-hand-rolling customSpecService fails its own row and leaves
// the other two green. Do not collapse them into a single assertion on the
// builder — that is the test this bug already got past.
func TestInternalTestServicesAuthenticate(t *testing.T) {
	tests := []struct {
		name  string
		build func(context.Context, *testing.T) *Service
	}{
		{
			// Called directly, as the event-handler test does.
			name: "NewTestService",
			build: func(ctx context.Context, t *testing.T) *Service {
				return newTestService(ctx, t, false)
			},
		},
		{
			name:  "CustomSpecService",
			build: customSpecService,
		},
		{
			// Its handler never runs here; only the service it hands back is
			// under inspection.
			name: "DisagreeingService",
			build: func(ctx context.Context, t *testing.T) *Service {
				return disagreeingService(ctx, t, "/never-requested", func(http.ResponseWriter, *http.Request) {})
			},
		},
	}

	// An ambient token is used as-is rather than overridden: against an
	// authenticated node a token of this test's own invention would be rejected
	// by New()'s ping, so the assertion would fail for the wrong reason. One is
	// synthesised only when the suite is running unauthenticated, where the node
	// ignores it.
	token := os.Getenv("HTTP_BEARER_TOKEN")
	if token == "" {
		token = "token-of-this-test's-own"
		t.Setenv("HTTP_BEARER_TOKEN", token)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			s := test.build(ctx, t)

			require.Equal(t, "Bearer "+token, s.extraHeaders["Authorization"])
		})
	}
}

// TestNewTestServiceOmitsAnAbsentBearerToken guards the other arm: no token in the
// environment must mean no Authorization header at all, not a bare "Bearer ".
//
// No node takes part — a closed port plus WithAllowDelayedStart is enough, and it
// avoids emptying HTTP_BEARER_TOKEN against a node that would reject the ping.
func TestNewTestServiceOmitsAnAbsentBearerToken(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Setenv("HTTP_ADDRESS", "http://127.0.0.1:1")
	t.Setenv("HTTP_BEARER_TOKEN", "")

	s := newTestService(ctx, t, false, WithAllowDelayedStart(true))

	require.NotContains(t, s.extraHeaders, "Authorization")
}
