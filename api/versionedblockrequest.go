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

package api

import (
	"errors"

	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedBlockRequest contains a versioned signed beacon block request.
type VersionedBlockRequest struct {
	Version   spec.DataVersion
	Bellatrix *bellatrix.SignedBeaconBlock
	Capella   *capella.SignedBeaconBlock
	Deneb     *apiv1deneb.SignedBlockContents
}

// Slot returns the slot of the signed beacon block.
func (v *VersionedBlockRequest) Slot() (phase0.Slot, error) {
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
		if v.Deneb == nil || v.Deneb.SignedBlock == nil || v.Deneb.SignedBlock.Message == nil {
			return 0, errors.New("no denb block")
		}
		return v.Deneb.SignedBlock.Message.Slot, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExecutionBlockHash returns the block hash of the beacon block.
func (v *VersionedBlockRequest) ExecutionBlockHash() (phase0.Hash32, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil || v.Bellatrix.Message.Body.ExecutionPayload == nil {
			return phase0.Hash32{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Body.ExecutionPayload.BlockHash, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil || v.Capella.Message.Body.ExecutionPayload == nil {
			return phase0.Hash32{}, errors.New("no capella block")
		}
		return v.Capella.Message.Body.ExecutionPayload.BlockHash, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlock == nil || v.Deneb.SignedBlock.Message == nil || v.Deneb.SignedBlock.Message.Body == nil || v.Deneb.SignedBlock.Message.Body.ExecutionPayload == nil {
			return phase0.Hash32{}, errors.New("no denb block")
		}
		return v.Deneb.SignedBlock.Message.Body.ExecutionPayload.BlockHash, nil
	default:
		return phase0.Hash32{}, errors.New("unknown version")
	}
}

// Attestations returns the attestations of the beacon block.
func (v *VersionedBlockRequest) Attestations() ([]*phase0.Attestation, error) {
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
		if v.Deneb == nil || v.Deneb.SignedBlock == nil || v.Deneb.SignedBlock.Message == nil || v.Deneb.SignedBlock.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlock.Message.Body.Attestations, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// Root returns the root of the beacon block.
func (v *VersionedBlockRequest) Root() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil {
			return phase0.Root{}, errors.New("no capella block")
		}
		return v.Capella.Message.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlock == nil || v.Deneb.SignedBlock.Message == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlock.Message.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// BodyRoot returns the body root of the beacon block.
func (v *VersionedBlockRequest) BodyRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Body.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return phase0.Root{}, errors.New("no capella block")
		}
		return v.Capella.Message.Body.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlock == nil || v.Deneb.SignedBlock.Message == nil || v.Deneb.SignedBlock.Message.Body == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlock.Message.Body.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// ParentRoot returns the parent root of the beacon block.
func (v *VersionedBlockRequest) ParentRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.ParentRoot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil {
			return phase0.Root{}, errors.New("no capella block")
		}
		return v.Capella.Message.ParentRoot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlock == nil || v.Deneb.SignedBlock.Message == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlock.Message.ParentRoot, nil
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// StateRoot returns the state root of the beacon block.
func (v *VersionedBlockRequest) StateRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.StateRoot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil {
			return phase0.Root{}, errors.New("no capella block")
		}
		return v.Capella.Message.StateRoot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlock == nil || v.Deneb.SignedBlock.Message == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlock.Message.StateRoot, nil
	default:
		return phase0.Root{}, errors.New("unknown version")
	}
}

// AttesterSlashings returns the attester slashings of the beacon block.
func (v *VersionedBlockRequest) AttesterSlashings() ([]*phase0.AttesterSlashing, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil {
			return nil, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Body.AttesterSlashings, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}
		return v.Capella.Message.Body.AttesterSlashings, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlock == nil || v.Deneb.SignedBlock.Message == nil || v.Deneb.SignedBlock.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlock.Message.Body.AttesterSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ProposerSlashings returns the proposer slashings of the beacon block.
func (v *VersionedBlockRequest) ProposerSlashings() ([]*phase0.ProposerSlashing, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil {
			return nil, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Body.ProposerSlashings, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}
		return v.Capella.Message.Body.ProposerSlashings, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlock == nil || v.Deneb.SignedBlock.Message == nil || v.Deneb.SignedBlock.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlock.Message.Body.ProposerSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// SyncAggregate returns the sync aggregate of the beacon block.
func (v *VersionedBlockRequest) SyncAggregate() (*altair.SyncAggregate, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil {
			return nil, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Body.SyncAggregate, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil {
			return nil, errors.New("no capella block")
		}
		return v.Capella.Message.Body.SyncAggregate, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlock == nil || v.Deneb.SignedBlock.Message == nil || v.Deneb.SignedBlock.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlock.Message.Body.SyncAggregate, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedBlockRequest) String() string {
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
	default:
		return "unsupported version"
	}
}
