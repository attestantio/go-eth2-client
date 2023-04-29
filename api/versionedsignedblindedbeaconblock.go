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
	"errors"

	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedSignedBlindedBeaconBlock contains a versioned signed blinded beacon block.
type VersionedSignedBlindedBeaconBlock struct {
	Version   spec.DataVersion
	Bellatrix *apiv1bellatrix.SignedBlindedBeaconBlock
	Capella   *apiv1capella.SignedBlindedBeaconBlock
	Deneb     *apiv1deneb.SignedBlindedBeaconBlock
}

// Slot returns the slot of the signed beacon block.
func (v *VersionedSignedBlindedBeaconBlock) Slot() (phase0.Slot, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil {
			return 0, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Slot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil {
			return 0, errors.New("no capella block")
		}
		return v.Capella.Message.Slot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil {
			return 0, errors.New("no deneb block")
		}
		return v.Deneb.Message.Slot, nil
	default:
		return 0, errors.New("unsupported version")
	}
}

// Attestations returns the attestations of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) Attestations() ([]*phase0.Attestation, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil {
			return nil, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Body.Attestations, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}
		return v.Capella.Message.Body.Attestations, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Message == nil || v.Deneb.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}
		return v.Deneb.Message.Body.Attestations, nil
	default:
		return nil, errors.New("unsupported version")
	}
}

// Root returns the root of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) Root() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}
		return v.Capella.Message.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}
		return v.Deneb.Message.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// BodyRoot returns the body root of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) BodyRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Body.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}
		return v.Capella.Message.Body.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}
		return v.Deneb.Message.Body.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// ParentRoot returns the parent root of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) ParentRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.ParentRoot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}
		return v.Capella.Message.ParentRoot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}
		return v.Deneb.Message.ParentRoot, nil
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// StateRoot returns the state root of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) StateRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.StateRoot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}
		return v.Capella.Message.StateRoot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}
		return v.Deneb.Message.StateRoot, nil
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// AttesterSlashings returns the attester slashings of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) AttesterSlashings() ([]*phase0.AttesterSlashing, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Body.AttesterSlashings, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella block")
		}
		return v.Capella.Message.Body.AttesterSlashings, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no deneb block")
		}
		return v.Deneb.Message.Body.AttesterSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ProposerSlashings returns the proposer slashings of the beacon block.
func (v *VersionedSignedBlindedBeaconBlock) ProposerSlashings() ([]*phase0.ProposerSlashing, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Body.ProposerSlashings, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella block")
		}
		return v.Capella.Message.Body.ProposerSlashings, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no deneb block")
		}
		return v.Deneb.Message.Body.ProposerSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}
