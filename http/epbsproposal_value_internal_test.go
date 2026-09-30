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

	err := validateEPBSProposalExecutionValue(proposal, &gloas.BuilderConfig{MinBid: 10}, "", staticBuilderIndexSelfBuild)
	require.Error(t, err)
	require.ErrorIs(t, err, client.ErrInconsistentResult)
}

func TestValidateEPBSProposalExecutionValueRejectsP2PBelowMinimum(t *testing.T) {
	proposal := epbsProposalWithBid(1, 9, 0, 9_000_000_000)

	err := validateEPBSProposalExecutionValue(proposal, &gloas.BuilderConfig{MinBid: 10}, "", staticBuilderIndexSelfBuild)
	require.ErrorIs(t, err, client.ErrInconsistentResult)
}

func TestValidateEPBSProposalExecutionValueRejectsMismatchedDirectHeader(t *testing.T) {
	proposal := epbsProposalWithBid(1, 10, 5, 15_000_000_000)
	config := &gloas.BuilderConfig{Builders: []*gloas.BuilderEntry{{
		URL:                 []byte("https://builder.example"),
		MaxExecutionPayment: 2,
		MinBid:              12,
	}}}

	err := validateEPBSProposalExecutionValue(proposal, config, "https://builder.example", staticBuilderIndexSelfBuild)
	require.ErrorIs(t, err, client.ErrInconsistentResult)
}

func TestValidateEPBSProposalExecutionValueUnknownAndSelfBuild(t *testing.T) {
	duplicateURLs := &gloas.BuilderConfig{Builders: []*gloas.BuilderEntry{
		{URL: []byte("https://duplicate.example"), MaxExecutionPayment: 2},
		{URL: []byte("https://duplicate.example"), MaxExecutionPayment: 3},
	}}
	direct := func(maxPayment, minBid phase0.Gwei) *gloas.BuilderConfig {
		return &gloas.BuilderConfig{Builders: []*gloas.BuilderEntry{
			{URL: []byte("https://direct.example"), MaxExecutionPayment: maxPayment, MinBid: minBid},
		}}
	}
	tests := []struct {
		name     string
		proposal *api.VersionedEPBSProposal
		config   *gloas.BuilderConfig
		url      string
		err      bool
		unknown  bool
	}{
		{
			name:     "NoHeader",
			proposal: epbsProposalWithoutValue(1, 10, 0),
			config:   &gloas.BuilderConfig{MinBid: 10},
			unknown:  true,
		},
		{
			name:     "NoHeaderBelowMinimum",
			proposal: epbsProposalWithoutValue(1, 9, 0),
			config:   &gloas.BuilderConfig{MinBid: 10},
			err:      true,
		},
		{
			name:     "P2PPaymentUncheckable",
			proposal: epbsProposalWithBid(1, 10, 3, 13_000_000_000),
			config:   &gloas.BuilderConfig{MinBid: 10},
			unknown:  true,
		},
		{
			name:     "DuplicateBuilderURL",
			proposal: epbsProposalWithBid(1, 10, 0, 10_000_000_000),
			config:   duplicateURLs,
			url:      "https://duplicate.example",
			unknown:  true,
		},
		{
			name:     "SelfBuild",
			proposal: epbsProposalWithBid(staticBuilderIndexSelfBuild, 0, 7, 7_000_000_000),
			config:   duplicateURLs,
			unknown:  true,
		},
		{
			name:     "SelfBuildNonZeroBid",
			proposal: epbsProposalWithBid(staticBuilderIndexSelfBuild, 1, 0, 1_000_000_000),
			config:   duplicateURLs,
			err:      true,
		},
		{
			name:     "DirectBelowMinimum",
			proposal: epbsProposalWithBid(1, 10, 2, 12_000_000_000),
			config:   direct(2, 13),
			url:      "https://direct.example",
			err:      true,
		},
		// beacon-APIs does not say whether the reported value includes the
		// execution payment, so both ends of the bid-bound bracket are accepted
		// and anything outside it is not.
		{
			name:     "P2PValueMatches",
			proposal: epbsProposalWithBid(1, 10, 0, 10_000_000_000),
			config:   &gloas.BuilderConfig{MinBid: 10},
		},
		{
			name:     "DirectValueWithoutPayment",
			proposal: epbsProposalWithBid(1, 10, 3, 10_000_000_000),
			config:   direct(5, 10),
			url:      "https://direct.example",
		},
		{
			name:     "DirectValueWithPayment",
			proposal: epbsProposalWithBid(1, 10, 3, 13_000_000_000),
			config:   direct(5, 10),
			url:      "https://direct.example",
		},
		{
			name:     "DirectValueWithUncappedPayment",
			proposal: epbsProposalWithBid(1, 10, 3, 13_000_000_000),
			config:   direct(2, 10),
			url:      "https://direct.example",
			err:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := test.proposal.ExecutionValue
			err := validateEPBSProposalExecutionValue(test.proposal, test.config, test.url, staticBuilderIndexSelfBuild)
			if test.err {
				require.ErrorIs(t, err, client.ErrInconsistentResult)

				return
			}
			require.NoError(t, err)
			if test.unknown {
				require.Nil(t, test.proposal.ExecutionValue)

				return
			}
			require.Equal(t, value, test.proposal.ExecutionValue)
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

func epbsProposalWithoutValue(builderIndex gloas.BuilderIndex,
	value phase0.Gwei,
	executionPayment phase0.Gwei,
) *api.VersionedEPBSProposal {
	proposal := epbsProposalWithBid(builderIndex, value, executionPayment, 0)
	proposal.ExecutionValue = nil

	return proposal
}
