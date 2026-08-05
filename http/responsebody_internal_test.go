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
