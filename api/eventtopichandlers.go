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

package api

import (
	"context"
	"fmt"
	"slices"

	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
)

// eventTopicHandler is the handling in EventsOpts of the events of a single topic.
type eventTopicHandler struct {
	// isSet reports whether the options carry a handler specific to the topic.
	isSet func(opts *EventsOpts) bool
	// handle decodes the data of an event of the topic and passes it to the topic's handler in
	// the options, or failing that to their generic handler.
	handle func(ctx context.Context, opts *EventsOpts, data []byte) error
	// filter sets the topic's handler in dst to one that passes events to the handler in src
	// when forward allows it, leaving it nil if src has none.
	filter func(dst *EventsOpts, src *EventsOpts, forward func(topic string) bool)
}

// eventTopicHandlers binds each topic in apiv1 to the field of EventsOpts that holds its
// handler, by topic.  It is the only list of those fields.
var eventTopicHandlers = eventTopicHandlersByName(
	bindEventTopic(apiv1.AttestationEventTopic,
		func(o *EventsOpts) *AttestationEventHandlerFunc { return &o.AttestationHandler }),
	bindEventTopic(apiv1.AttesterSlashingEventTopic,
		func(o *EventsOpts) *AttesterSlashingEventHandlerFunc { return &o.AttesterSlashingHandler }),
	bindEventTopic(apiv1.BlobSidecarEventTopic,
		func(o *EventsOpts) *BlobSidecarEventHandlerFunc { return &o.BlobSidecarHandler }),
	bindEventTopic(apiv1.BlockEventTopic,
		func(o *EventsOpts) *BlockEventHandlerFunc { return &o.BlockHandler }),
	bindEventTopic(apiv1.BlockGossipEventTopic,
		func(o *EventsOpts) *BlockGossipEventHandlerFunc { return &o.BlockGossipHandler }),
	bindEventTopic(apiv1.BLSToExecutionChangeEventTopic,
		func(o *EventsOpts) *BLSToExecutionChangeEventHandlerFunc { return &o.BLSToExecutionChangeHandler }),
	bindEventTopic(apiv1.ChainReorgEventTopic,
		func(o *EventsOpts) *ChainReorgEventHandlerFunc { return &o.ChainReorgHandler }),
	bindEventTopic(apiv1.ContributionAndProofEventTopic,
		func(o *EventsOpts) *ContributionAndProofEventHandlerFunc { return &o.ContributionAndProofHandler }),
	bindEventTopic(apiv1.DataColumnSidecarEventTopic,
		func(o *EventsOpts) *DataColumnSidecarEventHandlerFunc { return &o.DataColumnSidecarHandler }),
	bindEventTopic(apiv1.ExecutionPayloadEventTopic,
		func(o *EventsOpts) *ExecutionPayloadEventHandlerFunc { return &o.ExecutionPayloadHandler }),
	bindEventTopic(apiv1.ExecutionPayloadAvailableEventTopic,
		func(o *EventsOpts) *ExecutionPayloadAvailableEventHandlerFunc {
			return &o.ExecutionPayloadAvailableHandler
		}),
	bindEventTopic(apiv1.ExecutionPayloadBidEventTopic,
		func(o *EventsOpts) *ExecutionPayloadBidEventHandlerFunc { return &o.ExecutionPayloadBidHandler }),
	bindEventTopic(apiv1.ExecutionPayloadGossipEventTopic,
		func(o *EventsOpts) *ExecutionPayloadGossipEventHandlerFunc { return &o.ExecutionPayloadGossipHandler }),
	bindEventTopic(apiv1.FastConfirmationEventTopic,
		func(o *EventsOpts) *FastConfirmationEventHandlerFunc { return &o.FastConfirmationHandler }),
	bindEventTopic(apiv1.FinalizedCheckpointEventTopic,
		func(o *EventsOpts) *FinalizedCheckpointEventHandlerFunc { return &o.FinalizedCheckpointHandler }),
	bindEventTopic(apiv1.HeadEventTopic,
		func(o *EventsOpts) *HeadEventHandlerFunc { return &o.HeadHandler }),
	bindEventTopic(apiv1.PayloadAttestationMessageEventTopic,
		func(o *EventsOpts) *PayloadAttestationMessageEventHandlerFunc {
			return &o.PayloadAttestationMessageHandler
		}),
	bindEventTopic(apiv1.PayloadAttributesEventTopic,
		func(o *EventsOpts) *PayloadAttributesEventHandlerFunc { return &o.PayloadAttributesHandler }),
	bindEventTopic(apiv1.ProposerPreferencesEventTopic,
		func(o *EventsOpts) *ProposerPreferencesEventHandlerFunc { return &o.ProposerPreferencesHandler }),
	bindEventTopic(apiv1.ProposerSlashingEventTopic,
		func(o *EventsOpts) *ProposerSlashingEventHandlerFunc { return &o.ProposerSlashingHandler }),
	bindEventTopic(apiv1.SingleAttestationEventTopic,
		func(o *EventsOpts) *SingleAttestationEventHandlerFunc { return &o.SingleAttestationHandler }),
	bindEventTopic(apiv1.VoluntaryExitEventTopic,
		func(o *EventsOpts) *VoluntaryExitEventHandlerFunc { return &o.VoluntaryExitHandler }),
)

