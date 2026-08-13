package testclients

import (
	"context"
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"
)

type forkutilTestService struct {
	specData map[string]any
	version  spec.DataVersion
}

func (s *forkutilTestService) Spec(context.Context, *api.SpecOpts) (*api.Response[map[string]any], error) {
	return &api.Response[map[string]any]{Data: s.specData}, nil
}

func (s *forkutilTestService) SignedBeaconBlock(context.Context, *api.SignedBeaconBlockOpts) (*api.Response[*spec.VersionedSignedBeaconBlock], error) {
	return &api.Response[*spec.VersionedSignedBeaconBlock]{Data: &spec.VersionedSignedBeaconBlock{Version: s.version}}, nil
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
			service:    &forkutilTestService{specData: map[string]any{"GLOAS_FORK_EPOCH": uint64(10)}, version: spec.DataVersionGloas},
			knowsGloas: true,
			onGloas:    true,
			head:       spec.DataVersionGloas.String(),
		},
		{
			name:       "PreGloasHeadStillKnows",
			service:    &forkutilTestService{specData: map[string]any{"GLOAS_FORK_EPOCH": uint64(10)}, version: spec.DataVersionFulu},
			knowsGloas: true,
			onGloas:    false,
			head:       spec.DataVersionFulu.String(),
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
