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

package tekuhttp

import (
	"bytes"
	"context"
	"time"

	client "github.com/attestantio/go-eth2-client"
	"github.com/pkg/errors"
)

// AddOnBeaconChainHeadUpdatedHandler adds a handler provided with beacon chain head updates.
func (s *Service) AddOnBeaconChainHeadUpdatedHandler(ctx context.Context, handler client.BeaconChainHeadUpdatedHandler) error {
	if handler == nil {
		return errors.New("no handler supplied")
	}
	s.beaconChainHeadUpdatedMutex.Lock()
	if s.beaconChainHeadUpdatedHandlers == nil {
		log.Trace().Msg("Adding first handler; starting poll")
		s.beaconChainHeadUpdatedHandlers = make([]client.BeaconChainHeadUpdatedHandler, 1, 16)
		s.beaconChainHeadUpdatedHandlers[0] = handler
		go s.pollBeaconChainHead()
	} else {
		s.beaconChainHeadUpdatedHandlers = append(s.beaconChainHeadUpdatedHandlers, handler)
	}
	s.beaconChainHeadUpdatedMutex.Unlock()
	return nil
}

// pollBeaconChainHead polls beacon chain head to feed beacon chain head update events.
// Teku does not broadcast chain head updates, so simulate with polling.
func (s *Service) pollBeaconChainHead() {
	pollPeriod := time.Duration(200) * time.Millisecond

	var lastBlockRoot []byte
	var lastSlot uint64
	for {
		head, err := s.beaconHead(s.ctx)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to poll for /beacon/head")
		} else if !bytes.Equal(head.BlockRoot, lastBlockRoot) {
			log.Trace().Uint64("slot", head.Slot).Msg("Received beacon chain head")

			slotsPerEpoch, err := s.SlotsPerEpoch(s.ctx)
			if err != nil {
				log.Warn().Err(err).Msg("Failed to obtain slots per epoch")
			}

			lastBlockRoot = head.BlockRoot
			s.beaconChainHeadUpdatedMutex.RLock()
			for i := range s.beaconChainHeadUpdatedHandlers {
				go func(handler client.BeaconChainHeadUpdatedHandler, stateRoot []byte, blockRoot []byte, slot uint64, lastSlot uint64, slotsPerEpoch uint64) {
					handler.OnBeaconChainHeadUpdated(s.ctx, slot, blockRoot, stateRoot, lastSlot/slotsPerEpoch != slot/slotsPerEpoch)
				}(s.beaconChainHeadUpdatedHandlers[i], head.StateRoot, head.BlockRoot, head.Slot, lastSlot, slotsPerEpoch)
			}
			lastSlot = head.Slot
			s.beaconChainHeadUpdatedMutex.RUnlock()
		}
		select {
		case <-time.After(pollPeriod):
			continue
		case <-s.ctx.Done():
			log.Info().Msg("Context done; stopping poll for /beacon/head")
			return
		}
	}
}
