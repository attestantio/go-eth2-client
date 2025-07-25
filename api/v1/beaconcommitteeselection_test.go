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

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestBeaconCommitteeSelectionJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.beaconCommitteeSelectionJSON",
		},
		{
			name:  "ValidatorIndexMissing",
			input: []byte(`{"slot":"1","selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "validator index missing",
		},
		{
			name:  "ValidatorIndexWrongType",
			input: []byte(`{"validator_index":true,"slot":"1","selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconCommitteeSelectionJSON.validator_index of type string",
		},
		{
			name:  "ValidatorIndexInvalid",
			input: []byte(`{"validator_index":"invalid","slot":"1","selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid value for validator index: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"validator_index":"10","selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"validator_index":"10","slot":true,"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconCommitteeSelectionJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"validator_index":"10","slot":"-1","selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "SelectionProofMissing",
			input: []byte(`{"validator_index":"10","slot":"1"}`),
			err:   "selection proof missing",
		},
		{
			name:  "SelectionProofWrongType",
			input: []byte(`{"validator_index":"10","slot":"1","selection_proof":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconCommitteeSelectionJSON.selection_proof of type string",
		},
		{
			name:  "SelectionProofInvalid",
			input: []byte(`{"validator_index":"10","slot":"1","selection_proof":"invalid"}`),
			err:   "invalid value for selection proof: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "SelectionProofShort",
			input: []byte(`{"validator_index":"10","slot":"1","selection_proof":"0x5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "incorrect length for selection proof",
		},
		{
			name:  "SelectionProofLong",
			input: []byte(`{"validator_index":"10","slot":"1","selection_proof":"0x8b8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
			err:   "incorrect length for selection proof",
		},
		{
			name:  "Good",
			input: []byte(`{"validator_index":"10","slot":"1","selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.BeaconCommitteeSelection
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
