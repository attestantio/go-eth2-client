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
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/attestantio/go-eth2-client/spec"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/stretchr/testify/require"
)

func TestDecodeEPBSProposalJSONAcceptsBareBuilderWin(t *testing.T) {
	service, err := mock.New(context.Background())
	require.NoError(t, err)
	includePayload := false
	response, err := service.EPBSProposal(context.Background(), &api.EPBSProposalOpts{
		Slot:           123,
		IncludePayload: &includePayload,
	})
	require.NoError(t, err)
	block := response.Data.Gloas
	data, err := json.Marshal(block)
	require.NoError(t, err)

	body := append([]byte(`{"execution_payload_included":false,"data":`), data...)
	body = append(body, '}')

	proposal := &api.VersionedEPBSProposal{Version: spec.DataVersionGloas}
	metadata, err := decodeEPBSProposalJSON(body, proposal)
	require.NoError(t, err)
	require.Equal(t, false, metadata["execution_payload_included"])
	require.False(t, proposal.ExecutionPayloadIncluded)
	require.Equal(t, block, proposal.Gloas)
	require.Nil(t, proposal.GloasContents)
}

func TestDecodeEPBSProposalJSONRejectsNullPayloadInclusion(t *testing.T) {
	service, err := mock.New(context.Background())
	require.NoError(t, err)
	includePayload := false
	response, err := service.EPBSProposal(context.Background(), &api.EPBSProposalOpts{
		Slot:           123,
		IncludePayload: &includePayload,
	})
	require.NoError(t, err)
	data, err := json.Marshal(response.Data.Gloas)
	require.NoError(t, err)
	body := append([]byte(`{"execution_payload_included":null,"data":`), data...)
	body = append(body, '}')

	tests := []struct {
		name string
		err  string
	}{
		{
			name: "NullPayloadInclusion",
			err:  "execution_payload_included cannot be null",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proposal := &api.VersionedEPBSProposal{Version: spec.DataVersionGloas}
			_, err := decodeEPBSProposalJSON(body, proposal)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestDecodeEPBSProposalJSONRejectsMissingOrMalformedMetadata(t *testing.T) {
	tests := []struct {
		name string
		body string
		err  string
	}{
		{
			name: "MissingPayloadInclusion",
			body: `{"data":{}}`,
			err:  "no execution_payload_included in epbs proposal response",
		},
		{
			name: "MalformedPayloadInclusion",
			body: `{"execution_payload_included":"false","data":{}}`,
			err:  "failed to unmarshal execution_payload_included\njson: cannot unmarshal string into Go value of type bool",
		},
		{
			name: "MissingData",
			body: `{"execution_payload_included":false}`,
			err:  "no gloas epbs proposal in response",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proposal := &api.VersionedEPBSProposal{Version: spec.DataVersionGloas}
			_, err := decodeEPBSProposalJSON([]byte(test.body), proposal)
			require.EqualError(t, err, test.err)
		})
	}
}

func TestAssertIncludedEPBSProposalEnvelopeMatchesBlockRejectsMismatchedBuilderIndex(t *testing.T) {
	mockService, err := mock.New(context.Background())
	require.NoError(t, err)

	includePayload := true
	response, err := mockService.EPBSProposal(context.Background(), &api.EPBSProposalOpts{
		Slot:           123,
		IncludePayload: &includePayload,
	})
	require.NoError(t, err)
	blockRoot, err := response.Data.Root()
	require.NoError(t, err)
	response.Data.GloasContents.ExecutionPayloadEnvelope.BeaconBlockRoot = blockRoot
	response.Data.GloasContents.Block.Body.SignedExecutionPayloadBid.Message.BuilderIndex = 1
	response.Data.GloasContents.ExecutionPayloadEnvelope.BuilderIndex = 2

	tests := []struct {
		name string
		err  string
	}{
		{
			name: "MismatchedBuilderIndex",
			err:  "execution payload envelope builder index does not match bid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := assertIncludedEPBSProposalEnvelopeMatchesBlock(response.Data, dynssz.GetGlobalDynSsz())
			require.ErrorContains(t, err, test.err)
			require.ErrorIs(t, err, client.ErrInconsistentResult)
		})
	}
}

func TestAssertIncludedEPBSProposalEnvelopeMatchesBlockRejectsMismatchedBlockHash(t *testing.T) {
	mockService, err := mock.New(context.Background())
	require.NoError(t, err)

	includePayload := true
	response, err := mockService.EPBSProposal(context.Background(), &api.EPBSProposalOpts{
		Slot:           123,
		IncludePayload: &includePayload,
	})
	require.NoError(t, err)
	blockRoot, err := response.Data.Root()
	require.NoError(t, err)
	response.Data.GloasContents.ExecutionPayloadEnvelope.BeaconBlockRoot = blockRoot
	response.Data.GloasContents.Block.Body.SignedExecutionPayloadBid.Message.BlockHash[0] = 1
	response.Data.GloasContents.ExecutionPayloadEnvelope.Payload.BlockHash[0] = 2

	tests := []struct {
		name string
		err  string
	}{
		{
			name: "MismatchedBlockHash",
			err:  "execution payload block hash does not match bid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := assertIncludedEPBSProposalEnvelopeMatchesBlock(response.Data, dynssz.GetGlobalDynSsz())
			require.ErrorContains(t, err, test.err)
			require.ErrorIs(t, err, client.ErrInconsistentResult)
		})
	}
}

