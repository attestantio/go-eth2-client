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
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
)

// SubmitBeaconCommitteeSubscriptions subscribes to beacon committees.
func (s *Service) SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*client.BeaconCommitteeSubscription) error {
	client := ethpb.NewBeaconNodeValidatorClient(s.conn)
	slots := make([]uint64, len(subscriptions))
	committeeIds := make([]uint64, len(subscriptions))
	isAggregator := make([]bool, len(subscriptions))
	for i, subscription := range subscriptions {
		slots[i] = subscription.Slot
		committeeIds[i] = subscription.CommitteeIndex
		isAggregator[i] = subscription.Aggregate
	}

	log.Trace().Msg("Calling SubscribeCommitteeSubnets()")
	_, err := client.SubscribeCommitteeSubnets(ctx, &ethpb.CommitteeSubnetsSubscribeRequest{
		Slots:        slots,
		CommitteeIds: committeeIds,
		IsAggregator: isAggregator,
	})
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to beacon committees")
	}

	return nil
}
