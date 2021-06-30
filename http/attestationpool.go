// Copyright © 2021 Attestant Limited.
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

type attestationPoolJSON struct {
	Data []*spec.Attestation `json:"data"`
}

// AttestationPool obtains the attestation pool for a given slot.
func (s *Service) AttestationPool(ctx context.Context, slot spec.Slot) ([]*spec.Attestation, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/beacon/pool/attestations?slot=%d", slot))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request attestation pool")
	}
	if respBodyReader == nil {
		return nil, errors.New("failed to obtain attestation pool")
	}

	var attestationPoolJSON attestationPoolJSON
	if err := json.NewDecoder(respBodyReader).Decode(&attestationPoolJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse attestation pool")
	}

	// Ensure the data returned to us is as expected given our input.
	if attestationPoolJSON.Data == nil {
		return nil, errors.New("attestation pool not returned")
	}
	for i := range attestationPoolJSON.Data {
		if attestationPoolJSON.Data[i].Data.Slot != slot {
			return nil, errors.New("attestation pool entry not for requested slot")
		}
	}

	return attestationPoolJSON.Data, nil
}
