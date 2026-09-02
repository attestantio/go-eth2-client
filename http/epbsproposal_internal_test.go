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
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"

	bitfield "github.com/OffchainLabs/go-bitfield"
	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1gloas "github.com/attestantio/go-eth2-client/api/v1/gloas"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

// validEPBSBeaconBlock returns a gloas beacon block populated just far enough
// for its own JSON marshaler to run: every pointer the marshaler dereferences
// is non-nil.
func validEPBSBeaconBlock() *gloas.BeaconBlock {
	return &gloas.BeaconBlock{
		Slot:          60300,
		ProposerIndex: 403,
		Body: &gloas.BeaconBlockBody{
			RANDAOReveal:              phase0.BLSSignature{0x0a, 0x0b, 0x0c},
			ETH1Data:                  &phase0.ETH1Data{BlockHash: make([]byte, phase0.Hash32Length)},
			SyncAggregate:             &altair.SyncAggregate{SyncCommitteeBits: bitfield.NewBitvector512()},
			SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{}},
			ParentExecutionRequests:   &gloas.ExecutionRequests{},
		},
	}
}

// validEPBSBlockContents returns block contents carrying one KZG proof and no
// blobs.  Blobs are left empty deliberately: deneb.Blob is a 128KiB array, so a
// populated one turns any assertion failure into a 256KiB hex diff.  The
// payload's BaseFeePerGas must be non-nil, as its own marshaler dereferences it.
func validEPBSBlockContents() *apiv1gloas.BlockContents {
	contents := &apiv1gloas.BlockContents{
		Block: validEPBSBeaconBlock(),
		ExecutionPayloadEnvelope: &gloas.ExecutionPayloadEnvelope{
			Payload:           &gloas.ExecutionPayload{BaseFeePerGas: uint256.NewInt(7)},
			ExecutionRequests: &gloas.ExecutionRequests{},
			BuilderIndex:      12,
		},
		KZGProofs: []deneb.KZGProof{{0x01, 0x02, 0x03}},
		Blobs:     []deneb.Blob{},
	}

	contents.Block.Body.SignedExecutionPayloadBid.Message.BuilderIndex = contents.ExecutionPayloadEnvelope.BuilderIndex
	executionRequestsRoot, err := contents.ExecutionPayloadEnvelope.ExecutionRequests.HashTreeRoot()
	if err != nil {
		panic(err)
	}
	contents.Block.Body.SignedExecutionPayloadBid.Message.ExecutionRequestsRoot = phase0.Root(executionRequestsRoot)

	root, err := contents.Block.HashTreeRoot()
	if err != nil {
		panic(err)
	}
	contents.ExecutionPayloadEnvelope.BeaconBlockRoot = root

	return contents
}

// epbsProposalJSONBody wraps a marshaled datum in the produceBlockV4 response
// object, whose value fields and payload-inclusion flag sit alongside the data.
// It returns the marshaled datum as well, so a test can compare what came back
// against the exact bytes that went in.
func epbsProposalJSONBody(t *testing.T, included bool, datum any) (body, data []byte) {
	t.Helper()

	data, err := json.Marshal(datum)
	require.NoError(t, err)

	return fmt.Appendf(nil,
		`{"version":"gloas","consensus_block_value":"3476149000000000","execution_payload_value":"12345","execution_payload_included":%t,"data":%s}`,
		included, data,
	), data
}

// requireJSONEquals asserts that a decoded datum re-marshals to the bytes it was
// decoded from.  Comparing bytes rather than using require.Equal on the structs
// is deliberate: the nested gloas codecs normalise a nil slice to an empty one
// on the way in, so struct equality would assert those packages' conventions
// rather than this decoder's behaviour.  Byte equality still catches what
// matters — drop a field and its key re-marshals differently.
func requireJSONEquals(t *testing.T, want []byte, got any) {
	t.Helper()

	reencoded, err := json.Marshal(got)
	require.NoError(t, err)
	require.JSONEq(t, string(want), string(reencoded))
}

