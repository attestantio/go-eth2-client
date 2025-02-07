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

package spec_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/stretchr/testify/assert"
)

func TestVersionedAttestation_CommitteeIndex(t *testing.T) {
	// Test cases
	tests := []struct {
		name            string
		expectedIndices []phase0.CommitteeIndex
		errorMsg        string
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
			expectedIndices: []phase0.CommitteeIndex{64},
			errorMsg:        "no committee index found in committee bits",
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
			for _, expectedIndex := range tt.expectedIndices {
				committeeBits.SetBitAt(uint64(expectedIndex), true)
			}
			attestation := spec.VersionedAttestation{
				Version: spec.DataVersionElectra,
				Electra: &electra.Attestation{
					CommitteeBits: committeeBits,
				},
			}
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
