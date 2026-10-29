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

	err := validateEPBSProposalExecutionValue(proposal, &gloas.BuilderConfig{MinBid: 10}, "", staticBuilderIndexSelfBuild, nil)
	require.Error(t, err)
	require.ErrorIs(t, err, client.ErrInconsistentResult)
}

func TestValidateEPBSProposalExecutionValueRejectsP2PBelowMinimum(t *testing.T) {
	proposal := epbsProposalWithBid(1, 9, 0, 9_000_000_000)

	err := validateEPBSProposalExecutionValue(proposal, &gloas.BuilderConfig{MinBid: 10}, "", staticBuilderIndexSelfBuild, nil)
	require.ErrorIs(t, err, client.ErrInconsistentResult)
}

func TestValidateEPBSProposalExecutionValueRejectsMismatchedDirectHeader(t *testing.T) {
	proposal := epbsProposalWithBid(1, 10, 5, 15_000_000_000)
	config := &gloas.BuilderConfig{Builders: []*gloas.BuilderEntry{{
		URL:                 []byte("https://builder.example"),
		MaxExecutionPayment: 2,
		MinBid:              12,
	}}}

	err := validateEPBSProposalExecutionValue(proposal, config, "https://builder.example", staticBuilderIndexSelfBuild, nil)
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
		// A self-build has no signed bid the reported value could contradict,
		// and bid.Value is zero by construction, so the node's figure is the
		// only signal of what the block is worth and is kept.
		{
			name:     "SelfBuildKeepsReportedValue",
			proposal: epbsProposalWithBid(staticBuilderIndexSelfBuild, 0, 7, 7_000_000_000),
			config:   duplicateURLs,
		},
		{
			name:     "SelfBuildNonZeroBid",
			proposal: epbsProposalWithBid(staticBuilderIndexSelfBuild, 1, 0, 1_000_000_000),
			config:   duplicateURLs,
			err:      true,
		},
		// An Eth-Builder-Url the config does not name is the node's to choose,
		// so it must not buy the bid an exemption from the configured floor.
		{
			name:     "UnmatchedBuilderURLStillChecksMinimum",
			proposal: epbsProposalWithBid(1, 1, 0, 1_000_000_000),
			config:   &gloas.BuilderConfig{MinBid: 1_000_000_000},
			url:      "https://unknown.example",
			err:      true,
		},
		{
			name:     "DuplicateBuilderURLStillChecksMinimum",
			proposal: epbsProposalWithBid(1, 1, 0, 1_000_000_000),
			config: &gloas.BuilderConfig{MinBid: 1_000_000_000, Builders: []*gloas.BuilderEntry{
				{URL: []byte("https://duplicate.example")},
				{URL: []byte("https://duplicate.example")},
			}},
			url: "https://duplicate.example",
			err: true,
		},
		// A header sent twice reaches the validator joined with ";".
		{
			name:     "JoinedBuilderURLStillChecksMinimum",
			proposal: epbsProposalWithBid(1, 1, 0, 1_000_000_000),
			config: &gloas.BuilderConfig{MinBid: 1_000_000_000, Builders: []*gloas.BuilderEntry{
				{URL: []byte("https://direct.example")},
			}},
			url: "https://direct.example;https://direct.example",
			err: true,
		},
		// Nor must a p2p builder disable the floor by attaching a token
		// execution payment that cannot be capped.
		{
			name:     "P2PPaymentStillChecksMinimum",
			proposal: epbsProposalWithBid(1, 1, 1, 1_000_000_000),
			config:   &gloas.BuilderConfig{MinBid: 1_000_000_000},
			err:      true,
		},
		// An attributed bid is judged against its entry's minimum only; the
		// config-wide minimum is the floor for unattributed bids.
		{
			name:     "AttributedBidUsesOnlyEntryMinimum",
			proposal: epbsProposalWithBid(1, 5, 0, 5_000_000_000),
			config: &gloas.BuilderConfig{MinBid: 1_000_000_000, Builders: []*gloas.BuilderEntry{
				{URL: []byte("https://direct.example"), MinBid: 5},
			}},
			url: "https://direct.example",
		},
		// Where the bid cannot be attributed the payment counts up to the most
		// permissive cap the config allows, which can only admit bids that a
		// resolvable cap would also have admitted.
		{
			name:     "UnattributedBidUsesMostPermissiveCap",
			proposal: epbsProposalWithBid(1, 6, 5, 11_000_000_000),
			config: &gloas.BuilderConfig{MinBid: 10, Builders: []*gloas.BuilderEntry{
				{URL: []byte("https://a.example"), MaxExecutionPayment: 1},
				{URL: []byte("https://b.example"), MaxExecutionPayment: 5},
			}},
			unknown: true,
		},
		// A matched entry's own cap applies, not the largest one in the config.
		{
			name:     "MatchedEntryUsesItsOwnCap",
			proposal: epbsProposalWithBid(1, 6, 5, 11_000_000_000),
			config: &gloas.BuilderConfig{MinBid: 10, Builders: []*gloas.BuilderEntry{
				{URL: []byte("https://a.example"), MaxExecutionPayment: 1},
				{URL: []byte("https://b.example"), MaxExecutionPayment: 5},
			}},
			url: "https://a.example",
			err: true,
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
			// Metadata holds every response header and, for JSON, every top-level
			// body key, so a withheld figure has to go from all of them.
			metadata := map[string]any{}
			if value != nil {
				metadata["Eth-Execution-Payload-Value"] = value.String()
				metadata["execution_payload_value"] = value.String()
			}
			err := validateEPBSProposalExecutionValue(test.proposal, test.config, test.url, staticBuilderIndexSelfBuild, metadata)
			if test.err {
				require.ErrorIs(t, err, client.ErrInconsistentResult)

				return
			}
			require.NoError(t, err)
			if test.unknown {
				require.Nil(t, test.proposal.ExecutionValue)
				require.NotContains(t, metadata, "Eth-Execution-Payload-Value")
				require.NotContains(t, metadata, "execution_payload_value")

				return
			}
			require.Equal(t, value, test.proposal.ExecutionValue)
			if value != nil {
				require.Contains(t, metadata, "Eth-Execution-Payload-Value")
			}
		})
	}
}

// TestValidateEPBSProposalExecutionValueRejectsAnAbsentBid covers the walk to the
// bid, which reaches this function from one caller today but panicked rather than
// errored if it were ever reached from another.
func TestValidateEPBSProposalExecutionValueRejectsAnAbsentBid(t *testing.T) {
	tests := []struct {
		name     string
		proposal *api.VersionedEPBSProposal
	}{
		{
			name:     "NoBlock",
			proposal: &api.VersionedEPBSProposal{Version: spec.DataVersionGloas},
		},
		{
			name: "NoBody",
			proposal: &api.VersionedEPBSProposal{
				Version: spec.DataVersionGloas,
				Gloas:   &gloas.BeaconBlock{},
			},
		},
		{
			name: "NoBid",
			proposal: &api.VersionedEPBSProposal{
				Version: spec.DataVersionGloas,
				Gloas:   &gloas.BeaconBlock{Body: &gloas.BeaconBlockBody{}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				err := validateEPBSProposalExecutionValue(test.proposal,
					&gloas.BuilderConfig{}, "", staticBuilderIndexSelfBuild, nil)
				require.ErrorIs(t, err, client.ErrInconsistentResult)
			})
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
