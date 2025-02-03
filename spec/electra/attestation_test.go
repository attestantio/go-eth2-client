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

package electra_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/stretchr/testify/assert"
)

func TestAttestation_CommitteeIndex(t *testing.T) {
	// Test cases
	tests := []struct {
		name            string
		expectedIndices []phase0.CommitteeIndex
		errorMsg        string
		doNotSet        bool
	}{
		{
			name:            "Valid index 0",
			expectedIndices: []phase0.CommitteeIndex{0},
		},
		{
			name:            "Valid index 4",
			expectedIndices: []phase0.CommitteeIndex{4},
		},
		{
			name:            "Valid index 40",
			expectedIndices: []phase0.CommitteeIndex{40},
		},
		{
			name:            "Invalid index 64",
			expectedIndices: []phase0.CommitteeIndex{64},
			errorMsg:        "no committee index found in committee bits",
		},
		{
			name:            "Invalid no index set",
			expectedIndices: []phase0.CommitteeIndex{4},
			errorMsg:        "no committee index found in committee bits",
			doNotSet:        true,
		},
		{
			name:            "Invalid multiple index set",
			expectedIndices: []phase0.CommitteeIndex{4, 40},
			errorMsg:        "multiple committee indices found in committee bits",
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			committeeBits := bitfield.NewBitvector64()
			if !tt.doNotSet {
				for _, expectedIndex := range tt.expectedIndices {
					committeeBits.SetBitAt(uint64(expectedIndex), true)
				}
			}
			attestation := &electra.Attestation{CommitteeBits: committeeBits}

			committeeIndex, err := attestation.CommitteeIndex()
			if tt.errorMsg == "" {
				require.NoError(t, err)
				for _, expectedIndex := range tt.expectedIndices {
					assert.Equal(t, expectedIndex, committeeIndex)
				}
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestAttestation_AggregateValidatorIndex(t *testing.T) {
	// Test cases
	tests := []struct {
		name            string
		expectedIndices []phase0.ValidatorIndex
		errorMsg        string
		doNotSet        bool
	}{
		{
			name:            "Valid index 0",
			expectedIndices: []phase0.ValidatorIndex{0},
		},
		{
			name:            "Valid index 4",
			expectedIndices: []phase0.ValidatorIndex{4},
		},
		{
			name:            "Valid index 140",
			expectedIndices: []phase0.ValidatorIndex{140},
		},
		{
			name:            "Invalid index 160",
			expectedIndices: []phase0.ValidatorIndex{160},
			errorMsg:        "no validator index found in aggregation bits",
		},
		{
			name:            "Invalid no index set",
			expectedIndices: []phase0.ValidatorIndex{64},
			errorMsg:        "no validator index found in aggregation bits",
			doNotSet:        true,
		},
		{
			name:            "Invalid multiple index set",
			expectedIndices: []phase0.ValidatorIndex{4, 40},
			errorMsg:        "multiple validator indices found in aggregation bits",
		},
	}
	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aggregateBits := bitfield.NewBitlist(160)
			if !tt.doNotSet {
				for _, expectedIndex := range tt.expectedIndices {
					aggregateBits.SetBitAt(uint64(expectedIndex), true)
				}
			}
			attestation := &electra.Attestation{AggregationBits: aggregateBits}

			aggregateValidatorIndex, err := attestation.AggregateValidatorIndex()
			if tt.errorMsg == "" {
				require.NoError(t, err)
				for _, expectedIndex := range tt.expectedIndices {
					assert.Equal(t, expectedIndex, aggregateValidatorIndex)
				}
				return
			}
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestAttestation_SSZ(t *testing.T) {
	aggregateSize := uint64(131072)
	aggregateBits := bitfield.NewBitlist(aggregateSize)
	committeeSize := uint64(64)
	committeeBits := bitfield.NewBitvector64()
	attestation := electra.Attestation{
		AggregationBits: aggregateBits,
		CommitteeBits:   committeeBits,
	}
	// Set a bit beyond the bit list.
	aggregateBits.SetBitAt(aggregateSize, true)
	// Set a bit in the last bit of the list.
	aggregateBits.SetBitAt(aggregateSize-1, true)

	// Should only have the bit that was set on the last bit.
	require.Equal(t, 1, len(aggregateBits.BitIndices()))
	require.Equal(t, aggregateSize, aggregateBits.Len())

	// Set a bit beyond the bit vector.
	committeeBits.SetBitAt(committeeSize, true)
	// Set a bit in the last bit of the vector.
	committeeBits.SetBitAt(committeeSize-1, true)

	// Should only have the bit that was set on the last bit.
	require.Equal(t, 1, len(committeeBits.BitIndices()))

	// Ensure we can actually serialise.
	_, err := attestation.MarshalSSZ()
	require.NoError(t, err)
}
