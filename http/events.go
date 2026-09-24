// Copyright © 2020 - 2026 Attestant Limited.
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
	"errors"
	"fmt"
	"maps"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/r3labs/sse/v2"
	"github.com/rs/zerolog"
)

// Events feeds requested events with the given topics to the supplied handler.
func (s *Service) Events(ctx context.Context, opts *api.EventsOpts) error {
	if err := s.assertIsActive(ctx); err != nil {
		return err
	}

	if err := ValidateEventsOpts(opts); err != nil {
		return err
	}

	// #nosec G404
	log := s.log.With().Str("id", fmt.Sprintf("%02x", rand.Int31())).Str("address", s.address).Logger()
	ctx = log.WithContext(ctx)

	endpoint := "/eth/v1/events"
	query := "topics=" + strings.Join(opts.Topics, "&topics=")
	callURL := urlForCall(s.base, endpoint, query)
	log.Trace().Str("url", callURL.String()).Msg("GET request to events stream")

	sseClient := sse.NewClient(callURL.String())
	maps.Copy(sseClient.Headers, s.extraHeaders)

	if _, exists := sseClient.Headers["User-Agent"]; !exists {
		sseClient.Headers["User-Agent"] = defaultUserAgent
	}

	sseClient.Headers["Accept"] = "text/event-stream"
	sseClient.Connection.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 2 * time.Second,
		}).Dial,
	}

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				log.Trace().Msg("Connecting to events stream")

				if err := sseClient.SubscribeRawWithContext(ctx, func(msg *sse.Event) {
					s.handleEvent(ctx, msg, opts)
				}); err != nil {
					log.Error().Err(err).Msg("Failed to subscribe to event stream")
				}

				log.Trace().Msg("Events stream disconnected")
			case <-ctx.Done():
				log.Debug().Msg("Context done")

				return
			}
		}
	}()

	return nil
}

// ValidateEventsOpts checks the options for an events subscription.
func ValidateEventsOpts(opts *api.EventsOpts) error {
	if opts == nil {
		return client.ErrNoOptions
	}
	if len(opts.Topics) == 0 {
		return errors.Join(errors.New("no topics supplied"), client.ErrInvalidOptions)
	}

	for _, topic := range opts.Topics {
		if _, exists := apiv1.SupportedEventTopics[topic]; !exists {
			return fmt.Errorf("unsupported event topic %s: %w", topic, client.ErrInvalidOptions)
		}

		if opts.Handler == nil && !opts.HasTopicHandler(topic) {
			return fmt.Errorf("no handler for %s event: %w", topic, client.ErrInvalidOptions)
		}
	}

	return nil
}

// handleEvent handles all events.
func (*Service) handleEvent(ctx context.Context,
	msg *sse.Event,
	opts *api.EventsOpts,
) {
	log := zerolog.Ctx(ctx)

	if msg == nil {
		log.Debug().Msg("No message supplied; ignoring")

		return
	}

	topic := string(msg.Event)

	switch {
	case len(topic) == 0:
		// Used as keepalive.  Ignore.
	case !apiv1.SupportedEventTopics[topic]:
		log.Warn().Str("topic", topic).Msg("Received message with unhandled topic; ignoring")
	default:
		if err := opts.HandleEvent(ctx, topic, msg.Data); err != nil {
			log.Error().Err(err).Str("topic", topic).RawJSON("data", msg.Data).Msg("Failed to parse event")
		}
	}
}
