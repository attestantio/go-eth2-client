// Copyright © 2020, 2021 Attestant Limited.
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

type proposerDutiesJSON struct {
	Data []*api.ProposerDuty `json:"data"`
}

// ProposerDuties obtains proposer duties for the given epoch.
// If validators is empty all duties are returned, otherwise only matching duties are returned.
func (s *Service) ProposerDuties(ctx context.Context, epoch phase0.Epoch, validatorIndices []phase0.ValidatorIndex) ([]*api.ProposerDuty, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/validator/duties/proposer/%d", epoch))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request proposer duties")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain proposer duties")
	}

	var resp proposerDutiesJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse proposer duties response")
	}

	// Validate the duties.
	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain slots per epoch")
	}
	startSlot := phase0.Slot(uint64(epoch) * slotsPerEpoch)
	endSlot := phase0.Slot(uint64(epoch)*slotsPerEpoch + slotsPerEpoch - 1)
	for _, duty := range resp.Data {
		if duty.Slot < startSlot || duty.Slot > endSlot {
			return nil, fmt.Errorf("received proposal for slot %d outside of range [%d,%d]", duty.Slot, startSlot, endSlot)
		}
	}

	if len(validatorIndices) == 0 {
		// Return all duties.
		return resp.Data, nil
	}

	// Filter duties based on supplied validators.
	validatorIndexMap := make(map[phase0.ValidatorIndex]bool, len(validatorIndices))
	for _, index := range validatorIndices {
		validatorIndexMap[index] = true
	}
	duties := make([]*api.ProposerDuty, 0, len(resp.Data))
	for _, duty := range resp.Data {
		if _, exists := validatorIndexMap[duty.ValidatorIndex]; exists {
			duties = append(duties, duty)
		}
	}

	return duties, nil
}
