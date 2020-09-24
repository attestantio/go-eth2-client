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
	"regexp"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SubmitAttestation submits an attestation.
func (s *Service) SubmitAttestation(ctx context.Context, specAttestation *spec.Attestation) error {
	specJSON, err := json.Marshal(specAttestation)
	if err != nil {
		return errors.Wrap(err, "failed to marshal JSON")
	}
	lhReader, err := specToLH(ctx, bytes.NewReader(specJSON))
	if err != nil {
		return errors.Wrap(err, "failed to convert spec JSON to lighthouse JSON")
	}

	attestation, err := ioutil.ReadAll(lhReader)
	if err != nil {
		return errors.Wrap(err, "failed to read lighthouse JSON")
	}

	// We require a subnet ID, which we don't have.  This has been raised at
	// https://github.com/sigp/lighthouse/issues/1550
	// In the meantime, the helpful error message tells us which subnet it should be so we:
	// - try with subnet 0
	// - if it succeeds, done
	// - if not, fetch the requested subnet ID from the error message and try again.
	attestations := fmt.Sprintf("[[%s,0]]", string(attestation))
	respBodyReader, cancel, err := s.post(ctx, "/validator/attestations", bytes.NewReader([]byte(attestations)))
	if err != nil {
		return errors.Wrap(err, "failed to post attestation for 0-subnet attestation")
	}
	defer cancel()

	resp, err := ioutil.ReadAll(respBodyReader)
	if err != nil {
		return errors.Wrap(err, "failed to obtain error message for 0-subnet attestation")
	}
	if resp == nil || bytes.Equal(resp, []byte("null")) {
		// Means we succeeded.
		return nil
	}
	re := regexp.MustCompile(`.*expected: SubnetId\(([0-9]+)\)`)
	match := re.FindSubmatch(resp)
	if len(match) != 2 {
		log.Warn().Str("resp", string(resp)).Msg("No subnet ID supplied in error message")
		return errors.New("no subnet ID supplied in error message")
	}
	//subnetID, err := strconv.ParseInt(string(match[1]),10,64)
	//if err != nil {
	//		return errors.New("invalid subnet ID supplied in error message")
	//}
	// Go again with the "borrowed" subnet ID.
	attestations = fmt.Sprintf("[[%s,%s]]", string(attestation), string(match[1]))
	respBodyReader, cancel, err = s.post(ctx, "/validator/attestations", bytes.NewReader([]byte(attestations)))
	if err != nil {
		return errors.Wrap(err, "failed to POST to /validator/attestations")
	}
	defer cancel()

	resp, err = ioutil.ReadAll(respBodyReader)
	if err != nil {
		return errors.Wrap(err, "failed to obtain error message for attestation")
	}
	if resp != nil && !bytes.Equal(resp, []byte("null")) {
		return fmt.Errorf("failed to submit attestation: %s", string(resp))
	}
	return nil

	// This is for when we no longer need to send the subnet ID.
	//	log.Trace().Msg("Sending to /validator/attestations")
	//	respBodyReader, err := s.post(ctx, "/validator/attestations", bytes.NewReader([]byte(attestations)))
	//	var resp []byte
	//	if respBodyReader != nil {
	//		resp, err = ioutil.ReadAll(respBodyReader)
	//		if err != nil {
	//			resp = nil
	//		}
	//	}
	//
	//	if err != nil {
	//		log.Debug().Err(err).Str("resp", string(resp)).Msg("POST to /validator/attestations failed")
	//		return errors.Wrap(err, fmt.Sprintf("failed to submit attestation: %s", resp))
	//	}
	//
	// 	return nil
}
