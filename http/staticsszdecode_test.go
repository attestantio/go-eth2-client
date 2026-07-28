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

package http_test

import (
	"context"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/stretchr/testify/require"
)

// TestBeaconStateSSZSpecModes fetches the same SSZ beacon state through both of the
// suite's spec modes, and pins what separates them.
//
// The custom-spec service decodes against the spec the node reports, so it must succeed
// whatever preset that node runs. The default-spec service takes the static generated
// codecs instead, which are built at the mainnet preset: it must succeed against a
// mainnet-preset node and must fail against any other. That failure is the reason this
// suite opts in to custom spec support at all (see ADR-0003), so it is asserted rather
// than skipped over.
//
// The beacon state is the only endpoint used here because it is the one SSZ endpoint the
// interim Glamsterdam devnet serves for both modes: the block endpoints have no gloas arm
// yet, and the devnet's head block carries no blobs. The static branch of every SSZ
// decoder, including those, is covered without a node in
// staticsszdecode_internal_test.go.
func TestBeaconStateSSZSpecModes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	customSpecService := testService(ctx, t).(client.Service)

	// Only testService takes a coordinator token, which is per test, not per service.
	defaultSpecService, err := newTestService(ctx, false)
	require.NoError(t, err)

	specs, err := customSpecService.(client.SpecProvider).Spec(ctx, &api.SpecOpts{})
	require.NoError(t, err)

	presetBase, isString := specs.Data["PRESET_BASE"].(string)
	require.True(t, isString, "spec has no PRESET_BASE")

	opts := &api.BeaconStateOpts{State: "head"}

	t.Run("CustomSpec", func(t *testing.T) {
		response, err := customSpecService.(client.BeaconStateProvider).BeaconState(ctx, opts)
		require.NoError(t, err)
		require.NotNil(t, response.Data)
	})

	t.Run("DefaultSpec", func(t *testing.T) {
		response, err := defaultSpecService.(client.BeaconStateProvider).BeaconState(ctx, opts)
		if presetBase == "mainnet" {
			require.NoError(t, err)
			require.NotNil(t, response.Data)

			return
		}

		// Any other preset sizes its containers differently from the generated codecs,
		// so the static decode must reject the body rather than mis-read it.
		require.ErrorContains(t, err, "failed to decode")
	})
}