// TestEPBSProposalFromResponse covers the decode half of the block-production
// endpoint, which takes an already-fetched response and so needs no beacon
// node.
func TestEPBSProposalFromResponse(t *testing.T) {
	ctx := context.Background()
	s := &Service{}

	// With the execution payload excluded the data is a bare beacon block, and
	// the caller must fetch the envelope from the producing node to publish it.
	t.Run("JSONPayloadExcluded", func(t *testing.T) {
		body, data := epbsProposalJSONBody(t, false, validEPBSBeaconBlock())

		response, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body,
			headers:          map[string]string{},
		})
		require.NoError(t, err)
		require.Equal(t, spec.DataVersionGloas, response.Data.Version)
		require.False(t, response.Data.ExecutionPayloadIncluded)
		require.Nil(t, response.Data.GloasContents)
		require.NotNil(t, response.Data.Gloas)
		require.Equal(t, phase0.Slot(60300), response.Data.Gloas.Slot)
		requireJSONEquals(t, data, response.Data.Gloas)
	})

	// With the execution payload included the data is block contents: the block
	// plus the envelope, blobs and proofs a caller needs to publish it from any
	// node.  The flag is what selects between the two containers, so decoding
	// this body as a bare block would silently discard everything but the block.
	t.Run("JSONPayloadIncluded", func(t *testing.T) {
		body, data := epbsProposalJSONBody(t, true, validEPBSBlockContents())

		response, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body,
			headers:          map[string]string{},
		})
		require.NoError(t, err)
		require.True(t, response.Data.ExecutionPayloadIncluded)
		require.Nil(t, response.Data.Gloas)
		require.NotNil(t, response.Data.GloasContents)
		requireJSONEquals(t, data, response.Data.GloasContents)

		envelope, err := response.Data.ExecutionPayloadEnvelope()
		require.NoError(t, err)
		require.Equal(t, gloas.BuilderIndex(12), envelope.BuilderIndex)

		proofs, err := response.Data.KZGProofs()
		require.NoError(t, err)
		require.Len(t, proofs, 1)
	})

	// SSZ carries no field names and no response wrapper, so the header is the
	// only place the payload-inclusion flag can be read from.
	t.Run("SSZPayloadExcluded", func(t *testing.T) {
		block := validEPBSBeaconBlock()
		sszBody, err := block.MarshalSSZ()
		require.NoError(t, err)

		response, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeSSZ,
			consensusVersion: spec.DataVersionGloas,
			body:             sszBody,
			headers:          map[string]string{"Eth-Execution-Payload-Included": "false"},
		})
		require.NoError(t, err)
		require.False(t, response.Data.ExecutionPayloadIncluded)
		require.Nil(t, response.Data.GloasContents)
		require.NotNil(t, response.Data.Gloas)

		want, err := block.HashTreeRoot()
		require.NoError(t, err)
		got, err := response.Data.Gloas.HashTreeRoot()
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	// A gloas block carries preset-derived fixed-size fields — a sync committee
	// bitvector is 512 bits at the mainnet preset and 32 at minimal — and SSZ
	// addresses variable-size fields by offsets stored after them, so every
	// offset in the body shifts by 60 bytes between the two presets.  The
	// generated codec has the mainnet preset baked in and cannot read a body
	// encoded at another one, which is the whole reason this arm decodes through
	// the request-scoped dynamic codec instead.
	t.Run("SSZAtTheNodesPreset", func(t *testing.T) {
		custom := customSpecService(ctx, t)

		ds, err := custom.dynSSZForRequest(ctx)
		require.NoError(t, err)

		// The sync committee bitvector is the preset-derived field that shifts
		// the offsets, so the fixture has to be sized from the node's spec
		// rather than from the compiled-in preset to encode at all.
		specResponse, err := custom.Spec(ctx, &api.SpecOpts{})
		require.NoError(t, err)
		syncCommitteeSize, isUint64 := specResponse.Data["SYNC_COMMITTEE_SIZE"].(uint64)
		require.True(t, isUint64, "SYNC_COMMITTEE_SIZE missing from spec")

		block := validEPBSBeaconBlock()
		block.Body.SyncAggregate.SyncCommitteeBits = make(bitfield.Bitvector512, syncCommitteeSize/8)

		sszBody, err := ds.MarshalSSZ(block)
		require.NoError(t, err)

		response, err := custom.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeSSZ,
			consensusVersion: spec.DataVersionGloas,
			body:             sszBody,
			headers:          map[string]string{"Eth-Execution-Payload-Included": "false"},
		})
		require.NoError(t, err)
		require.NotNil(t, response.Data.Gloas)
		require.Equal(t, phase0.Slot(60300), response.Data.Gloas.Slot)
	})

	// Both value components arrive in headers.  Value() sums them, so a caller
	// choosing between proposals is comparing consensus rewards plus execution
	// payload value.
	t.Run("ValuesFromHeaders", func(t *testing.T) {
		body, _ := epbsProposalJSONBody(t, false, validEPBSBeaconBlock())

		response, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body,
			headers: map[string]string{
				"Eth-Consensus-Block-Value":   "3476149000000000",
				"Eth-Execution-Payload-Value": "999",
			},
		})
		require.NoError(t, err)
		require.Equal(t, big.NewInt(3476149000000000), response.Data.ConsensusValue)
		require.Equal(t, big.NewInt(999), response.Data.ExecutionValue)
		require.Equal(t, big.NewInt(3476149000000999), response.Data.Value())
	})

	// The JSON decoder fills in only the keys it finds, so a body carrying no
	// data key at all, or an explicit null, leaves the seeded nil pointer
	// untouched.  That has to be an error rather than a success wrapping
	// nothing: a caller handed a proposal whose arms are both nil gets errors
	// only later, from whichever accessor it happens to reach first.
	t.Run("NoDatum", func(t *testing.T) {
		tests := []struct {
			name string
			body string
		}{
			{name: "MissingDataKeyExcluded", body: `{"version":"gloas","execution_payload_included":false}`},
			{name: "NullDataExcluded", body: `{"version":"gloas","execution_payload_included":false,"data":null}`},
			{name: "MissingDataKeyIncluded", body: `{"version":"gloas","execution_payload_included":true}`},
			{name: "NullDataIncluded", body: `{"version":"gloas","execution_payload_included":true,"data":null}`},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				_, err := s.epbsProposalFromResponse(ctx, &httpResponse{
					statusCode:       http.StatusOK,
					contentType:      ContentTypeJSON,
					consensusVersion: spec.DataVersionGloas,
					body:             []byte(test.body),
					headers:          map[string]string{},
				})
				require.ErrorContains(t, err, "no gloas epbs proposal in response")
			})
		}
	})

	// The endpoint is gloas-onwards, so a response claiming any earlier fork is
	// a node answering a request it should have rejected.
	t.Run("PreGloasVersion", func(t *testing.T) {
		body, _ := epbsProposalJSONBody(t, false, validEPBSBeaconBlock())

		_, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionFulu,
			body:             body,
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "epbs proposal not available for version fulu")
	})

	t.Run("UnknownContentType", func(t *testing.T) {
		body, _ := epbsProposalJSONBody(t, false, validEPBSBeaconBlock())

		_, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeUnknown,
			consensusVersion: spec.DataVersionGloas,
			body:             body,
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "unhandled content type")
	})

	// The flag has no default in the spec, and the two containers are not
	// interchangeable, so a body that omits it cannot be decoded at all.
	t.Run("JSONWithoutInclusionFlag", func(t *testing.T) {
		_, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             []byte(`{"version":"gloas","data":{}}`),
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "no execution_payload_included in epbs proposal response")
	})

	// On the SSZ path the header is the only source for the flag.  Absent, the
	// body cannot be attributed to either container.
	t.Run("SSZWithoutInclusionHeader", func(t *testing.T) {
		_, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeSSZ,
			consensusVersion: spec.DataVersionGloas,
			body:             []byte{},
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "no Eth-Execution-Payload-Included header in epbs proposal response")
	})

	// A malformed flag must be rejected rather than read as false.  Treating an
	// unparseable value as false would decode block contents as a bare block and
	// silently drop the envelope, blobs and proofs.
	t.Run("SSZWithMalformedInclusionHeader", func(t *testing.T) {
		_, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeSSZ,
			consensusVersion: spec.DataVersionGloas,
			body:             []byte{},
			headers:          map[string]string{"Eth-Execution-Payload-Included": "yes"},
		})
		require.ErrorContains(t, err, "Eth-Execution-Payload-Included is not a valid boolean")

		// The offending value is named but not quoted back.  Its length is the
		// node's to choose, so an error carrying it puts however much the node
		// sent into a log line.
		_, err = s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeSSZ,
			consensusVersion: spec.DataVersionGloas,
			body:             []byte{},
			headers:          map[string]string{"Eth-Execution-Payload-Included": strings.Repeat("y", 4096)},
		})
		require.Less(t, len(err.Error()), 256, "the error must not echo an unbounded header value")
	})

	// A value header's length is chosen by the node, and big.Int.SetString is
	// quadratic in it.  The parse also happens after the request's deadline has
	// been released, so nothing on the caller's side interrupts it: a
	// multi-megabyte value burns minutes of CPU on the one path that has a slot
	// to meet.  Rejected on length before parsing.
	t.Run("OverlongValueHeader", func(t *testing.T) {
		body, _ := epbsProposalJSONBody(t, false, validEPBSBeaconBlock())

		_, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body,
			headers: map[string]string{
				"Eth-Consensus-Block-Value": strings.Repeat("9", maxProposalValueDigits+1),
			},
		})
		require.ErrorContains(t, err, "more than the 40 a value can need")

		// A value at the limit is still accepted, so the bound rejects only what
		// no real node sends.
		response, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body,
			headers: map[string]string{
				"Eth-Consensus-Block-Value": strings.Repeat("9", maxProposalValueDigits),
			},
		})
		require.NoError(t, err)
		require.Equal(t, strings.Repeat("9", maxProposalValueDigits), response.Data.ConsensusValue.String())
	})

	// A negative value is not a value.  big.Int.SetString accepts a leading
	// minus, and Value() would then subtract it from the total, making a proposal
	// look cheaper than a rival's.
	t.Run("NegativeValueHeader", func(t *testing.T) {
		body, _ := epbsProposalJSONBody(t, false, validEPBSBeaconBlock())

		_, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body,
			headers:          map[string]string{"Eth-Execution-Payload-Value": "-1"},
		})
		require.ErrorContains(t, err, "is negative")
	})

	t.Run("MalformedValueHeader", func(t *testing.T) {
		body, _ := epbsProposalJSONBody(t, false, validEPBSBeaconBlock())

		_, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body,
			headers:          map[string]string{"Eth-Consensus-Block-Value": "lots"},
		})
		require.ErrorContains(t, err, "not a valid integer")
	})

	// An omitted execution value is unknown, not zero.  A caller must not use
	// it as evidence that the execution bid has no value.
	t.Run("AbsentExecutionValueIsUnknown", func(t *testing.T) {
		body, _ := epbsProposalJSONBody(t, false, validEPBSBeaconBlock())

		response, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body,
			headers:          map[string]string{},
		})
		require.NoError(t, err)
		require.Nil(t, response.Data.ConsensusValue)
		require.Nil(t, response.Data.ExecutionValue)
		require.Nil(t, response.Data.Value())
	})

	t.Run("CorruptJSON", func(t *testing.T) {
		body, _ := epbsProposalJSONBody(t, false, validEPBSBeaconBlock())

		_, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body[:len(body)-1],
			headers:          map[string]string{},
		})
		require.ErrorContains(t, err, "failed to parse epbs proposal response")
	})

	t.Run("CorruptSSZ", func(t *testing.T) {
		block := validEPBSBeaconBlock()
		sszBody, err := block.MarshalSSZ()
		require.NoError(t, err)

		_, err = s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeSSZ,
			consensusVersion: spec.DataVersionGloas,
			body:             sszBody[:len(sszBody)-1],
			headers:          map[string]string{"Eth-Execution-Payload-Included": "false"},
		})
		require.ErrorContains(t, err, "failed to decode gloas SSZ epbs proposal")
	})
}

