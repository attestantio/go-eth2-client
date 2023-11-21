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

// SignedBeaconBlock is a signed beacon block.
type SignedBeaconBlock struct {
	Message   *BeaconBlock
	Signature phase0.BLSSignature `ssz-size:"96"`
}

// signedBeaconBlockJSON is the spec representation of the struct.
type signedBeaconBlockJSON struct {
	Message   *BeaconBlock `json:"message"`
	Signature string       `json:"signature"`
}

// signedBeaconBlockYAML is the spec representation of the struct.
type signedBeaconBlockYAML struct {
	Message   *BeaconBlock `yaml:"message"`
	Signature string       `yaml:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedBeaconBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedBeaconBlockJSON{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedBeaconBlock) UnmarshalJSON(input []byte) error {
	var signedBeaconBlockJSON signedBeaconBlockJSON
	if err := json.Unmarshal(input, &signedBeaconBlockJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return s.unpack(&signedBeaconBlockJSON)
}

func (s *SignedBeaconBlock) unpack(signedBeaconBlockJSON *signedBeaconBlockJSON) error {
	if signedBeaconBlockJSON.Message == nil {
		return errors.New("message missing")
	}
	s.Message = signedBeaconBlockJSON.Message
	if signedBeaconBlockJSON.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(signedBeaconBlockJSON.Signature, "0x"))
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
func (s *SignedBeaconBlock) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&signedBeaconBlockYAML{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SignedBeaconBlock) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var signedBeaconBlockJSON signedBeaconBlockJSON
	if err := yaml.Unmarshal(input, &signedBeaconBlockJSON); err != nil {
		return err
	}

	return s.unpack(&signedBeaconBlockJSON)
}

// String returns a string version of the structure.
func (s *SignedBeaconBlock) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
