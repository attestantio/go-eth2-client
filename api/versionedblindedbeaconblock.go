// Copyright © 2022, 2023 Attestant Limited.
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

package api

import (
	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	apiv1electra "github.com/attestantio/go-eth2-client/api/v1/electra"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedBlindedBeaconBlock contains a versioned blinded beacon block.
type VersionedBlindedBeaconBlock struct {
	Version   spec.DataVersion
	Bellatrix *apiv1bellatrix.BlindedBeaconBlock
	Capella   *apiv1capella.BlindedBeaconBlock
	Deneb     *apiv1deneb.BlindedBeaconBlock
	Electra   *apiv1electra.BlindedBeaconBlock
	Fulu      *apiv1electra.BlindedBeaconBlock
	Gloas     *apiv1electra.BlindedBeaconBlock
}

// IsEmpty returns true if there is no block.
func (v *VersionedBlindedBeaconBlock) IsEmpty() bool {
	return v.Bellatrix == nil && v.Capella == nil && v.Deneb == nil && v.Electra == nil && v.Fulu == nil && v.Gloas == nil
}

// Slot returns the slot of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) Slot() (phase0.Slot, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, ErrDataMissing
		}

		return v.Bellatrix.Slot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return 0, ErrDataMissing
		}

		return v.Capella.Slot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return 0, ErrDataMissing
		}

		return v.Deneb.Slot, nil
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return 0, ErrDataMissing
		}

		return v.Electra.Slot, nil
	case spec.DataVersionFulu:
		if v.Fulu == nil {
			return 0, ErrDataMissing
		}

		return v.Fulu.Slot, nil
	case spec.DataVersionGloas:
		if v.Gloas == nil {
			return 0, ErrDataMissing
		}

		return v.Gloas.Slot, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// ProposerIndex returns the proposer index of the beacon block.
func (v *VersionedBlindedBeaconBlock) ProposerIndex() (phase0.ValidatorIndex, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, ErrDataMissing
		}

		return v.Bellatrix.ProposerIndex, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return 0, ErrDataMissing
		}

		return v.Capella.ProposerIndex, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return 0, ErrDataMissing
		}

		return v.Deneb.ProposerIndex, nil
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return 0, ErrDataMissing
		}

		return v.Electra.ProposerIndex, nil
	case spec.DataVersionFulu:
		if v.Fulu == nil {
			return 0, ErrDataMissing
		}

		return v.Fulu.ProposerIndex, nil
	case spec.DataVersionGloas:
		if v.Gloas == nil {
			return 0, ErrDataMissing
		}

		return v.Gloas.ProposerIndex, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// RandaoReveal returns the RANDAO reveal of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) RandaoReveal() (phase0.BLSSignature, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}

		return v.Bellatrix.Body.RANDAOReveal, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}

		return v.Capella.Body.RANDAOReveal, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Body == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}

		return v.Deneb.Body.RANDAOReveal, nil
	case spec.DataVersionElectra:
		if v.Electra == nil ||
			v.Electra.Body == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}

		return v.Electra.Body.RANDAOReveal, nil
	case spec.DataVersionFulu:
		if v.Fulu == nil ||
			v.Fulu.Body == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}

		return v.Fulu.Body.RANDAOReveal, nil
	case spec.DataVersionGloas:
		if v.Gloas == nil ||
			v.Gloas.Body == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}

		return v.Gloas.Body.RANDAOReveal, nil
	default:
		return phase0.BLSSignature{}, ErrUnsupportedVersion
	}
}

// Graffiti returns the graffiti of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) Graffiti() ([32]byte, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil {
			return [32]byte{}, ErrDataMissing
		}

		return v.Bellatrix.Body.Graffiti, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil {
			return [32]byte{}, ErrDataMissing
		}

		return v.Capella.Body.Graffiti, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Body == nil {
			return [32]byte{}, ErrDataMissing
		}

		return v.Deneb.Body.Graffiti, nil
	case spec.DataVersionElectra:
		if v.Electra == nil ||
			v.Electra.Body == nil {
			return [32]byte{}, ErrDataMissing
		}

		return v.Electra.Body.Graffiti, nil
	case spec.DataVersionFulu:
		if v.Fulu == nil ||
			v.Fulu.Body == nil {
			return [32]byte{}, ErrDataMissing
		}

		return v.Fulu.Body.Graffiti, nil
	case spec.DataVersionGloas:
		if v.Gloas == nil ||
			v.Gloas.Body == nil {
			return [32]byte{}, ErrDataMissing
		}

		return v.Gloas.Body.Graffiti, nil
	default:
		return [32]byte{}, ErrUnsupportedVersion
	}
}

