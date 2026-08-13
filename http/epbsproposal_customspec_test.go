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

	bitfield "github.com/OffchainLabs/go-bitfield"
	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	eth2http "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestEPBSProposalCustomSpec verifies that an ePBS proposal decoded against a
// custom (minimal) preset is accepted and yields the correct roots.
//
// Generated HashTreeRoot methods inline mainnet preset sizes (a 512-bit sync
// committee bitvector here) as Go literals and never consult the request's
// codec, so on a preset with a smaller SYNC_COMMITTEE_SIZE they compute a
// wrong body root.  The response's execution payload envelope carries the
// correct root, computed by this test with the same request-scoped dynssz
// codec the client uses, via header composition rather than the block's own
// generated HashTreeRoot.  Before the fix, the transport compared that
// correct envelope root against the generated (wrong) block root and
// rejected every such response as inconsistent.
func TestEPBSProposalCustomSpec(t *testing.T) {
	ctx := context.Background()
	dynamicSSZ := dynssz.NewDynSsz(map[string]any{
		"SYNC_COMMITTEE_SIZE": uint64(32),
	})

	slot := phase0.Slot(123)
	randaoReveal := phase0.BLSSignature{0x01, 0x02}
	graffiti := [32]byte{0x03, 0x04}

	block := &gloas.BeaconBlock{
		Slot:          slot,
		ProposerIndex: 7,
		ParentRoot:    phase0.Root{0x11},
		StateRoot:     phase0.Root{0x22},
		Body: &gloas.BeaconBlockBody{
			RANDAOReveal: randaoReveal,
			Graffiti:     graffiti,
			SyncAggregate: &altair.SyncAggregate{
				// Sized for the minimal preset (32 sync committee members,
				// 4 bytes); the mainnet-sized generated hasher pads this to
				// 64 bytes regardless of what is actually here.
				SyncCommitteeBits: bitfield.Bitvector512(bitfield.NewBitvector32()),
			},
		},
	}

	wantBodyRootRaw, err := dynamicSSZ.HashTreeRoot(block.Body)
	require.NoError(t, err)
	wantBodyRoot := phase0.Root(wantBodyRootRaw)

	generatedBodyRootRaw, err := block.Body.HashTreeRoot()
	require.NoError(t, err)
	generatedBodyRoot := phase0.Root(generatedBodyRootRaw)

	// The fixture is only useful if the generated (mainnet-preset) root
	// actually differs from the minimal-preset one; otherwise this test
	// could pass by accident on either code path.
	require.NotEqual(t, wantBodyRoot, generatedBodyRoot)

	wantBlockRootRaw, err := (&phase0.BeaconBlockHeader{
		Slot:          block.Slot,
		ProposerIndex: block.ProposerIndex,
		ParentRoot:    block.ParentRoot,
		StateRoot:     block.StateRoot,
		BodyRoot:      wantBodyRoot,
	}).HashTreeRoot()
	require.NoError(t, err)
	wantBlockRoot := phase0.Root(wantBlockRootRaw)

	contents := &apiv1gloas.BlockContents{
		Block: block,
		ExecutionPayloadEnvelope: &gloas.ExecutionPayloadEnvelope{
			BeaconBlockRoot: wantBlockRoot,
		},
	}

	encoded, err := dynamicSSZ.MarshalSSZ(contents)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/syncing":
			require.Equal(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"data":{"head_slot":"0","sync_distance":"0","is_optimistic":false,"is_syncing":false}}`))
			require.NoError(t, err)
		case "/eth/v1/node/version":
			require.Equal(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"data":{"version":"test"}}`))
			require.NoError(t, err)
		case "/eth/v1/config/spec":
			require.Equal(t, http.MethodGet, r.Method)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{"data":{"SYNC_COMMITTEE_SIZE":"32"}}`))
			require.NoError(t, err)
		case "/eth/v4/validator/blocks/123":
			require.Equal(t, http.MethodGet, r.Method)
			require.Contains(t, r.Header.Get("Accept"), "application/octet-stream")
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Eth-Consensus-Version", "gloas")
			w.Header().Set("Eth-Execution-Payload-Included", "true")
			_, err := w.Write(encoded)
			require.NoError(t, err)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	service, err := eth2http.New(ctx,
		eth2http.WithAddress(server.URL),
		eth2http.WithCustomSpecSupport(true),
		eth2http.WithLogLevel(zerolog.Disabled),
	)
	require.NoError(t, err)

	includePayload := true
	response, err := service.(client.EPBSProposalProvider).EPBSProposal(ctx, &api.EPBSProposalOpts{
		Slot:           slot,
		RandaoReveal:   randaoReveal,
		Graffiti:       graffiti,
		IncludePayload: &includePayload,
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.True(t, response.Data.ExecutionPayloadIncluded)

	gotBodyRoot, err := response.Data.BodyRoot()
	require.NoError(t, err)
	require.Equal(t, wantBodyRoot, gotBodyRoot)
	require.NotEqual(t, generatedBodyRoot, gotBodyRoot)

	gotRoot, err := response.Data.Root()
	require.NoError(t, err)
	require.Equal(t, wantBlockRoot, gotRoot)
}
