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
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// SignedAggregateAndProof provides information about a signed aggregate and proof.
type SignedAggregateAndProof struct {
	Message   *AggregateAndProof
	Signature BLSSignature `ssz-size:"96"`
}

// signedAggregateAndProofJSON is a raw representation of the struct.
type signedAggregateAndProofJSON struct {
	Message   *AggregateAndProof `json:"message"`
	Signature string             `json:"signature"`
}

// signedAggregateAndProofYAML is a raw representation of the struct.
type signedAggregateAndProofYAML struct {
	Message   *AggregateAndProof `yaml:"message"`
	Signature string             `yaml:"signature"`
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
	var signedAggregateAndProofJSON signedAggregateAndProofJSON
	if err := json.Unmarshal(input, &signedAggregateAndProofJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return s.unpack(&signedAggregateAndProofJSON)
}

func (s *SignedAggregateAndProof) unpack(signedAggregateAndProofJSON *signedAggregateAndProofJSON) error {
	if signedAggregateAndProofJSON.Message == nil {
		return errors.New("message missing")
	}
	s.Message = signedAggregateAndProofJSON.Message
	if signedAggregateAndProofJSON.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(signedAggregateAndProofJSON.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(signature) != SignatureLength {
		return errors.New("incorrect length for signature")
	}
	copy(s.Signature[:], signature)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (s *SignedAggregateAndProof) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&signedAggregateAndProofYAML{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SignedAggregateAndProof) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var signedAggregateAndProofJSON signedAggregateAndProofJSON
	if err := yaml.Unmarshal(input, &signedAggregateAndProofJSON); err != nil {
		return err
	}

	return s.unpack(&signedAggregateAndProofJSON)
}

// String returns a string version of the structure.
func (s *SignedAggregateAndProof) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
