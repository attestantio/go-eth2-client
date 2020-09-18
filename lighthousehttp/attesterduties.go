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
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	client "github.com/attestantio/go-eth2-client"
	api "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/pkg/errors"
)

// AttesterDuties obtains attester duties.
func (s *Service) AttesterDuties(ctx context.Context, epoch uint64, validators []client.ValidatorIDProvider) ([]*api.AttesterDuty, error) {
	if validators == nil {
		// Best handled by a different call.
		committees, err := s.BeaconCommittees(ctx, fmt.Sprintf("%d", epoch*(*s.slotsPerEpoch)))
		if err != nil {
			return nil, errors.Wrap(err, "failed to obtain committees")
		}
		res := make([]*api.AttesterDuty, 0)
		for committeeIndex, committee := range committees {
			for index, entry := range committee.Validators {
				res = append(res, &api.AttesterDuty{
					Slot:                    committee.Slot,
					CommitteeIndex:          uint64(committeeIndex),
					CommitteeLength:         uint64(len(committee.Validators)),
					ValidatorIndex:          entry,
					ValidatorCommitteeIndex: uint64(index),
				})
			}
		}
		return res, nil
	}
	var reqBodyReader bytes.Buffer
	if _, err := reqBodyReader.WriteString(`{"epoch":`); err != nil {
		return nil, errors.Wrap(err, "failed to write header")
	}
	if _, err := reqBodyReader.WriteString(fmt.Sprintf("%d", epoch)); err != nil {
		return nil, errors.Wrap(err, "failed to write epoch")
	}
	if _, err := reqBodyReader.WriteString(`,"pubkeys":[`); err != nil {
		return nil, errors.Wrap(err, "failed to write pubkeys")
	}
	for i := range validators {
		pubKey, err := validators[i].PubKey(ctx)
		if err != nil {
			// Warn but continue.
			log.Warn().Err(err).Msg("Failed to obtain public key for validator; skipping")
			continue
		}
		if _, err := reqBodyReader.WriteString(fmt.Sprintf(`"%#x"`, pubKey)); err != nil {
			return nil, errors.Wrap(err, "failed to write pubkey")
		}
		if i != len(validators)-1 {
			if _, err := reqBodyReader.WriteString(`,`); err != nil {
				return nil, errors.Wrap(err, "failed to write separator")
			}
		}
	}
	if _, err := reqBodyReader.WriteString(`]}`); err != nil {
		return nil, errors.Wrap(err, "failed to write footer")
	}

	respBodyReader, err := s.post(ctx, "/validator/duties", &reqBodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request attester duties")
	}
	defer func() {
		if err := respBodyReader.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close HTTP body")
		}
	}()

	var resp []*dutyJSON
	if err := json.NewDecoder(respBodyReader).Decode(&resp); err != nil {
		return nil, errors.Wrap(err, "failed to parse attester duties response")
	}

	attesterDuties := make([]*api.AttesterDuty, 0, len(resp))
	for i := range resp {
		attesterDuties = append(attesterDuties, &api.AttesterDuty{
			Slot:                    resp[i].AttestationSlot,
			ValidatorIndex:          resp[i].ValidatorIndex,
			CommitteeIndex:          resp[i].AttestationCommitteeIndex,
			ValidatorCommitteeIndex: resp[i].AttestationCommitteePosition,
		})
	}

	// Need to obtain the committee size; comes from a different call.
	// Fetch the data.
	respBodyReader, err = s.get(ctx, fmt.Sprintf("/beacon/committees?epoch=%d", epoch))
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain committees")
	}
	defer func() {
		if err := respBodyReader.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close HTTP body")
		}
	}()
	type committeeData struct {
		Slot      uint64   `json:"slot"`
		Index     uint64   `json:"index"`
		Committee []uint64 `json:"committee"`
	}
	var committees []*committeeData
	if err := json.NewDecoder(respBodyReader).Decode(&committees); err != nil {
		return nil, errors.Wrap(err, "failed to parse committees")
	}
	// Parse the data.
	// Map is slot -> committee -> size.
	committeeSizes := make(map[uint64]map[uint64]int)
	for _, committee := range committees {
		slotSizes, exists := committeeSizes[committee.Slot]
		if !exists {
			slotSizes = make(map[uint64]int)
			committeeSizes[committee.Slot] = slotSizes
		}
		slotSizes[committee.Index] = len(committee.Committee)
	}
	// Integrate the committee sizes.
	for _, duty := range attesterDuties {
		duty.CommitteeLength = uint64(committeeSizes[duty.Slot][duty.CommitteeIndex])
	}

	return attesterDuties, nil
}

// dutyJSON handles the JSON returned from lighthouse.
type dutyJSON struct {
	PubKey                       string   `json:"validator_pubkey"`
	ValidatorIndex               uint64   `json:"validator_index"`
	AttestationSlot              uint64   `json:"attestation_slot"`
	AttestationCommitteeIndex    uint64   `json:"attestation_committee_index"`
	AttestationCommitteePosition uint64   `json:"attestation_committee_position"`
	BlockProposalSlots           []uint64 `json:"block_proposal_slots"`
	AggregatorModulo             uint64   `json:"aggregator_modulo"`
}
