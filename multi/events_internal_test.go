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
	"errors"
	"sync/atomic"
	"testing"
	"time"

	consensusclient "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/mock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestEventsDoesNotWriteServiceClientLists confirms that an active client failing to subscribe
// is queued for retry on a list of Events' own.  activateClient and deactivateClient leave the
// service's lists with spare capacity, so appending to a copy of the slice header would write
// into the service's backing array outside the lock, racing with those functions' own appends.
func TestEventsDoesNotWriteServiceClientLists(t *testing.T) {
	// The failing client is retried until the context is done, so end it with the test.  The
	// mocks are given a context of their own, as ending theirs has them race on mock's logger.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	failing, err := mock.New(context.Background(), mock.WithName("failing"))
	require.NoError(t, err)
	failing.EventsFunc = func(context.Context, *api.EventsOpts) error {
		return errors.New("failed to subscribe")
	}

	inactive, err := mock.New(context.Background(), mock.WithName("inactive"))
	require.NoError(t, err)

	inactiveClients := make([]consensusclient.Service, 1, 2)
	inactiveClients[0] = inactive

	s := &Service{
		log:                 zerolog.Nop(),
		activeClients:       []consensusclient.Service{failing},
		inactiveClients:     inactiveClients,
		eventsRetryInterval: time.Millisecond,
	}

	require.NoError(t, s.Events(ctx, &api.EventsOpts{
		Topics:  []string{"head"},
		Handler: func(*apiv1.Event) {},
	}))

	s.clientsMu.RLock()
	defer s.clientsMu.RUnlock()
	require.Nil(t, s.inactiveClients[:cap(s.inactiveClients)][1], "Events wrote into the service's inactive client list")
}

// TestEventsRetriesWithDefaultIntervalWhenUnset confirms that a Service whose retry interval is
// unset, as one not built by New, waits between attempts rather than retrying without pause.
func TestEventsRetriesWithDefaultIntervalWhenUnset(t *testing.T) {
	// As in TestEventsDoesNotWriteServiceClientLists, the mock has a context of its own.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var attempts atomic.Int32
	failing, err := mock.New(context.Background(), mock.WithName("failing"))
	require.NoError(t, err)
	failing.EventsFunc = func(context.Context, *api.EventsOpts) error {
		attempts.Add(1)

		return errors.New("failed to subscribe")
	}

	s := &Service{
		log:             zerolog.Nop(),
		inactiveClients: []consensusclient.Service{failing},
	}

	require.NoError(t, s.Events(ctx, &api.EventsOpts{
		Topics:  []string{"head"},
		Handler: func(*apiv1.Event) {},
	}))

	time.Sleep(50 * time.Millisecond)
	require.LessOrEqual(t, attempts.Load(), int32(1), "retried without waiting for an interval")
}
