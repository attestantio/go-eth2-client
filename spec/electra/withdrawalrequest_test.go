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

// TestWithdrawalRequestUnmarshalJSON_CaplinBareAmount confirms that the
// lenient phase0.Gwei.UnmarshalJSON propagates through to WithdrawalRequest,
// which delegates field unmarshaling via codecs.RawJSON.
func TestWithdrawalRequestUnmarshalJSON_CaplinBareAmount(t *testing.T) {
	const (
		sourceAddr = "0x0000000000000000000000000000000000000001"
		valPubkey  = "0xa99a76ed7796f7be22d5b7e85deeb7c5677e88e511e0b337618f8c4eb61349b4bf2d153f649f7b53359fe8b94a38e44c"
	)
	input := `{
		"source_address":"` + sourceAddr + `",
		"validator_pubkey":"` + valPubkey + `",
		"amount":32000000000
	}`

	var wr electra.WithdrawalRequest
	require.NoError(t, json.Unmarshal([]byte(input), &wr))
	assert.Equal(t, phase0.Gwei(32000000000), wr.Amount)
}
