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

// Package eventdispatch passes the events of each topic to its handler in api.EventsOpts.
package eventdispatch

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/internal/eventtopic"
	"github.com/rs/zerolog"
)

// topicHandler is the handling in api.EventsOpts of the events of a single topic.
type topicHandler struct {
	// isSet reports whether the options carry a handler specific to the topic.
	isSet func(opts *api.EventsOpts) bool
	// handle decodes the data of an event of the topic and passes it to the topic's handler in
	// the options, or failing that to their generic handler.
	handle func(ctx context.Context, opts *api.EventsOpts, data []byte) error
	// filter sets the topic's handler in dst to one that passes events to the handler in src
	// when forward allows it, leaving it nil if src has none.
	filter func(dst *api.EventsOpts, src *api.EventsOpts, forward func(topic string) bool)
}

// topicHandlers binds each field of api.EventsOpts holding the handler for a topic to that
// topic, by the name of the topic.  It is the only list of those fields.  Each field's topic is
// the one whose data decodes into the type its handler takes, so no topic is named here.
//
// Building it looks the topics up in eventtopic's registry, into which api/v1 registers them as
// it is initialised.  Go initialises api/v1 first because this package imports it; were that
// import dropped, every lookup here would panic.
var topicHandlers = topicHandlersByName(
	bind(func(o *api.EventsOpts) *api.AttestationEventHandlerFunc { return &o.AttestationHandler }),
	bind(func(o *api.EventsOpts) *api.AttesterSlashingEventHandlerFunc { return &o.AttesterSlashingHandler }),
	bind(func(o *api.EventsOpts) *api.BlobSidecarEventHandlerFunc { return &o.BlobSidecarHandler }),
	bind(func(o *api.EventsOpts) *api.BlockEventHandlerFunc { return &o.BlockHandler }),
	bind(func(o *api.EventsOpts) *api.BlockGossipEventHandlerFunc { return &o.BlockGossipHandler }),
	bind(func(o *api.EventsOpts) *api.BLSToExecutionChangeEventHandlerFunc {
		return &o.BLSToExecutionChangeHandler
	}),
	bind(func(o *api.EventsOpts) *api.ChainReorgEventHandlerFunc { return &o.ChainReorgHandler }),
	bind(func(o *api.EventsOpts) *api.ContributionAndProofEventHandlerFunc {
		return &o.ContributionAndProofHandler
	}),
	bind(func(o *api.EventsOpts) *api.DataColumnSidecarEventHandlerFunc { return &o.DataColumnSidecarHandler }),
	bind(func(o *api.EventsOpts) *api.ExecutionPayloadEventHandlerFunc { return &o.ExecutionPayloadHandler }),
	bind(func(o *api.EventsOpts) *api.ExecutionPayloadAvailableEventHandlerFunc {
		return &o.ExecutionPayloadAvailableHandler
	}),
	bind(func(o *api.EventsOpts) *api.ExecutionPayloadBidEventHandlerFunc { return &o.ExecutionPayloadBidHandler }),
	bind(func(o *api.EventsOpts) *api.ExecutionPayloadGossipEventHandlerFunc {
		return &o.ExecutionPayloadGossipHandler
	}),
	bind(func(o *api.EventsOpts) *api.FastConfirmationEventHandlerFunc { return &o.FastConfirmationHandler }),
	bind(func(o *api.EventsOpts) *api.FinalizedCheckpointEventHandlerFunc { return &o.FinalizedCheckpointHandler }),
	bind(func(o *api.EventsOpts) *api.HeadEventHandlerFunc { return &o.HeadHandler }),
	bind(func(o *api.EventsOpts) *api.HeadV2EventHandlerFunc { return &o.HeadV2Handler }),
	bindGeneric[apiv1.LightClientFinalityUpdateEvent](),
	bindGeneric[apiv1.LightClientOptimisticUpdateEvent](),
	bind(func(o *api.EventsOpts) *api.PayloadAttestationMessageEventHandlerFunc {
		return &o.PayloadAttestationMessageHandler
	}),
	bind(func(o *api.EventsOpts) *api.PayloadAttributesEventHandlerFunc { return &o.PayloadAttributesHandler }),
	bind(func(o *api.EventsOpts) *api.ProposerPreferencesEventHandlerFunc { return &o.ProposerPreferencesHandler }),
	bind(func(o *api.EventsOpts) *api.ProposerSlashingEventHandlerFunc { return &o.ProposerSlashingHandler }),
	bind(func(o *api.EventsOpts) *api.SingleAttestationEventHandlerFunc { return &o.SingleAttestationHandler }),
	bind(func(o *api.EventsOpts) *api.VoluntaryExitEventHandlerFunc { return &o.VoluntaryExitHandler }),
)

// namedTopicHandler is a topicHandler with the name of its topic.
type namedTopicHandler struct {
	name    string
	handler topicHandler
}

