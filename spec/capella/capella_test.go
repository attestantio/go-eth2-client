// Copyright © 2022, 2023 Attestant Limited.
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

package capella_test

// Create a test to verify withdrawalindex.unmarshalJSON
import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/capella"
	require "github.com/stretchr/testify/require"
)

func TestWithdrawalIndexUnmarshalJSON(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		input    []byte
		expected capella.WithdrawalIndex
		wantErr  bool
	}{
		{
			name:     "Valid input 100000",
			input:    []byte("\"100000\""),
			expected: capella.WithdrawalIndex(100000),
			wantErr:  false,
		},

		{
			name:     "Valid input",
			input:    []byte("\"1\""),
			expected: capella.WithdrawalIndex(1),
			wantErr:  false,
		},

		{
			name:     "Invalid input text",
			input:    []byte("not-a-number"),
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "Invalid input single quote",
			input:    []byte("\""),
			expected: 0,
			wantErr:  true,
		},
	}

	// Run tests
	for _, tt := range tests {
		// Run the test
		t.Run(tt.name, func(t *testing.T) {
			var withdrawalIndex capella.WithdrawalIndex
			err := withdrawalIndex.UnmarshalJSON(tt.input)
			if tt.wantErr {
				require.NotNil(t, err, "UnmarshalJSON() did not return an error")
				require.Equal(t, tt.expected, withdrawalIndex, "UnmarshalJSON() returned incorrect value")
			} else {
				require.Nil(t, err, "UnmarshalJSON() returned an error")
				require.Equal(t, tt.expected, withdrawalIndex, "UnmarshalJSON() returned incorrect value")
			}
		})
	}
}
