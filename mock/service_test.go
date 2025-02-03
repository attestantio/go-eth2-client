// Copyright © 2025 Attestant Limited.
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

package mock_test

import (
	"context"
	"testing"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/stretchr/testify/require"
)

func TestMockFunc(t *testing.T) {
	ctx := context.Background()

	m, err := mock.New(ctx)
	require.NoError(t, err)

	m.SpecFunc = testSpec

	specResponse, err := m.Spec(ctx, &api.SpecOpts{})
	require.NoError(t, err)
	require.NotNil(t, specResponse)

	require.Contains(t, specResponse.Data, "MOCK_VALUE")
}

func testSpec(ctx context.Context, opts *api.SpecOpts) (*api.Response[map[string]any], error) {
	data := map[string]any{
		"MOCK_VALUE": "foo",
	}

	return &api.Response[map[string]any]{
		Data:     data,
		Metadata: make(map[string]any),
	}, nil
}
