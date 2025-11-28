// Copyright © 2022 Attestant Limited.
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

package http_test

import (
	"context"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"
)

func TestSubmitBLSToExecutionChanges(t *testing.T) {
	tests := []struct {
		name string
		ops  []*capella.SignedBLSToExecutionChange
	}{
		{
			name: "InvalidSignature",
			ops: []*capella.SignedBLSToExecutionChange{
				{
					Message: &capella.BLSToExecutionChange{
						ValidatorIndex:     12345,
						FromBLSPubkey:      phase0.BLSPubKey{},
						ToExecutionAddress: bellatrix.ExecutionAddress{},
					},
					Signature: phase0.BLSSignature{},
				},
			},
		},
	}

	service := testService(context.Background(), t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.(client.BLSToExecutionChangesSubmitter).SubmitBLSToExecutionChanges(context.Background(), test.ops)
			require.Contains(t, err.Error(), "400")
		})
	}
}
