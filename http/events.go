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
	"encoding/json"
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

		if opts.Handler != nil {
			continue
		}

		if err := checkEventSpecificHandler(opts, topic); err != nil {
			return fmt.Errorf("%w: %w", err, client.ErrInvalidOptions)
		}
	}

	return nil
}

func checkEventSpecificHandler(opts *api.EventsOpts, topic string) error {
	handling, exists := eventTopics[topic]
	if !exists {
		return fmt.Errorf("unsupported event %s", topic)
	}

	if !handling.hasHandler(opts) {
		return fmt.Errorf("no handler for %s event", topic)
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

	handling, exists := eventTopics[string(msg.Event)]

	switch {
	case exists:
		handling.handle(ctx, msg, opts)
	case len(msg.Event) == 0:
		// Used as keepalive.  Ignore.
	default:
		log.Warn().Str("topic", string(msg.Event)).Msg("Received message with unhandled topic; ignoring")
	}
}

// eventTopic is the handling of the events of a single topic.
type eventTopic struct {
	// hasHandler reports whether the options carry a handler specific to the topic.
	hasHandler func(opts *api.EventsOpts) bool
	// handle decodes an event of the topic and passes it to the handler for it in the options.
	handle func(ctx context.Context, msg *sse.Event, opts *api.EventsOpts)
}

// eventTopics is the handling of each event topic the client can decode, by topic.
var eventTopics = map[string]eventTopic{
	"attestation": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.AttestationEventHandlerFunc { return o.AttestationHandler }),
	"attester_slashing": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.AttesterSlashingEventHandlerFunc { return o.AttesterSlashingHandler }),
	"blob_sidecar": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.BlobSidecarEventHandlerFunc { return o.BlobSidecarHandler }),
	"block": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.BlockEventHandlerFunc { return o.BlockHandler }),
	"block_gossip": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.BlockGossipEventHandlerFunc { return o.BlockGossipHandler }),
	"bls_to_execution_change": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.BLSToExecutionChangeEventHandlerFunc { return o.BLSToExecutionChangeHandler }),
	"chain_reorg": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.ChainReorgEventHandlerFunc { return o.ChainReorgHandler }),
	"contribution_and_proof": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.ContributionAndProofEventHandlerFunc { return o.ContributionAndProofHandler }),
	"data_column_sidecar": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.DataColumnSidecarEventHandlerFunc { return o.DataColumnSidecarHandler }),
	"execution_payload": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.ExecutionPayloadEventHandlerFunc { return o.ExecutionPayloadHandler }),
	"execution_payload_available": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.ExecutionPayloadAvailableEventHandlerFunc {
			return o.ExecutionPayloadAvailableHandler
		}),
	"execution_payload_bid": newEventTopic(unmarshalVersionedEventData,
		func(o *api.EventsOpts) api.ExecutionPayloadBidEventHandlerFunc { return o.ExecutionPayloadBidHandler }),
	"execution_payload_gossip": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.ExecutionPayloadGossipEventHandlerFunc {
			return o.ExecutionPayloadGossipHandler
		}),
	"fast_confirmation": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.FastConfirmationEventHandlerFunc { return o.FastConfirmationHandler }),
	"finalized_checkpoint": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.FinalizedCheckpointEventHandlerFunc { return o.FinalizedCheckpointHandler }),
	"head": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.HeadEventHandlerFunc { return o.HeadHandler }),
	"payload_attestation_message": newEventTopic(unmarshalVersionedEventData,
		func(o *api.EventsOpts) api.PayloadAttestationMessageEventHandlerFunc {
			return o.PayloadAttestationMessageHandler
		}),
	"payload_attributes": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.PayloadAttributesEventHandlerFunc { return o.PayloadAttributesHandler }),
	"proposer_preferences": newEventTopic(unmarshalVersionedEventData,
		func(o *api.EventsOpts) api.ProposerPreferencesEventHandlerFunc { return o.ProposerPreferencesHandler }),
	"proposer_slashing": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.ProposerSlashingEventHandlerFunc { return o.ProposerSlashingHandler }),
	"single_attestation": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.SingleAttestationEventHandlerFunc { return o.SingleAttestationHandler }),
	"voluntary_exit": newEventTopic(json.Unmarshal,
		func(o *api.EventsOpts) api.VoluntaryExitEventHandlerFunc { return o.VoluntaryExitHandler }),
}

// newEventTopic creates the handling of a topic whose events decode, with the given decoder,
// into a T.  handler selects the topic's specific handler from the options; an event is passed
// to it if it is set, and otherwise to the generic handler.
func newEventTopic[T any, H ~func(context.Context, *T)](decode func([]byte, any) error,
	handler func(opts *api.EventsOpts) H,
) eventTopic {
	return eventTopic{
		hasHandler: func(opts *api.EventsOpts) bool {
			return handler(opts) != nil
		},
		handle: func(ctx context.Context, msg *sse.Event, opts *api.EventsOpts) {
			log := zerolog.Ctx(ctx)

			data := new(T)
			if err := decode(msg.Data, data); err != nil {
				log.Error().Err(err).Str("topic", string(msg.Event)).RawJSON("data", msg.Data).Msg("Failed to parse event")

				return
			}

			switch specific := handler(opts); {
			case specific != nil:
				specific(ctx, data)
			case opts.Handler != nil:
				opts.Handler(&apiv1.Event{
					Topic: string(msg.Event),
					Data:  data,
				})
			default:
				log.Debug().Msg("No specific or generic handler supplied; ignoring")
			}
		},
	}
}

// unmarshalVersionedEventData decodes the payload of an SSE event that the beacon-API spec wraps
// as {"version": "...", "data": {...}}, being those of the execution_payload_bid,
// payload_attestation_message and proposer_preferences topics.  A bare, unwrapped object is
// accepted as well, for nodes that do not wrap it.
func unmarshalVersionedEventData(raw []byte, v any) error {
	var wrapper struct {
		Version string          `json:"version"`
		Data    json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(raw, &wrapper); err == nil && len(wrapper.Data) > 0 && wrapper.Version != "" {
		return json.Unmarshal(wrapper.Data, v)
	}

	return json.Unmarshal(raw, v)
}
