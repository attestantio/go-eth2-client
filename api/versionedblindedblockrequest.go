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

	apiv1bellatrix "github.com/attestantio/go-eth2-client/api/v1/bellatrix"
	apiv1capella "github.com/attestantio/go-eth2-client/api/v1/capella"
	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedBlindedBlockRequest contains a versioned signed blinded beacon block.
type VersionedBlindedBlockRequest struct {
	Version   spec.DataVersion
	Bellatrix *apiv1bellatrix.SignedBlindedBeaconBlock
	Capella   *apiv1capella.SignedBlindedBeaconBlock
	Deneb     *apiv1deneb.SignedBlindedBlockContents
}

// Slot returns the slot of the signed beacon block.
func (v *VersionedBlindedBlockRequest) Slot() (phase0.Slot, error) {
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
		if v.Deneb == nil || v.Deneb.SignedBlindedBlock == nil || v.Deneb.SignedBlindedBlock.Message == nil {
			return 0, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlindedBlock.Message.Slot, nil
	default:
		return 0, errors.New("unsupported version")
	}
}

// Attestations returns the attestations of the beacon block.
func (v *VersionedBlindedBlockRequest) Attestations() ([]*phase0.Attestation, error) {
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
		if v.Deneb == nil || v.Deneb.SignedBlindedBlock == nil || v.Deneb.SignedBlindedBlock.Message == nil || v.Deneb.SignedBlindedBlock.Message.Body == nil {
			return nil, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlindedBlock.Message.Body.Attestations, nil
	default:
		return nil, errors.New("unsupported version")
	}
}

// Root returns the root of the beacon block.
func (v *VersionedBlindedBlockRequest) Root() (phase0.Root, error) {
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
		return v.Deneb.SignedBlindedBlock.Message.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// BodyRoot returns the body root of the beacon block.
func (v *VersionedBlindedBlockRequest) BodyRoot() (phase0.Root, error) {
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
		return v.Deneb.SignedBlindedBlock.Message.Body.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// ParentRoot returns the parent root of the beacon block.
func (v *VersionedBlindedBlockRequest) ParentRoot() (phase0.Root, error) {
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
		return v.Deneb.SignedBlindedBlock.Message.ParentRoot, nil
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// StateRoot returns the state root of the beacon block.
func (v *VersionedBlindedBlockRequest) StateRoot() (phase0.Root, error) {
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
		return v.Deneb.SignedBlindedBlock.Message.StateRoot, nil
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// AttesterSlashings returns the attester slashings of the beacon block.
func (v *VersionedBlindedBlockRequest) AttesterSlashings() ([]*phase0.AttesterSlashing, error) {
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
		return v.Deneb.SignedBlindedBlock.Message.Body.AttesterSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ProposerSlashings returns the proposer slashings of the beacon block.
func (v *VersionedBlindedBlockRequest) ProposerSlashings() ([]*phase0.ProposerSlashing, error) {
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
		return v.Deneb.SignedBlindedBlock.Message.Body.ProposerSlashings, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// ProposerIndex returns the proposer index of the beacon block.
func (v *VersionedBlindedBlockRequest) ProposerIndex() (phase0.ValidatorIndex, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil {
			return 0, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.ProposerIndex, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil {
			return 0, errors.New("no capella block")
		}
		return v.Capella.Message.ProposerIndex, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlindedBlock.Message == nil {
			return 0, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlindedBlock.Message.ProposerIndex, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// ExecutionBlockHash returns the hash of the beacon block.
func (v *VersionedBlindedBlockRequest) ExecutionBlockHash() (phase0.Hash32, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil || v.Bellatrix.Message.Body.ExecutionPayloadHeader == nil {
			return phase0.Hash32{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Body.ExecutionPayloadHeader.BlockHash, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil || v.Capella.Message.Body.ExecutionPayloadHeader == nil {
			return phase0.Hash32{}, errors.New("no capella block")
		}
		return v.Capella.Message.Body.ExecutionPayloadHeader.BlockHash, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlindedBlock.Message == nil || v.Deneb.SignedBlindedBlock.Message.Body == nil || v.Deneb.SignedBlindedBlock.Message.Body.ExecutionPayloadHeader == nil {
			return phase0.Hash32{}, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlindedBlock.Message.Body.ExecutionPayloadHeader.BlockHash, nil
	default:
		return phase0.Hash32{}, errors.New("unknown version")
	}
}

// ExecutionBlockNumber returns the block number of the beacon block.
func (v *VersionedBlindedBlockRequest) ExecutionBlockNumber() (uint64, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Message == nil || v.Bellatrix.Message.Body == nil || v.Bellatrix.Message.Body.ExecutionPayloadHeader == nil {
			return 0, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Message.Body.ExecutionPayloadHeader.BlockNumber, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Message == nil || v.Capella.Message.Body == nil || v.Capella.Message.Body.ExecutionPayloadHeader == nil {
			return 0, errors.New("no capella block")
		}
		return v.Capella.Message.Body.ExecutionPayloadHeader.BlockNumber, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlindedBlock.Message == nil || v.Deneb.SignedBlindedBlock.Message.Body == nil || v.Deneb.SignedBlindedBlock.Message.Body.ExecutionPayloadHeader == nil {
			return 0, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlindedBlock.Message.Body.ExecutionPayloadHeader.BlockNumber, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// Signature returns the signature of the beacon block.
func (v *VersionedBlindedBlockRequest) BeaconBlockSignature() (phase0.BLSSignature, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.BLSSignature{}, errors.New("no bellatrix block")
		}
		return v.Bellatrix.Signature, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.BLSSignature{}, errors.New("no capella block")
		}
		return v.Capella.Signature, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.SignedBlindedBlock == nil {
			return phase0.BLSSignature{}, errors.New("no deneb block")
		}
		return v.Deneb.SignedBlindedBlock.Signature, nil
	default:
		return phase0.BLSSignature{}, errors.New("unknown version")
	}
}
