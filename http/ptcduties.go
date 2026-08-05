// Copyright © 2026 Attestant Limited.
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
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// PTCDuties obtains payload timeliness committee duties.
func (s *Service) PTCDuties(ctx context.Context,
	opts *api.PTCDutiesOpts,
) (
	*api.Response[[]*apiv1.PTCDuty],
	error,
) {
	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	if len(opts.Indices) == 0 {
		return nil, errors.Join(errors.New("no validator indices specified"), client.ErrInvalidOptions)
	}

	// The endpoint takes an array of validator indices, which the beacon API
	// encodes as quoted decimal strings.
	indices := make([]string, len(opts.Indices))
	for i := range opts.Indices {
		indices[i] = strconv.FormatUint(uint64(opts.Indices[i]), 10)
	}

	reqBody, err := json.Marshal(indices)
	if err != nil {
		return nil, errors.Join(errors.New("failed to marshal validator indices"), err)
	}

	endpoint := fmt.Sprintf("/eth/v1/validator/duties/ptc/%d", opts.Epoch)

	httpResponse, err := s.post(ctx,
		endpoint,
		"",
		&opts.Common,
		bytes.NewReader(reqBody),
		ContentTypeJSON,
		map[string]string{},
	)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request PTC duties"), err)
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.PTCDuty{})
	if err != nil {
		return nil, err
	}

	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return nil, errors.Join(errors.New("failed to obtain slots per epoch"), err)
	}

	if err := verifyPTCDuties(opts.Epoch, slotsPerEpoch, data); err != nil {
		return nil, err
	}

	return &api.Response[[]*apiv1.PTCDuty]{
		Metadata: metadata,
		Data:     data,
	}, nil
}

// verifyPTCDuties confirms that every duty returned by the node falls in the
// epoch that was requested.
func verifyPTCDuties(epoch phase0.Epoch, slotsPerEpoch uint64, duties []*apiv1.PTCDuty) error {
	if slotsPerEpoch == 0 {
		// Without this the end of the range underflows to MaxUint64 and the
		// check below accepts every duty the node cares to return.
		return errors.New("invalid slots per epoch 0")
	}

	startSlot := phase0.Slot(uint64(epoch) * slotsPerEpoch)
	endSlot := startSlot + phase0.Slot(slotsPerEpoch) - 1

	for _, duty := range duties {
		if duty == nil {
			return errors.New("received nil PTC duty")
		}

		if duty.Slot < startSlot || duty.Slot > endSlot {
			return fmt.Errorf("received PTC duty for slot %d outside of range [%d,%d]", duty.Slot, startSlot, endSlot)
		}
	}

	return nil
}
