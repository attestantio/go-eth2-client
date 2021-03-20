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

package v2_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	standardhttp "github.com/attestantio/go-eth2-client/standardhttp/v2"
	"github.com/stretchr/testify/require"
)

func TestStateRoot(t *testing.T) {
	tests := []struct {
		name              string
		stateID           string
		expectedErrorCode int
	}{
		{
			name:              "Invalid",
			stateID:           "current",
			expectedErrorCode: 400,
		},
		{
			name:    "Zero",
			stateID: "0",
		},
		{
			name:    "Head",
			stateID: "head",
		},
		{
			name:    "Finalized",
			stateID: "finalized",
		},
		{
			name:    "Justified",
			stateID: "justified",
		},
	}

	service, err := standardhttp.New(context.Background(),
		standardhttp.WithTimeout(timeout),
		standardhttp.WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stateRoot, err := service.StateRoot(context.Background(), test.stateID)
			if test.expectedErrorCode != 0 {
				require.Contains(t, err.Error(), fmt.Sprintf("%d", test.expectedErrorCode))
			} else {
				require.NoError(t, err)
				require.NotNil(t, stateRoot)
			}
		})
	}
}
