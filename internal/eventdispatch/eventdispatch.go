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
	"fmt"
	"slices"

	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/internal/eventtopic"
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

				if opts.Handler != nil {
					opts.Handler(&apiv1.Event{
						Topic: name,
						Data:  data,
					})
				}

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
		return fmt.Errorf("unsupported event topic %s", topic)
	}

	return handler.handle(ctx, opts, data)
}

// Filtered returns new options with the same topics and common options as those given, in which
// each handler passes an event on to the matching handler in the given options only if forward
// allows the event's topic.  A handler the given options do not have is left nil, because events
// of a topic whose specific handler is nil go to the generic handler, and setting one would
// starve it.
//
// The returned options are always a struct of their own and never the given options modified in
// place, so filtering the same options more than once never has one filter wrap another.
func Filtered(opts *api.EventsOpts, forward func(topic string) bool) *api.EventsOpts {
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
