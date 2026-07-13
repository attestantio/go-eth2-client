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

package altair_test

import (
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseParticipationFlags(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []altair.ParticipationFlags
		wantErr string
	}{
		{
			name:  "Spec form: array of decimal strings",
			input: `["0","3","2","7"]`,
			want:  []altair.ParticipationFlags{0, 3, 2, 7},
		},
		{
			name:  "Spec form: empty array",
			input: `[]`,
			want:  []altair.ParticipationFlags{},
		},
		{
			// Caplin's actual format: SSZ-encoded bytes as a single hex string.
			name:  "Caplin form: 0x-prefixed hex string",
			input: `"0x000302070100"`,
			want:  []altair.ParticipationFlags{0x00, 0x03, 0x02, 0x07, 0x01, 0x00},
		},
		{
			name:  "Caplin form: empty hex",
			input: `"0x"`,
			want:  []altair.ParticipationFlags{},
		},
		{
			name:  "Null",
			input: `null`,
			want:  nil,
		},
		{
			name:    "Spec form: invalid element",
			input:   `["0","abc"]`,
			wantErr: "invalid value",
		},
		{
			name:    "Spec form: out-of-range element (>255)",
			input:   `["0","999"]`,
			wantErr: "invalid value",
		},
		{
			name:    "Caplin form: invalid hex",
			input:   `"0xzz"`,
			wantErr: "invalid hex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := altair.ParseParticipationFlags(json.RawMessage(tt.input), "field")
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)

				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMarshalParticipationFlags_Roundtrip(t *testing.T) {
	original := []altair.ParticipationFlags{0, 3, 2, 7, 1, 0}

	marshaled, err := altair.MarshalParticipationFlags(original)
	require.NoError(t, err)
	assert.Equal(t, `["0","3","2","7","1","0"]`, string(marshaled))

	parsed, err := altair.ParseParticipationFlags(marshaled, "field")
	require.NoError(t, err)
	assert.Equal(t, original, parsed)
}
