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

package v1_test

import (
	"encoding/json"
	"testing"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

// builderExample is a beacon-APIs Gloas.Builder as served by Prysm, Lighthouse
// and Teku.
const builderExample = `{"pubkey":"0x863917c482adda67fba8305a8ef2cce39cfcd0d13311c86c0cb5b648102b6a481d0b25931d7cfbb5fa3e21c114ec0951","version":"0","execution_address":"0x67565aa29bcf60c3fdcd0b79a36cd85cc13c7a91","balance":"39550000000","deposit_epoch":"1434","withdrawable_epoch":"18446744073709551615"}`

const builderResponseExample = `{"index":"1","status":"active","builder":` + builderExample + `}`

func TestBuilderJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   string
	}{
		{name: "Good", input: builderResponseExample},
		{name: "JSONBad", input: `[]`, err: "invalid JSON: json: cannot unmarshal array into Go value of type v1.builderJSON"},
		{name: "IndexMissing", input: `{"status":"active","builder":` + builderExample + `}`, err: "index missing"},
		{name: "IndexInvalid", input: `{"index":"-1","status":"active","builder":` + builderExample + `}`, err: "invalid value for index: strconv.ParseUint: parsing \"-1\": invalid syntax"},
		{name: "StatusMissing", input: `{"index":"1","builder":` + builderExample + `}`, err: "status missing"},
		{name: "StatusInvalid", input: `{"index":"1","status":"active_ongoing","builder":` + builderExample + `}`, err: "invalid JSON: unrecognised builder status \"active_ongoing\""},
		{name: "BuilderMissing", input: `{"index":"1","status":"active"}`, err: "builder missing"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.Builder
			err := json.Unmarshal([]byte(test.input), &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, gloas.BuilderIndex(1), res.Index)
			require.Equal(t, api.BuilderStatusActive, res.Status)
			require.Equal(t, uint8(0), res.Builder.Version)

			rt, err := json.Marshal(&res)
			require.NoError(t, err)
			var decoded api.Builder
			require.NoError(t, json.Unmarshal(rt, &decoded))
			require.Equal(t, res, decoded)
			require.Equal(t, string(rt), res.String())
		})
	}
}
