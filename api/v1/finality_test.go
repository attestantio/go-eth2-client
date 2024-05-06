// Copyright © 2020 Attestant Limited.
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

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	require "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestFinalityJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.finalityJSON",
		},
		{
			name:  "FinalizedMissing",
			input: []byte(`{"current_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"},"previous_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"}}`),
			err:   "finalized checkpoint missing",
		},
		{
			name:  "FinalizedWrongType",
			input: []byte(`{"finalized":true,"current_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"},"previous_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"}}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.checkpointJSON",
		},
		{
			name:  "FinalizedInvalid",
			input: []byte(`{"finalized":{},"current_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"},"previous_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"}}`),
			err:   "invalid JSON: epoch missing",
		},
		{
			name:  "CurrentJustifiedMissing",
			input: []byte(`{"finalized":{"epoch":"15614","root":"0xb3806428b52a802fb9c4355b6e93a6afde02ecbd27a9f4723eb427c27cadb440"},"previous_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"}}`),
			err:   "justified checkpoint missing",
		},
		{
			name:  "CurrentJustifiedWrongType",
			input: []byte(`{"finalized":{"epoch":"15614","root":"0xb3806428b52a802fb9c4355b6e93a6afde02ecbd27a9f4723eb427c27cadb440"},"current_justified":true,"previous_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"}}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.checkpointJSON",
		},
		{
			name:  "CurrentJustifiedInvalid",
			input: []byte(`{"finalized":{"epoch":"15614","root":"0xb3806428b52a802fb9c4355b6e93a6afde02ecbd27a9f4723eb427c27cadb440"},"current_justified":{},"previous_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"}}`),
			err:   "invalid JSON: epoch missing",
		},
		{
			name:  "PreviousJustifiedMissing",
			input: []byte(`{"finalized":{"epoch":"15614","root":"0xb3806428b52a802fb9c4355b6e93a6afde02ecbd27a9f4723eb427c27cadb440"},"current_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"}}`),
			err:   "previous justified checkpoint missing",
		},
		{
			name:  "PreviousJustifiedWrongType",
			input: []byte(`{"finalized":{"epoch":"15614","root":"0xb3806428b52a802fb9c4355b6e93a6afde02ecbd27a9f4723eb427c27cadb440"},"current_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"},"previous_justified":true}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.checkpointJSON",
		},
		{
			name:  "PreviousJustifiedInvalid",
			input: []byte(`{"finalized":{"epoch":"15614","root":"0xb3806428b52a802fb9c4355b6e93a6afde02ecbd27a9f4723eb427c27cadb440"},"current_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"},"previous_justified":{}}`),
			err:   "invalid JSON: epoch missing",
		},
		{
			name:  "Good",
			input: []byte(`{"finalized":{"epoch":"15614","root":"0xb3806428b52a802fb9c4355b6e93a6afde02ecbd27a9f4723eb427c27cadb440"},"current_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"},"previous_justified":{"epoch":"15705","root":"0x66ba71dfb29bada27c3f99e9823dac4272ff1a057814d0672353358571cb0142"}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.Finality
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}
