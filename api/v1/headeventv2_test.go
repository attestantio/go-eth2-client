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
	"github.com/attestantio/go-eth2-client/spec/phase0"
	require "github.com/stretchr/testify/require"
)

func TestHeadEventV2JSONErrors(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name:  "InvalidJSON",
			input: []byte(`[]`),
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type v1.headEventV2JSON",
		},
		{
			name:  "VersionMissing",
			input: []byte(`{"data":{}}`),
			err:   "version missing",
		},
		{
			name:  "DataMissing",
			input: []byte(`{"version":"gloas"}`),
			err:   "data missing",
		},
		{
			name:  "PayloadStatusMissing",
			input: []byte(`{"version":"gloas","data":{"slot":"10","block":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","state":"0x600e852a08c1200654ddf11025f1ceacb3c2e74bdd5c630cde0838b2591b69f9","current_epoch_dependent_root":"0x5e0043f107cb57913498fbf2f99ff55e730bf1e151f02f221e977c91a90a0e91","next_epoch_dependent_root":"0x7f1154a218dc77addb1fe144266e83c04f67c8dfd03aacf955df05308cd8b27c"}}`),
			err:   "payload status missing",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"version":"gloas","data":{"payload_status":"empty"}}`),
			err:   "slot missing",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"version":"gloas","data":{"slot":"invalid","payload_status":"empty"}}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"invalid\": invalid syntax",
		},
		{
			name:  "BlockMissing",
			input: []byte(`{"version":"gloas","data":{"slot":"10","payload_status":"empty"}}`),
			err:   "block missing",
		},
		{
			name:  "BlockInvalid",
			input: []byte(`{"version":"gloas","data":{"slot":"10","block":"0x01","payload_status":"empty"}}`),
			err:   "incorrect length 1 for block",
		},
		{
			name:  "StateMissing",
			input: []byte(`{"version":"gloas","data":{"slot":"10","block":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","payload_status":"empty"}}`),
			err:   "state missing",
		},
		{
			name:  "CurrentEpochDependentRootMissing",
			input: []byte(`{"version":"gloas","data":{"slot":"10","block":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","state":"0x600e852a08c1200654ddf11025f1ceacb3c2e74bdd5c630cde0838b2591b69f9","payload_status":"empty"}}`),
			err:   "current epoch dependent root missing",
		},
		{
			name:  "NextEpochDependentRootMissing",
			input: []byte(`{"version":"gloas","data":{"slot":"10","block":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","state":"0x600e852a08c1200654ddf11025f1ceacb3c2e74bdd5c630cde0838b2591b69f9","payload_status":"empty","current_epoch_dependent_root":"0x5e0043f107cb57913498fbf2f99ff55e730bf1e151f02f221e977c91a90a0e91"}}`),
			err:   "next epoch dependent root missing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var event api.HeadEventV2
			require.EqualError(t, json.Unmarshal(test.input, &event), test.err)
		})
	}
}

func TestHeadEventV2JSON(t *testing.T) {
	input := []byte(`{"version":"gloas","data":{"slot":"10","block":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","state":"0x600e852a08c1200654ddf11025f1ceacb3c2e74bdd5c630cde0838b2591b69f9","payload_status":"empty","epoch_transition":false,"current_epoch_dependent_root":"0x5e0043f107cb57913498fbf2f99ff55e730bf1e151f02f221e977c91a90a0e91","next_epoch_dependent_root":"0x7f1154a218dc77addb1fe144266e83c04f67c8dfd03aacf955df05308cd8b27c","execution_optimistic":true}}`)

	var actual api.HeadEventV2
	require.NoError(t, json.Unmarshal(input, &actual))
	require.Equal(t, spec.DataVersionGloas, actual.Version)
	require.Equal(t, phase0.Slot(10), actual.Slot)
	require.Equal(t, phase0.Root{0x9a, 0x2f, 0xef, 0xd2, 0xfd, 0xb5, 0x7f, 0x74, 0x99, 0x3c, 0x77, 0x80, 0xea, 0x5b, 0x90, 0x30, 0xd2, 0x89, 0x7b, 0x61, 0x5b, 0x89, 0xf8, 0x08, 0x01, 0x1c, 0xa5, 0xae, 0xbe, 0xd5, 0x4e, 0xaf}, actual.Block)
	require.Equal(t, phase0.Root{0x60, 0x0e, 0x85, 0x2a, 0x08, 0xc1, 0x20, 0x06, 0x54, 0xdd, 0xf1, 0x10, 0x25, 0xf1, 0xce, 0xac, 0xb3, 0xc2, 0xe7, 0x4b, 0xdd, 0x5c, 0x63, 0x0c, 0xde, 0x08, 0x38, 0xb2, 0x59, 0x1b, 0x69, 0xf9}, actual.State)
	require.Equal(t, "empty", actual.PayloadStatus)
	require.False(t, actual.EpochTransition)
	require.Equal(t, phase0.Root{0x5e, 0x00, 0x43, 0xf1, 0x07, 0xcb, 0x57, 0x91, 0x34, 0x98, 0xfb, 0xf2, 0xf9, 0x9f, 0xf5, 0x5e, 0x73, 0x0b, 0xf1, 0xe1, 0x51, 0xf0, 0x2f, 0x22, 0x1e, 0x97, 0x7c, 0x91, 0xa9, 0x0a, 0x0e, 0x91}, actual.CurrentEpochDependentRoot)
	require.Equal(t, phase0.Root{0x7f, 0x11, 0x54, 0xa2, 0x18, 0xdc, 0x77, 0xad, 0xdb, 0x1f, 0xe1, 0x44, 0x26, 0x6e, 0x83, 0xc0, 0x4f, 0x67, 0xc8, 0xdf, 0xd0, 0x3a, 0xac, 0xf9, 0x55, 0xdf, 0x05, 0x30, 0x8c, 0xd8, 0xb2, 0x7c}, actual.NextEpochDependentRoot)
	require.True(t, actual.ExecutionOptimistic)

	output, err := json.Marshal(&actual)
	require.NoError(t, err)
	require.JSONEq(t, string(input), string(output))
	require.JSONEq(t, string(input), actual.String())
}
