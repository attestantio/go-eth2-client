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
	"fmt"
	"io"

	client "github.com/attestantio/go-eth2-client"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// AddOnBeaconChainHeadUpdatedHandler adds a handler provided with beacon chain head updates.
func (s *Service) AddOnBeaconChainHeadUpdatedHandler(ctx context.Context, handler client.BeaconChainHeadUpdatedHandler) error {
	if handler == nil {
		return errors.New("no handler supplied")
	}
	s.beaconChainHeadUpdatedMutex.Lock()
	if s.beaconChainHeadUpdatedHandlers == nil {
		log.Trace().Msg("Adding first handler; starting stream")
		s.beaconChainHeadUpdatedHandlers = make([]client.BeaconChainHeadUpdatedHandler, 1, 16)
		s.beaconChainHeadUpdatedHandlers[0] = handler
		go s.streamBeaconChainHead(ctx)
	} else {
		s.beaconChainHeadUpdatedHandlers = append(s.beaconChainHeadUpdatedHandlers, handler)
	}
	s.beaconChainHeadUpdatedMutex.Unlock()
	return nil
}

// streamBeaconChainHead streams beacon chain head to feed beacon chain head update events.
func (s *Service) streamBeaconChainHead(ctx context.Context) {
	conn := ethpb.NewBeaconChainClient(s.conn)
	log.Trace().Msg("Calling StreamChainHead()")
	stream, err := conn.StreamChainHead(s.ctx, &types.Empty{})
	if err != nil {
		log.Warn().Err(err).Msg("failed to open chain head stream")
		return
	}
	defer func() {
		if err := stream.CloseSend(); err != nil {
			log.Warn().Err(err).Msg("failed to close chain head stream")
		}
	}()
	lastEpoch := uint64(0)
	for {
		beaconChainHead, err := stream.Recv()
		if err == io.EOF {
			// Natural EOF.
			return
		}
		if err != nil {
			// Unnatural error.
			log.Warn().Err(err).Msg("received error from blocks stream")
			return
		}
		if beaconChainHead != nil {
			log.Trace().Uint64("slot", beaconChainHead.HeadSlot).Msg("Received beacon chain head")

			// Need the state root for this slot.
			signedBeaconBlock, err := s.SignedBeaconBlock(ctx, fmt.Sprintf("%d", beaconChainHead.HeadSlot))
			if err != nil {
				log.Warn().Err(err).Msg("failed to obtain block for slot")
				return
			}
			if signedBeaconBlock == nil {
				log.Warn().Err(err).Msg("obtained nil block for slot")
				return
			}

			s.beaconChainHeadUpdatedMutex.RLock()
			for i := range s.beaconChainHeadUpdatedHandlers {
				go func(handler client.BeaconChainHeadUpdatedHandler) {
					handler.OnBeaconChainHeadUpdated(s.ctx, beaconChainHead.HeadSlot, beaconChainHead.HeadBlockRoot, signedBeaconBlock.Message.StateRoot[:], beaconChainHead.HeadEpoch != lastEpoch)
				}(s.beaconChainHeadUpdatedHandlers[i])
			}
			lastEpoch = beaconChainHead.HeadEpoch
			s.beaconChainHeadUpdatedMutex.RUnlock()
		}
	}
}
