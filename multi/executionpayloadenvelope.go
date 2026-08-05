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

package multi

import (
	"context"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec"
)

// ExecutionPayloadEnvelope obtains the cached execution payload envelope for the
// given slot and beacon block root.
//
// Worth knowing when using this through a multi-client: only the node that
// produced the block holds its envelope, so every other client in the set
// answers ErrNoExecutionPayloadEnvelope.  A caller that excluded the payload
// from its proposal should address the producing node directly rather than rely
// on failover here.
func (s *Service) ExecutionPayloadEnvelope(ctx context.Context,
	opts *api.ExecutionPayloadEnvelopeOpts,
) (
	*api.Response[*spec.VersionedExecutionPayloadEnvelope],
	error,
) {
	res, err := s.doCall(ctx, func(ctx context.Context, client consensusclient.Service) (any, error) {
		envelope, err := client.(consensusclient.ExecutionPayloadEnvelopeProvider).ExecutionPayloadEnvelope(ctx, opts)
		if err != nil {
			return nil, err
		}

		return envelope, nil
	}, nil)
	if err != nil {
		return nil, err
	}

	response, isResponse := res.(*api.Response[*spec.VersionedExecutionPayloadEnvelope])
	if !isResponse {
		return nil, ErrIncorrectType
	}

	return response, nil
}