// TestAssertEPBSProposalMatchesRequest covers the checks made on what the node
// returned.  The faults cannot be provoked against a live node — one that
// honours the request cannot be made to answer with the wrong slot, mode or
// RANDAO reveal — so this is the only place they are exercised.
// BuilderBidExcludesRequestedPayload is the other side of that: a divergence
// the guard has to let through rather than reject.
func TestAssertEPBSProposalMatchesRequest(t *testing.T) {
	s := &Service{}

	includePayload := true
	reveal := phase0.BLSSignature{0x0a, 0x0b, 0x0c}

	// Matches validEPBSBeaconBlock, so the request and the response agree.
	matchingOpts := func() *api.EPBSProposalOpts {
		want := includePayload

		return &api.EPBSProposalOpts{
			Slot:           60300,
			RandaoReveal:   reveal,
			IncludePayload: &want,
		}
	}

	matchingProposal := func() *api.VersionedEPBSProposal {
		return &api.VersionedEPBSProposal{
			Version:                  spec.DataVersionGloas,
			ExecutionPayloadIncluded: true,
			GloasContents:            validEPBSBlockContents(),
		}
	}

	t.Run("Matching", func(t *testing.T) {
		require.NoError(t, s.assertEPBSProposalMatchesRequest(matchingProposal(), matchingOpts()))
	})

	t.Run("WrongSlot", func(t *testing.T) {
		opts := matchingOpts()
		opts.Slot = 60301

		err := s.assertEPBSProposalMatchesRequest(matchingProposal(), opts)
		require.ErrorIs(t, err, client.ErrInconsistentResult)
		require.ErrorContains(t, err, "for slot 60300; expected 60301")
	})

	// The mode decides where the block can be published, so a node that includes
	// a payload nobody asked for has changed that without saying so.  This is the
	// only direction that is a fault; the other one is what an external builder's
	// bid legitimately looks like.
	t.Run("WrongPayloadInclusion", func(t *testing.T) {
		exclude := false
		opts := matchingOpts()
		opts.IncludePayload = &exclude

		err := s.assertEPBSProposalMatchesRequest(matchingProposal(), opts)
		require.ErrorIs(t, err, client.ErrInconsistentResult)
		require.ErrorContains(t, err, "execution payload included true; expected false")
	})

	// An external builder's bid comes back as the beacon block alone whatever
	// include_payload asked for, because the node does not hold the builder's
	// payload and so has none to send.  Refusing it would leave the caller with
	// nothing to sign on every slot a builder wins.  The proposal is decoded from
	// a response body rather than built as a struct literal so that what is proved
	// is a real bid surviving the whole path, not just the check at the end of it
	// agreeing with a hand-made value.
	t.Run("BuilderBidExcludesRequestedPayload", func(t *testing.T) {
		ctx := context.Background()
		body, _ := epbsProposalJSONBody(t, false, validEPBSBeaconBlock())

		response, err := s.epbsProposalFromResponse(ctx, &httpResponse{
			statusCode:       http.StatusOK,
			contentType:      ContentTypeJSON,
			consensusVersion: spec.DataVersionGloas,
			body:             body,
			headers:          map[string]string{},
		})
		require.NoError(t, err)
		require.False(t, response.Data.ExecutionPayloadIncluded)

		require.NoError(t, s.assertEPBSProposalMatchesRequest(response.Data, matchingOpts()))
	})

	t.Run("WrongRandaoReveal", func(t *testing.T) {
		opts := matchingOpts()
		opts.RandaoReveal = phase0.BLSSignature{0xff}

		err := s.assertEPBSProposalMatchesRequest(matchingProposal(), opts)
		require.ErrorIs(t, err, client.ErrInconsistentResult)
		require.ErrorContains(t, err, "RANDAO reveal")
	})

	// Behind DVT middleware the reveal is the middleware's to decide, so a
	// mismatch there is expected rather than a fault.  The slot and the mode are
	// still checked.
	t.Run("RandaoRevealIgnoredBehindDVT", func(t *testing.T) {
		dvt := &Service{connectedToDVTMiddleware: true}

		opts := matchingOpts()
		opts.RandaoReveal = phase0.BLSSignature{0xff}

		require.NoError(t, dvt.assertEPBSProposalMatchesRequest(matchingProposal(), opts))

		opts.Slot = 60301
		require.ErrorIs(t,
			dvt.assertEPBSProposalMatchesRequest(matchingProposal(), opts),
			client.ErrInconsistentResult,
		)
	})
}

// customSpecService returns a service configured the way a caller on a
// non-mainnet preset must configure one: with custom spec support, so SSZ is
// decoded against the spec the node reports rather than the compiled-in preset.
func customSpecService(ctx context.Context, t *testing.T) *Service {
	t.Helper()

	return newTestService(ctx, t, true)
}
