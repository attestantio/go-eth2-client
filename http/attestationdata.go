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

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type attestationDataJSON struct {
	Data *spec.AttestationData `json:"data"`
}

// AttestationData obtains attestation data for a slot.
func (s *Service) AttestationData(ctx context.Context, slot spec.Slot, committeeIndex spec.CommitteeIndex) (*spec.AttestationData, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/validator/attestation_data?slot=%d&committee_index=%d", slot, committeeIndex))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request attestation data")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain attestation data")
	}

	var attestationDataJSON attestationDataJSON
	if err := json.NewDecoder(respBodyReader).Decode(&attestationDataJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse attestation data")
	}

	// Ensure the data returned to us is as expected given our input.
	if attestationDataJSON.Data == nil {
		return nil, errors.New("attestation not returned")
	}
	if attestationDataJSON.Data.Slot != slot {
		return nil, errors.New("attestation data not for requested slot")
	}
	if attestationDataJSON.Data.Index != committeeIndex {
		return nil, errors.New("attestation data not for requested committee index")
	}

	return attestationDataJSON.Data, nil
}
