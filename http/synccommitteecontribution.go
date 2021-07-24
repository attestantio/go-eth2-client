// Copyright © 2021 Attestant Limited.
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
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type syncCommitteeContributionJSON struct {
	Data *altair.SyncCommitteeContribution `json:"data"`
}

// SyncCommitteeContribution provides a sync committee contribution.
func (s *Service) SyncCommitteeContribution(ctx context.Context,
	slot phase0.Slot,
	subcommitteeIndex uint64,
	beaconBlockRoot phase0.Root,
) (
	*altair.SyncCommitteeContribution,
	error,
) {
	url := fmt.Sprintf("/eth/v1/validator/sync_committee_contribution?slot=%d&subcommittee_index=%d&beacon_block_root=%#x", slot, subcommitteeIndex, beaconBlockRoot)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		log.Trace().Str("url", url).Err(err).Msg("Request failed")
		return nil, errors.Wrap(err, "failed to request sync committee contribution")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain sync committee contribution")
	}

	var resp syncCommitteeContributionJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse sync committee contribution")
	}

	return resp.Data, nil
}
