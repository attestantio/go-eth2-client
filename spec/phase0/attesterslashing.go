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
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// AttesterSlashing provides information about an attester slashing.
type AttesterSlashing struct {
	Attestation1 *IndexedAttestation
	Attestation2 *IndexedAttestation
}

// attesterSlashingJSON is the spec representation of the struct.
type attesterSlashingJSON struct {
	Attestation1 *IndexedAttestation `json:"attestation_1"`
	Attestation2 *IndexedAttestation `json:"attestation_2"`
}

// MarshalJSON implements json.Marshaler.
func (a *AttesterSlashing) MarshalJSON() ([]byte, error) {
	return json.Marshal(&attesterSlashingJSON{
		Attestation1: a.Attestation1,
		Attestation2: a.Attestation2,
	})
}

// UnmarshalJSON implements json.Unmarshaler.
func (a *AttesterSlashing) UnmarshalJSON(input []byte) error {
	var err error

	var attesterSlashingJSON attesterSlashingJSON
	if err = json.Unmarshal(input, &attesterSlashingJSON); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}
	if attesterSlashingJSON.Attestation1 == nil {
		return errors.New("attestation 1 missing")
	}
	a.Attestation1 = attesterSlashingJSON.Attestation1
	if attesterSlashingJSON.Attestation2 == nil {
		return errors.New("attestation 2 missing")
	}
	a.Attestation2 = attesterSlashingJSON.Attestation2

	return nil
}

func (a *AttesterSlashing) String() string {
	data, err := json.Marshal(a)
	if err != nil {
		return fmt.Sprintf("ERR: %v", err)
	}
	return string(data)
}
