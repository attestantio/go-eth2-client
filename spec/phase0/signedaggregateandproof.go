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

package phase0

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// SignedAggregateAndProof provides information about a signed aggregate and proof.
type SignedAggregateAndProof struct {
	Message   *AggregateAndProof
	Signature []byte `ssz-size:"96"`
}

// signedAggregateAndProofJSON is the spec representation of the struct.
type signedAggregateAndProofJSON struct {
	Message   *AggregateAndProof `json:"message"`
	Signature string             `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedAggregateAndProof) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedAggregateAndProofJSON{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedAggregateAndProof) UnmarshalJSON(input []byte) error {
	var err error

	var signedAggregateAndProofJSON signedAggregateAndProofJSON
	if err = json.Unmarshal(input, &signedAggregateAndProofJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if signedAggregateAndProofJSON.Message == nil {
		return errors.New("message missing")
	}
	s.Message = signedAggregateAndProofJSON.Message
	if signedAggregateAndProofJSON.Signature == "" {
		return errors.New("signature missing")
	}
	if s.Signature, err = hex.DecodeString(strings.TrimPrefix(signedAggregateAndProofJSON.Signature, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(s.Signature) != signatureLength {
		return errors.New("incorrect length for signature")
	}

	return nil
}

// String returns a string version of the structure.
func (s *SignedAggregateAndProof) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
