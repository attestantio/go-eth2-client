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

package bellatrix_test

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	require "github.com/stretchr/testify/require"
)

func TestExecutionAddressString(t *testing.T) {
	tests := []struct {
		name   string
		input  bellatrix.ExecutionAddress
		output string
	}{
		{
			name:   "Zero",
			input:  bellatrix.ExecutionAddress{},
			output: "0x0000000000000000000000000000000000000000",
		},
		{
			name:   "Ten",
			input:  bellatrix.ExecutionAddress{0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a, 0x0a},
			output: "0x0A0A0a0a0a0a0a0A0a0a0A0a0A0A0A0a0a0a0a0a",
		},
		{
			name:   "A",
			input:  bellatrix.ExecutionAddress{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
			output: "0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.output, test.input.String())
		})
	}
}
