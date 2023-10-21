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

// SignedVoluntaryExit provides information about a signed voluntary exit.
type SignedVoluntaryExit struct {
	Message   *VoluntaryExit
	Signature BLSSignature `ssz-size:"96"`
}

// signedVoluntaryExitJSON is the spec representation of the struct.
type signedVoluntaryExitJSON struct {
	Message   *VoluntaryExit `json:"message"`
	Signature string         `json:"signature"`
}

// signedVoluntaryExitYAML is the spec representation of the struct.
type signedVoluntaryExitYAML struct {
	Message   *VoluntaryExit `yaml:"message"`
	Signature string         `yaml:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedVoluntaryExit) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedVoluntaryExitJSON{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedVoluntaryExit) UnmarshalJSON(input []byte) error {
	var signedVoluntaryExitJSON signedVoluntaryExitJSON
	if err := json.Unmarshal(input, &signedVoluntaryExitJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return s.unpack(&signedVoluntaryExitJSON)
}

func (s *SignedVoluntaryExit) unpack(signedVoluntaryExitJSON *signedVoluntaryExitJSON) error {
	s.Message = signedVoluntaryExitJSON.Message
	if s.Message == nil {
		return errors.New("message missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(signedVoluntaryExitJSON.Signature, "0x"))
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
func (s *SignedVoluntaryExit) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&signedVoluntaryExitYAML{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SignedVoluntaryExit) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var signedVoluntaryExitJSON signedVoluntaryExitJSON
	if err := yaml.Unmarshal(input, &signedVoluntaryExitJSON); err != nil {
		return err
	}

	return s.unpack(&signedVoluntaryExitJSON)
}

// String returns a string version of the structure.
func (s *SignedVoluntaryExit) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
