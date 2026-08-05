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
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/stretchr/testify/require"
)

// requireGloas turns every fork gate below from a skip into a failure.  It is set
// from HTTP_REQUIRE_GLOAS in TestMain, beside the other HTTP_* knobs.
//
// It exists because the two endpoints this suite runs against want opposite things.
// On the regression endpoint, which is not on Gloas and will not be for years, a
// skip is the correct outcome forever and asserting on it would be noise.  On the
// validation devnet a skip is always a bug — a stale address, a predicate typo, an
// enclave port a docker restart moved — and without this the job degrades to an
// all-skip run and reports green having verified nothing.
var requireGloas bool

// gateT is the part of *testing.T that a fork gate uses.
//
// It is an interface so that the choice below can be tested.  The real Skip and
// Fatalf both end the calling goroutine, so a test cannot observe what was done to
// it; a recorder standing in for *testing.T can.
type gateT interface {
	Helper()
	Skip(args ...any)
	Fatalf(format string, args ...any)
}

// requireKnowsGloas leaves the calling test out unless the node implements the ePBS
// endpoints, which is to say unless its configuration mentions the fork at all.
//
// This is the weaker of the two levels and the right one for a test of pre-Gloas
// behaviour — asking an endpoint to refuse a pre-fork slot needs the endpoint to
// exist, not the chain to have reached the fork.  Gating such a test on the head's
// version instead would keep it dark on the regression endpoint forever, where the
// stronger level can never be satisfied.
func requireKnowsGloas(ctx context.Context, t gateT, service any) {
	t.Helper()

	if testclients.KnowsGloas(ctx, service) {
		return
	}

	gate(t, "node does not know gloas (GLOAS_FORK_EPOCH absent)")
}

// requireOnGloas leaves the calling test out unless the node's chain head is a Gloas
// block, which is what a test operating on head needs.
//
// The head's fork is fetched a second time to name it in the message.  That only
// happens on the path where the test is not going to run, and it buys the one thing
// the message is for: a reader can tell a node that has never heard of Gloas from
// one that knows it and has not reached it.
func requireOnGloas(ctx context.Context, t gateT, service any) {
	t.Helper()

	if testclients.OnGloas(ctx, service) {
		return
	}

	gate(t, fmt.Sprintf("chain is not on gloas (head is %s)", testclients.HeadVersion(ctx, service)))
}

// gate declines to run the calling test, or fails it if a Gloas node was required.
//
// reason names which of the two levels was not met, and is the whole diagnostic:
// on the gating endpoint it is the only visibility anyone has into that node's
// Gloas posture.
func gate(t gateT, reason string) {
	t.Helper()

	if requireGloas {
		t.Fatalf("HTTP_REQUIRE_GLOAS is set, so this test must not be skipped: %s", reason)

		return
	}

	t.Skip(reason)
}

// recordingT stands in for *testing.T so that what a gate does to the test it is
// gating can be asserted instead of happening.
type recordingT struct {
	skipped string
	failed  string
}

func (r *recordingT) Helper() {}

func (r *recordingT) Skip(args ...any) {
	r.skipped = fmt.Sprint(args...)
}

func (r *recordingT) Fatalf(format string, args ...any) {
	r.failed = fmt.Sprintf(format, args...)
}

// TestGateSkipsOrFails is the guard on the one property this whole mechanism rests
// on: HTTP_REQUIRE_GLOAS turns a skip into a failure.
//
// Its value is in what it forbids rather than what it proves.  A gated test that
// skips where it should have run is invisible, which is the failure mode the
// inversion exists to catch; leaving that to a maintainer noticing a SKIP line is
// how the unconditional skip in attestationrewards_test.go went unnoticed for an
// unknown period.  If a later change softens the failure arm to a warning, or lets
// both arms fire, this is what stops it.
func TestGateSkipsOrFails(t *testing.T) {
	tests := []struct {
		name    string
		require bool
		skipped string
		failed  string
	}{
		{
			// The regression endpoint: a skip here is correct forever.
			name:    "SkipsWhenNotRequired",
			require: false,
			skipped: "chain is not on gloas (head is fulu)",
		},
		{
			// The validation devnet: a skip here is always a bug, so it is not
			// available as an outcome.
			name:    "FailsWhenRequired",
			require: true,
			failed: "HTTP_REQUIRE_GLOAS is set, so this test must not be skipped: " +
				"chain is not on gloas (head is fulu)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := requireGloas
			t.Cleanup(func() { requireGloas = restore })
			requireGloas = test.require

			recorder := &recordingT{}

			gate(recorder, "chain is not on gloas (head is fulu)")

			require.Equal(t, test.skipped, recorder.skipped, "skip")
			require.Equal(t, test.failed, recorder.failed, "failure")
		})
	}
}

// TestGloasGatesNameTheirLevel pins that the two gates are told apart by their
// message alone, which is what makes a bucket-B test lighting up on the regression
// endpoint observable rather than inferred.
func TestGloasGatesNameTheirLevel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	restore := requireGloas
	t.Cleanup(func() { requireGloas = restore })
	requireGloas = false

	// Not a service, so neither question can be answered and both gates decline,
	// each naming its own level.
	notAService := "http://a-string-is-not-a-service"

	knows := &recordingT{}
	requireKnowsGloas(ctx, knows, notAService)
	require.Equal(t, "node does not know gloas (GLOAS_FORK_EPOCH absent)", knows.skipped)

	on := &recordingT{}
	requireOnGloas(ctx, on, notAService)
	require.Equal(t, "chain is not on gloas (head is unknown)", on.skipped)
}
