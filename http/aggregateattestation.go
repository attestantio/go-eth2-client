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
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

type aggregateAttestationDataJSON struct {
	Data *phase0.Attestation `json:"data"`
}

// AggregateAttestation fetches the aggregate attestation given an attestation.
// N.B if an aggregate attestation for the attestation is not available this will return nil without an error.
func (s *Service) AggregateAttestation(ctx context.Context, slot phase0.Slot, attestationDataRoot phase0.Root) (*phase0.Attestation, error) {
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/eth/v1/validator/aggregate_attestation?slot=%d&attestation_data_root=%#x", slot, attestationDataRoot))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request aggregate attestation")
	}
	if respBodyReader == nil {
		return nil, nil
	}

	var aggregateAttestationDataJSON aggregateAttestationDataJSON
	if err := json.NewDecoder(respBodyReader).Decode(&aggregateAttestationDataJSON); err != nil {
		return nil, errors.Wrap(err, "failed to parse aggregate attestation")
	}

	// Ensure the data returned to us is as expected given our input.
	if aggregateAttestationDataJSON.Data.Data.Slot != slot {
		return nil, errors.New("aggregate attestation not for requested slot")
	}
	dataRoot, err := aggregateAttestationDataJSON.Data.Data.HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain hash tree root of aggregate attestation data")
	}
	if !bytes.Equal(dataRoot[:], attestationDataRoot[:]) {
		return nil, errors.New("aggregate attestation not for requested data root")
	}

	return aggregateAttestationDataJSON.Data, nil
}
