// Copyright © 2024 Attestant Limited.
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

package electra

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/codecs"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// signedConsolidationJSON is the spec representation of the struct.
type signedConsolidationJSON struct {
	Message   Consolidation       `json:"message"`
	Signature phase0.BLSSignature `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedConsolidation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedConsolidationJSON{
		Message:   s.Message,
		Signature: s.Signature,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedConsolidation) UnmarshalJSON(input []byte) error {
	raw, err := codecs.RawJSON(&signedConsolidationJSON{}, input)
	if err != nil {
		return err
	}

	if err := s.Message.UnmarshalJSON(raw["message"]); err != nil {
		return errors.Wrap(err, "message")
	}
	if err := s.Signature.UnmarshalJSON(raw["signature"]); err != nil {
		return errors.Wrap(err, "signature")
	}

	return nil
}