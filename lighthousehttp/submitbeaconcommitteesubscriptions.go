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

package lighthousehttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	client "github.com/attestantio/go-eth2-client"
	"github.com/pkg/errors"
)

// SubmitBeaconCommitteeSubscriptions subscribes to beacon committees.
func (s *Service) SubmitBeaconCommitteeSubscriptions(ctx context.Context, subscriptions []*client.BeaconCommitteeSubscription) error {
	type subscriptionReq struct {
		Slot           uint64 `json:"slot"`
		CommitteeIndex uint64 `json:"attestation_committee_index"`
		CommitteeSize  uint64 `json:"committee_count_at_slot"`
		ValidatorIndex uint64 `json:"validator_index"`
		Aggregate      bool   `json:"is_aggregator"`
	}
	reqBody := make([]*subscriptionReq, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		reqBody = append(reqBody, &subscriptionReq{
			Slot:           subscription.Slot,
			CommitteeIndex: subscription.CommitteeIndex,
			CommitteeSize:  subscription.CommitteeSize,
			ValidatorIndex: subscription.ValidatorIndex,
			Aggregate:      subscription.Aggregate,
		})
	}

	var reqBodyReader bytes.Buffer
	if err := json.NewEncoder(&reqBodyReader).Encode(reqBody); err != nil {
		return errors.Wrap(err, "failed to encode beacon committee subscriptions request")
	}

	respBodyReader, cancel, err := s.post(ctx, "/validator/subscribe", &reqBodyReader)
	if err != nil {
		return errors.Wrap(err, "failed to request beacon committee subscriptions")
	}
	defer cancel()

	body, err := ioutil.ReadAll(respBodyReader)
	if err != nil {
		return errors.Wrap(err, "failed to read beacon committee subscriptions response")
	}
	if string(body) != "null" {
		log.Warn().Str("error", string(body)).Msg("Bad response from server on beacon committee subscription request")
		return fmt.Errorf("server rejected beacon committee subscriptions request")
	}

	return nil
}
