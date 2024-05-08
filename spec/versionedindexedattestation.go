// Copyright © 2024 Attestant Limited.
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
	"errors"

	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedIndexedAttestation contains a versioned indexed attestation.
type VersionedIndexedAttestation struct {
	Version   DataVersion
	Phase0    *phase0.IndexedAttestation
	Altair    *phase0.IndexedAttestation
	Bellatrix *phase0.IndexedAttestation
	Capella   *phase0.IndexedAttestation
	Deneb     *phase0.IndexedAttestation
	Electra   *electra.IndexedAttestation
}

// IsEmpty returns true if there is no block.
func (v *VersionedIndexedAttestation) IsEmpty() bool {
	return v.Phase0 == nil && v.Altair == nil && v.Bellatrix == nil && v.Capella == nil && v.Deneb == nil && v.Electra == nil
}

// AttestingIndices returns the attesting indices of the indexed attestation.
func (v *VersionedIndexedAttestation) AttestingIndices() ([]uint64, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no Phase0 indexed attestation")
		}

		return v.Phase0.AttestingIndices, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no Altair indexed attestation")
		}

		return v.Altair.AttestingIndices, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no Bellatrix indexed attestation")
		}

		return v.Bellatrix.AttestingIndices, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella indexed attestation")
		}

		return v.Capella.AttestingIndices, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no Deneb indexed attestation")
		}

		return v.Deneb.AttestingIndices, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra indexed attestation")
		}

		return v.Electra.AttestingIndices, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Data returns the data of the indexed attestation.
func (v *VersionedIndexedAttestation) Data() (*phase0.AttestationData, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no Phase0 indexed attestation")
		}

		return v.Phase0.Data, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no Altair indexed attestation")
		}

		return v.Altair.Data, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no Bellatrix indexed attestation")
		}

		return v.Bellatrix.Data, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella indexed attestation")
		}

		return v.Capella.Data, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no Deneb indexed attestation")
		}

		return v.Deneb.Data, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra indexed attestation")
		}

		return v.Electra.Data, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Signature returns the signature of the indexed attestation.
func (v *VersionedIndexedAttestation) Signature() (phase0.BLSSignature, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.BLSSignature{}, errors.New("no Phase0 indexed attestation")
		}

		return v.Phase0.Signature, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.BLSSignature{}, errors.New("no Altair indexed attestation")
		}

		return v.Altair.Signature, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.BLSSignature{}, errors.New("no Bellatrix indexed attestation")
		}

		return v.Bellatrix.Signature, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.BLSSignature{}, errors.New("no Capella indexed attestation")
		}

		return v.Capella.Signature, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.BLSSignature{}, errors.New("no Deneb indexed attestation")
		}

		return v.Deneb.Signature, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return phase0.BLSSignature{}, errors.New("no Electra indexed attestation")
		}

		return v.Electra.Signature, nil
	default:
		return phase0.BLSSignature{}, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedIndexedAttestation) String() string {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return ""
		}

		return v.Phase0.String()
	case DataVersionAltair:
		if v.Altair == nil {
			return ""
		}

		return v.Altair.String()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	case DataVersionElectra:
		if v.Electra == nil {
			return ""
		}

		return v.Electra.String()
	default:
		return "unknown version"
	}
}
