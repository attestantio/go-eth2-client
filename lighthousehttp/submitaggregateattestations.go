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
	"io/ioutil"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SubmitAggregateAttestations submits aggregate attestations.
func (s *Service) SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*spec.SignedAggregateAndProof) error {
	specJSON, err := json.Marshal(aggregateAndProofs)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}
	lhReader, err := specToLH(ctx, bytes.NewReader(specJSON))
	if err != nil {
		return errors.Wrap(err, "failed to convert lighthouse response to spec response")
	}

	aggregateAttestations, err := ioutil.ReadAll(lhReader)
	if err != nil {
		return errors.Wrap(err, "failed to read lighthouse JSON")
	}

	respBodyReader, err := s.post(ctx, "/validator/aggregate_and_proofs", bytes.NewReader(aggregateAttestations))
	if err != nil {
		return errors.Wrap(err, "error submitting aggregate attestation")
	}

	resp, err := ioutil.ReadAll(respBodyReader)
	if err != nil {
		return errors.Wrap(err, "failed to read response")
	}
	if resp != nil && !bytes.Equal(resp, []byte("null")) {
		return fmt.Errorf("failed to submit aggregate attestation: %s", string(resp))
	}
	return nil
}
