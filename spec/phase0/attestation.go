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
	"strings"

	"github.com/pkg/errors"
	bitfield "github.com/prysmaticlabs/go-bitfield"
)

// Attestation is the Ethereum 2 attestation structure.
type Attestation struct {
	AggregationBits bitfield.Bitlist `ssz-max:"2048"`
	Data            *AttestationData
	Signature       []byte `ssz-size:"96"`
}

// attestationJSON is the spec representation of the struct.
type attestationJSON struct {
	AggregationBits string           `json:"aggregation_bits"`
	Data            *AttestationData `json:"data"`
	Signature       string           `json:"signature"`
}

// MarshalJSON implements json.Marshaler.
func (a *Attestation) MarshalJSON() ([]byte, error) {
	return json.Marshal(&attestationJSON{
		AggregationBits: fmt.Sprintf("%#x", []byte(a.AggregationBits)),
		Data:            a.Data,
		Signature:       fmt.Sprintf("%#x", a.Signature),
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *Attestation) UnmarshalJSON(input []byte) error {
	var err error

	var attestationJSON attestationJSON
	if err = json.Unmarshal(input, &attestationJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if a.AggregationBits, err = hex.DecodeString(strings.TrimPrefix(attestationJSON.AggregationBits, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for beacon block root")
	}
	a.Data = attestationJSON.Data
	if a.Data == nil {
		return errors.New("data missing")
	}
	if a.Signature, err = hex.DecodeString(strings.TrimPrefix(attestationJSON.Signature, "0x")); err != nil {
		return errors.Wrap(err, "invalid value for signature")
	}
	if len(a.Signature) != signatureLength {
		return errors.New("incorrect length for signature")
	}

	return nil
}

// String returns a string version of the structure.
func (a *Attestation) String() string {
	data, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
