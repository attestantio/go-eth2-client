// Copyright © 2021 Attestant Limited.
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

package altair

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// SignedContributionAndProof provides information about a signed contribution and proof.
type SignedContributionAndProof struct {
	Message   *ContributionAndProof
	Signature phase0.BLSSignature `ssz-size:"96"`
}

// signedContributionAndProofJSON is a raw representation of the struct.
type signedContributionAndProofJSON struct {
	Message   *ContributionAndProof `json:"message"`
	Signature string                `json:"signature"`
}

// signedContributionAndProofYAML is a raw representation of the struct.
type signedContributionAndProofYAML struct {
	Message   *ContributionAndProof `yaml:"message"`
	Signature string                `yaml:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedContributionAndProof) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedContributionAndProofJSON{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedContributionAndProof) UnmarshalJSON(input []byte) error {
	var signedContributionAndProofJSON signedContributionAndProofJSON
	if err := json.Unmarshal(input, &signedContributionAndProofJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return s.unpack(&signedContributionAndProofJSON)
}

func (s *SignedContributionAndProof) unpack(signedContributionAndProofJSON *signedContributionAndProofJSON) error {
	if signedContributionAndProofJSON.Message == nil {
		return errors.New("message missing")
	}
	s.Message = signedContributionAndProofJSON.Message
	if signedContributionAndProofJSON.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(signedContributionAndProofJSON.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(signature) != phase0.SignatureLength {
		return errors.New("incorrect length for signature")
	}
	copy(s.Signature[:], signature)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (s *SignedContributionAndProof) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&signedContributionAndProofYAML{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SignedContributionAndProof) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var signedContributionAndProofJSON signedContributionAndProofJSON
	if err := yaml.Unmarshal(input, &signedContributionAndProofJSON); err != nil {
		return err
	}

	return s.unpack(&signedContributionAndProofJSON)
}

// String returns a string version of the structure.
func (s *SignedContributionAndProof) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
