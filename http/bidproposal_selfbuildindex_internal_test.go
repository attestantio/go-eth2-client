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
	"context"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

// TestBuilderIndexSelfBuild covers how the self-build sentinel is resolved.  It
// is a spec value rather than a fixed constant, but it arrived with Gloas, so a
// node on an earlier release serves a spec without it; that has to fall back
// rather than fail the proposal it is being read for.
func TestBuilderIndexSelfBuild(t *testing.T) {
	tests := []struct {
		name              string
		customSpecSupport bool
		spec              map[string]any
		expected          gloas.BuilderIndex
		err               error
	}{
		{
			name:     "StaticWithoutCustomSpecSupport",
			spec:     map[string]any{"BUILDER_INDEX_SELF_BUILD": uint64(5)},
			expected: staticBuilderIndexSelfBuild,
		},
		{
			name:              "SpecValue",
			customSpecSupport: true,
			spec:              map[string]any{"BUILDER_INDEX_SELF_BUILD": uint64(5)},
			expected:          gloas.BuilderIndex(5),
		},
		{
			name:              "AbsentFallsBackToMainnet",
			customSpecSupport: true,
			spec:              map[string]any{"SYNC_COMMITTEE_SIZE": uint64(32)},
			expected:          staticBuilderIndexSelfBuild,
		},
		{
			name:              "WrongType",
			customSpecSupport: true,
			spec:              map[string]any{"BUILDER_INDEX_SELF_BUILD": "self"},
			err:               ErrIncorrectType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &Service{
				customSpecSupport: test.customSpecSupport,
				connectionActive:  true,
				spec:              test.spec,
			}

			index, err := s.builderIndexSelfBuild(context.Background())
			if test.err != nil {
				require.ErrorIs(t, err, test.err)

				return
			}
			require.NoError(t, err)
			require.Equal(t, test.expected, index)
		})
	}
}
