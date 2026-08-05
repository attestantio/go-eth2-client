// Copyright © 2026 Attestant Limited.
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

package gloas

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/pkg/errors"
)

// signedExecutionPayloadEnvelopeContentsJSON is the spec representation of the struct.
type signedExecutionPayloadEnvelopeContentsJSON struct {
	SignedExecutionPayloadEnvelope *gloas.SignedExecutionPayloadEnvelope `json:"signed_execution_payload_envelope"`
	KZGProofs                      []deneb.KZGProof                      `json:"kzg_proofs"`
	Blobs                          []deneb.Blob                          `json:"blobs"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedExecutionPayloadEnvelopeContents) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedExecutionPayloadEnvelopeContentsJSON{
		SignedExecutionPayloadEnvelope: s.SignedExecutionPayloadEnvelope,
		KZGProofs:                      s.KZGProofs,
		Blobs:                          s.Blobs,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedExecutionPayloadEnvelopeContents) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&signedExecutionPayloadEnvelopeContentsJSON{}, input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw["signed_execution_payload_envelope"], &s.SignedExecutionPayloadEnvelope); err != nil {
		return errors.Wrap(err, "signed_execution_payload_envelope")
	}

	if err := json.Unmarshal(raw["kzg_proofs"], &s.KZGProofs); err != nil {
		return errors.Wrap(err, "kzg_proofs")
	}

	if err := json.Unmarshal(raw["blobs"], &s.Blobs); err != nil {
		return errors.Wrap(err, "blobs")
	}

	return nil
}
