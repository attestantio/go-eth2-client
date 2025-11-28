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

package http_test

import (
	"context"
	"sync"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/stretchr/testify/require"
)

func TestEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tests := []struct {
		name   string
		topics []string
	}{
		{
			name:   "Good",
			topics: []string{"head", "chain_reorg"},
		},
	}

	service := testService(ctx, t).(client.Service)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			eventsMu := sync.Mutex{}
			events := 0
			err := service.(client.EventsProvider).Events(ctx, &api.EventsOpts{
				Topics: test.topics,
				Handler: func(*apiv1.Event) {
					eventsMu.Lock()
					events++
					eventsMu.Unlock()
				},
			})
			require.NoError(t, err)
			time.Sleep(30 * time.Second)
			eventsMu.Lock()
			defer eventsMu.Unlock()
			require.NotEqual(t, 0, events)
			cancel()
		})
	}
}
