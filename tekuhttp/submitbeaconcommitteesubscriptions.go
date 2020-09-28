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
	"encoding/json"
	"fmt"
	"io/ioutil"

	client "github.com/attestantio/go-eth2-client"
	"github.com/pkg/errors"
)

type subscriptionReq struct {
	AggregationSlot string `json:"aggregation_slot"`
	CommitteeIndex  string `json:"committee_index"`
}

// SubmitBeaconCommitteeSubscriptions subscribes to beacon committees.
func (s *Service) SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*client.BeaconCommitteeSubscription) error {
	hasErrors := false
	for _, subscription := range subscriptions {
		reqBody := &subscriptionReq{
			AggregationSlot: fmt.Sprintf("%d", subscription.Slot),
			CommitteeIndex:  fmt.Sprintf("%d", subscription.CommitteeIndex),
		}
		if err := s.submitBeaconCommitteeSubscription(ctx, reqBody); err != nil {
			// We want to subscribe to as many subnets as we can.
			// Rather than exit, log the error and set a flag.
			log.Error().Err(err).Msg("Failed to subscribe to beacon committee")
			hasErrors = true
		}
	}

	if hasErrors {
		return errors.New("submitted with errors")
	}
	return nil
}

func (s *Service) submitBeaconCommitteeSubscription(ctx context.Context, reqBody *subscriptionReq) error {
	var reqBodyReader bytes.Buffer
	if err := json.NewEncoder(&reqBodyReader).Encode(reqBody); err != nil {
		return errors.Wrap(err, "failed to encode beacon committee subscription request")
	}

	respBodyReader, err := s.post(ctx, "/validator/beacon_committee_subscription", &reqBodyReader)
	if err != nil {
		return errors.Wrap(err, "failed to request beacon committee subscription")
	}

	var resp []byte
	if respBodyReader != nil {
		resp, err = ioutil.ReadAll(respBodyReader)
		if err != nil {
			resp = nil
		}
	}
	if err != nil {
		return errors.Wrap(err, "failed to obtain error message for beacon committee subscription")
	}
	if len(resp) > 0 {
		return fmt.Errorf("failed to submit beacon committee subscription: %s", string(resp))
	}

	return nil
}
