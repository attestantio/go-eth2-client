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
	"github.com/stretchr/testify/require"
)

func TestBuilderStatus(t *testing.T) {
	for _, status := range []api.BuilderStatus{api.BuilderStatusUnknown, api.BuilderStatusPending, api.BuilderStatusActive, api.BuilderStatusExited} {
		data, err := json.Marshal(status)
		require.NoError(t, err)
		require.Equal(t, `"`+status.String()+`"`, string(data))

		var decoded api.BuilderStatus
		require.NoError(t, json.Unmarshal(data, &decoded))
		require.Equal(t, status, decoded)
	}

	require.Equal(t, "pending", api.BuilderStatusPending.String())
	require.Equal(t, "active", api.BuilderStatusActive.String())
	require.Equal(t, "exited", api.BuilderStatusExited.String())
	require.Equal(t, "unknown", api.BuilderStatus(99).String())

	var decoded api.BuilderStatus
	require.NoError(t, json.Unmarshal([]byte(`"ACTIVE"`), &decoded))
	require.Equal(t, api.BuilderStatusActive, decoded)
	require.EqualError(t, json.Unmarshal([]byte(`"active_ongoing"`), &decoded), `unrecognised builder status "active_ongoing"`)
}
