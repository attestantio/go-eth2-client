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

package phase0_test

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

func TestEpochUnmarshalJSON(t *testing.T) {
	tests := []struct {
		label    string
		input    []byte
		expected phase0.Epoch
		wantErr  bool
	}{
		{
			label:    "Valid quoted",
			input:    []byte("\"42\""),
			expected: phase0.Epoch(42),
		},
		{
			label:    "Invalid text",
			input:    []byte("not-a-number"),
			expected: 0,
			wantErr:  true,
		},
		{
			label:    "Invalid single quote",
			input:    []byte("\""),
			expected: 0,
			wantErr:  true,
		},
		// Caplin emits bare uint64 epoch fields (e.g.
		// PendingPartialWithdrawal.withdrawable_epoch); accept both.
		{
			label:    "Caplin bare number",
			input:    []byte("42"),
			expected: phase0.Epoch(42),
		},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			var e phase0.Epoch
			err := e.UnmarshalJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if e != tt.expected {
				t.Errorf("UnmarshalJSON() got = %v, expected %v", e, tt.expected)
			}
		})
	}
}
