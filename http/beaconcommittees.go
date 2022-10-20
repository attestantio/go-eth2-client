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

package http

import (
	"context"
	"encoding/json"
	"fmt"

	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type beaconCommitteesJSON struct {
	Data []*api.BeaconCommittee `json:"data"`
}

// BeaconCommittees fetches all beacon committees for the epoch at the given state.
func (s *Service) BeaconCommittees(ctx context.Context, stateID string) ([]*api.BeaconCommittee, error) {
	url := fmt.Sprintf("/eth/v1/beacon/states/%s/committees", stateID)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon committees")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain beacon committees")
	}

	var resp beaconCommitteesJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon committees")
	}

	return resp.Data, nil
}

// BeaconCommitteesAtEpoch fetches all beacon committees for the given epoch at the given state.
func (s *Service) BeaconCommitteesAtEpoch(ctx context.Context, stateID string, epoch phase0.Epoch) ([]*api.BeaconCommittee, error) {
	url := fmt.Sprintf("/eth/v1/beacon/states/%s/committees?epoch=%d", stateID, epoch)
	respBodyReader, err := s.get(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request beacon committees")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain beacon committees")
	}

	var resp beaconCommitteesJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse beacon committees")
	}

	return resp.Data, nil
}
