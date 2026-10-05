package testclients

import (
	"context"
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

type forkutilTestService struct {
	specData map[string]any
	head     phase0.Version
}

func (s *forkutilTestService) Spec(context.Context, *api.SpecOpts) (*api.Response[map[string]any], error) {
	return &api.Response[map[string]any]{Data: s.specData}, nil
}

func (s *forkutilTestService) Fork(context.Context, *api.ForkOpts) (*api.Response[*phase0.Fork], error) {
	return &api.Response[*phase0.Fork]{Data: &phase0.Fork{CurrentVersion: s.head}}, nil
}

// gloasConfig is the shape a node that knows Gloas reports.
func gloasConfig() map[string]any {
	return map[string]any{
		"GLOAS_FORK_EPOCH":   uint64(10),
		"GLOAS_FORK_VERSION": phase0.Version{0x70, 0x00, 0x00, 0x00},
		"FULU_FORK_VERSION":  phase0.Version{0x60, 0x00, 0x00, 0x00},
	}
}

func TestForkPosture(t *testing.T) {
	tests := []struct {
		name       string
		service    any
		knowsGloas bool
		onGloas    bool
		head       string
	}{
		{
			name:       "GloasHeadAndConfig",
			service:    &forkutilTestService{specData: gloasConfig(), head: phase0.Version{0x70, 0x00, 0x00, 0x00}},
			knowsGloas: true,
			onGloas:    true,
			head:       "gloas",
		},
		{
			name:       "PreGloasHeadStillKnows",
			service:    &forkutilTestService{specData: gloasConfig(), head: phase0.Version{0x60, 0x00, 0x00, 0x00}},
			knowsGloas: true,
			onGloas:    false,
			head:       "fulu",
		},
		{
			// A node whose head is on a fork its own config does not name.
			name:    "UnnamedHeadFork",
			service: &forkutilTestService{specData: gloasConfig(), head: phase0.Version{0xff}},
			// GLOAS_FORK_EPOCH is still published, so the config knows Gloas.
			knowsGloas: true,
			head:       "unknown",
		},
		{
			name:    "UnknownService",
			service: "not a service",
			head:    "unknown",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.knowsGloas, KnowsGloas(context.Background(), test.service))
			require.Equal(t, test.onGloas, OnGloas(context.Background(), test.service))
			require.Equal(t, test.head, HeadVersion(context.Background(), test.service))
		})
	}
}
