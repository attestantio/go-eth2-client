// Copyright © 2020 - 2024 Attestant Limited.
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
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// AttesterDuties obtains attester duties.
func (s *Service) AttesterDuties(ctx context.Context,
	opts *api.AttesterDutiesOpts,
) (
	*api.Response[[]*apiv1.AttesterDuty],
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

	var reqBodyReader bytes.Buffer
	if _, err := reqBodyReader.WriteString(`[`); err != nil {
		return nil, errors.Join(errors.New("failed to write validator index array start"), err)
	}
	for i := range opts.Indices {
		if _, err := reqBodyReader.WriteString(fmt.Sprintf(`"%d"`, opts.Indices[i])); err != nil {
			return nil, errors.Join(errors.New("failed to write index"), err)
		}
		if i != len(opts.Indices)-1 {
			if _, err := reqBodyReader.WriteString(`,`); err != nil {
				return nil, errors.Join(errors.New("failed to write separator"), err)
			}
		}
	}
	if _, err := reqBodyReader.WriteString(`]`); err != nil {
		return nil, errors.Join(errors.New("failed to write end of validator index array"), err)
	}

	endpoint := fmt.Sprintf("/eth/v1/validator/duties/attester/%d", opts.Epoch)
	query := ""

	httpResponse, err := s.post(ctx,
		endpoint,
		query,
		&api.CommonOpts{},
		&reqBodyReader,
		ContentTypeJSON,
		map[string]string{},
	)
	if err != nil {
		return nil, errors.Join(errors.New("failed to request attester duties"), err)
	}

	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), []*apiv1.AttesterDuty{})
	if err != nil {
		return nil, err
	}

	// Confirm that duties are for the requested epoch.
	slotsPerEpoch, err := s.SlotsPerEpoch(ctx)
	if err != nil {
		return nil, errors.Join(errors.New("failed to obtain slots per epoch"), err)
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