// Attestations returns the attestations of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) Attestations() ([]spec.VersionedAttestation, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Body == nil {
			return nil, ErrDataMissing
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Bellatrix.Body.Attestations))
		for i, attestation := range v.Bellatrix.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version:   spec.DataVersionBellatrix,
				Bellatrix: attestation,
			}
		}

		return versionedAttestations, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return nil, ErrDataMissing
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Capella.Body.Attestations))
		for i, attestation := range v.Capella.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version: spec.DataVersionCapella,
				Capella: attestation,
			}
		}

		return versionedAttestations, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Body == nil {
			return nil, ErrDataMissing
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Deneb.Body.Attestations))
		for i, attestation := range v.Deneb.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version: spec.DataVersionDeneb,
				Deneb:   attestation,
			}
		}

		return versionedAttestations, nil
	case spec.DataVersionElectra:
		if v.Electra == nil || v.Electra.Body == nil {
			return nil, ErrDataMissing
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Electra.Body.Attestations))
		for i, attestation := range v.Electra.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version: spec.DataVersionElectra,
				Electra: attestation,
			}
		}

		return versionedAttestations, nil
	case spec.DataVersionFulu:
		if v.Fulu == nil || v.Fulu.Body == nil {
			return nil, ErrDataMissing
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Fulu.Body.Attestations))
		for i, attestation := range v.Fulu.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version: spec.DataVersionFulu,
				Fulu:    attestation,
			}
		}

		return versionedAttestations, nil
	case spec.DataVersionGloas:
		if v.Gloas == nil || v.Gloas.Body == nil {
			return nil, ErrDataMissing
		}

		versionedAttestations := make([]spec.VersionedAttestation, len(v.Gloas.Body.Attestations))
		for i, attestation := range v.Gloas.Body.Attestations {
			versionedAttestations[i] = spec.VersionedAttestation{
				Version: spec.DataVersionGloas,
				Gloas:   attestation,
			}
		}

		return versionedAttestations, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// Root returns the root of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) Root() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Bellatrix.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Capella.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.HashTreeRoot()
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Electra.HashTreeRoot()
	case spec.DataVersionFulu:
		if v.Fulu == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Fulu.HashTreeRoot()
	case spec.DataVersionGloas:
		if v.Gloas == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Gloas.HashTreeRoot()
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// BodyRoot returns the body root of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) BodyRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Bellatrix.Body.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Capella.Body.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.Body.HashTreeRoot()
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Electra.Body.HashTreeRoot()
	case spec.DataVersionFulu:
		if v.Fulu == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Fulu.Body.HashTreeRoot()
	case spec.DataVersionGloas:
		if v.Gloas == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Gloas.Body.HashTreeRoot()
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// ParentRoot returns the parent root of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) ParentRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Bellatrix.ParentRoot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Capella.ParentRoot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.ParentRoot, nil
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Electra.ParentRoot, nil
	case spec.DataVersionFulu:
		if v.Fulu == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Fulu.ParentRoot, nil
	case spec.DataVersionGloas:
		if v.Gloas == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Gloas.ParentRoot, nil
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// StateRoot returns the state root of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) StateRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Bellatrix.StateRoot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Capella.StateRoot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.StateRoot, nil
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Electra.StateRoot, nil
	case spec.DataVersionFulu:
		if v.Fulu == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Fulu.StateRoot, nil
	case spec.DataVersionGloas:
		if v.Gloas == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Gloas.StateRoot, nil
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// TransactionsRoot returns the transactions root of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) TransactionsRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil ||
			v.Bellatrix.Body.ExecutionPayloadHeader == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Bellatrix.Body.ExecutionPayloadHeader.TransactionsRoot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayloadHeader == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayloadHeader.TransactionsRoot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Body == nil ||
			v.Deneb.Body.ExecutionPayloadHeader == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.Body.ExecutionPayloadHeader.TransactionsRoot, nil
	case spec.DataVersionElectra:
		if v.Electra == nil ||
			v.Electra.Body == nil ||
			v.Electra.Body.ExecutionPayloadHeader == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Electra.Body.ExecutionPayloadHeader.TransactionsRoot, nil
	case spec.DataVersionFulu:
		if v.Fulu == nil ||
			v.Fulu.Body == nil ||
			v.Fulu.Body.ExecutionPayloadHeader == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Fulu.Body.ExecutionPayloadHeader.TransactionsRoot, nil
	case spec.DataVersionGloas:
		if v.Gloas == nil ||
			v.Gloas.Body == nil ||
			v.Gloas.Body.ExecutionPayloadHeader == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Gloas.Body.ExecutionPayloadHeader.TransactionsRoot, nil
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// FeeRecipient returns the fee recipient of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil ||
			v.Bellatrix.Body.ExecutionPayloadHeader == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Bellatrix.Body.ExecutionPayloadHeader.FeeRecipient, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayloadHeader == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayloadHeader.FeeRecipient, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Body == nil ||
			v.Deneb.Body.ExecutionPayloadHeader == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Deneb.Body.ExecutionPayloadHeader.FeeRecipient, nil
	case spec.DataVersionElectra:
		if v.Electra == nil ||
			v.Electra.Body == nil ||
			v.Electra.Body.ExecutionPayloadHeader == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Electra.Body.ExecutionPayloadHeader.FeeRecipient, nil
	case spec.DataVersionFulu:
		if v.Fulu == nil ||
			v.Fulu.Body == nil ||
			v.Fulu.Body.ExecutionPayloadHeader == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Fulu.Body.ExecutionPayloadHeader.FeeRecipient, nil
	case spec.DataVersionGloas:
		if v.Gloas == nil ||
			v.Gloas.Body == nil ||
			v.Gloas.Body.ExecutionPayloadHeader == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Gloas.Body.ExecutionPayloadHeader.FeeRecipient, nil
	default:
		return bellatrix.ExecutionAddress{}, ErrUnsupportedVersion
	}
}

