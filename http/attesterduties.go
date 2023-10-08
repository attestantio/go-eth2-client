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

// AttesterDuties obtains attester duties.
func (s *Service) AttesterDuties(ctx context.Context,
	opts *api.AttesterDutiesOpts,
) (
	*api.Response[[]*apiv1.AttesterDuty],
	error,
) {
	if opts == nil {
		return nil, errors.New("no options specified")
	}
	if len(opts.Indices) == 0 {
		return nil, errors.New("no validator indices specified")
	}

	var reqBodyReader bytes.Buffer
	if _, err := reqBodyReader.WriteString(`[`); err != nil {
		return nil, errors.Wrap(err, "failed to write validator index array start")
	}
	for i := range opts.Indices {
		if _, err := reqBodyReader.WriteString(fmt.Sprintf(`"%d"`, opts.Indices[i])); err != nil {
			return nil, errors.Wrap(err, "failed to write index")
		}
		if i != len(opts.Indices)-1 {
			if _, err := reqBodyReader.WriteString(`,`); err != nil {
				return nil, errors.Wrap(err, "failed to write separator")
			}
		}
	}
	if _, err := reqBodyReader.WriteString(`]`); err != nil {
		return nil, errors.Wrap(err, "failed to write end of validator index array")
	}

	url := fmt.Sprintf("/eth/v1/validator/duties/attester/%d", opts.Epoch)
	respBodyReader, err := s.post(ctx, url, &reqBodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request attester duties")
	}

	data, metadata, err := decodeJSONResponse(respBodyReader, []*apiv1.AttesterDuty{})
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
			return nil, fmt.Errorf("received attester duty for slot %d outside of range [%d,%d]", duty.Slot, startSlot, endSlot)
		}
	}

	return &api.Response[[]*apiv1.AttesterDuty]{
		Metadata: metadata,
		Data:     data,
	}, nil
}
