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
	"context"
	"testing"

	v1 "github.com/attestantio/go-eth2-client/api/v1"
	require "github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(ctx context.Context) context.Context
		expected string
	}{
		{
			name: "Empty",
			setup: func(ctx context.Context) context.Context {
				return ctx
			},
			expected: "",
		},
		{
			name: "Set",
			setup: func(ctx context.Context) context.Context {
				return v1.WithProvider(ctx, "http://localhost:5052")
			},
			expected: "http://localhost:5052",
		},
		{
			name: "Override",
			setup: func(ctx context.Context) context.Context {
				ctx = v1.WithProvider(ctx, "http://localhost:5052")
				return v1.WithProvider(ctx, "http://localhost:5053")
			},
			expected: "http://localhost:5053",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = test.setup(ctx)
			result := v1.Provider(ctx)
			require.Equal(t, test.expected, result)
		})
	}
}
