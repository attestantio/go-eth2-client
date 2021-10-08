// Copyright © 2020, 2021 Attestant Limited.
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
	"os"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/r3labs/sse/v2"
	"github.com/stretchr/testify/require"
)

// timeout for tests.
var timeout = 60 * time.Second

func TestEventHandler(t *testing.T) {
	handled := false
	handler := func(*api.Event) {
		handled = true
	}

	tests := []struct {
		name    string
		message *sse.Event
		handler client.EventHandlerFunc
		handled bool
	}{
		{
			name:    "MessageNil",
			handler: handler,
			handled: false,
		},
		{
			name:    "MessageEmpty",
			message: &sse.Event{},
			handler: handler,
			handled: false,
		},
		{
			name: "EventUnknown",
			message: &sse.Event{
				Event: []byte("unknown"),
			},
			handler: handler,
			handled: false,
		},
		{
			name: "HandlerNil",
			message: &sse.Event{
				Event: []byte("head"),
			},
			handled: false,
		},
		{
			name: "HeadGood",
			message: &sse.Event{
				Event: []byte("head"),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "BlockGood",
			message: &sse.Event{
				Event: []byte("block"),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "AttestationGood",
			message: &sse.Event{
				Event: []byte("attestation"),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "VoluntaryExitGood",
			message: &sse.Event{
				Event: []byte("voluntary_exit"),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "FinalizedCheckpointGood",
			message: &sse.Event{
				Event: []byte("finalized_checkpoint"),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "ChainReorgGood",
			message: &sse.Event{
				Event: []byte("chain_reorg"),
			},
			handler: handler,
			handled: true,
		},
		{
			name: "ContributionAndProofGood",
			message: &sse.Event{
				Event: []byte("contribution_and_proof"),
			},
			handler: handler,
			handled: true,
		},
	}

	s, err := New(context.Background(),
		WithTimeout(timeout),
		WithAddress(os.Getenv("HTTP_ADDRESS")),
	)
	require.NoError(t, err)

	h, isHTTPService := s.(*Service)
	require.True(t, isHTTPService)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handled = false
			h.handleEvent(test.message, test.handler)
			require.Equal(t, test.handled, handled)
		})
	}
}
