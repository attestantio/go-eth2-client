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

package http

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadResponseBody(t *testing.T) {
	tests := []struct {
		name  string
		body  string
		limit int
		err   string
	}{
		{name: "WithinLimit", body: "1234", limit: 4},
		{name: "ExceedsLimit", body: "12345", limit: 4, err: "response body exceeds 4 bytes"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := readResponseBody(strings.NewReader(test.body), test.limit)
			if test.err != "" {
				require.EqualError(t, err, test.err)
				require.Nil(t, body)
			} else {
				require.NoError(t, err)
				require.Equal(t, []byte(test.body), body)
			}
		})
	}
}

// TestEPBSProposalResponseLimitsAreSurvivable pins the property that makes the
// block-production limits a defence rather than a comment.  readResponseBody
// reaches its length check only after io.ReadAll has buffered the whole body,
// and ReadAll grows its buffer geometrically, so peak allocation is close to
// twice the limit.  A limit taken from the protocol's theoretical maximum
// (~2GiB of SSZ, ~5GiB of JSON hex) is therefore never reached -- the process
// dies first -- and a body admitted just under one is equally fatal.
//
// The ceiling below is an allocation budget, not a protocol bound: every live
// preset produces a payload-included response orders of magnitude smaller, so
// the JSON limit sitting exactly at half the budget is the deliberate ceiling
// rather than a coincidence.  Raising either limit past it is the regression
// this test exists to catch.
func TestEPBSProposalResponseLimitsAreSurvivable(t *testing.T) {
	const maxProposalResponsePeakAllocation = 1024 * 1024 * 1024

	tests := []struct {
		name  string
		limit int
	}{
		{name: "SSZ", limit: maxEPBSProposalResponseSize},
		{name: "JSON", limit: maxEPBSProposalJSONResponseSize},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.LessOrEqual(t, 2*test.limit, maxProposalResponsePeakAllocation)
		})
	}
}
