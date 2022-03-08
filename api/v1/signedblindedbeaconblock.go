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

// SignedBlindedBeaconBlock is a signed beacon block.
type SignedBlindedBeaconBlock struct {
	Message   *BlindedBeaconBlock
	Signature phase0.BLSSignature `ssz-size:"96"`
}

// signedBlindedBeaconBlockJSON is the spec representation of the struct.
type signedBlindedBeaconBlockJSON struct {
	Message   *BlindedBeaconBlock `json:"message"`
	Signature string              `json:"signature"`
}

// signedBlindedBeaconBlockYAML is the spec representation of the struct.
type signedBlindedBeaconBlockYAML struct {
	Message   *BlindedBeaconBlock `yaml:"message"`
	Signature string              `yaml:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (s *SignedBlindedBeaconBlock) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signedBlindedBeaconBlockJSON{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SignedBlindedBeaconBlock) UnmarshalJSON(input []byte) error {
	var data signedBlindedBeaconBlockJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	return s.unpack(&data)
}

func (s *SignedBlindedBeaconBlock) unpack(data *signedBlindedBeaconBlockJSON) error {
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
func (s *SignedBlindedBeaconBlock) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&signedBlindedBeaconBlockYAML{
		Message:   s.Message,
		Signature: fmt.Sprintf("%#x", s.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}
	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SignedBlindedBeaconBlock) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var data signedBlindedBeaconBlockJSON
	if err := yaml.Unmarshal(input, &data); err != nil {
		return err
	}
	return s.unpack(&data)
}

// String returns a string version of the structure.
func (s *SignedBlindedBeaconBlock) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
