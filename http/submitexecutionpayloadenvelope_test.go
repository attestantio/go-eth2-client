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
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/stretchr/testify/require"
)

func TestSubmitExecutionPayloadEnvelope(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	tests := []struct {
		name string
		opts *api.SubmitExecutionPayloadEnvelopeOpts
		err  string
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name: "NilEnvelope",
			opts: &api.SubmitExecutionPayloadEnvelopeOpts{},
			err:  "no envelope supplied",
		},
		{
			name: "UnsupportedVersion",
			opts: &api.SubmitExecutionPayloadEnvelopeOpts{
				SignedExecutionPayloadEnvelope: &spec.VersionedSignedExecutionPayloadEnvelope{
					Version: spec.DataVersionPhase0,
				},
			},
			err: "unsupported envelope version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.(client.ExecutionPayloadEnvelopeSubmitter).SubmitExecutionPayloadEnvelope(ctx, test.opts)
			require.ErrorContains(t, err, test.err)
		})
	}

	// A structurally-valid Gloas envelope (with a garbage signature) passes
	// client-side validation and exercises the SSZ marshal + POST path (with
	// the Eth-Blob-Data-Included header) that the other cases never reach. The
	// beacon node rejects it, so we tolerate any returned error rather than
	// asserting a node-specific message.
	t.Run("ReachesServer", func(t *testing.T) {
		err := service.(client.ExecutionPayloadEnvelopeSubmitter).SubmitExecutionPayloadEnvelope(ctx, &api.SubmitExecutionPayloadEnvelopeOpts{
			SignedExecutionPayloadEnvelope: &spec.VersionedSignedExecutionPayloadEnvelope{
				Version: spec.DataVersionGloas,
				Gloas: &gloas.SignedExecutionPayloadEnvelope{
					Message: &gloas.ExecutionPayloadEnvelope{
						Payload:           &gloas.ExecutionPayload{},
						ExecutionRequests: &gloas.ExecutionRequests{},
					},
				},
			},
		})
		require.Error(t, err)
	})
}
