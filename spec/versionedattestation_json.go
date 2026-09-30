// Copyright © 2025 - 2026 Attestant Limited.
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

// MarshalJSON implements json.Marshaler.
//
// It emits the attestation of the populated fork, which is the shape
// UnmarshalJSON reads and the shape a node sends on the attestation event
// topic.  Without it, marshalling falls back to the struct's own layout and
// produces {"Altair":null,...,"Electra":{...},"Version":"electra"} -- Go field
// names, one key per fork -- so an attestation decoded from an event and
// re-marshalled is no longer the beacon-API shape it arrived as.
func (v *VersionedAttestation) MarshalJSON() ([]byte, error) {
	switch v.Version {
	case DataVersionPhase0:
		return marshalVersionedAttestation(v.Phase0)
	case DataVersionAltair:
		return marshalVersionedAttestation(v.Altair)
	case DataVersionBellatrix:
		return marshalVersionedAttestation(v.Bellatrix)
	case DataVersionCapella:
		return marshalVersionedAttestation(v.Capella)
	case DataVersionDeneb:
		return marshalVersionedAttestation(v.Deneb)
	case DataVersionElectra:
		return marshalVersionedAttestation(v.Electra)
	case DataVersionFulu:
		return marshalVersionedAttestation(v.Fulu)
	case DataVersionGloas:
		return marshalVersionedAttestation(v.Gloas)
	case DataVersionUnknown:
		return nil, errors.New("unknown version")
	default:
		return nil, errors.New("unsupported version")
	}
}

// marshalVersionedAttestation marshals the attestation of the populated fork,
// rejecting a version whose arm is nil rather than emitting "null".
func marshalVersionedAttestation[T any](attestation *T) ([]byte, error) {
	if attestation == nil {
		return nil, errors.New("no attestation")
	}

	return json.Marshal(attestation)
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
