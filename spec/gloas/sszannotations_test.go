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

package gloas_test

import (
	"reflect"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

// TestGloasSSZFieldAnnotations pins SSZ annotations whose values are not
// otherwise apparent from the surrounding field declaration.
func TestGloasSSZFieldAnnotations(t *testing.T) {
	tests := []struct {
		name     string
		typ      reflect.Type
		field    string
		tag      string
		expected string
	}{
		{
			name:     "HistoricalRootsDynamicLimit",
			typ:      reflect.TypeFor[gloas.BeaconState](),
			field:    "HistoricalRoots",
			tag:      "dynssz-max",
			expected: "HISTORICAL_ROOTS_LIMIT",
		},
		{
			name:     "HistoricalSummariesDynamicLimit",
			typ:      reflect.TypeFor[gloas.BeaconState](),
			field:    "HistoricalSummaries",
			tag:      "dynssz-max",
			expected: "HISTORICAL_ROOTS_LIMIT",
		},
		{
			name:     "RANDAORevealSize",
			typ:      reflect.TypeFor[gloas.BeaconBlockBody](),
			field:    "RANDAOReveal",
			tag:      "ssz-size",
			expected: "96",
		},
		{
			name:     "GraffitiSize",
			typ:      reflect.TypeFor[gloas.BeaconBlockBody](),
			field:    "Graffiti",
			tag:      "ssz-size",
			expected: "32",
		},
		{
			name:     "ParentBlockHashSize",
			typ:      reflect.TypeFor[gloas.ExecutionPayloadBid](),
			field:    "ParentBlockHash",
			tag:      "ssz-size",
			expected: "32",
		},
		{
			name:     "ParentBlockRootSize",
			typ:      reflect.TypeFor[gloas.ExecutionPayloadBid](),
			field:    "ParentBlockRoot",
			tag:      "ssz-size",
			expected: "32",
		},
		{
			name:     "BlockHashSize",
			typ:      reflect.TypeFor[gloas.ExecutionPayloadBid](),
			field:    "BlockHash",
			tag:      "ssz-size",
			expected: "32",
		},
		{
			name:     "PrevRandaoSize",
			typ:      reflect.TypeFor[gloas.ExecutionPayloadBid](),
			field:    "PrevRandao",
			tag:      "ssz-size",
			expected: "32",
		},
		{
			name:     "FeeRecipientSize",
			typ:      reflect.TypeFor[gloas.ExecutionPayloadBid](),
			field:    "FeeRecipient",
			tag:      "ssz-size",
			expected: "20",
		},
		{
			name:     "BeaconBlockRootSize",
			typ:      reflect.TypeFor[gloas.ExecutionPayloadEnvelope](),
			field:    "BeaconBlockRoot",
			tag:      "ssz-size",
			expected: "32",
		},
		{
			name:     "ParentBeaconBlockRootSize",
			typ:      reflect.TypeFor[gloas.ExecutionPayloadEnvelope](),
			field:    "ParentBeaconBlockRoot",
			tag:      "ssz-size",
			expected: "32",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			field, found := test.typ.FieldByName(test.field)
			require.True(t, found)
			require.Equal(t, test.expected, field.Tag.Get(test.tag))
		})
	}
}
