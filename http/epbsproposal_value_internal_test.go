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
	"math/big"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestValidateEPBSProposalExecutionValueRejectsMismatchedP2PHeader(t *testing.T) {
	proposal := epbsProposalWithBid(1, 10, 0, 11_000_000_000)

	err := validateEPBSProposalExecutionValue(proposal, &gloas.BuilderConfig{MinBid: 10}, "")
	require.Error(t, err)
	require.ErrorIs(t, err, client.ErrInconsistentResult)
}

func TestValidateEPBSProposalExecutionValueRejectsP2PBelowMinimum(t *testing.T) {
	proposal := epbsProposalWithBid(1, 9, 0, 9_000_000_000)

	err := validateEPBSProposalExecutionValue(proposal, &gloas.BuilderConfig{MinBid: 10}, "")
	require.ErrorIs(t, err, client.ErrInconsistentResult)
}

func TestValidateEPBSProposalExecutionValueRejectsMismatchedDirectHeader(t *testing.T) {
	proposal := epbsProposalWithBid(1, 10, 5, 15_000_000_000)
	config := &gloas.BuilderConfig{Builders: []*gloas.BuilderEntry{{
		URL:                 []byte("https://builder.example"),
		MaxExecutionPayment: 2,
		MinBid:              12,
	}}}

	err := validateEPBSProposalExecutionValue(proposal, config, "https://builder.example")
	require.ErrorIs(t, err, client.ErrInconsistentResult)
}

func TestValidateEPBSProposalExecutionValueUnknownAndSelfBuild(t *testing.T) {
	config := &gloas.BuilderConfig{Builders: []*gloas.BuilderEntry{
		{URL: []byte("https://duplicate.example"), MaxExecutionPayment: 2},
		{URL: []byte("https://duplicate.example"), MaxExecutionPayment: 3},
	}}
	tests := []struct {
		name     string
		proposal *api.VersionedEPBSProposal
		url      string
		err      bool
		unknown  bool
	}{
		{name: "UnknownHeader", proposal: epbsProposalWithBid(1, 10, 0, 0)},
		{name: "DuplicateBuilderURL", proposal: epbsProposalWithBid(1, 10, 0, 10_000_000_000), url: "https://duplicate.example", unknown: true},
		{name: "SelfBuild", proposal: epbsProposalWithBid(gloas.BuilderIndex(^uint64(0)), 0, 7, 7_000_000_000), unknown: true},
		{name: "SelfBuildNonZeroBid", proposal: epbsProposalWithBid(gloas.BuilderIndex(^uint64(0)), 1, 0, 1_000_000_000), err: true},
		{name: "DirectBelowMinimum", proposal: epbsProposalWithBid(1, 10, 2, 12_000_000_000), url: "https://direct.example", err: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configured := config
			if test.name == "DirectBelowMinimum" {
				configured = &gloas.BuilderConfig{Builders: []*gloas.BuilderEntry{{URL: []byte(test.url), MaxExecutionPayment: 2, MinBid: 13}}}
			}
			if test.name == "UnknownHeader" {
				test.proposal.ExecutionValue = nil
			}
			err := validateEPBSProposalExecutionValue(test.proposal, configured, test.url)
			if test.err {
				require.ErrorIs(t, err, client.ErrInconsistentResult)
			} else {
				require.NoError(t, err)
			}
			if test.unknown {
				require.Nil(t, test.proposal.ExecutionValue)
			}
		})
	}
}

func epbsProposalWithBid(builderIndex gloas.BuilderIndex,
	value phase0.Gwei,
	executionPayment phase0.Gwei,
	executionValue uint64,
) *api.VersionedEPBSProposal {
	return &api.VersionedEPBSProposal{
		Version:        spec.DataVersionGloas,
		ExecutionValue: new(big.Int).SetUint64(executionValue),
		Gloas: &gloas.BeaconBlock{Body: &gloas.BeaconBlockBody{
			SignedExecutionPayloadBid: &gloas.SignedExecutionPayloadBid{Message: &gloas.ExecutionPayloadBid{
				BuilderIndex:     builderIndex,
				Value:            value,
				ExecutionPayment: executionPayment,
			}},
		}},
	}
}
