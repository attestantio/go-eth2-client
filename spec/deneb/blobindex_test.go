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

// Create a test to verify blobindex.unmarshalJSON
import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/deneb"
	require "github.com/stretchr/testify/require"
)

func TestBlobIndexUnmarshalJSON(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		input    []byte
		expected deneb.BlobIndex
		wantErr  bool
	}{
		{
			name:     "Valid input",
			input:    []byte("\"1\""),
			expected: deneb.BlobIndex(1),
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
		err := tt.expected.UnmarshalJSON(tt.input)
		if tt.wantErr {
			require.Error(t, err, tt.name)
		} else {
			require.NoError(t, err, tt.name)
			require.Equal(t, tt.expected, tt.expected, tt.name)
		}
	}
}
