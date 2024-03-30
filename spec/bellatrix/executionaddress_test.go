// Copyright © 2024 Attestant Limited.
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
	"github.com/stretchr/testify/require"
)

func TestZeroExecutionAddress(t *testing.T) {
	zeroAddress := &bellatrix.ExecutionAddress{}
	require.True(t, zeroAddress.IsZero())

	nonZeroAddress := &bellatrix.ExecutionAddress{0x01}
	require.False(t, nonZeroAddress.IsZero())
}
