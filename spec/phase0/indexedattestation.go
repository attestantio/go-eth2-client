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
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
)

// IndexedAttestation provides a signed attestation with a list of attesting indices.
type IndexedAttestation struct {
	// Currently using primitives as sszgen does not handle []ValidatorIndex
	AttestingIndices []uint64 `ssz-max:"2048"`
	Data             *AttestationData
	Signature        BLSSignature `ssz-size:"96"`
}

// indexedAttestationJSON is the spec representation of the struct.
type indexedAttestationJSON struct {
	AttestingIndices []string         `json:"attesting_indices"`
	Data             *AttestationData `json:"data"`
	Signature        string           `json:"signature"`
}

// indexedAttestationYAML is a raw representation of the struct.
type indexedAttestationYAML struct {
	AttestingIndices []uint64         `yaml:"attesting_indices"`
	Data             *AttestationData `yaml:"data"`
	Signature        string           `yaml:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (i *IndexedAttestation) MarshalJSON() ([]byte, error) {
	attestingIndices := make([]string, len(i.AttestingIndices))
	for j := range i.AttestingIndices {
		attestingIndices[j] = strconv.FormatUint(i.AttestingIndices[j], 10)
	}

	return json.Marshal(&indexedAttestationJSON{
		AttestingIndices: attestingIndices,
		Data:             i.Data,
		Signature:        fmt.Sprintf("%#x", i.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *IndexedAttestation) UnmarshalJSON(input []byte) error {
	var indexedAttestationJSON indexedAttestationJSON
	if err := json.Unmarshal(input, &indexedAttestationJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	return i.unpack(&indexedAttestationJSON)
}

func (i *IndexedAttestation) unpack(indexedAttestationJSON *indexedAttestationJSON) error {
	var err error
	// Spec tests contain indexed attestations with empty attesting indices.
	// if indexedAttestationJSON.AttestingIndices == nil {
	// 	return errors.New("attesting indices missing")
	// }
	// if len(indexedAttestationJSON.AttestingIndices) == 0 {
	// 	return errors.New("attesting indices missing")
	// }
	i.AttestingIndices = make([]uint64, len(indexedAttestationJSON.AttestingIndices))
	for j := range indexedAttestationJSON.AttestingIndices {
		if i.AttestingIndices[j], err = strconv.ParseUint(indexedAttestationJSON.AttestingIndices[j], 10, 64); err != nil {
			return errors.Wrap(err, "failed to parse attesting index")
		}
	}
	if indexedAttestationJSON.Data == nil {
		return errors.New("data missing")
	}
	i.Data = indexedAttestationJSON.Data
	if indexedAttestationJSON.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(indexedAttestationJSON.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(signature) != SignatureLength {
		return errors.New("incorrect length for signature")
	}
	copy(i.Signature[:], signature)

	return nil
}

// MarshalYAML implements yaml.Marshaler.
func (i *IndexedAttestation) MarshalYAML() ([]byte, error) {
	yamlBytes, err := yaml.MarshalWithOptions(&indexedAttestationYAML{
		AttestingIndices: i.AttestingIndices,
		Data:             i.Data,
		Signature:        fmt.Sprintf("%#x", i.Signature),
	}, yaml.Flow(true))
	if err != nil {
		return nil, err
	}

	return bytes.ReplaceAll(yamlBytes, []byte(`"`), []byte(`'`)), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (i *IndexedAttestation) UnmarshalYAML(input []byte) error {
	// We unmarshal to the JSON struct to save on duplicate code.
	var indexedAttestationJSON indexedAttestationJSON
	if err := yaml.Unmarshal(input, &indexedAttestationJSON); err != nil {
		return err
	}

	return i.unpack(&indexedAttestationJSON)
}

// String returns a string version of the structure.
func (i *IndexedAttestation) String() string {
	data, err := yaml.Marshal(i)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}

	return string(data)
}
