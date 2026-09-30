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

package mock_test

import (
	"context"
	"encoding/json"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestEPBSProposal verifies that the mock's default proposal follows the
// payload-inclusion option and that what it hands back can actually be used.
//
// The marshaling assertions are a smoke check that the mock's output is
// serialisable at all, which is easy to break by leaving one of the many
// pointers a gloas block body holds unset.  They do not cover BaseFeePerGas:
// the payload's JSON marshaler substitutes a zero for a nil one, so a nil there
// serialises quietly rather than failing.  That is asserted directly instead,
// where a consumer reading the field rather than the JSON would trip over it.
func TestEPBSProposal(t *testing.T) {
	ctx := context.Background()

	service, err := mock.New(ctx)
	require.NoError(t, err)

	includePayload := true
	excludePayload := false

	t.Run("PayloadIncluded", func(t *testing.T) {
		response, err := service.EPBSProposal(ctx, &api.EPBSProposalOpts{
			Slot:           123,
			RandaoReveal:   phase0.BLSSignature{0x01, 0x02},
			IncludePayload: &includePayload,
			BuilderConfig:  &gloas.BuilderConfig{},
		})
		require.NoError(t, err)
		require.Equal(t, spec.DataVersionGloas, response.Data.Version)
		require.True(t, response.Data.ExecutionPayloadIncluded)
		require.NotNil(t, response.Data.GloasContents)
		require.Nil(t, response.Data.Gloas)

		// The proposal must be consistent with what was asked for, since a
		// caller's own checks reject it otherwise.
		slot, err := response.Data.Slot()
		require.NoError(t, err)
		require.Equal(t, phase0.Slot(123), slot)

		reveal, err := response.Data.RandaoReveal()
		require.NoError(t, err)
		require.Equal(t, phase0.BLSSignature{0x01, 0x02}, reveal)

		// Everything needed to publish from any node must be reachable.
		envelope, err := response.Data.ExecutionPayloadEnvelope()
		require.NoError(t, err)
		require.NotNil(t, envelope.Payload)
		require.NotNil(t, envelope.Payload.BaseFeePerGas)

		// The mock serves mainnet-preset data, so the generated
		// Body.HashTreeRoot() is the correct root here; BodyRoot() must
		// return exactly that rather than erroring for want of a retained
		// root.
		wantBodyRoot, err := response.Data.GloasContents.Block.Body.HashTreeRoot()
		require.NoError(t, err)
		bodyRoot, err := response.Data.BodyRoot()
		require.NoError(t, err)
		require.Equal(t, phase0.Root(wantBodyRoot), bodyRoot)

		_, err = response.Data.Blobs()
		require.NoError(t, err)
		_, err = response.Data.KZGProofs()
		require.NoError(t, err)

		var marshaled []byte
		require.NotPanics(t, func() {
			marshaled, err = json.Marshal(response.Data.GloasContents)
		})
		require.NoError(t, err)
		require.NotEmpty(t, marshaled)
	})

	t.Run("PayloadExcluded", func(t *testing.T) {
		response, err := service.EPBSProposal(ctx, &api.EPBSProposalOpts{
			Slot:           123,
			IncludePayload: &excludePayload,
			BuilderConfig:  &gloas.BuilderConfig{},
		})
		require.NoError(t, err)
		require.False(t, response.Data.ExecutionPayloadIncluded)
		require.NotNil(t, response.Data.Gloas)
		require.Nil(t, response.Data.GloasContents)

		// There is no envelope to hand over in this mode; the caller has to
		// fetch it from the node that produced the block.
		_, err = response.Data.ExecutionPayloadEnvelope()
		require.ErrorContains(t, err, "the execution payload was not included")

		wantBodyRoot, err := response.Data.Gloas.Body.HashTreeRoot()
		require.NoError(t, err)
		bodyRoot, err := response.Data.BodyRoot()
		require.NoError(t, err)
		require.Equal(t, phase0.Root(wantBodyRoot), bodyRoot)

		var marshaled []byte
		require.NotPanics(t, func() {
			marshaled, err = json.Marshal(response.Data.Gloas)
		})
		require.NoError(t, err)
		require.NotEmpty(t, marshaled)
	})

	t.Run("OverrideFunc", func(t *testing.T) {
		mockService, err := mock.New(ctx)
		require.NoError(t, err)

		mockService.EPBSProposalFunc = func(context.Context, *api.EPBSProposalOpts) (
			*api.Response[*api.VersionedEPBSProposal], error,
		) {
			return &api.Response[*api.VersionedEPBSProposal]{
				Data: &api.VersionedEPBSProposal{Version: spec.DataVersionGloas},
			}, nil
		}

		response, err := mockService.EPBSProposal(ctx, &api.EPBSProposalOpts{
			Slot:           123,
			IncludePayload: &includePayload,
		})
		require.NoError(t, err)
		require.True(t, response.Data.IsEmpty())
	})
}

// TestExecutionPayloadEnvelope verifies that the mock's default envelope echoes
// the requested block root, which is what a caller checks to be sure the
// envelope belongs to the block it proposed.
func TestExecutionPayloadEnvelope(t *testing.T) {
	ctx := context.Background()

	service, err := mock.New(ctx)
	require.NoError(t, err)

	root := phase0.Root{0xaa, 0xbb}

	response, err := service.ExecutionPayloadEnvelope(ctx, &api.ExecutionPayloadEnvelopeOpts{
		Slot:            123,
		BeaconBlockRoot: root,
	})
	require.NoError(t, err)
	require.Equal(t, spec.DataVersionGloas, response.Data.Version)
	require.False(t, response.Data.IsEmpty())

	got, err := response.Data.BeaconBlockRoot()
	require.NoError(t, err)
	require.Equal(t, root, got)

	payload, err := response.Data.Payload()
	require.NoError(t, err)
	require.NotNil(t, payload.BaseFeePerGas)

	var marshaled []byte
	require.NotPanics(t, func() {
		marshaled, err = json.Marshal(response.Data.Gloas)
	})
	require.NoError(t, err)
	require.NotEmpty(t, marshaled)
}

// TestEPBSProposalMirrorsTheHTTPPreconditions keeps the double honest.  A mock
// that accepts options the real client refuses lets a downstream project's unit
// tests pass against a call that fails against a node.
func TestEPBSProposalMirrorsTheHTTPPreconditions(t *testing.T) {
	ctx := context.Background()

	service, err := mock.New(ctx)
	require.NoError(t, err)

	includePayload := true

	tests := []struct {
		name string
		opts *api.EPBSProposalOpts
		err  error
		msg  string
	}{
		{
			name: "NoOptions",
			err:  client.ErrNoOptions,
		},
		{
			name: "NoSlot",
			opts: &api.EPBSProposalOpts{IncludePayload: &includePayload, BuilderConfig: &gloas.BuilderConfig{}},
			err:  client.ErrInvalidOptions,
			msg:  "no slot specified",
		},
		{
			name: "NoPayloadInclusion",
			opts: &api.EPBSProposalOpts{Slot: 123, BuilderConfig: &gloas.BuilderConfig{}},
			err:  client.ErrInvalidOptions,
			msg:  "no payload inclusion specified",
		},
		{
			name: "NoBuilderConfig",
			opts: &api.EPBSProposalOpts{Slot: 123, IncludePayload: &includePayload},
			err:  client.ErrInvalidOptions,
			msg:  "no builder config supplied",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := service.EPBSProposal(ctx, test.opts)
			require.ErrorIs(t, err, test.err)
			if test.msg != "" {
				require.ErrorContains(t, err, test.msg)
			}
		})
	}
}

// TestEPBSProposalPopulatesBuilderIndex covers the other half of the drift: the
// http client always sets BuilderIndex, so code reading it nil-dereferenced
// under the mock and not in production.
func TestEPBSProposalPopulatesBuilderIndex(t *testing.T) {
	ctx := context.Background()

	service, err := mock.New(ctx)
	require.NoError(t, err)

	for _, included := range []bool{true, false} {
		includePayload := included
		response, err := service.EPBSProposal(ctx, &api.EPBSProposalOpts{
			Slot:           123,
			IncludePayload: &includePayload,
			BuilderConfig:  &gloas.BuilderConfig{},
		})
		require.NoError(t, err)
		require.NotNil(t, response.Data.BuilderIndex)
	}
}
