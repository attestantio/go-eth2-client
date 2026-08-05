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

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
)

// KnowsGloas reports whether the node's configuration includes the Gloas fork,
// which is to say whether GLOAS_FORK_EPOCH is present in /eth/v1/config/spec.
//
// Presence alone, deliberately: a client that knows Gloas publishes the epoch as
// soon as it has one, typically far in the future, and its ePBS endpoints exist to
// refuse pre-fork requests from that same moment. A far-future epoch therefore
// counts as knowing, and no arithmetic is done on the value.
func KnowsGloas(ctx context.Context, service any) bool {
	specProvider, isProvider := service.(interface {
		Spec(ctx context.Context, opts *api.SpecOpts) (*api.Response[map[string]any], error)
	})
	if !isProvider {
		return false
	}

	response, err := specProvider.Spec(ctx, &api.SpecOpts{})
	if err != nil || response == nil || response.Data == nil {
		return false
	}

	_, exists := response.Data["GLOAS_FORK_EPOCH"]

	return exists
}

// OnGloas reports whether the node's chain head is a Gloas block.
//
// The head's own version rather than current_epoch >= GLOAS_FORK_EPOCH arithmetic:
// one call instead of several spec reads, it is what the gated tests actually
// depend on since they all operate on head, and it cannot disagree with reality
// when a node's configuration and its chain diverge.
func OnGloas(ctx context.Context, service any) bool {
	return HeadVersion(ctx, service) == spec.DataVersionGloas.String()
}

// HeadVersion returns the name of the fork the node's head block belongs to, or
// "unknown" if the node cannot be asked.
//
// It exists so that a caller declining to run can say which fork it found instead
// of only that the fork was wrong, and so that OnGloas above and that message are
// answered by one piece of code rather than two that can disagree.
//
// Because it reads the head block rather than a header, the service must be able to
// decode one: against a non-mainnet preset that means custom spec support, without
// which the compiled-in codec fails on the offsets. A service that cannot decode
// the node's blocks reports "unknown" rather than a fork name, so the caller says
// it could not tell rather than naming a fork the node is not on.
func HeadVersion(ctx context.Context, service any) string {
	blockProvider, isProvider := service.(interface {
		SignedBeaconBlock(ctx context.Context, opts *api.SignedBeaconBlockOpts) (
			*api.Response[*spec.VersionedSignedBeaconBlock],
			error,
		)
	})
	if !isProvider {
		return "unknown"
	}

	response, err := blockProvider.SignedBeaconBlock(ctx, &api.SignedBeaconBlockOpts{Block: "head"})
	if err != nil || response == nil || response.Data == nil {
		return "unknown"
	}

	return response.Data.Version.String()
}
