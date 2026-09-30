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
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/require"
)

func TestLightClientEventsJSON(t *testing.T) {
	tests := []struct {
		name     string
		newEvent func() any
	}{
		{name: "Finality", newEvent: func() any { return &api.LightClientFinalityUpdateEvent{} }},
		{name: "Optimistic", newEvent: func() any { return &api.LightClientOptimisticUpdateEvent{} }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, input := range []string{`null`, `{"version":"altair"}`, `{"version":"altair","data":null}`, `{"data":{}}`} {
				t.Run(input, func(t *testing.T) {
					require.Error(t, json.Unmarshal([]byte(input), test.newEvent()))
				})
			}
		})
	}

	var event api.LightClientFinalityUpdateEvent
	require.NoError(t, json.Unmarshal([]byte(`{"version":"altair","data":{"signature_slot":"1"}}`), &event))
	require.Equal(t, spec.DataVersionAltair, event.Version)
	require.JSONEq(t, `{"signature_slot":"1"}`, string(event.Data))
}
