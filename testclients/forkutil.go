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

package testclients

import (
	"context"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
)

// KnowsGloas reports whether the node's configuration includes the Gloas fork.
func KnowsGloas(ctx context.Context, service any) bool {
	provider, ok := service.(client.SpecProvider)
	if !ok {
		return false
	}
	response, err := provider.Spec(ctx, &api.SpecOpts{})
	if err != nil || response == nil || response.Data == nil {
		return false
	}
	_, ok = response.Data["GLOAS_FORK_EPOCH"]
	return ok
}

// OnGloas reports whether the node's head block is a Gloas block.
func OnGloas(ctx context.Context, service any) bool {
	return HeadVersion(ctx, service) == spec.DataVersionGloas.String()
}

// HeadVersion returns the fork name of the node's head block, or "unknown".
func HeadVersion(ctx context.Context, service any) string {
	provider, ok := service.(client.SignedBeaconBlockProvider)
	if !ok {
		return "unknown"
	}
	response, err := provider.SignedBeaconBlock(ctx, &api.SignedBeaconBlockOpts{Block: "head"})
	if err != nil || response == nil || response.Data == nil {
		return "unknown"
	}
	return response.Data.Version.String()
}
