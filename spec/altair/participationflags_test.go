// Copyright © 2021 Attestant Limited.
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

package altair_test

// create a test for participationflags.UnmarshalJSON
import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	require "github.com/stretchr/testify/require"
)

func TestParticipationFlagsUnmarshalJSON(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		input    []byte
		expected altair.ParticipationFlags
		wantErr  bool
	}{
		{
			name:     "Valid input 1",
			input:    []byte("\"1\""),
			expected: altair.ParticipationFlags(1),
			wantErr:  false,
		},

		{
			name:     "Valid input 2",
			input:    []byte("\"2\""),
			expected: altair.ParticipationFlags(2),
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
		// Run test
		t.Run(tt.name, func(t *testing.T) {
			var participationFlags altair.ParticipationFlags
			err := participationFlags.UnmarshalJSON(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, participationFlags)
			}
		})
	}
}
