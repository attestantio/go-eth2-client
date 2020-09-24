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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// Teku has 'index' for aggregate attestations; spec is 'aggregator_index'.
var aggregateAttestationsRe1 = regexp.MustCompile(`"aggregator_index"`)

// Teku has 'attestation' for aggregate attestations; spec is 'aggregate'.
var aggregateAttestationsRe2 = regexp.MustCompile(`"aggregate"`)

// SubmitAggregateAttestations submits aggregate attestations.
func (s *Service) SubmitAggregateAttestations(ctx context.Context, aggregateAndProofs []*spec.SignedAggregateAndProof) error {
	for _, aggregateAndProof := range aggregateAndProofs {
		specJSON, err := json.Marshal(aggregateAndProof)
		if err != nil {
			return errors.Wrap(err, "failed to marshal JSON")
		}
		specJSON = aggregateAttestationsRe1.ReplaceAll(specJSON, []byte(`"index"`))
		specJSON = aggregateAttestationsRe2.ReplaceAll(specJSON, []byte(`"attestation"`))

		respBodyReader, cancel, err := s.post(ctx, "/validator/aggregate_and_proofs", bytes.NewReader(specJSON))
		if err != nil {
			return errors.Wrap(err, "error submitting aggregate attestation")
		}
		defer cancel()

		var resp []byte
		if respBodyReader != nil {
			resp, err = ioutil.ReadAll(respBodyReader)
			if err != nil {
				resp = nil
			}
		}
		if err != nil {
			return errors.Wrap(err, "failed to obtain error message for beacon block proposal")
		}
		if len(resp) > 0 {
			return fmt.Errorf("failed to submit beacon block proposal: %s", string(resp))
		}
	}

	return nil
}
