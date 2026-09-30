// Copyright © 2020 - 2026 Attestant Limited.
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
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// AttestationData obtains attestation data given the options.
func (s *Service) AttestationData(ctx context.Context,
	opts *api.AttestationDataOpts,
) (
	*api.Response[*phase0.AttestationData],
	error,
) {
	if err := s.assertIsSynced(ctx); err != nil {
		return nil, err
	}

	if opts == nil {
		return nil, client.ErrNoOptions
	}

	endpoint := "/eth/v1/validator/attestation_data"
	query := fmt.Sprintf("slot=%d&committee_index=%d", opts.Slot, opts.CommitteeIndex)

	httpResponse, err := s.get(ctx, endpoint, query, &opts.Common, false)
	if err != nil {
		return nil, err
	}

	switch httpResponse.contentType {
	case ContentTypeJSON:
		return s.attestationDataFromJSON(ctx, opts, httpResponse)
	default:
		return nil, fmt.Errorf("unhandled content type %v", httpResponse.contentType)
	}
}

func (s *Service) attestationDataFromJSON(ctx context.Context,
	opts *api.AttestationDataOpts,
	httpResponse *httpResponse,
) (
	*api.Response[*phase0.AttestationData],
	error,
) {
	data, metadata, err := decodeJSONResponse(bytes.NewReader(httpResponse.body), phase0.AttestationData{})
	if err != nil {
		return nil, err
	}

	if err := s.verifyAttestationData(ctx, opts, &data); err != nil {
		return nil, err
	}

	return &api.Response[*phase0.AttestationData]{
		Metadata: metadata,
		Data:     &data,
	}, nil
}

func (s *Service) verifyAttestationData(ctx context.Context, opts *api.AttestationDataOpts, data *phase0.AttestationData) error {
	if data.Slot != opts.Slot {
		return errors.Join(
			fmt.Errorf("attestation data for slot %d; expected %d", data.Slot, opts.Slot),
			client.ErrInconsistentResult,
		)
	}

	onGloas, err := s.isGloasSlot(ctx, opts.Slot)
	if err != nil {
		return errors.Join(errors.New("failed to determine whether the slot is in the gloas era"), err)
	}

	// Gloas repurposes data.Index as a one-bit vote on the availability of the attested
	// block's execution payload: 0 for a payload the attester does not see in the
	// canonical chain, 1 for one it does.  Holding it to the electra-era value of 0 would
	// discard every payload-FULL answer a conformant node gives, so bound it the way
	// process_attestation does rather than pinning it.
	if onGloas {
		if data.Index > 1 {
			return errors.Join(
				fmt.Errorf("attestation data payload availability vote %d; expected 0 or 1", data.Index),
				client.ErrInconsistentResult,
			)
		}

		return nil
	}

	electraSlot, err := s.calculateElectraSlot(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to calculate electra slot"), err)
	}

	// When in the electra era the data.Index is hardcoded to 0.
	index := opts.CommitteeIndex
	if opts.Slot >= electraSlot {
		index = 0
	}

	if data.Index != index {
		return errors.Join(
			fmt.Errorf("attestation data for committee index %d; expected %d", data.Index, index),
			client.ErrInconsistentResult,
		)
	}

	return nil
}

// isGloasSlot reports whether the given slot is at or after the node's gloas fork.
//
// A node that does not know the fork answers false rather than failing, unlike its
// electra counterpart's treatment of a missing key: a client publishes GLOAS_FORK_EPOCH
// only once it has one, so on every node that predates the fork the key is simply
// missing, and reporting that as a failure would break the endpoint everywhere.
//
// A predicate on a slot rather than the fork's first slot, because the two facts that
// answer would rest on -- the epoch, and whether there is one at all -- are only
// meaningful together.  A caller handed both separately can compare against the slot
// without consulting the flag, and an unknown fork has no first slot to return, so that
// mistake would admit every slot rather than none.
func (s *Service) isGloasSlot(ctx context.Context, slot phase0.Slot) (bool, error) {
	response, err := s.Spec(ctx, &api.SpecOpts{})
	if err != nil {
		return false, err
	}

	value, exists := response.Data["GLOAS_FORK_EPOCH"]
	if !exists {
		return false, nil
	}

	gloasEpoch, isCorrectType := value.(uint64)
	if !isCorrectType {
		return false, ErrIncorrectType
	}

	slotsPerEpoch, isCorrectType := response.Data["SLOTS_PER_EPOCH"].(uint64)
	if !isCorrectType {
		return false, ErrIncorrectType
	}

	return slot >= phase0.Slot(slotsPerEpoch*gloasEpoch), nil
}

func (s *Service) calculateElectraSlot(ctx context.Context) (phase0.Slot, error) {
	response, err := s.Spec(ctx, &api.SpecOpts{})
	if err != nil {
		return 0, err
	}

	slotsPerEpoch, isCorrectType := response.Data["SLOTS_PER_EPOCH"].(uint64)
	if !isCorrectType {
		return 0, ErrIncorrectType
	}

	electraEpoch, isCorrectType := response.Data["ELECTRA_FORK_EPOCH"].(uint64)
	if !isCorrectType {
		return 0, ErrIncorrectType
	}

	electraSlot := phase0.Slot(slotsPerEpoch * electraEpoch)

	return electraSlot, nil
}
