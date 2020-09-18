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

package tekuhttp

import (
	"context"
	"encoding/json"
	"fmt"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// NonSpecAggregateAttestation fetches the aggregate attestation given an attestation.
func (s *Service) NonSpecAggregateAttestation(ctx context.Context, attestation *spec.Attestation, validatorPubKey []byte, slotSignature []byte) (*spec.Attestation, error) {
	root, err := attestation.Data.HashTreeRoot()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain hash tree root for attestation data")
	}
	respBodyReader, err := s.get(ctx, fmt.Sprintf("/validator/aggregate_attestation?attestation_data_root=%#x&slot=%d", root, attestation.Data.Slot))
	if err != nil {
		return nil, errors.Wrap(err, "failed to request aggregate attestation")
	}
	defer func() {
		if err := respBodyReader.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close HTTP body")
		}
	}()

	var aggregateAttestation *spec.Attestation
	if err := json.NewDecoder(respBodyReader).Decode(&aggregateAttestation); err != nil {
		return nil, errors.Wrap(err, "failed to parse aggregate attestation")
	}

	return aggregateAttestation, nil
}