// namedEventTopicHandler is an eventTopicHandler with the name of its topic.
type namedEventTopicHandler struct {
	name    string
	handler eventTopicHandler
}

func eventTopicHandlersByName(handlers ...namedEventTopicHandler) map[string]eventTopicHandler {
	byName := make(map[string]eventTopicHandler, len(handlers))
	for _, handler := range handlers {
		byName[handler.name] = handler.handler
	}

	return byName
}

// bindEventTopic binds a topic to the field of EventsOpts, returned by field, holding its handler.
// The handler must take the topic's decode type, so that binding a topic to a handler for data of
// another type does not compile.
func bindEventTopic[T any, H ~func(context.Context, *T)](topic apiv1.EventTopic[T],
	field func(opts *EventsOpts) *H,
) namedEventTopicHandler {
	name := topic.Name()

	return namedEventTopicHandler{
		name: name,
		handler: eventTopicHandler{
			isSet: func(opts *EventsOpts) bool {
				return *field(opts) != nil
			},
			handle: func(ctx context.Context, opts *EventsOpts, input []byte) error {
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
			filter: func(dst *EventsOpts, src *EventsOpts, forward func(topic string) bool) {
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
func (o *EventsOpts) HasTopicHandler(topic string) bool {
	handler, exists := eventTopicHandlers[topic]

	return exists && handler.isSet(o)
}

// HandleEvent decodes the data of an event with the given topic, as sent on the events stream,
// and passes it to the topic's specific handler, or failing that to the generic handler.  It
// returns an error if the topic is not supported or the data does not decode.
func (o *EventsOpts) HandleEvent(ctx context.Context, topic string, data []byte) error {
	handler, exists := eventTopicHandlers[topic]
	if !exists {
		return fmt.Errorf("unsupported event topic %s", topic)
	}

	return handler.handle(ctx, o, data)
}

// Filtered returns new options with the same topics and common options, in which each handler
// passes an event on to the matching handler in these options only if forward allows the event's
// topic.  A handler these options do not have is left nil, because events of a topic whose
// specific handler is nil go to the generic handler, and setting one would starve it.
//
// The returned options are always a struct of their own and never these options modified in
// place, so filtering the same options more than once never has one filter wrap another.
func (o *EventsOpts) Filtered(forward func(topic string) bool) *EventsOpts {
	filtered := &EventsOpts{
		Common: o.Common,
		Topics: slices.Clone(o.Topics),
	}

	if handler := o.Handler; handler != nil {
		filtered.Handler = func(event *apiv1.Event) {
			if forward(event.Topic) {
				handler(event)
			}
		}
	}

	for _, topicHandler := range eventTopicHandlers {
		topicHandler.filter(filtered, o, forward)
	}

	return filtered
}
