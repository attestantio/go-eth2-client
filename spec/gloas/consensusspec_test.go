// Copyright © 2025 Attestant Limited.
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

package gloas_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/golang/snappy"
	clone "github.com/huandu/go-clone/generic"
	"github.com/pk910/dynamic-ssz/sszutils"
	require "github.com/stretchr/testify/require"
)

// TestConsensusSpec tests the types against the Ethereum consensus spec tests.
func TestConsensusSpec(t *testing.T) {
	if os.Getenv("CONSENSUS_SPEC_TESTS_DIR") == "" {
		t.Skip("CONSENSUS_SPEC_TESTS_DIR not supplied, not running spec tests")
	}

	tests := []struct {
		name string
		s    any
	}{
		{
			name: "AggregateAndProof",
			s:    &gloas.AggregateAndProof{},
		},
		{
			name: "Attestation",
			s:    &gloas.Attestation{},
		},
		{
			name: "AttestationData",
			s:    &phase0.AttestationData{},
		},
		{
			name: "AttesterSlashing",
			s:    &gloas.AttesterSlashing{},
		},
		{
			name: "BeaconBlock",
			s:    &gloas.BeaconBlock{},
		},
		{
			name: "BeaconBlockBody",
			s:    &gloas.BeaconBlockBody{},
		},
		{
			name: "BeaconBlockHeader",
			s:    &phase0.BeaconBlockHeader{},
		},
		{
			name: "BeaconState",
			s:    &gloas.BeaconState{},
		},
		{
			name: "BLSToExecutionChange",
			s:    &capella.BLSToExecutionChange{},
		},
		{
			name: "Builder",
			s:    &gloas.Builder{},
		},
		{
			name: "BuilderDepositRequest",
			s:    &gloas.BuilderDepositRequest{},
		},
		{
			name: "BuilderExitRequest",
			s:    &gloas.BuilderExitRequest{},
		},
		{
			name: "BuilderPendingPayment",
			s:    &gloas.BuilderPendingPayment{},
		},
		{
			name: "BuilderPendingWithdrawal",
			s:    &gloas.BuilderPendingWithdrawal{},
		},
		{
			name: "Checkpoint",
			s:    &phase0.Checkpoint{},
		},
		{
			name: "ConsolidationRequest",
			s:    &electra.ConsolidationRequest{},
		},
		{
			name: "ContributionAndProof",
			s:    &altair.ContributionAndProof{},
		},
		{
			name: "Deposit",
			s:    &phase0.Deposit{},
		},
		{
			name: "DepositData",
			s:    &phase0.DepositData{},
		},
		{
			name: "DepositMessage",
			s:    &phase0.DepositMessage{},
		},
		{
			name: "DepositRequest",
			s:    &electra.DepositRequest{},
		},
		{
			name: "Eth1Data",
			s:    &phase0.ETH1Data{},
		},
		{
			name: "ExecutionPayload",
			s:    &gloas.ExecutionPayload{},
		},
		{
			name: "ExecutionPayloadBid",
			s:    &gloas.ExecutionPayloadBid{},
		},
		{
			name: "ExecutionPayloadEnvelope",
			s:    &gloas.ExecutionPayloadEnvelope{},
		},
		{
			name: "ExecutionRequests",
			s:    &gloas.ExecutionRequests{},
		},
		{
			name: "Fork",
			s:    &phase0.Fork{},
		},
		{
			name: "ForkData",
			s:    &phase0.ForkData{},
		},
		{
			name: "HistoricalSummary",
			s:    &capella.HistoricalSummary{},
		},
		{
			name: "IndexedAttestation",
			s:    &gloas.IndexedAttestation{},
		},
		{
			name: "IndexedPayloadAttestation",
			s:    &gloas.IndexedPayloadAttestation{},
		},
		{
			name: "PayloadAttestation",
			s:    &gloas.PayloadAttestation{},
		},
		{
			name: "PayloadAttestationData",
			s:    &gloas.PayloadAttestationData{},
		},
		{
			name: "PayloadAttestationMessage",
			s:    &gloas.PayloadAttestationMessage{},
		},
		{
			name: "PendingConsolidation",
			s:    &electra.PendingConsolidation{},
		},
		{
			name: "PendingDeposit",
			s:    &electra.PendingDeposit{},
		},
		{
			name: "PendingPartialWithdrawal",
			s:    &electra.PendingPartialWithdrawal{},
		},
		{
			name: "ProposerPreferences",
			s:    &gloas.ProposerPreferences{},
		},
		{
			name: "ProposerSlashing",
			s:    &phase0.ProposerSlashing{},
		},
		{
			name: "SignedAggregateAndProof",
			s:    &gloas.SignedAggregateAndProof{},
		},
		{
			name: "SignedBeaconBlock",
			s:    &gloas.SignedBeaconBlock{},
		},
		{
			name: "SignedBeaconBlockHeader",
			s:    &phase0.SignedBeaconBlockHeader{},
		},
		{
			name: "SignedBLSToExecutionChange",
			s:    &capella.SignedBLSToExecutionChange{},
		},
		{
			name: "SignedContributionAndProof",
			s:    &altair.SignedContributionAndProof{},
		},
		{
			name: "SignedExecutionPayloadBid",
			s:    &gloas.SignedExecutionPayloadBid{},
		},
		{
			name: "SignedExecutionPayloadEnvelope",
			s:    &gloas.SignedExecutionPayloadEnvelope{},
		},
		{
			name: "SignedProposerPreferences",
			s:    &gloas.SignedProposerPreferences{},
		},
		{
			name: "SignedVoluntaryExit",
			s:    &phase0.SignedVoluntaryExit{},
		},
		{
			name: "SingleAttestation",
			s:    &electra.SingleAttestation{},
		},
		{
			name: "SyncAggregate",
			s:    &altair.SyncAggregate{},
		},
		{
			name: "SyncAggregatorSelectionData",
			s:    &altair.SyncAggregatorSelectionData{},
		},
		{
			name: "SyncCommittee",
			s:    &altair.SyncCommittee{},
		},
		{
			name: "SyncCommitteeContribution",
			s:    &altair.SyncCommitteeContribution{},
		},
		{
			name: "SyncCommitteeMessage",
			s:    &altair.SyncCommitteeMessage{},
		},
		{
			name: "Validator",
			s:    &phase0.Validator{},
		},
		{
			name: "VoluntaryExit",
			s:    &phase0.VoluntaryExit{},
		},
		{
			name: "Withdrawal",
			s:    &capella.Withdrawal{},
		},
		{
			name: "WithdrawalRequest",
			s:    &electra.WithdrawalRequest{},
		},
	}

	baseDir := filepath.Join(os.Getenv("CONSENSUS_SPEC_TESTS_DIR"), "tests", "mainnet", "gloas", "ssz_static")
	for _, test := range tests {
		dir := filepath.Join(baseDir, test.name, "ssz_random")
		require.NoError(t, filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if path == dir {
				// Only interested in subdirectories.
				return nil
			}
			require.NoError(t, err)
			if info.IsDir() {
				t.Run(fmt.Sprintf("%s/%s", test.name, info.Name()), func(t *testing.T) {
					s1 := clone.Clone(test.s)
					// Obtain the struct from the YAML.
					specYAML, err := os.ReadFile(filepath.Join(path, "value.yaml"))
					require.NoError(t, err)
					require.NoError(t, yaml.Unmarshal(specYAML, s1))
					// Confirm we can return to the YAML.
					remarshalledSpecYAML, err := yaml.Marshal(s1)
					require.NoError(t, err)
					require.YAMLEq(t, testYAMLFormat(specYAML), testYAMLFormat(remarshalledSpecYAML))

					// Obtain the struct from the SSZ.
					s2 := clone.Clone(test.s)
					compressedSpecSSZ, err := os.ReadFile(filepath.Join(path, "serialized.ssz_snappy"))
					require.NoError(t, err)
					var specSSZ []byte
					specSSZ, err = snappy.Decode(specSSZ, compressedSpecSSZ)
					require.NoError(t, err)
					require.NoError(t, s2.(sszutils.FastsszUnmarshaler).UnmarshalSSZ(specSSZ))
					// Confirm we can return to the SSZ.
					remarshalledSpecSSZ, err := s2.(sszutils.FastsszMarshaler).MarshalSSZ()
					require.NoError(t, err)
					require.Equal(t, specSSZ, remarshalledSpecSSZ)

					// Obtain the hash tree root from the YAML.
					specYAMLRoot, err := os.ReadFile(filepath.Join(path, "roots.yaml"))
					require.NoError(t, err)
					// Confirm we calculate the same root.
					generatedRootBytes, err := s2.(sszutils.FastsszHashRoot).HashTreeRoot()
					require.NoError(t, err)
					generatedRoot := fmt.Sprintf("{root: '%#x'}\n", string(generatedRootBytes[:]))
					require.YAMLEq(t, string(specYAMLRoot), generatedRoot)
				})
			}

			return nil
		}))
	}
}

func testYAMLFormat(input []byte) string {
	val := make(map[string]any)
	if err := yaml.UnmarshalWithOptions(input, &val, yaml.UseOrderedMap()); err != nil {
		panic(err)
	}

	res, err := yaml.MarshalWithOptions(val, yaml.Flow(true))
	if err != nil {
		panic(err)
	}

	replacements := [][][]byte{
		{[]byte(`"`), []byte(`'`)},
		// Field 'extra_data' in ExecutionPayloadHeader/case_1 has a non-standard format, fix here.
		{[]byte(`extra_data: 0,`), []byte(`extra_data: '0x',`)},
	}
	for _, replacement := range replacements {
		res = bytes.ReplaceAll(res, replacement[0], replacement[1])
	}

	return string(bytes.ToLower(res))
}
