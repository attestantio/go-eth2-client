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

package version_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/version"
	"github.com/stretchr/testify/require"
)

// TestDataVersionMarshalJSON pins the wire format of every version in the enum
// and requires an out-of-range version to marshal as "unknown" rather than
// panicking, since a panic here aborts the marshal of the whole enclosing
// response.
//
// The address is taken deliberately: MarshalJSON has a pointer receiver, so
// json.Marshal(test.version) would skip the marshaler and emit the number.
func TestDataVersionMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		version  version.DataVersion
		expected string
	}{
		{name: "Unknown", version: version.DataVersionUnknown, expected: `"unknown"`},
		{name: "Phase0", version: version.DataVersionPhase0, expected: `"phase0"`},
		{name: "Altair", version: version.DataVersionAltair, expected: `"altair"`},
		{name: "Bellatrix", version: version.DataVersionBellatrix, expected: `"bellatrix"`},
		{name: "Capella", version: version.DataVersionCapella, expected: `"capella"`},
		{name: "Deneb", version: version.DataVersionDeneb, expected: `"deneb"`},
		{name: "Electra", version: version.DataVersionElectra, expected: `"electra"`},
		{name: "Fulu", version: version.DataVersionFulu, expected: `"fulu"`},
		{name: "Gloas", version: version.DataVersionGloas, expected: `"gloas"`},
		// A fork added to the enum but not to the string table lands here.
		{name: "PastLastVersion", version: version.DataVersionGloas + 1, expected: `"unknown"`},
		// A bounds check written against int, not uint64, wraps negative here.
		{name: "SignedRangeBoundary", version: 1 << 63, expected: `"unknown"`},
		{name: "MaxUint64", version: math.MaxUint64, expected: `"unknown"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				res, err := json.Marshal(&test.version)
				require.NoError(t, err)
				require.Equal(t, test.expected, string(res))
			})
		})
	}
}
