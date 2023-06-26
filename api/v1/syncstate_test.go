// Copyright © 2020 Attestant Limited.
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
	require "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestSyncStateJSON(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name: "Empty",
			err:  "unexpected end of JSON input",
		},
		{
			name:  "JSONBad",
			input: []byte("[]"),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.syncStateJSON",
		},
		{
			name:  "HeadSlotMissing",
			input: []byte(`{"sync_distance":"2","is_syncing":true}`),
			err:   "head slot missing",
		},
		{
			name:  "HeadSlotWrongType",
			input: []byte(`{"head_slot":true,"sync_distance":"2","is_syncing":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncStateJSON.head_slot of type string",
		},
		{
			name:  "HeadSlotInvalid",
			input: []byte(`{"head_slot":"-1","sync_distance":"2","is_syncing":true}`),
			err:   "invalid value for head slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "SyncDistanceMissing",
			input: []byte(`{"head_slot":"1","is_syncing":true}`),
			err:   "sync distance missing",
		},
		{
			name:  "SyncDistanceWrongType",
			input: []byte(`{"head_slot":"1","sync_distance":true,"is_syncing":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncStateJSON.sync_distance of type string",
		},
		{
			name:  "SyncDistanceInvalid",
			input: []byte(`{"head_slot":"1","sync_distance":"-1","is_syncing":true}`),
			err:   "invalid value for sync distance: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Good",
			input: []byte(`{"head_slot":"1","sync_distance":"2","is_optimistic":false,"is_syncing":true}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res api.SyncState
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
				assert.Equal(t, string(rt), res.String())
			}
		})
	}
}
