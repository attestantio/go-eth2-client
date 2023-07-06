// Copyright © 2023 Attestant Limited.
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

package bellatrix_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/goccy/go-yaml"
	"github.com/golang/snappy"
	clone "github.com/huandu/go-clone/generic"
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
			s:    &phase0.AggregateAndProof{},
		},
		{
			name: "Attestation",
			s:    &phase0.Attestation{},
		},
		{
			name: "AttestationData",
			s:    &phase0.AttestationData{},
		},
		{
			name: "AttesterSlashing",
			s:    &phase0.AttesterSlashing{},
		},
		{
			name: "BeaconBlock",
			s:    &bellatrix.BeaconBlock{},
		},
		{
			name: "BeaconBlockBody",
			s:    &bellatrix.BeaconBlockBody{},
		},
		{
			name: "BeaconBlockHeader",
			s:    &phase0.BeaconBlockHeader{},
		},
		{
			name: "BeaconState",
			s:    &bellatrix.BeaconState{},
		},
		{
			name: "Checkpoint",
			s:    &phase0.Checkpoint{},
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
			name: "Eth1Data",
			s:    &phase0.ETH1Data{},
		},
		{
			name: "ExecutionPayload",
			s:    &bellatrix.ExecutionPayload{},
		},
		{
			name: "ExecutionPayloadHeader",
			s:    &bellatrix.ExecutionPayloadHeader{},
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
			name: "IndexedAttestation",
			s:    &phase0.IndexedAttestation{},
		},
		{
			name: "PendingAttestation",
			s:    &phase0.PendingAttestation{},
		},
		{
			name: "ProposerSlashing",
			s:    &phase0.ProposerSlashing{},
		},
		{
			name: "SignedAggregateAndProof",
			s:    &phase0.SignedAggregateAndProof{},
		},
		{
			name: "SignedBeaconBlock",
			s:    &bellatrix.SignedBeaconBlock{},
		},
		{
			name: "SignedBeaconBlockHeader",
			s:    &phase0.SignedBeaconBlockHeader{},
		},
		{
			name: "SignedContributionAndproof",
			s:    &altair.SignedContributionAndProof{},
		},
		{
			name: "SignedVoluntaryExit",
			s:    &phase0.SignedVoluntaryExit{},
		},
		{
			name: "SyncAggregate",
			s:    &altair.SyncAggregate{},
		},
		{
			name: "SyncCommitteeContribuion",
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
	}

	baseDir := filepath.Join(os.Getenv("CONSENSUS_SPEC_TESTS_DIR"), "tests", "mainnet", "bellatrix", "ssz_static")
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
					require.Equal(t, testYAMLFormat(specYAML), testYAMLFormat(remarshalledSpecYAML))

					// Obtain the struct from the SSZ.
					s2 := clone.Clone(test.s)
					compressedSpecSSZ, err := os.ReadFile(filepath.Join(path, "serialized.ssz_snappy"))
					require.NoError(t, err)
					var specSSZ []byte
					specSSZ, err = snappy.Decode(specSSZ, compressedSpecSSZ)
					require.NoError(t, err)
					require.NoError(t, s2.(ssz.Unmarshaler).UnmarshalSSZ(specSSZ))
					// Confirm we can return to the SSZ.
					remarshalledSpecSSZ, err := s2.(ssz.Marshaler).MarshalSSZ()
					require.NoError(t, err)
					require.Equal(t, specSSZ, remarshalledSpecSSZ)

					// Obtain the hash tree root from the YAML.
					specYAMLRoot, err := os.ReadFile(filepath.Join(path, "roots.yaml"))
					require.NoError(t, err)
					// Confirm we calculate the same root.
					generatedRootBytes, err := s2.(ssz.HashRoot).HashTreeRoot()
					require.NoError(t, err)
					generatedRoot := fmt.Sprintf("{root: '%#x'}\n", string(generatedRootBytes[:]))
					require.Equal(t, string(specYAMLRoot), generatedRoot)
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
	}
	for _, replacement := range replacements {
		res = bytes.ReplaceAll(res, replacement[0], replacement[1])
	}

	return string(bytes.ToLower(res))
}
