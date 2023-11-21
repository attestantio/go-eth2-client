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

// SigningData provides information about signing data.
type SigningData struct {
	ObjectRoot Root   `ssz-size:"32"`
	Domain     Domain `ssz-size:"32"`
}

// signingDataJSON is the spec representation of the struct.
type signingDataJSON struct {
	ObjectRoot string `json:"object_root"`
	Domain     string `json:"domain"`
}

// signingDataYAML is the spec representation of the struct.
type signingDataYAML struct {
	ObjectRoot string `yaml:"object_root"`
	Domain     string `yaml:"domain"`
}

// MarshalJSON implements json.Marshaler.
func (s *SigningData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&signingDataJSON{
		ObjectRoot: fmt.Sprintf("%#x", s.ObjectRoot),
		Domain:     fmt.Sprintf("%#x", s.Domain),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *SigningData) UnmarshalJSON(input []byte) error {
	var signingDataJSON signingDataJSON
	if err := json.Unmarshal(input, &signingDataJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return s.unpack(&signingDataJSON)
}

func (s *SigningData) unpack(signingDataJSON *signingDataJSON) error {
	if signingDataJSON.ObjectRoot == "" {
		return errors.New("object root missing")
	}
	objectRoot, err := hex.DecodeString(strings.TrimPrefix(signingDataJSON.ObjectRoot, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for object root")
	}
	if len(objectRoot) != RootLength {
		return errors.New("incorrect length for object root")
	}
	copy(s.ObjectRoot[:], objectRoot)
	if signingDataJSON.Domain == "" {
		return errors.New("domain missing")
	}
	domain, err := hex.DecodeString(strings.TrimPrefix(signingDataJSON.Domain, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for domain")
	}
	if len(domain) != DomainLength {
		return errors.New("incorrect length for domain")
	}
	copy(s.Domain[:], domain)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (s *SigningData) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&signingDataYAML{
		ObjectRoot: fmt.Sprintf("%#x", s.ObjectRoot),
		Domain:     fmt.Sprintf("%#x", s.Domain),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (s *SigningData) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var signingDataJSON signingDataJSON
	if err := yaml.Unmarshal(input, &signingDataJSON); err != nil {
		return err
	}

	return s.unpack(&signingDataJSON)
}

// String returns a string version of the structure.
func (s *SigningData) String() string {
	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
