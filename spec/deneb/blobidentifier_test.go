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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/goccy/go-yaml"
	"github.com/golang/snappy"
	require "github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestBlobIdentifierJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		output []byte
		err    string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type deneb.blobIdentifierJSON",
		},
		{
			name:  "BlockRootMissing",
			input: []byte(`{"index":"17189299593882149153"}`),
			err:   "block root missing",
		},
		{
			name:  "BlockRootWrongType",
			input: []byte(`{"block_root":true,"index":"17189299593882149153"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blobIdentifierJSON.block_root of type string",
		},
		{
			name:  "BlockRootInvalid",
			input: []byte(`{"block_root":"true","index":"17189299593882149153"}`),
			err:   "invalid value for block root: encoding/hex: invalid byte: U+0074 't'",
		},
		{
			name:  "BlockRootIncorrectLength",
			input: []byte(`{"block_root":"0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec7","index":"17189299593882149153"}`),
			err:   "incorrect length for block root",
		},
		{
			name:  "IndexMissing",
			input: []byte(`{"block_root":"0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec799"}`),
			err:   "index missing",
		},
		{
			name:  "IndexWrongType",
			input: []byte(`{"block_root":"0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec799","index":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field blobIdentifierJSON.index of type string",
		},
		{
			name:  "IndexInvalid",
			input: []byte(`{"block_root":"0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec799","index":"-1"}`),
			err:   "invalid value for index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"block_root":"0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec799","index":"17189299593882149153"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res deneb.BlobIdentifier
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				if len(test.output) > 0 {
					assert.Equal(t, string(test.output), string(rt))
				} else {
					assert.Equal(t, string(test.input), string(rt))
				}
			}
		})
	}
}

func TestBlobIdentifierYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{block_root: '0x813b05d7c10dc4bdf45201a3539ec805ff4e016fbadd98a8b24cbf1f428ec799', index: 17189299593882149153}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res deneb.BlobIdentifier
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, res.String(), string(rt))
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestBlobIdentifierSpec(t *testing.T) {
	if os.Getenv("ETH2_SPEC_TESTS_DIR") == "" {
		t.Skip("ETH2_SPEC_TESTS_DIR not suppplied, not running spec tests")
	}
	baseDir := filepath.Join(os.Getenv("ETH2_SPEC_TESTS_DIR"), "tests", "mainnet", "deneb", "ssz_static", "BlobIdentifier", "ssz_random")
	require.NoError(t, filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if path == baseDir {
			// Only interested in subdirectories.
			return nil
		}
		require.NoError(t, err)
		if info.IsDir() {
			t.Run(info.Name(), func(t *testing.T) {
				specYAML, err := os.ReadFile(filepath.Join(path, "value.yaml"))
				require.NoError(t, err)
				var res deneb.BlobIdentifier
				require.NoError(t, yaml.Unmarshal(specYAML, &res))

				compressedSpecSSZ, err := os.ReadFile(filepath.Join(path, "serialized.ssz_snappy"))
				require.NoError(t, err)
				var specSSZ []byte
				specSSZ, err = snappy.Decode(specSSZ, compressedSpecSSZ)
				require.NoError(t, err)

				unmarshalled := &deneb.BlobIdentifier{}
				require.NoError(t, unmarshalled.UnmarshalSSZ(specSSZ))
				remarshalled, err := unmarshalled.MarshalSSZ()
				require.NoError(t, err)
				require.Equal(t, specSSZ, remarshalled)

				ssz, err := res.MarshalSSZ()
				require.NoError(t, err)
				require.Equal(t, specSSZ, ssz)

				root, err := res.HashTreeRoot()
				require.NoError(t, err)
				rootsYAML, err := os.ReadFile(filepath.Join(path, "roots.yaml"))
				require.NoError(t, err)
				require.Equal(t, string(rootsYAML), fmt.Sprintf("{root: '%#x'}\n", root))
			})
		}
		return nil
	}))
}
