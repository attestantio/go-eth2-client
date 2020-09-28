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

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SubmitAttestation submits an attestation.
func (s *Service) SubmitAttestation(ctx context.Context, specAttestation *spec.Attestation) error {
	specJSON, err := json.Marshal(specAttestation)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}

	respBodyReader, err := s.post(ctx, "/validator/attestation", bytes.NewReader(specJSON))
	if err != nil {
		return errors.Wrap(err, "failed to POST to /validator/attestation")
	}

	var resp []byte
	if respBodyReader != nil {
		resp, err = ioutil.ReadAll(respBodyReader)
		if err != nil {
			resp = nil
		}
	}
	if err != nil {
		return errors.Wrap(err, "failed to obtain error message for attestation")
	}
	if len(resp) > 0 {
		return fmt.Errorf("failed to submit attestation: %s", string(resp))
	}
	return nil
}
