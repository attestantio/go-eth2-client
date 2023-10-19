// Copyright © 2021 - 2023 Attestant Limited.
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

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedBeaconBlock contains a versioned beacon block.
type VersionedBeaconBlock struct {
	Version   DataVersion
	Phase0    *phase0.BeaconBlock
	Altair    *altair.BeaconBlock
	Bellatrix *bellatrix.BeaconBlock
	Capella   *capella.BeaconBlock
	Deneb     *deneb.BeaconBlock
}

// IsEmpty returns true if there is no block.
func (v *VersionedBeaconBlock) IsEmpty() bool {
	return v.Phase0 == nil && v.Altair == nil && v.Bellatrix == nil && v.Capella == nil && v.Deneb == nil
}

// Slot returns the slot of the beacon block.
func (v *VersionedBeaconBlock) Slot() (phase0.Slot, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, errors.New("no phase0 block")
		}

		return v.Phase0.Slot, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return 0, errors.New("no altair block")
		}

		return v.Altair.Slot, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Slot, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.Slot, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no deneb block contents")
		}

		return v.Deneb.Slot, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// RandaoReveal returns the RANDAO reveal of the beacon block.
func (v *VersionedBeaconBlock) RandaoReveal() (phase0.BLSSignature, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.BLSSignature{}, errors.New("no phase0 block")
		}
		if v.Phase0.Body == nil {
			return phase0.BLSSignature{}, errors.New("no phase0 block body")
		}

		return v.Phase0.Body.RANDAOReveal, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.BLSSignature{}, errors.New("no altair block")
		}
		if v.Altair.Body == nil {
			return phase0.BLSSignature{}, errors.New("no altair block body")
		}

		return v.Altair.Body.RANDAOReveal, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.BLSSignature{}, errors.New("no bellatrix block")
		}
		if v.Bellatrix.Body == nil {
			return phase0.BLSSignature{}, errors.New("no bellatrix block body")
		}

		return v.Bellatrix.Body.RANDAOReveal, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.BLSSignature{}, errors.New("no capella block")
		}
		if v.Capella.Body == nil {
			return phase0.BLSSignature{}, errors.New("no capella block body")
		}

		return v.Capella.Body.RANDAOReveal, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.BLSSignature{}, errors.New("no deneb block contents")
		}
		if v.Deneb.Body == nil {
			return phase0.BLSSignature{}, errors.New("no deneb block body")
		}

		return v.Deneb.Body.RANDAOReveal, nil
	default:
		return phase0.BLSSignature{}, errors.New("unknown version")
	}
}

// Graffiti returns the graffiti of the beacon block.
func (v *VersionedBeaconBlock) Graffiti() ([32]byte, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return [32]byte{}, errors.New("no phase0 block")
		}
		if v.Phase0.Body == nil {
			return [32]byte{}, errors.New("no phase0 block body")
		}

		return v.Phase0.Body.Graffiti, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return [32]byte{}, errors.New("no altair block")
		}
		if v.Altair.Body == nil {
			return [32]byte{}, errors.New("no altair block body")
		}

		return v.Altair.Body.Graffiti, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return [32]byte{}, errors.New("no bellatrix block")
		}
		if v.Bellatrix.Body == nil {
			return [32]byte{}, errors.New("no bellatrix block body")
		}

		return v.Bellatrix.Body.Graffiti, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return [32]byte{}, errors.New("no capella block")
		}
		if v.Capella.Body == nil {
			return [32]byte{}, errors.New("no capella block body")
		}

		return v.Capella.Body.Graffiti, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return [32]byte{}, errors.New("no deneb block contents")
		}
		if v.Deneb.Body == nil {
			return [32]byte{}, errors.New("no deneb block body")
		}

		return v.Deneb.Body.Graffiti, nil
	default:
		return [32]byte{}, errors.New("unknown version")
	}
}

