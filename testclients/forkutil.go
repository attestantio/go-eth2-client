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
	"strings"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// unknownFork is what the helpers below report when the node cannot be asked.
const unknownFork = "unknown"

// nodeSpec fetches the node's configuration, or nil if it cannot be obtained.
func nodeSpec(ctx context.Context, service any) map[string]any {
	provider, isProvider := service.(client.SpecProvider)
	if !isProvider {
		return nil
	}
	response, err := provider.Spec(ctx, &api.SpecOpts{})
	if err != nil || response == nil || response.Data == nil {
		return nil
	}

	return response.Data
}

// headForkVersion fetches the fork version in force at the chain head.
func headForkVersion(ctx context.Context, service any) (phase0.Version, bool) {
	provider, isProvider := service.(client.ForkProvider)
	if !isProvider {
		return phase0.Version{}, false
	}
	response, err := provider.Fork(ctx, &api.ForkOpts{State: "head"})
	if err != nil || response == nil || response.Data == nil {
		return phase0.Version{}, false
	}

	return response.Data.CurrentVersion, true
}

// KnowsGloas reports whether the node's configuration includes the Gloas fork.
func KnowsGloas(ctx context.Context, service any) bool {
	config := nodeSpec(ctx, service)
	if config == nil {
		return false
	}
	_, isPresent := config["GLOAS_FORK_EPOCH"]

	return isPresent
}

// OnGloas reports whether the chain head is in the Gloas fork.
//
// The head's fork version is compared against GLOAS_FORK_VERSION rather than the
// head block being fetched and its version read: a block fetched over SSZ by a
// service without custom spec support cannot be decoded at all on a non-mainnet
// preset, which would report every Gloas chain as not being on Gloas.
func OnGloas(ctx context.Context, service any) bool {
	head, haveHead := headForkVersion(ctx, service)
	if !haveHead {
		return false
	}

	gloas, isVersion := nodeSpec(ctx, service)["GLOAS_FORK_VERSION"].(phase0.Version)

	return isVersion && head == gloas
}

// HeadVersion returns the name of the fork in force at the chain head, or
// "unknown".  It is diagnostic only: use OnGloas to decide.
func HeadVersion(ctx context.Context, service any) string {
	head, haveHead := headForkVersion(ctx, service)
	if !haveHead {
		return unknownFork
	}

	for key, value := range nodeSpec(ctx, service) {
		version, isVersion := value.(phase0.Version)
		if !isVersion || version != head {
			continue
		}
		if name, isForkVersion := strings.CutSuffix(key, "_FORK_VERSION"); isForkVersion {
			return strings.ToLower(name)
		}
	}

	return unknownFork
}
