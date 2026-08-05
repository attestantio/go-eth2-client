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
	"github.com/stretchr/testify/require"
)

func TestSignedExecutionPayloadEnvelope(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := testService(ctx, t).(client.Service)

	tests := []struct {
		name string
		opts *api.SignedExecutionPayloadEnvelopeOpts
		err  string
	}{
		{
			name: "NilOpts",
			err:  "no options specified",
		},
		{
			name: "NoBlock",
			opts: &api.SignedExecutionPayloadEnvelopeOpts{},
			err:  "no block specified",
		},
		{
			name: "Head",
			opts: &api.SignedExecutionPayloadEnvelopeOpts{Block: "head"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := service.(client.ExecutionPayloadProvider).SignedExecutionPayloadEnvelope(ctx, test.opts)
			switch {
			case test.err != "":
				require.ErrorContains(t, err, test.err)
			default:
				if err != nil {
					// The execution payload envelope endpoint is Gloas-only and
					// the envelope for a given block may not yet be revealed, so
					// tolerate a node that cannot serve it.
					t.Skipf("execution payload envelope not available: %v", err)
				}

				require.NotNil(t, response)
				require.Equal(t, spec.DataVersionGloas, response.Data.Version)
				require.NotNil(t, response.Data.Gloas)
			}
		})
	}
}
