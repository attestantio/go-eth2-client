// Copyright © 2023 Attestant Limited.
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

package deneb

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/pkg/errors"
)

// signedBlindedBlobSidecarJSON is the spec representation of the struct.
type signedBlindedBlobSidecarJSON struct {
	Message   *BlindedBlobSidecar `json:"message"`
	Signature string              `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedBlindedBlobSidecar) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedBlindedBlobSidecarJSON{
		Message:   s.Message,
		Signature: s.Signature.String(),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedBlindedBlobSidecar) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&signedBlindedBlobSidecarJSON{}, input)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw["message"], &s.Message); err != nil {
		return errors.Wrap(err, "message")
	}

	if err := json.Unmarshal(raw["signature"], &s.Signature); err != nil {
		return errors.Wrap(err, "signature")
	}

	return nil
}
