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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	eth2http "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestEPBSProposalCustomSpec verifies that an ePBS proposal decoded against a
// custom (minimal) preset is accepted and yields the correct roots.
//
// The request wire contract is pinned to beacon-APIs 159622d983a703eb03a8a37bb1edeab7ffc3b6bc:
// its required POST body, fields, static SSZ bounds, query parameters and response headers match
// Prysm 393e9b3c8c7b58bee90da9ae03d0f89988427afb's BuilderConfig proto and handler. Prysm retains
// a legacy GET handler, but this client deliberately sends only the pinned POST contract.
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
	builderConfig := &gloas.BuilderConfig{
		BuilderBoostFactor: 100,
		Builders:           []*gloas.BuilderEntry{},
	}

	block := &gloas.BeaconBlock{
		Slot:          slot,
		ProposerIndex: 7,
		ParentRoot:    phase0.Root{0x11},
		StateRoot:     phase0.Root{0x22},
		Body: &gloas.BeaconBlockBody{
			RANDAOReveal: randaoReveal,
			Graffiti:     graffiti,
			SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{
				Message: &gloas.ExecutionPayloadBid{BuilderIndex: 3},
			},
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
			Payload: &gloas.ExecutionPayload{
				BaseFeePerGas:   uint256.NewInt(0),
				ExtraData:       []byte{},
				Transactions:    []bellatrix.Transaction{},
				Withdrawals:     []*capella.Withdrawal{},
				BlockAccessList: []byte{},
			},
			ExecutionRequests: &gloas.ExecutionRequests{
				Deposits:        []*electra.DepositRequest{},
				Withdrawals:     []*electra.WithdrawalRequest{},
				Consolidations:  []*electra.ConsolidationRequest{},
				BuilderDeposits: []*gloas.BuilderDepositRequest{},
				BuilderExits:    []*gloas.BuilderExitRequest{},
			},
			BuilderIndex:    3,
			BeaconBlockRoot: wantBlockRoot,
		},
	}

	executionRequestsRoot, err := dynamicSSZ.HashTreeRoot(contents.ExecutionPayloadEnvelope.ExecutionRequests)
	require.NoError(t, err)
	block.Body.SignedExecutionPayloadBid.Message.ExecutionRequestsRoot = phase0.Root(executionRequestsRoot)
	wantBodyRootRaw, err = dynamicSSZ.HashTreeRoot(block.Body)
	require.NoError(t, err)
	wantBodyRoot = phase0.Root(wantBodyRootRaw)
	wantBlockRootRaw, err = (&phase0.BeaconBlockHeader{
		Slot:          block.Slot,
		ProposerIndex: block.ProposerIndex,
		ParentRoot:    block.ParentRoot,
		StateRoot:     block.StateRoot,
		BodyRoot:      wantBodyRoot,
	}).HashTreeRoot()
	require.NoError(t, err)
	wantBlockRoot = phase0.Root(wantBlockRootRaw)
	contents.ExecutionPayloadEnvelope.BeaconBlockRoot = wantBlockRoot

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
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "gloas", r.Header.Get("Eth-Consensus-Version"))
			require.Equal(t, "true", r.URL.Query().Get("include_payload"))
			require.Equal(t, fmt.Sprintf("%#x", randaoReveal), r.URL.Query().Get("randao_reveal"))
			require.Equal(t, fmt.Sprintf("%#x", graffiti), r.URL.Query().Get("graffiti"))
			require.Contains(t, r.Header.Get("Content-Type"), "application/octet-stream")
			body, readErr := io.ReadAll(r.Body)
			require.NoError(t, readErr)
			var gotBuilderConfig gloas.BuilderConfig
			require.NoError(t, gotBuilderConfig.UnmarshalSSZ(body))
			require.Equal(t, builderConfig, &gotBuilderConfig)
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Eth-Consensus-Version", "gloas")
			w.Header().Set("Eth-Execution-Payload-Included", "true")
			w.Header().Set("Eth-Builder-Url", "https://builder.example")
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
		BuilderConfig:  builderConfig,
		IncludePayload: &includePayload,
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, "https://builder.example", response.Metadata["Eth-Builder-Url"])
	require.NotNil(t, response.Data.BuilderIndex)
	require.Equal(t, gloas.BuilderIndex(3), *response.Data.BuilderIndex)
	require.Nil(t, response.Data.ExecutionValue)
	require.True(t, response.Data.ExecutionPayloadIncluded)

	gotBodyRoot, err := response.Data.BodyRoot()
	require.NoError(t, err)
	require.Equal(t, wantBodyRoot, gotBodyRoot)
	require.NotEqual(t, generatedBodyRoot, gotBodyRoot)

	gotRoot, err := response.Data.Root()
	require.NoError(t, err)
	require.Equal(t, wantBlockRoot, gotRoot)
}
