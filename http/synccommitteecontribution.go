// Copyright © 2021 - 2024 Attestant Limited.
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
	"errors"
	"fmt"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/altair"
)

// SyncCommitteeContribution provides a sync committee contribution.
func (s *Service) SyncCommitteeContribution(ctx context.Context,
	opts *api.SyncCommitteeContributionOpts,
) (
	*api.Response[*altair.SyncCommitteeContribution],
	error,
) {
	if err := s.assertIsActive(ctx); err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, client.ErrNoOptions
	}
	if opts.BeaconBlockRoot.IsZero() {
		return nil, errors.Join(errors.New("no beacon block root specified"), client.ErrInvalidOptions)
	}

	endpoint := "/eth/v1/validator/sync_committee_contribution"
	query := fmt.Sprintf("slot=%d&subcommittee_index=%d&beacon_block_root=%#x", opts.Slot, opts.SubcommitteeIndex, opts.BeaconBlockRoot)
	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), altair.SyncCommitteeContribution{})
	if err != nil {
		return nil, err
	}

	// Confirm the contribution is for the requested slot.
	if data.Slot != opts.Slot {
		return nil, errors.Join(fmt.Errorf("sync committee contiribution for slot %d; expected %d", data.Slot, opts.Slot), client.ErrInconsistentResult)
	}

	// Confirm the beacon block root is correct.
	if !bytes.Equal(data.BeaconBlockRoot[:], opts.BeaconBlockRoot[:]) {
		return nil, errors.Join(fmt.Errorf("sync committee proposal has beacon bock root %#x; expected %#x", data.BeaconBlockRoot[:], opts.BeaconBlockRoot[:]), client.ErrInconsistentResult)
	}

	return &api.Response[*altair.SyncCommitteeContribution]{
		Metadata: metadata,
		Data:     &data,
	}, nil
}
