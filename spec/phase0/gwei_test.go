// Copyright © 2020, 2021 Attestant Limited.
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

package phase0_test

// Create a test to verify gwei.unmarshalJSON
import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

func TestGweiUnmarshalJSON(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		input    []byte
		expected phase0.Gwei
		wantErr  bool
	}{
		{
			name:     "Valid input 1000000000",
			input:    []byte("\"1000000000\""),
			expected: phase0.Gwei(1000000000),
			wantErr:  false,
		},

		{
			name:     "Valid input",
			input:    []byte("\"1\""),
			expected: phase0.Gwei(1),
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
		t.Run(tt.name, func(t *testing.T) {
			var g phase0.Gwei
			err := g.UnmarshalJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if g != tt.expected {
				t.Errorf("UnmarshalJSON() got = %v, expected %v", g, tt.expected)
			}
		})
	}
}