func TestAssertIncludedEPBSProposalEnvelopeMatchesBlockRejectsMismatchedParentRoot(t *testing.T) {
	mockService, err := mock.New(context.Background())
	require.NoError(t, err)

	includePayload := true
	response, err := mockService.EPBSProposal(context.Background(), &api.EPBSProposalOpts{
		Slot:           123,
		IncludePayload: &includePayload,
	})
	require.NoError(t, err)
	blockRoot, err := response.Data.Root()
	require.NoError(t, err)
	response.Data.GloasContents.ExecutionPayloadEnvelope.BeaconBlockRoot = blockRoot
	response.Data.GloasContents.ExecutionPayloadEnvelope.ParentBeaconBlockRoot[0] = 1

	tests := []struct {
		name string
		err  string
	}{
		{
			name: "MismatchedParentRoot",
			err:  "execution payload envelope parent root does not match bid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := assertIncludedEPBSProposalEnvelopeMatchesBlock(response.Data, dynssz.GetGlobalDynSsz())
			require.ErrorContains(t, err, test.err)
			require.ErrorIs(t, err, client.ErrInconsistentResult)
		})
	}
}

func TestAssertIncludedEPBSProposalEnvelopeMatchesBlockRejectsMismatchedExecutionRequests(t *testing.T) {
	mockService, err := mock.New(context.Background())
	require.NoError(t, err)

	includePayload := true
	response, err := mockService.EPBSProposal(context.Background(), &api.EPBSProposalOpts{
		Slot:           123,
		IncludePayload: &includePayload,
	})
	require.NoError(t, err)
	blockRoot, err := response.Data.Root()
	require.NoError(t, err)
	response.Data.GloasContents.ExecutionPayloadEnvelope.BeaconBlockRoot = blockRoot
	response.Data.GloasContents.Block.Body.SignedExecutionPayloadBid.Message.ExecutionRequestsRoot[0] = 1

	tests := []struct {
		name string
		err  string
	}{
		{
			name: "MismatchedExecutionRequests",
			err:  "execution payload envelope requests do not match bid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := assertIncludedEPBSProposalEnvelopeMatchesBlock(response.Data, dynssz.GetGlobalDynSsz())
			require.ErrorContains(t, err, test.err)
			require.ErrorIs(t, err, client.ErrInconsistentResult)
		})
	}
}

func TestAssertEPBSProposalMatchesRequestAllowsBareBuilderWin(t *testing.T) {
	mockService, err := mock.New(context.Background())
	require.NoError(t, err)

	payloadIncluded := false
	response, err := mockService.EPBSProposal(context.Background(), &api.EPBSProposalOpts{
		Slot:           123,
		IncludePayload: &payloadIncluded,
	})
	require.NoError(t, err)

	requestedIncluded := true
	tests := []struct {
		name string
		err  string
	}{
		{
			name: "BareBuilderWin",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := (&Service{}).assertEPBSProposalMatchesRequest(response.Data, &api.EPBSProposalOpts{
				Slot:           123,
				IncludePayload: &requestedIncluded,
			})
			if test.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.err)
			}
		})
	}
}

func TestAssertEPBSProposalMatchesRequestRejectsUnexpectedContents(t *testing.T) {
	mockService, err := mock.New(context.Background())
	require.NoError(t, err)

	payloadIncluded := true
	response, err := mockService.EPBSProposal(context.Background(), &api.EPBSProposalOpts{
		Slot:           123,
		IncludePayload: &payloadIncluded,
	})
	require.NoError(t, err)

	requestedIncluded := false
	tests := []struct {
		name string
		err  string
	}{
		{
			name: "UnexpectedContents",
			err:  "epbs beacon block proposal has execution payload included true; expected false\ninconsistent result",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := (&Service{}).assertEPBSProposalMatchesRequest(response.Data, &api.EPBSProposalOpts{
				Slot:           123,
				IncludePayload: &requestedIncluded,
			})
			require.EqualError(t, err, test.err)
		})
	}
}

func TestEPBSPayloadIncludedFromHeaders(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		included bool
		err      string
	}{
		{
			name:     "Included",
			headers:  map[string]string{"Eth-Execution-Payload-Included": "true"},
			included: true,
		},
		{
			name:    "Missing",
			headers: map[string]string{},
			err:     "no Eth-Execution-Payload-Included header in epbs proposal response",
		},
		{
			name:    "Malformed",
			headers: map[string]string{"Eth-Execution-Payload-Included": "sometimes"},
			err:     "proposal header Eth-Execution-Payload-Included is not a valid boolean",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			included, err := epbsPayloadIncludedFromHeaders(test.headers)
			if test.err == "" {
				require.NoError(t, err)
				require.Equal(t, test.included, included)
			} else {
				require.EqualError(t, err, test.err)
			}
		})
	}
}
