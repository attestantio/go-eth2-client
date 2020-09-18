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
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// IndexedAttestation provides a signed attestation with a list of attesting indices.
type IndexedAttestation struct {
	AttestingIndices []uint64
	Data             *AttestationData
	Signature        []byte
}

// indexedAttestationJSON is the spec representation of the struct.
type indexedAttestationJSON struct {
	AttestingIndices []string         `json:"attesting_indices"`
	Data             *AttestationData `json:"data"`
	Signature        string           `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (i *IndexedAttestation) MarshalJSON() ([]byte, error) {
	attestingIndices := make([]string, len(i.AttestingIndices))
	for j := range i.AttestingIndices {
		attestingIndices[j] = fmt.Sprintf("%d", i.AttestingIndices[j])
	}
	return json.Marshal(&indexedAttestationJSON{
		AttestingIndices: attestingIndices,
		Data:             i.Data,
		Signature:        fmt.Sprintf("%#x", i.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *IndexedAttestation) UnmarshalJSON(input []byte) error {
	var err error

	var indexedAttestationJSON indexedAttestationJSON
	if err = json.Unmarshal(input, &indexedAttestationJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if indexedAttestationJSON.AttestingIndices == nil {
		return errors.New("attesting indices missing")
	}
	if len(indexedAttestationJSON.AttestingIndices) == 0 {
		return errors.New("attesting indices missing")
	}
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
	if i.Signature, err = hex.DecodeString(strings.TrimPrefix(indexedAttestationJSON.Signature, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(i.Signature) != signatureLength {
		return errors.New("incorrect length for signature")
	}

	return nil
}

// String returns a string version of the structure.
func (i *IndexedAttestation) String() string {
	data, err := json.Marshal(i)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
