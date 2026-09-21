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

package client_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	eth2http "github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog"
	require "github.com/stretchr/testify/require"
)

func TestEventsLightClientUpdates(t *testing.T) {
	server := newEventsTestServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = fmt.Fprint(w, "event: light_client_finality_update\ndata: {\"version\":\"altair\",\"data\":{\"signature_slot\":\"1\"}}\n\n")
		_, _ = fmt.Fprint(w, "event: light_client_optimistic_update\ndata: {\"version\":\"altair\",\"data\":{\"signature_slot\":\"2\"}}\n\n")
		w.(http.Flusher).Flush()
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	service, err := eth2http.New(ctx, eth2http.WithAddress(server.URL), eth2http.WithLogLevel(zerolog.Disabled))
	require.NoError(t, err)

	received := make(chan *apiv1.Event, 2)
	err = service.(client.EventsProvider).Events(ctx, &api.EventsOpts{
		Topics: []string{"light_client_finality_update", "light_client_optimistic_update"},
		Handler: func(event *apiv1.Event) {
			select {
			case received <- event:
			default:
			}
		},
	})
	require.NoError(t, err)

	events := map[string]*apiv1.Event{}
	for len(events) < 2 {
		select {
		case event := <-received:
			events[event.Topic] = event
		case <-ctx.Done():
			t.Fatal("timed out waiting for light client events")
		}
	}
	cancel()

	finality, ok := events["light_client_finality_update"].Data.(*apiv1.LightClientFinalityUpdateEvent)
	require.True(t, ok)
	require.Equal(t, spec.DataVersionAltair, finality.Version)
	require.JSONEq(t, `{"signature_slot":"1"}`, string(finality.Data))

	optimistic, ok := events["light_client_optimistic_update"].Data.(*apiv1.LightClientOptimisticUpdateEvent)
	require.True(t, ok)
	require.Equal(t, spec.DataVersionAltair, optimistic.Version)
	require.JSONEq(t, `{"signature_slot":"2"}`, string(optimistic.Data))
}

func TestEventsTopicValidation(t *testing.T) {
	server := newEventsTestServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service, err := eth2http.New(ctx, eth2http.WithAddress(server.URL), eth2http.WithLogLevel(zerolog.Disabled))
	require.NoError(t, err)
	provider := service.(client.EventsProvider)

	require.NoError(t, provider.Events(ctx, &api.EventsOpts{
		Topics:  []string{"blob_sidecar"},
		Handler: func(*apiv1.Event) {},
	}))
	require.EqualError(t, provider.Events(ctx, &api.EventsOpts{
		Topics: []string{"light_client_finality_update"},
	}), "no handler for light_client_finality_update event")
	require.EqualError(t, provider.Events(ctx, &api.EventsOpts{
		Topics: []string{"light_client_optimistic_update"},
	}), "no handler for light_client_optimistic_update event")
}

func TestEventsHeadV2(t *testing.T) {
	fixture := `{"version":"gloas","data":{"slot":"10","block":"0x9a2fefd2fdb57f74993c7780ea5b9030d2897b615b89f808011ca5aebed54eaf","state":"0x600e852a08c1200654ddf11025f1ceacb3c2e74bdd5c630cde0838b2591b69f9","payload_status":"empty","epoch_transition":false,"current_epoch_dependent_root":"0x5e0043f107cb57913498fbf2f99ff55e730bf1e151f02f221e977c91a90a0e91","next_epoch_dependent_root":"0x7f1154a218dc77addb1fe144266e83c04f67c8dfd03aacf955df05308cd8b27c","execution_optimistic":true}}`

	server := newEventsTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("topics") != "head_v2" {
			http.Error(w, "incorrect topic", http.StatusBadRequest)

			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = fmt.Fprintf(w, "event: head_v2\ndata: %s\n\n", fixture)
		w.(http.Flusher).Flush()
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service, err := eth2http.New(ctx, eth2http.WithAddress(server.URL), eth2http.WithLogLevel(zerolog.Disabled))
	require.NoError(t, err)

	received := make(chan *apiv1.HeadEventV2, 1)
	err = service.(client.EventsProvider).Events(ctx, &api.EventsOpts{
		Topics: []string{"head_v2"},
		HeadV2Handler: func(_ context.Context, event *apiv1.HeadEventV2) {
			select {
			case received <- event:
			default:
			}
		},
	})
	require.NoError(t, err)

	select {
	case event := <-received:
		cancel()
		require.Equal(t, spec.DataVersionGloas, event.Version)
		require.Equal(t, phase0.Slot(10), event.Slot)
		require.Equal(t, phase0.Root{0x9a, 0x2f, 0xef, 0xd2, 0xfd, 0xb5, 0x7f, 0x74, 0x99, 0x3c, 0x77, 0x80, 0xea, 0x5b, 0x90, 0x30, 0xd2, 0x89, 0x7b, 0x61, 0x5b, 0x89, 0xf8, 0x08, 0x01, 0x1c, 0xa5, 0xae, 0xbe, 0xd5, 0x4e, 0xaf}, event.Block)
		require.Equal(t, phase0.Root{0x60, 0x0e, 0x85, 0x2a, 0x08, 0xc1, 0x20, 0x06, 0x54, 0xdd, 0xf1, 0x10, 0x25, 0xf1, 0xce, 0xac, 0xb3, 0xc2, 0xe7, 0x4b, 0xdd, 0x5c, 0x63, 0x0c, 0xde, 0x08, 0x38, 0xb2, 0x59, 0x1b, 0x69, 0xf9}, event.State)
		require.Equal(t, "empty", event.PayloadStatus)
		require.False(t, event.EpochTransition)
		require.Equal(t, phase0.Root{0x5e, 0x00, 0x43, 0xf1, 0x07, 0xcb, 0x57, 0x91, 0x34, 0x98, 0xfb, 0xf2, 0xf9, 0x9f, 0xf5, 0x5e, 0x73, 0x0b, 0xf1, 0xe1, 0x51, 0xf0, 0x2f, 0x22, 0x1e, 0x97, 0x7c, 0x91, 0xa9, 0x0a, 0x0e, 0x91}, event.CurrentEpochDependentRoot)
		require.Equal(t, phase0.Root{0x7f, 0x11, 0x54, 0xa2, 0x18, 0xdc, 0x77, 0xad, 0xdb, 0x1f, 0xe1, 0x44, 0x26, 0x6e, 0x83, 0xc0, 0x4f, 0x67, 0xc8, 0xdf, 0xd0, 0x3a, 0xac, 0xf9, 0x55, 0xdf, 0x05, 0x30, 0x8c, 0xd8, 0xb2, 0x7c}, event.NextEpochDependentRoot)
		require.True(t, event.ExecutionOptimistic)
	case <-ctx.Done():
		t.Fatal("timed out waiting for head_v2 event")
	}
}

func newEventsTestServer(eventsHandler http.Handler) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/eth/v1/node/syncing":
			_, _ = fmt.Fprint(w, `{"data":{"head_slot":"0","sync_distance":"0","is_syncing":false,"is_optimistic":false}}`)
		case "/eth/v1/node/version":
			_, _ = fmt.Fprint(w, `{"data":{"version":"test/v1"}}`)
		case "/eth/v1/events":
			eventsHandler.ServeHTTP(w, r)
		}
	}))
}
