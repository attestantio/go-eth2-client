// Copyright © 2022 Attestant Limited.
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

package v1

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

// SignedValidatorRegistration is a signed ValidatorRegistrationV1.
type SignedValidatorRegistration struct {
	Message   *ValidatorRegistration
	Signature phase0.BLSSignature `ssz-size:"96"`
}

// signedValidatorRegistrationJSON is the spec representation of the struct.
type signedValidatorRegistrationJSON struct {
	Message   *ValidatorRegistration `json:"message"`
	Signature string                 `json:"signature"`
}

// signedValidatorRegistrationYAML is the spec representation of the struct.
type signedValidatorRegistrationYAML struct {
	Message   *ValidatorRegistration `yaml:"message"`
	Signature string                 `yaml:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedValidatorRegistration) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedValidatorRegistrationJSON{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedValidatorRegistration) UnmarshalJSON(input []byte) error {
	var data signedValidatorRegistrationJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return s.unpack(&data)
}

func (s *SignedValidatorRegistration) unpack(data *signedValidatorRegistrationJSON) error {
	if data.Message == nil {
		return errors.New("message missing")
	}
	s.Message = data.Message
	if data.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(data.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(signature) != phase0.SignatureLength {
		return fmt.Errorf("incorrect length %d for signature", len(signature))
	}
	copy(s.Signature[:], signature)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (s *SignedValidatorRegistration) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&signedValidatorRegistrationYAML{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SignedValidatorRegistration) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data signedValidatorRegistrationJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}

	return s.unpack(&data)
}

// String returns a string version of the structure.
func (s *SignedValidatorRegistration) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