// Timestamp returns the timestamp of the blinded beacon block.
func (v *VersionedBlindedBeaconBlock) Timestamp() (uint64, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil ||
			v.Bellatrix.Body.ExecutionPayloadHeader == nil {
			return 0, ErrDataMissing
		}

		return v.Bellatrix.Body.ExecutionPayloadHeader.Timestamp, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayloadHeader == nil {
			return 0, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayloadHeader.Timestamp, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Body == nil ||
			v.Deneb.Body.ExecutionPayloadHeader == nil {
			return 0, ErrDataMissing
		}

		return v.Deneb.Body.ExecutionPayloadHeader.Timestamp, nil
	case spec.DataVersionElectra:
		if v.Electra == nil ||
			v.Electra.Body == nil ||
			v.Electra.Body.ExecutionPayloadHeader == nil {
			return 0, ErrDataMissing
		}

		return v.Electra.Body.ExecutionPayloadHeader.Timestamp, nil
	case spec.DataVersionFulu:
		if v.Fulu == nil ||
			v.Fulu.Body == nil ||
			v.Fulu.Body.ExecutionPayloadHeader == nil {
			return 0, ErrDataMissing
		}

		return v.Fulu.Body.ExecutionPayloadHeader.Timestamp, nil
	case spec.DataVersionGloas:
		if v.Gloas == nil ||
			v.Gloas.Body == nil ||
			v.Gloas.Body.ExecutionPayloadHeader == nil {
			return 0, ErrDataMissing
		}

		return v.Gloas.Body.ExecutionPayloadHeader.Timestamp, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// String returns a string version of the structure.
func (v *VersionedBlindedBeaconBlock) String() string {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return ""
		}

		return v.Bellatrix.String()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return ""
		}

		return v.Capella.String()
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return ""
		}

		return v.Deneb.String()
	case spec.DataVersionElectra:
		if v.Electra == nil {
			return ""
		}

		return v.Electra.String()
	case spec.DataVersionFulu:
		if v.Fulu == nil {
			return ""
		}

		return v.Fulu.String()
	case spec.DataVersionGloas:
		if v.Gloas == nil {
			return ""
		}

		return v.Gloas.String()
	default:
		return "unknown version"
	}
}
