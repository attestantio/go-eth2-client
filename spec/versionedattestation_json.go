// Copyright © 2025 Attestant Limited.
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

package spec

import (
	"encoding/json"

	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// attestationIdentificationJSON contains fields that allow us to identify the attestation variant.
type attestationIdentificationJSON struct {
	CommitteeBits *string `json:"committee_bits"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *VersionedAttestation) UnmarshalJSON(input []byte) error {
	var id attestationIdentificationJSON
	if err := json.Unmarshal(input, &id); err != nil {
		return errors.Wrap(err, "invalid JSON")
	}

	switch {
	case id.CommitteeBits != nil:
		v.Version = DataVersionElectra
		v.Electra = &electra.Attestation{}

		return v.Electra.UnmarshalJSON(input)
	default:
		v.Version = DataVersionPhase0
		v.Phase0 = &phase0.Attestation{}

		return v.Phase0.UnmarshalJSON(input)
	}
}
