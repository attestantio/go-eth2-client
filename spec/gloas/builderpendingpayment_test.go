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

package gloas_test

import (
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	require "github.com/stretchr/testify/require"
)

func TestBuilderPendingPaymentJSON(t *testing.T) {
	payment := &gloas.BuilderPendingPayment{
		Weight:        1,
		ProposerIndex: 2,
		Withdrawal: &gloas.BuilderPendingWithdrawal{
			FeeRecipient: bellatrix.ExecutionAddress{0x03},
			Amount:       4,
			BuilderIndex: 5,
		},
	}

	data, err := json.Marshal(payment)
	require.NoError(t, err)
	require.Contains(t, string(data), `"weight":"1"`)
	require.Contains(t, string(data), `"withdrawal":{`)
	require.Contains(t, string(data), `"proposer_index":"2"`)
	require.Contains(t, string(data), `"fee_recipient":"0x0300000000000000000000000000000000000000"`)
	require.Contains(t, string(data), `"amount":"4"`)
	require.Contains(t, string(data), `"builder_index":"5"`)

	var decoded gloas.BuilderPendingPayment
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.Equal(t, payment, &decoded)
}