// ProposerIndex returns the proposer index of the beacon block.
func (v *VersionedBeaconBlock) ProposerIndex() (phase0.ValidatorIndex, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, errors.New("no phase0 block")
		}

		return v.Phase0.ProposerIndex, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return 0, errors.New("no altair block")
		}

		return v.Altair.ProposerIndex, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix block")
		}

		return v.Bellatrix.ProposerIndex, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.ProposerIndex, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return 0, errors.New("no deneb block contents")
		}

		return v.Deneb.ProposerIndex, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// Root returns the root of the beacon block.
func (v *VersionedBeaconBlock) Root() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, errors.New("no phase0 block")
		}

		return v.Phase0.HashTreeRoot()
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, errors.New("no altair block")
		}

		return v.Altair.HashTreeRoot()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.HashTreeRoot()
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.HashTreeRoot()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, errors.New("no deneb block contents")
		}

		return v.Deneb.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// BodyRoot returns the body root of the beacon block.
func (v *VersionedBeaconBlock) BodyRoot() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, errors.New("no phase0 block")
		}

		return v.Phase0.Body.HashTreeRoot()
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, errors.New("no altair block")
		}

		return v.Altair.Body.HashTreeRoot()
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Body.HashTreeRoot()
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.Body.HashTreeRoot()
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, errors.New("no deneb block contents")
		}

		return v.Deneb.Body.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// ParentRoot returns the parent root of the beacon block.
func (v *VersionedBeaconBlock) ParentRoot() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, errors.New("no phase0 block")
		}

		return v.Phase0.ParentRoot, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, errors.New("no altair block")
		}

		return v.Altair.ParentRoot, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.ParentRoot, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.ParentRoot, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, errors.New("no deneb block contents")
		}

		return v.Deneb.ParentRoot, nil
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// StateRoot returns the state root of the beacon block.
func (v *VersionedBeaconBlock) StateRoot() (phase0.Root, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, errors.New("no phase0 block")
		}

		return v.Phase0.StateRoot, nil
	case DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, errors.New("no altair block")
		}

		return v.Altair.StateRoot, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.StateRoot, nil
	case DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.StateRoot, nil
	case DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, errors.New("no deneb block contents")
		}

		return v.Deneb.StateRoot, nil
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// Attestations returns the attestations of the beacon block.
func (v *VersionedBeaconBlock) Attestations() ([]*phase0.Attestation, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Body == nil {
			return nil, errors.New("no phase0 block")
		}

		return v.Phase0.Body.Attestations, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Body == nil {
			return nil, errors.New("no altair block")
		}

		return v.Altair.Body.Attestations, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Body == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Body.Attestations, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Body.Attestations, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Body.Attestations, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// AttesterSlashings returns the attester slashings of the beacon block.
func (v *VersionedBeaconBlock) AttesterSlashings() ([]*phase0.AttesterSlashing, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Body == nil {
			return nil, errors.New("no phase0 block")
		}

		return v.Phase0.Body.AttesterSlashings, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Body == nil {
			return nil, errors.New("no altair block")
		}

		return v.Altair.Body.AttesterSlashings, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Body == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Body.AttesterSlashings, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Body.AttesterSlashings, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Body.AttesterSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ProposerSlashings returns the proposer slashings of the beacon block.
func (v *VersionedBeaconBlock) ProposerSlashings() ([]*phase0.ProposerSlashing, error) {
	switch v.Version {
	case DataVersionPhase0:
		if v.Phase0 == nil || v.Phase0.Body == nil {
			return nil, errors.New("no phase0 block")
		}

		return v.Phase0.Body.ProposerSlashings, nil
	case DataVersionAltair:
		if v.Altair == nil || v.Altair.Body == nil {
			return nil, errors.New("no altair block")
		}

		return v.Altair.Body.ProposerSlashings, nil
	case DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Body == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Body.ProposerSlashings, nil
	case DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Body.ProposerSlashings, nil
	case DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Body.ProposerSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedBeaconBlock) String() string {
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
	default:
		return "unknown version"
	}
}
