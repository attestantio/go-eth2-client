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

package altair_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
)

func TestSyncAggregatorSelectionDataJSON(t *testing.T) {
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
			err:   "invalid JSON: json: cannot unmarshal array into Go value of type altair.syncAggregatorSelectionDataJSON",
		},
		{
			name:  "SlotMissing",
			input: []byte(`{"subcommittee_index":"3"}`),
			err:   "slot missing",
		},
		{
			name:  "SlotWrongType",
			input: []byte(`{"slot":true,"subcommittee_index":"3"}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncAggregatorSelectionDataJSON.slot of type string",
		},
		{
			name:  "SlotInvalid",
			input: []byte(`{"slot":"-1","subcommittee_index":"3"}`),
			err:   "invalid value for slot: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "SubcommitteeIndexMissing",
			input: []byte(`{"slot":"1"}`),
			err:   "subcommittee index missing",
		},
		{
			name:  "SubcommitteeIndexWrongType",
			input: []byte(`{"slot":"1","subcommittee_index":true}`),
			err:   "invalid JSON: json: cannot unmarshal bool into Go struct field syncAggregatorSelectionDataJSON.subcommittee_index of type string",
		},
		{
			name:  "SubcommitteeIndexInvalid",
			input: []byte(`{"slot":"1","subcommittee_index":"-1"}`),
			err:   "invalid value for subcommittee index: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:  "Valid",
			input: []byte(`{"slot":"1","subcommittee_index":"3"}`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.SyncAggregatorSelectionData
			err := json.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := json.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}

func TestSyncAggregatorSelectionDataYAML(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   string
	}{
		{
			name:  "Good",
			input: []byte(`{slot: 1, subcommittee_index: 3}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res altair.SyncAggregatorSelectionData
			err := yaml.Unmarshal(test.input, &res)
			if test.err != "" {
				require.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				rt, err := yaml.Marshal(&res)
				require.NoError(t, err)
				assert.Equal(t, string(rt), res.String())
				rt = bytes.TrimSuffix(rt, []byte("\n"))
				assert.Equal(t, string(test.input), string(rt))
			}
		})
	}
}
