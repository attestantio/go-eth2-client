// Copyright © 2020 - 2023 Attestant Limited.
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
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// ProposerDuties obtains proposer duties for the given options.
func (s *Service) ProposerDuties(ctx context.Context,
	opts *api.ProposerDutiesOpts,
) (
	*api.Response[[]*apiv1.ProposerDuty],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}

	url := fmt.Sprintf("/eth/v1/validator/duties/proposer/%d", opts.Epoch)
	httpResponse, err := s.get(ctx, url, &opts.Common)
	if err != nil {
		return nil, err
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.ProposerDuty{})
	if err != nil {
		return nil, err
	}

	// Confirm that duties are for the requested epoch.
	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain slots per epoch")
	}
	startSlot := phase0.Slot(uint64(opts.Epoch) * slotsPerEpoch)
	endSlot := phase0.Slot(uint64(opts.Epoch)*slotsPerEpoch + slotsPerEpoch - 1)
	for _, duty := range data {
		if duty.Slot < startSlot || duty.Slot > endSlot {
			return nil, fmt.Errorf("received proposer duty for slot %d outside of range [%d,%d]", duty.Slot, startSlot, endSlot)
		}
	}

	if len(opts.Indices) > 0 {
		// Filter duties based on supplied validators.
		validatorIndexMap := make(map[phase0.ValidatorIndex]bool, len(opts.Indices))
		for _, index := range opts.Indices {
			validatorIndexMap[index] = true
		}

		filteredData := make([]*apiv1.ProposerDuty, 0, len(data))
		for _, duty := range data {
			if _, exists := validatorIndexMap[duty.ValidatorIndex]; exists {
				filteredData = append(filteredData, duty)
			}
		}
		data = filteredData
	}

	return &api.Response[[]*apiv1.ProposerDuty]{
		Metadata: metadata,
		Data:     data,
	}, nil
}
