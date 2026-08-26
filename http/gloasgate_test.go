package http_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/testclients"
	"github.com/stretchr/testify/require"
)

type gateT interface {
	Helper()
	Skip(args ...any)
	Fatalf(format string, args ...any)
}

func requireKnowsGloas(ctx context.Context, t gateT, service any) {
	t.Helper()
	if testclients.KnowsGloas(ctx, service) {
		return
	}
	// KnowsGloas cannot tell a node without the fork from a node it failed to
	// ask, so the reason must not claim the former.
	gate(t, "node does not report knowing gloas (no GLOAS_FORK_EPOCH in its spec, or the spec could not be read)")
}

func requireOnGloas(ctx context.Context, t gateT, service any) {
	t.Helper()
	if testclients.OnGloas(ctx, service) {
		return
	}
	gate(t, fmt.Sprintf("chain is not on gloas (head is %s)", testclients.HeadVersion(ctx, service)))
}

func gate(t gateT, reason string) {
	t.Helper()
	if requireGloas {
		t.Fatalf("HTTP_REQUIRE_GLOAS is set, so this test must not be skipped: %s", reason)
		return
	}
	t.Skip(reason)
}

type recordingT struct {
	skipped string
	failed  string
}

func (r *recordingT) Helper()                           {}
func (r *recordingT) Skip(args ...any)                  { r.skipped = fmt.Sprint(args...) }
func (r *recordingT) Fatalf(format string, args ...any) { r.failed = fmt.Sprintf(format, args...) }

func TestGateSkipsOrFails(t *testing.T) {
	tests := []struct {
		name    string
		require bool
		skipped string
		failed  string
	}{
		{name: "SkipsWhenNotRequired", skipped: "chain is not on gloas (head is fulu)"},
		{name: "FailsWhenRequired", require: true, failed: "HTTP_REQUIRE_GLOAS is set, so this test must not be skipped: chain is not on gloas (head is fulu)"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			old := requireGloas
			t.Cleanup(func() { requireGloas = old })
			requireGloas = test.require
			r := &recordingT{}
			gate(r, "chain is not on gloas (head is fulu)")
			require.Equal(t, test.skipped, r.skipped)
			require.Equal(t, test.failed, r.failed)
		})
	}
}

func TestGloasGatesNameTheirLevel(t *testing.T) {
	old := requireGloas
	t.Cleanup(func() { requireGloas = old })
	requireGloas = false
	ctx := context.Background()

	knows := &recordingT{}
	requireKnowsGloas(ctx, knows, "not a service")
	require.Equal(t,
		"node does not report knowing gloas (no GLOAS_FORK_EPOCH in its spec, or the spec could not be read)",
		knows.skipped)

	on := &recordingT{}
	requireOnGloas(ctx, on, "not a service")
	require.Equal(t, "chain is not on gloas (head is unknown)", on.skipped)
}
