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

// The consistency guards on the two new gloas endpoints are tested directly
// elsewhere, which establishes that they decide correctly but not that anything
// calls them.  A live node cannot establish it either: a node that honours the
// request never disagrees with it, so the guards are unreachable in a passing
// live test and deleting the call sites leaves the whole suite green.
//
// These tests close that by answering from a server that deliberately disagrees.
// A service is built against the live node first, so its spec and genesis are
// real and cached, then its base URL is repointed at the local server for the one
// call under test.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

// disagreeingService returns a service whose requests are answered by handler,
// and which is otherwise configured from the live node.
func disagreeingService(ctx context.Context, t *testing.T, handler http.HandlerFunc) *Service {
	t.Helper()

	// Built against the real node so that the spec fetch, the sync assertion and
	// the SSZ codec are all genuine; only the endpoint under test is faked.
	s := newTestService(ctx, t, true, WithEnforceJSON(true))

	// Warm the spec cache while the real address is still in place, since the
	// fake server does not serve it.
	_, err := s.Spec(ctx, &api.SpecOpts{})
	require.NoError(t, err)

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	base, err := url.Parse(server.URL)
	require.NoError(t, err)
	s.base = base

	return s
}

// TestEPBSProposalRejectsADisagreeingNode verifies that the request-consistency
// guard is reached from EPBSProposal, by asking a node that includes an
// execution payload that was not requested.
func TestEPBSProposalRejectsADisagreeingNode(t *testing.T) {
	ctx := context.Background()

	const slot = 60300

	reveal := phase0.BLSSignature{0x0a, 0x0b, 0x0c}

	// Always answers payload-included, whatever was asked for.  This is the
	// direction that remains a disagreement: a node volunteering a payload nobody
	// asked for changes where the block can be published.  The opposite answer is
	// what an external builder's bid legitimately looks like, so it would no
	// longer prove the guard was reached.
	s := disagreeingService(ctx, t, func(w http.ResponseWriter, _ *http.Request) {
		contents := validEPBSBlockContents()
		contents.Block.Slot = slot
		contents.Block.Body.RANDAOReveal = reveal
		root, err := contents.Block.HashTreeRoot()
		require.NoError(t, err)
		contents.ExecutionPayloadEnvelope.BeaconBlockRoot = root

		data, err := json.Marshal(contents)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Eth-Consensus-Version", "gloas")
		fmt.Fprintf(w, `{"version":"gloas","execution_payload_included":true,"data":%s}`, data)
	})

	includePayload := false

	_, err := s.EPBSProposal(ctx, &api.EPBSProposalOpts{
		Slot:           slot,
		RandaoReveal:   reveal,
		BuilderConfig:  &gloas.BuilderConfig{Builders: []*gloas.BuilderEntry{}},
		IncludePayload: &includePayload,
	})
	require.ErrorIs(t, err, client.ErrInconsistentResult)
	require.ErrorContains(t, err, "execution payload included true; expected false")
}

// TestExecutionPayloadEnvelopeRejectsADisagreeingNode verifies that the block-root
// guard is reached from ExecutionPayloadEnvelope, by asking a node that answers
// with an envelope committed to a different block.
func TestExecutionPayloadEnvelopeRejectsADisagreeingNode(t *testing.T) {
	ctx := context.Background()

	// Always answers with the fixture's own root, whatever was asked for.
	s := disagreeingService(ctx, t, func(w http.ResponseWriter, _ *http.Request) {
		data, err := json.Marshal(validExecutionPayloadEnvelope())
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Eth-Consensus-Version", "gloas")
		fmt.Fprintf(w, `{"version":"gloas","data":%s}`, data)
	})

	_, err := s.ExecutionPayloadEnvelope(ctx, &api.ExecutionPayloadEnvelopeOpts{
		Slot:            60300,
		BeaconBlockRoot: phase0.Root{0xde, 0xad, 0xbe, 0xef},
	})
	require.ErrorIs(t, err, client.ErrInconsistentResult)
	require.ErrorContains(t, err, "execution payload envelope for block")
}
