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
	"context"
	"encoding/json"
	"fmt"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
)

type beaconCommitteesResponse struct {
	Slot      uint64   `json:"slot"`
	Index     uint64   `json:"index"`
	Committee []uint64 `json:"committee"`
}

// BeaconCommittees fetches the chain's beacon committees given a state.
func (s *Service) BeaconCommittees(ctx context.Context, stateID string) ([]*api.BeaconCommittee, error) {
	slot, err := s.slotFromStateID(ctx, stateID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse state ID")
	}
	epoch := slot / (*s.slotsPerEpoch)

	respBodyReader, err := s.get(ctx, fmt.Sprintf("/beacon/committees?epoch=%d", epoch))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon committees")
	}
	defer func() {
		if err := respBodyReader.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close HTTP body")
		}
	}()

	beaconCommitteesResponse := make([]*beaconCommitteesResponse, 0)
	if err := json.NewDecoder(respBodyReader).Decode(&beaconCommitteesResponse); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon committees")
	}

	resp := make([]*api.BeaconCommittee, len(beaconCommitteesResponse))
	for i := range beaconCommitteesResponse {
		resp[i] = &api.BeaconCommittee{
			Slot:       beaconCommitteesResponse[i].Slot,
			Index:      beaconCommitteesResponse[i].Index,
			Validators: beaconCommitteesResponse[i].Committee,
		}
	}

	return resp, nil
}
