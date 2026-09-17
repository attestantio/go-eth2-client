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
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// TestEPBSProposalOptsValidation covers the options the endpoint refuses before
// it reaches the network.  The payload-inclusion case is the one that matters
// most: the spec marks the parameter required with no default, and the two modes
// carry different operational constraints, so an unset value has to be rejected
// rather than resolved on the caller's behalf.
func TestEPBSProposalOptsValidation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)
	provider := service.(client.EPBSProposalProvider)

	includePayload := true

	t.Run("NilOpts", func(t *testing.T) {
		_, err := provider.EPBSProposal(ctx, nil)
		require.ErrorIs(t, err, client.ErrNoOptions)
	})

	t.Run("NoSlot", func(t *testing.T) {
		_, err := provider.EPBSProposal(ctx, &api.EPBSProposalOpts{
			IncludePayload: &includePayload,
		})
		require.ErrorIs(t, err, client.ErrInvalidOptions)
		require.ErrorContains(t, err, "no slot specified")
	})

	// Leaving IncludePayload unset must not silently select either mode.  False
	// is the constraining one — the producing node caches the envelope and must
	// also be the publisher — so a Go bool's zero value would quietly commit the
	// caller to stateful operation.
	t.Run("NoIncludePayload", func(t *testing.T) {
		_, err := provider.EPBSProposal(ctx, &api.EPBSProposalOpts{Slot: 1})
		require.ErrorIs(t, err, client.ErrInvalidOptions)
		require.ErrorContains(t, err, "no payload inclusion specified")
	})

	t.Run("SkipRandaoVerificationWithoutInfinity", func(t *testing.T) {
		_, err := provider.EPBSProposal(ctx, &api.EPBSProposalOpts{
			Slot:                   1,
			IncludePayload:         &includePayload,
			SkipRandaoVerification: true,
			RandaoReveal:           phase0.BLSSignature{0x01},
		})
		require.ErrorIs(t, err, client.ErrInvalidOptions)
		require.ErrorContains(t, err, "randao reveal must be point at infinity")
	})
}

// TestEPBSProposal exercises the endpoint against a live node, across both
// payload-inclusion modes and both encodings.
//
// The RANDAO reveal is the point at infinity with skip_randao_verification set,
// which is what lets a test produce a block without holding a validator key: the
// node then returns the reveal it was given, so the endpoint's own consistency
// check still has something to verify.
func TestEPBSProposal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	// Every subtest asks the node to produce a gloas block, so the gate belongs
	// on the whole test rather than on each of them.
	requireOnGloas(ctx, t, service)

	// The node produces a block for the slot it is about to propose, so the
	// target is derived from its head rather than from the wall clock.
	slot := headSlot(ctx, t, service) + 1

	// A service that pins JSON, so the two encodings are both reached.  Custom
	// spec support is on either way: the validation devnet runs the minimal
	// preset, and a gloas block is undecodable against the compiled-in one.
	jsonService, err := newTestService(ctx, true, http.WithEnforceJSON(true))
	require.NoError(t, err)

	tests := []struct {
		name          string
		service       client.Service
		included      bool
		builderConfig *gloas.BuilderConfig
	}{
		{name: "SSZPayloadExcluded", service: service, included: false, builderConfig: localPreferredBuilderConfig()},
		{name: "SSZPayloadIncluded", service: service, included: true, builderConfig: localPreferredBuilderConfig()},
		{name: "JSONPayloadExcluded", service: jsonService, included: false, builderConfig: localPreferredBuilderConfig()},
		{name: "JSONPayloadIncluded", service: jsonService, included: true, builderConfig: localPreferredBuilderConfig()},
		// A populated builders list is a different body on the wire, in both
		// encodings, and nothing an empty one sends reaches it.
		{name: "SSZBuilderRequested", service: service, included: false, builderConfig: builderRequestedConfig(slot)},
		{name: "JSONBuilderRequested", service: jsonService, included: false, builderConfig: builderRequestedConfig(slot)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			includePayload := test.included
			infinity := infinitySignature()

			response, err := test.service.(client.EPBSProposalProvider).EPBSProposal(ctx,
				&api.EPBSProposalOpts{
					Slot:                   slot,
					RandaoReveal:           infinity,
					IncludePayload:         &includePayload,
					SkipRandaoVerification: true,
					BuilderConfig:          test.builderConfig,
				},
			)
			require.NoError(t, err)
			require.Equal(t, spec.DataVersionGloas, response.Data.Version)
			require.Equal(t, test.included, response.Data.ExecutionPayloadIncluded)
			require.False(t, response.Data.IsEmpty())

			proposalSlot, err := response.Data.Slot()
			require.NoError(t, err)
			require.Equal(t, slot, proposalSlot)

			reveal, err := response.Data.RandaoReveal()
			require.NoError(t, err)
			require.Equal(t, infinity, reveal)

			// The envelope, blobs and proofs exist only when the payload
			// travelled with the block.  Asking for them in the other mode is a
			// caller error, not an empty result: the envelope has to be fetched
			// from the node that produced the block.
			envelope, err := response.Data.ExecutionPayloadEnvelope()
			if !test.included {
				require.ErrorContains(t, err, "the execution payload was not included")

				return
			}

			require.NoError(t, err)
			require.NotNil(t, envelope.Payload)
			require.NotNil(t, envelope.Payload.BaseFeePerGas)
		})
	}
}

// localPreferredBuilderConfig asks for a build with no direct builder bids
// solicited.  The contract documents the empty list as exactly that -- "Empty
// means request none, so only p2p bids are considered" -- and a zero boost
// factor expresses no preference among those.
func localPreferredBuilderConfig() *gloas.BuilderConfig {
	return &gloas.BuilderConfig{
		Builders: []*gloas.BuilderEntry{},
	}
}

// builderRequestedConfig solicits a direct bid from one builder.
//
// The bid cannot be won from a test: auth.message carries a signature the
// slot's proposer makes, whose key a test does not hold, so the builder
// declines and the node falls back to a p2p or local build.  Everything up to
// that point is still exercised, and none of it is reachable with an empty
// list: a populated BuilderEntry through both encodings, and the node parsing
// the entry and routing on its URL.
//
// The URL is a refused loopback port rather than a devnet builder so the node's
// dial fails at once instead of holding the request open for a timeout.
// devnet-8 publishes no builder API to point at in any case -- its buildoor
// pairs expose only a beacon and an EL RPC endpoint.
//
// The authorization slot has to be the request's own; validateBuilderConfig
// refuses a mismatch before the request is sent.
func builderRequestedConfig(slot phase0.Slot) *gloas.BuilderConfig {
	return &gloas.BuilderConfig{
		BuilderBoostFactor: 100,
		Builders: []*gloas.BuilderEntry{{
			URL: []byte("http://127.0.0.1:1"),
			Auth: &gloas.SignedBuilderRequestAuth{
				Message: &gloas.BuilderRequestAuth{Data: []byte{0x01}, Slot: slot},
			},
			BuilderPubkeys: []phase0.BLSPubKey{},
		}},
	}
}

// infinitySignature returns the BLS point at infinity, which is what a caller
// passes as the RANDAO reveal when asking the node to skip verifying it.
func infinitySignature() phase0.BLSSignature {
	var signature phase0.BLSSignature
	signature[0] = 0xc0

	return signature
}
