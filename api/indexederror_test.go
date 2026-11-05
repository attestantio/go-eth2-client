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

package api_test

import (
	"encoding/json"
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/stretchr/testify/require"
)

// TestIndexedErrorUnmarshalJSONKeepsUsableFailures decodes directly, as a caller
// holding an *api.Error's Data would: a failure that does not decode is skipped
// rather than discarding the others.
func TestIndexedErrorUnmarshalJSONKeepsUsableFailures(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		message  string
		failures []api.IndexedFailure
	}{
		{
			name:     "StringIndex",
			input:    `{"code":400,"message":"m","failures":[{"index":"1","message":"a"}]}`,
			message:  "m",
			failures: []api.IndexedFailure{{Index: 1, Message: "a"}},
		},
		{
			name:     "MissingIndex",
			input:    `{"code":400,"message":"m","failures":[{"index":0,"message":"a"},{"message":"b"}]}`,
			message:  "m",
			failures: []api.IndexedFailure{{Index: 0, Message: "a"}},
		},
		{
			name:     "FloatIndex",
			input:    `{"code":400,"message":"m","failures":[{"index":1.5,"message":"a"},{"index":1,"message":"b"}]}`,
			message:  "m",
			failures: []api.IndexedFailure{{Index: 1, Message: "b"}},
		},
		{
			name:     "NonObjectFailure",
			input:    `{"code":400,"message":"m","failures":["x",{"index":2,"message":"b"}]}`,
			message:  "m",
			failures: []api.IndexedFailure{{Index: 2, Message: "b"}},
		},
		{
			name:     "NumericMessage",
			input:    `{"code":400,"message":1,"failures":[{"index":0,"message":"a"}]}`,
			failures: []api.IndexedFailure{{Index: 0, Message: "a"}},
		},
		{
			name:     "NumericFailureMessage",
			input:    `{"code":400,"message":"m","failures":[{"index":0,"message":5}]}`,
			message:  "m",
			failures: []api.IndexedFailure{{Index: 0}},
		},
		{
			name:     "NoFailures",
			input:    `{"code":400,"message":"No data submitted"}`,
			message:  "No data submitted",
			failures: []api.IndexedFailure{},
		},
		{
			name:     "StringCode",
			input:    `{"code":"400","message":"m","failures":[{"index":0,"message":"a"}]}`,
			message:  "m",
			failures: []api.IndexedFailure{{Index: 0, Message: "a"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var indexedErr api.IndexedError
			require.NoError(t, json.Unmarshal([]byte(test.input), &indexedErr))
			require.Equal(t, test.message, indexedErr.Message)
			require.Equal(t, test.failures, indexedErr.Failures)
		})
	}
}

// TestIndexedErrorUnmarshalJSONRejectsMalformedInput covers the only inputs
// that fail: text that is not JSON, and failures that is not an array.
func TestIndexedErrorUnmarshalJSONRejectsMalformedInput(t *testing.T) {
	for _, input := range []string{`bad request`, `[]`, `{"message":"m","failures":{}}`} {
		var indexedErr api.IndexedError
		require.Error(t, json.Unmarshal([]byte(input), &indexedErr), input)
	}
}
