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

package electra_test

import (
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDepositRequestUnmarshalJSON covers spec-compliant input as well as the
// bare-number `amount` / `index` form that Erigon's Caplin emits.
func TestDepositRequestUnmarshalJSON(t *testing.T) {
	const (
		pubkey   = "0xa99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c"
		credsHex = "0x0100000000000000000000000000000000000000000000000000000000000001"
		sigHex   = "0xb9d4d4d18b8ff1d8b8a0c4f4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4b9d4d4d18b8ff1d8b8a0c4f4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4d4"
	)

	tests := []struct {
		name       string
		input      string
		wantAmount phase0.Gwei
		wantIndex  uint64
		wantErr    string
	}{
		{
			name: "spec quoted strings",
			input: `{
				"pubkey":"` + pubkey + `",
				"withdrawal_credentials":"` + credsHex + `",
				"amount":"32000000000",
				"signature":"` + sigHex + `",
				"index":"5"
			}`,
			wantAmount: phase0.Gwei(32000000000),
			wantIndex:  5,
		},
		{
			// Caplin's bare-number form in both amount and index.
			name: "Caplin bare numbers",
			input: `{
				"pubkey":"` + pubkey + `",
				"withdrawal_credentials":"` + credsHex + `",
				"amount":32000000000,
				"signature":"` + sigHex + `",
				"index":5
			}`,
			wantAmount: phase0.Gwei(32000000000),
			wantIndex:  5,
		},
		{
			// Mixed: amount bare, index quoted — defensively accept both.
			name: "Mixed bare amount and quoted index",
			input: `{
				"pubkey":"` + pubkey + `",
				"withdrawal_credentials":"` + credsHex + `",
				"amount":32000000000,
				"signature":"` + sigHex + `",
				"index":"5"
			}`,
			wantAmount: phase0.Gwei(32000000000),
			wantIndex:  5,
		},
		{
			name: "Missing amount",
			input: `{
				"pubkey":"` + pubkey + `",
				"withdrawal_credentials":"` + credsHex + `",
				"signature":"` + sigHex + `",
				"index":"5"
			}`,
			wantErr: "amount missing",
		},
		{
			name: "Missing index",
			input: `{
				"pubkey":"` + pubkey + `",
				"withdrawal_credentials":"` + credsHex + `",
				"amount":"32000000000",
				"signature":"` + sigHex + `"
			}`,
			wantErr: "index missing",
		},
		{
			name: "Invalid bare amount",
			input: `{
				"pubkey":"` + pubkey + `",
				"withdrawal_credentials":"` + credsHex + `",
				"amount":-1,
				"signature":"` + sigHex + `",
				"index":5
			}`,
			wantErr: "invalid value for amount",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dr electra.DepositRequest
			err := json.Unmarshal([]byte(tt.input), &dr)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)

				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantAmount, dr.Amount)
			assert.Equal(t, tt.wantIndex, dr.Index)

			// Re-marshal and ensure the output is spec form (quoted strings)
			// regardless of how the input was shaped.
			out, err := json.Marshal(&dr)
			require.NoError(t, err)
			assert.Contains(t, string(out), `"amount":"32000000000"`)
			assert.Contains(t, string(out), `"index":"5"`)
		})
	}
}
