// Copyright © 2023 Attestant Limited.
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

package gloas

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// indexedPayloadAttestationJSON is the spec representation of the struct.
type indexedPayloadAttestationJSON struct {
	AttestingIndices []string                `json:"attesting_indices"`
	Data             *PayloadAttestationData `json:"data"`
	Signature        string                  `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (i *IndexedPayloadAttestation) MarshalJSON() ([]byte, error) {
	attestingIndices := make([]string, len(i.AttestingIndices))
	for index := range i.AttestingIndices {
		attestingIndices[index] = fmt.Sprintf("%d", i.AttestingIndices[index])
	}

	return json.Marshal(&indexedPayloadAttestationJSON{
		AttestingIndices: attestingIndices,
		Data:             i.Data,
		Signature:        fmt.Sprintf("%#x", i.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (i *IndexedPayloadAttestation) UnmarshalJSON(input []byte) error {
	var data indexedPayloadAttestationJSON
	if err := json.Unmarshal(input, &data); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	i.AttestingIndices = make([]phase0.ValidatorIndex, len(data.AttestingIndices))
	for index := range data.AttestingIndices {
		validatorIndex, err := strconv.ParseUint(data.AttestingIndices[index], 10, 64)
		if err != nil {
			return errors.Wrap(err, "invalid validator index")
		}
		i.AttestingIndices[index] = phase0.ValidatorIndex(validatorIndex)
	}

	if data.Data == nil {
		return errors.New("data missing")
	}
	i.Data = data.Data

	if data.Signature == "" {
		return errors.New("signature missing")
	}
	signature, err := hex.DecodeString(strings.TrimPrefix(data.Signature, "0x"))
	if err != nil {
		return errors.Wrap(err, "invalid signature")
	}
	copy(i.Signature[:], signature)

	return nil
}
