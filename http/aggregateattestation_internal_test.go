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

package http

import (
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAggregateAttestationDecode(t *testing.T) {
	responseJson := `{
  "version": "Electra",
  "data": {
    "aggregation_bits": "0x97aff9afffedbbfefedbdffdfbf5ffaebfffffffecfd03",
    "data": {
      "slot": "84434",
      "index": "0",
      "beacon_block_root": "0xaa95c9d1a4f380b4331378e92ba88f4c757c6d252e29d43e6c8ac804caccca9a",
      "source": {
        "epoch": "2637",
        "root": "0x22aa73e2e76e27404e4bf259d27012faa9f5a2e6e7c611fdfe32510b66470b82"
      },
      "target": {
        "epoch": "2638",
        "root": "0x4f6545fcd8b1e24daeb6872dfe42898ba0e4be917f459b9c42f6fbb56715699c"
      }
    },
    "signature": "0xad0ea974c685ff2c8d8971e456569c9b5f8ba830f86154172bcfc77142f65ce64d866cf1b15adbb2dbb8f2958bb952bb061b66aaacbb2e764193d2fee666a1d2fefe34c41932fbf07522be456c1b768ceb8899ba5bd107a1a3700f58d6414d54",
    "committee_bits": "0x0100000000000000"
  }
}`

	t.Run("ElectraAttestationAggregate", func(t *testing.T) {
		response := httpResponse{
			// consensusVersion: spec.DataVersionElectra,
			body: []byte(responseJson),
		}
		data, metadata, err := decodeAggregateAttestation(&response)
		require.NoError(t, err)

		require.Equal(t, data.Version, spec.DataVersionElectra)
		versionKey, ok := metadata["version"]
		require.True(t, ok)
		require.Equal(t, versionKey, "Electra")
	})
}