func topicHandlersByName(handlers ...namedTopicHandler) map[string]topicHandler {
	byName := make(map[string]topicHandler, len(handlers))
	for _, handler := range handlers {
		if _, exists := byName[handler.name]; exists {
			panic(fmt.Sprintf("event topic %s bound to more than one handler", handler.name))
		}

		byName[handler.name] = handler.handler
	}

	return byName
}

// bind binds the field of api.EventsOpts, returned by field, holding the handler for the topic
// whose data decodes into the type the handler takes.  It panics if there is no such topic, so
// that a binding with none fails as the package is initialised.
func bind[T any, H ~func(context.Context, *T)](field func(opts *api.EventsOpts) *H) namedTopicHandler {
	topic := eventtopic.Lookup[T]()
	name := topic.Name()

	return namedTopicHandler{
		name: name,
		handler: topicHandler{
			isSet: func(opts *api.EventsOpts) bool {
				return *field(opts) != nil
			},
			handle: func(ctx context.Context, opts *api.EventsOpts, input []byte) error {
				data, err := topic.Decode(input)
				if err != nil {
					return err
				}

				if handler := *field(opts); handler != nil {
					handler(ctx, data)

					return nil
				}

				if opts.Handler == nil {
					zerolog.Ctx(ctx).Debug().Str("topic", name).Msg("No specific or generic handler supplied; ignoring")

					return nil
				}

				opts.Handler(&apiv1.Event{
					Topic: name,
					Data:  data,
				})

				return nil
			},
			filter: func(dst *api.EventsOpts, src *api.EventsOpts, forward func(topic string) bool) {
				handler := *field(src)
				if handler == nil {
					return
				}

				*field(dst) = H(func(ctx context.Context, data *T) {
					if forward(name) {
						handler(ctx, data)
					}
				})
			},
		},
	}
}

// bindGeneric binds a topic that has only the generic handler.
//
// These topics have no typed handler field in api.EventsOpts, so isSet is
// always false and ValidateEventsOpts will not accept a subscription to one
// without opts.Handler also being set.  That is intended -- there is no other
// way to deliver them -- but it does make them the only topics that cannot be
// subscribed to with a typed handler alone.
func bindGeneric[T any]() namedTopicHandler {
	topic := eventtopic.Lookup[T]()
	name := topic.Name()

	return namedTopicHandler{
		name: name,
		handler: topicHandler{
			isSet: func(*api.EventsOpts) bool { return false },
			handle: func(ctx context.Context, opts *api.EventsOpts, input []byte) error {
				data, err := topic.Decode(input)
				if err != nil {
					return err
				}

				if opts.Handler == nil {
					// Matches bind: an unsolicited event from a node is worth a
					// line at debug level rather than being dropped silently.
					zerolog.Ctx(ctx).Debug().Str("topic", name).Msg("No generic handler supplied; ignoring")

					return nil
				}

				opts.Handler(&apiv1.Event{Topic: name, Data: data})

				return nil
			},
			filter: func(*api.EventsOpts, *api.EventsOpts, func(string) bool) {},
		},
	}
}

// ErrUnsupportedTopic is returned by Handle for an event of a topic that is not supported.
var ErrUnsupportedTopic = errors.New("unsupported event topic")

// Supports reports whether events of the given topic can be handled.  It is to be used in place
// of apiv1.SupportedEventTopics, which is an exported map that callers can change.
func Supports(topic string) bool {
	_, exists := topicHandlers[topic]

	return exists
}

// HasTopicHandler reports whether the options carry a handler specific to the given topic.
func HasTopicHandler(opts *api.EventsOpts, topic string) bool {
	handler, exists := topicHandlers[topic]

	return exists && handler.isSet(opts)
}

// Handle decodes the data of an event with the given topic, as sent on the events stream, and
// passes it to the topic's specific handler in the options, or failing that to their generic
// handler.  It returns an error if the topic is not supported or the data does not decode.
func Handle(ctx context.Context, opts *api.EventsOpts, topic string, data []byte) error {
	handler, exists := topicHandlers[topic]
	if !exists {
		return fmt.Errorf("%w %s", ErrUnsupportedTopic, topic)
	}

	return handler.handle(ctx, opts, data)
}

// ForwardingGuarded returns a copy of opts whose handlers forward events only when forward
// allows their topic.  It leaves missing topic handlers nil so events can reach the generic
// handler.  It copies the topics slice and never changes opts, so guards for different clients
// do not stack.
func ForwardingGuarded(opts *api.EventsOpts, forward func(topic string) bool) *api.EventsOpts {
	filtered := &api.EventsOpts{
		Common: opts.Common,
		Topics: slices.Clone(opts.Topics),
	}

	if handler := opts.Handler; handler != nil {
		filtered.Handler = func(event *apiv1.Event) {
			if forward(event.Topic) {
				handler(event)
			}
		}
	}

	for _, topicHandler := range topicHandlers {
		topicHandler.filter(filtered, opts, forward)
	}

	return filtered
}
