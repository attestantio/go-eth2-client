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

package deneb_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/deneb"
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
			name: "BeaconBlockBody",
			s:    &deneb.BeaconBlockBody{},
		},
		{
			name: "BeaconBlock",
			s:    &deneb.BeaconBlock{},
		},
		{
			name: "BlobIdentifier",
			s:    &deneb.BlobIdentifier{},
		},
		{
			name: "BlobSidecar",
			s:    &deneb.BlobSidecar{},
		},
		{
			name: "ExecutionPayload",
			s:    &deneb.ExecutionPayload{},
		},
		{
			name: "ExecutionPayloadHeader",
			s:    &deneb.ExecutionPayloadHeader{},
		},
		{
			name: "SignedBeaconBlock",
			s:    &deneb.SignedBeaconBlock{},
		},
	}

	baseDir := filepath.Join(os.Getenv("CONSENSUS_SPEC_TESTS_DIR"), "tests", "mainnet", "deneb", "ssz_static")
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
