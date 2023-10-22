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

func TestBeaconBlockHeaderJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.beaconBlockHeaderJSON",
		},
		{
			name:  "RootMissing",
			input: []byte(`{"canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"0xa8d684242ee025ee96e877b28433d93176072b8c8e8295609501863147bb1d174b8a16aed661d001f30859c9e42c0f9d18ea35786a9bdf115dff1877980046e19e0e4c9310e281f8129f2692ddc4680673ab78b7f8db72f91be7863dd9fe1e55"}}`),
			err:   "root missing",
		},
		{
			name:  "RootWrongType",
			input: []byte(`{"root":true,"canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"0xa8d684242ee025ee96e877b28433d93176072b8c8e8295609501863147bb1d174b8a16aed661d001f30859c9e42c0f9d18ea35786a9bdf115dff1877980046e19e0e4c9310e281f8129f2692ddc4680673ab78b7f8db72f91be7863dd9fe1e55"}}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field beaconBlockHeaderJSON.root of type string",
		},
		{
			name:  "RootInvalid",
			input: []byte(`{"root":"invalid","canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"0xa8d684242ee025ee96e877b28433d93176072b8c8e8295609501863147bb1d174b8a16aed661d001f30859c9e42c0f9d18ea35786a9bdf115dff1877980046e19e0e4c9310e281f8129f2692ddc4680673ab78b7f8db72f91be7863dd9fe1e55"}}`),
			err:   "invalid value for root: encoding/hex: invalid byte: U+0069 'i'",
		},
		{
			name:  "RootShort",
			input: []byte(`{"root":"0x354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"0xa8d684242ee025ee96e877b28433d93176072b8c8e8295609501863147bb1d174b8a16aed661d001f30859c9e42c0f9d18ea35786a9bdf115dff1877980046e19e0e4c9310e281f8129f2692ddc4680673ab78b7f8db72f91be7863dd9fe1e55"}}`),
			err:   "incorrect length 31 for root",
		},
		{
			name:  "RootLong",
			input: []byte(`{"root":"0xbcbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"0xa8d684242ee025ee96e877b28433d93176072b8c8e8295609501863147bb1d174b8a16aed661d001f30859c9e42c0f9d18ea35786a9bdf115dff1877980046e19e0e4c9310e281f8129f2692ddc4680673ab78b7f8db72f91be7863dd9fe1e55"}}`),
			err:   "incorrect length 33 for root",
		},
		{
			name:  "CanonicalWrongType",
			input: []byte(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":"true","header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"0xa8d684242ee025ee96e877b28433d93176072b8c8e8295609501863147bb1d174b8a16aed661d001f30859c9e42c0f9d18ea35786a9bdf115dff1877980046e19e0e4c9310e281f8129f2692ddc4680673ab78b7f8db72f91be7863dd9fe1e55"}}`),
			err:   "invalid JSON: json: cannot unmarshal string into Go struct field beaconBlockHeaderJSON.canonical of type bool",
		},
		{
			name:  "CanonicalInvalid",
			input: []byte(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":maybe,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"0xa8d684242ee025ee96e877b28433d93176072b8c8e8295609501863147bb1d174b8a16aed661d001f30859c9e42c0f9d18ea35786a9bdf115dff1877980046e19e0e4c9310e281f8129f2692ddc4680673ab78b7f8db72f91be7863dd9fe1e55"}}`),
			err:   "invalid character 'm' looking for beginning of value",
		},
		{
			name:  "HeaderMissing",
			input: []byte(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true}`),
			err:   "header missing",
		},
		{
			name:  "HeaderWrongType",
			input: []byte(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true,"header":true}`),
			err:   "invalid JSON: invalid JSON: json: cannot unmarshal bool into Go value of type phase0.signedBeaconBlockHeaderJSON",
		},
		{
			name:  "HeaderInvalid",
			input: []byte(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true,"header":{}}`),
			err:   "invalid JSON: message missing",
		},
		{
			name:  "Good",
			input: []byte(`{"root":"0xbc354f1a5f27f8d096eee9e6b6139e1b730385f9752513832a57c9849a149df7","canonical":true,"header":{"message":{"slot":"585321","proposer_index":"29787","parent_root":"0xba4d784293df28bab771a14df58cdbed9d8d64afd0ddf1c52dff3e25fcdd51df","state_root":"0x4e405274abd4f59c6a2268b4e6ca93dba01e15ae6b56401fb20a1ad9701b036d","body_root":"0x57bb79520694c132a35dc887cac2e4dad9acc5ded58b5ae66b491644ab8835c8"},"signature":"0xa8d684242ee025ee96e877b28433d93176072b8c8e8295609501863147bb1d174b8a16aed661d001f30859c9e42c0f9d18ea35786a9bdf115dff1877980046e19e0e4c9310e281f8129f2692ddc4680673ab78b7f8db72f91be7863dd9fe1e55"}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.BeaconBlockHeader
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
