// Copyright © 2021, 2023 Attestant Limited.
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
	"bytes"
	"context"
	"fmt"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/pkg/errors"
)

// SyncCommitteeContribution provides a sync committee contribution.
func (s *Service) SyncCommitteeContribution(ctx context.Context,
	opts *api.SyncCommitteeContributionOpts,
) (
	*api.Response[*altair.SyncCommitteeContribution],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if opts.BeaconBlockRoot.IsZero() {
		return nil, errors.New("no beacon block root specified")
	}

	url := fmt.Sprintf("/eth/v1/validator/sync_committee_contribution?slot=%d&subcommittee_index=%d&beacon_block_root=%#x", opts.Slot, opts.SubcommitteeIndex, opts.BeaconBlockRoot)
	httpResponse, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), altair.SyncCommitteeContribution{})
	if err != nil {
		return nil, err
	}

	// Confirm the contribution is for the requested slot.
	if data.Slot != opts.Slot {
		return nil, fmt.Errorf("received sync committee contribution for slot %d; expected %d", data.Slot, opts.Slot)
	}

	// Confirm the beacon block root is correct.
	if !bytes.Equal(data.BeaconBlockRoot[:], opts.BeaconBlockRoot[:]) {
		return nil, errors.New("sync committee contribution not for requested beacon block root")
	}

	return &api.Response[*altair.SyncCommitteeContribution]{
		Metadata: metadata,
		Data:     &data,
	}, nil
}
