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
	"testing"

	"github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestSignedBlindedBlobSidecarJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type map[string]json.RawMessage",
		},
		{
			name:  "MessageMissing",
			input: []byte(`{"signature":"0x8c3095fd9d3a18e43ceeb7648281e16bb03044839dffea796432c4e5a1372bef22c11a98a31e0c1c5389b98cc6d45917170a0f1634bcf152d896f360dc599fabba2ec4de77898b5dff080fa1628482bdbad5b37d2e64fea3d8721095186cfe50"}`),
			err:   "message: missing",
		},
		{
			name:  "MessageWrongType",
			input: []byte(`{"message":true,"signature":"0x8c3095fd9d3a18e43ceeb7648281e16bb03044839dffea796432c4e5a1372bef22c11a98a31e0c1c5389b98cc6d45917170a0f1634bcf152d896f360dc599fabba2ec4de77898b5dff080fa1628482bdbad5b37d2e64fea3d8721095186cfe50"}`),
			err:   "message: invalid JSON: json: cannot unmarshal bool into Go value of type map[string]json.RawMessage",
		},
		{
			name:  "MessageInvalid",
			input: []byte(`{"message":{},"signature":"0x8c3095fd9d3a18e43ceeb7648281e16bb03044839dffea796432c4e5a1372bef22c11a98a31e0c1c5389b98cc6d45917170a0f1634bcf152d896f360dc599fabba2ec4de77898b5dff080fa1628482bdbad5b37d2e64fea3d8721095186cfe50"}`),
			err:   "message: block_root: missing",
		},
		{
			name:  "SignatureMissing",
			input: []byte(`{"message":{"block_root":"0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29","index":"17762875709721895328","slot":"12231583639632491026","block_parent_root":"0x22de86edc38dc56c4255cba641c83251a2a2dcc7535e773c9a2fb2e8b73758a4","proposer_index":"16148839969926959295","blob_root":"0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29","kzg_commitment":"0x0748ac5c58e66b1fae24289f9014948876fbd78da88931bb6cbcd2e44a01bd07ab4f33e54ec9b9a2ada2e83c840dceb6","kzg_proof":"0xc6e27a3ae80243ba7ea88eab107a0675020e0745d75ab6a1553691007a50f7f99f597693ac33ae3cea63bf0b90a734ff"}}`),
			err:   "signature: missing",
		},
		{
			name:  "SignatureWrongType",
			input: []byte(`{"message":{"block_root":"0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29","index":"17762875709721895328","slot":"12231583639632491026","block_parent_root":"0x22de86edc38dc56c4255cba641c83251a2a2dcc7535e773c9a2fb2e8b73758a4","proposer_index":"16148839969926959295","blob_root":"0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29","kzg_commitment":"0x0748ac5c58e66b1fae24289f9014948876fbd78da88931bb6cbcd2e44a01bd07ab4f33e54ec9b9a2ada2e83c840dceb6","kzg_proof":"0xc6e27a3ae80243ba7ea88eab107a0675020e0745d75ab6a1553691007a50f7f99f597693ac33ae3cea63bf0b90a734ff"},"signature":true}`),
			err:   "signature: invalid prefix",
		},
		{
			name:  "SignatureInvalid",
			input: []byte(`{"message":{"block_root":"0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29","index":"17762875709721895328","slot":"12231583639632491026","block_parent_root":"0x22de86edc38dc56c4255cba641c83251a2a2dcc7535e773c9a2fb2e8b73758a4","proposer_index":"16148839969926959295","blob_root":"0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29","kzg_commitment":"0x0748ac5c58e66b1fae24289f9014948876fbd78da88931bb6cbcd2e44a01bd07ab4f33e54ec9b9a2ada2e83c840dceb6","kzg_proof":"0xc6e27a3ae80243ba7ea88eab107a0675020e0745d75ab6a1553691007a50f7f99f597693ac33ae3cea63bf0b90a734ff"},"signature":"0xic3095fd9d3a18e43ceeb7648281e16bb03044839dffea796432c4e5a1372bef22c11a98a31e0c1c5389b98cc6d45917170a0f1634bcf152d896f360dc599fabba2ec4de77898b5dff080fa1628482bdbad5b37d2e64fea3d8721095186cfe50"}`),
			err:   "signature: invalid value ic3095fd9d3a18e43ceeb7648281e16bb03044839dffea796432c4e5a1372bef22c11a98a31e0c1c5389b98cc6d45917170a0f1634bcf152d896f360dc599fabba2ec4de77898b5dff080fa1628482bdbad5b37d2e64fea3d8721095186cfe50: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SignatureIncorrectLength",
			input: []byte(`{"message":{"block_root":"0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29","index":"17762875709721895328","slot":"12231583639632491026","block_parent_root":"0x22de86edc38dc56c4255cba641c83251a2a2dcc7535e773c9a2fb2e8b73758a4","proposer_index":"16148839969926959295","blob_root":"0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29","kzg_commitment":"0x0748ac5c58e66b1fae24289f9014948876fbd78da88931bb6cbcd2e44a01bd07ab4f33e54ec9b9a2ada2e83c840dceb6","kzg_proof":"0xc6e27a3ae80243ba7ea88eab107a0675020e0745d75ab6a1553691007a50f7f99f597693ac33ae3cea63bf0b90a734ff"},"signature":"0x8c3095fd9d3a18e43ceeb7648281e16bb03044839dffea796432c4e5a1372bef22c11a98a31e0c1c5389b98cc6d45917170a0f1634bcf152d896f360dc599fabba2ec4de77898b5dff080fa1628482bdbad5b37d2e64fea3d8721095186cfe5"}`),
			err:   "signature: incorrect length",
		},
		{
			name:  "Good",
			input: []byte(`{"message":{"block_root":"0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29","index":"17762875709721895328","slot":"12231583639632491026","block_parent_root":"0x22de86edc38dc56c4255cba641c83251a2a2dcc7535e773c9a2fb2e8b73758a4","proposer_index":"16148839969926959295","blob_root":"0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29","kzg_commitment":"0x0748ac5c58e66b1fae24289f9014948876fbd78da88931bb6cbcd2e44a01bd07ab4f33e54ec9b9a2ada2e83c840dceb6","kzg_proof":"0xc6e27a3ae80243ba7ea88eab107a0675020e0745d75ab6a1553691007a50f7f99f597693ac33ae3cea63bf0b90a734ff"},"signature":"0x8c3095fd9d3a18e43ceeb7648281e16bb03044839dffea796432c4e5a1372bef22c11a98a31e0c1c5389b98cc6d45917170a0f1634bcf152d896f360dc599fabba2ec4de77898b5dff080fa1628482bdbad5b37d2e64fea3d8721095186cfe50"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res deneb.SignedBlindedBlobSidecar
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestSignedBlindedBlobSidecarYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		root  []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{message: {block_root: '0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29', index: 17762875709721895328, slot: 12231583639632491026, block_parent_root: '0x22de86edc38dc56c4255cba641c83251a2a2dcc7535e773c9a2fb2e8b73758a4', proposer_index: 16148839969926959295, blob_root: '0x3c1820c62034fc45c10abc983dbce08de28f303192dea32371a902b3e6a1fc29', kzg_commitment: '0x0748ac5c58e66b1fae24289f9014948876fbd78da88931bb6cbcd2e44a01bd07ab4f33e54ec9b9a2ada2e83c840dceb6', kzg_proof: '0xc6e27a3ae80243ba7ea88eab107a0675020e0745d75ab6a1553691007a50f7f99f597693ac33ae3cea63bf0b90a734ff'}, signature: '0x8c3095fd9d3a18e43ceeb7648281e16bb03044839dffea796432c4e5a1372bef22c11a98a31e0c1c5389b98cc6d45917170a0f1634bcf152d896f360dc599fabba2ec4de77898b5dff080fa1628482bdbad5b37d2e64fea3d8721095186cfe50'}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res deneb.SignedBlindedBlobSidecar
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
