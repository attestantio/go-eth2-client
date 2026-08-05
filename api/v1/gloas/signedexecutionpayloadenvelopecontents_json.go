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
	"github.com/pkg/errors"
)

// MarshalJSON implements json.Marshaler.
func (s *SignedExecutionPayloadEnvelopeContents) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedExecutionPayloadEnvelopeContents{
		SignedExecutionPayloadEnvelope: s.SignedExecutionPayloadEnvelope,
		KZGProofs:                      s.KZGProofs,
		Blobs:                          s.Blobs,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedExecutionPayloadEnvelopeContents) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&signedExecutionPayloadEnvelopeContents{}, input)
	if err != nil {
		return err
	}

	var unmarshaled signedExecutionPayloadEnvelopeContents
	if err := json.Unmarshal(raw["signed_execution_payload_envelope"], &unmarshaled.SignedExecutionPayloadEnvelope); err != nil {
		return errors.Wrap(err, "signed_execution_payload_envelope")
	}

	if err := json.Unmarshal(raw["kzg_proofs"], &unmarshaled.KZGProofs); err != nil {
		return errors.Wrap(err, "kzg_proofs")
	}

	if err := json.Unmarshal(raw["blobs"], &unmarshaled.Blobs); err != nil {
		return errors.Wrap(err, "blobs")
	}

	s.unpack(&unmarshaled)

	return nil
}
