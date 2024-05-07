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

// VersionedAttesterSlashing contains a versioned attestation.
type VersionedAttesterSlashing struct {
	Version   DataVersion
	Phase0    *phase0.AttesterSlashing
	Altair    *phase0.AttesterSlashing
	Bellatrix *phase0.AttesterSlashing
	Capella   *phase0.AttesterSlashing
	Deneb     *phase0.AttesterSlashing
	Electra   *electra.AttesterSlashing
}

// IsEmpty returns true if there is no block.
func (v *VersionedAttesterSlashing) IsEmpty() bool {
	return v.Phase0 == nil && v.Altair == nil && v.Bellatrix == nil && v.Capella == nil && v.Deneb == nil && v.Electra == nil
}

// Attestation1 returns the first indexed attestation.
func (v *VersionedAttesterSlashing) Attestation1() (*VersionedIndexedAttestation, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no Phase0 indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionPhase0,
			Phase0:  v.Phase0.Attestation1,
		}

		return &versionedIndexedAttestation, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no Altair indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionPhase0,
			Altair:  v.Altair.Attestation1,
		}

		return &versionedIndexedAttestation, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no Bellatrix indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version:   DataVersionPhase0,
			Bellatrix: v.Bellatrix.Attestation1,
		}

		return &versionedIndexedAttestation, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionCapella,
			Capella: v.Capella.Attestation1,
		}

		return &versionedIndexedAttestation, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no Deneb indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionDeneb,
			Deneb:   v.Deneb.Attestation1,
		}

		return &versionedIndexedAttestation, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionElectra,
			Electra: v.Electra.Attestation1,
		}

		return &versionedIndexedAttestation, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Attestation2 returns the second indexed attestation.
func (v *VersionedAttesterSlashing) Attestation2() (*VersionedIndexedAttestation, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return nil, errors.New("no Phase0 indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionPhase0,
			Phase0:  v.Phase0.Attestation2,
		}

		return &versionedIndexedAttestation, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return nil, errors.New("no Altair indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionPhase0,
			Altair:  v.Altair.Attestation2,
		}

		return &versionedIndexedAttestation, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no Bellatrix indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version:   DataVersionPhase0,
			Bellatrix: v.Bellatrix.Attestation2,
		}

		return &versionedIndexedAttestation, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no Capella indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionCapella,
			Capella: v.Capella.Attestation2,
		}

		return &versionedIndexedAttestation, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no Deneb indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionDeneb,
			Deneb:   v.Deneb.Attestation2,
		}

		return &versionedIndexedAttestation, nil
	case DataVersionElectra:
		if v.Electra == nil {
			return nil, errors.New("no Electra indexed attestation")
		}

		versionedIndexedAttestation := VersionedIndexedAttestation{
			Version: DataVersionElectra,
			Electra: v.Electra.Attestation2,
		}

		return &versionedIndexedAttestation, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedAttesterSlashing) String() string {
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
