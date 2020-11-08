// Copyright © 2020 Attestant Limited.
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

package prysmgrpc

import (
	"context"

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	spec "github.com/attestantio/go-eth2-client/spec/phase0"
)

type eventPackager struct {
	handler client.EventHandlerFunc
}

func (e *eventPackager) OnBeaconChainHeadUpdated(ctx context.Context, slot uint64, blockRoot []byte, stateRoot []byte, epochTransition bool) {
	event := &api.Event{
		Topic: "head",
		Data: &api.HeadEvent{
			Slot:            spec.Slot(slot),
			EpochTransition: epochTransition,
		},
	}
	copy(event.Data.(*api.HeadEvent).Block[:], blockRoot)
	copy(event.Data.(*api.HeadEvent).State[:], stateRoot)
	e.handler(event)
}

// Events feeds requested events with the given topics to the supplied handler.
func (s *Service) Events(ctx context.Context, topics []string, handler client.EventHandlerFunc) error {
	packager := &eventPackager{handler: handler}
	return s.AddOnBeaconChainHeadUpdatedHandler(ctx, packager)
}
