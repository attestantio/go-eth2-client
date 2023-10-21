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

	apiv1deneb "github.com/attestantio/go-eth2-client/api/v1/deneb"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedProposal contains a versioned proposal.
type VersionedProposal struct {
	Version   spec.DataVersion
	Phase0    *phase0.BeaconBlock
	Altair    *altair.BeaconBlock
	Bellatrix *bellatrix.BeaconBlock
	Capella   *capella.BeaconBlock
	Deneb     *apiv1deneb.BlockContents
}

// IsEmpty returns true if there is no proposal.
func (v *VersionedProposal) IsEmpty() bool {
	return v.Phase0 == nil &&
		v.Altair == nil &&
		v.Bellatrix == nil &&
		v.Capella == nil &&
		v.Deneb == nil
}

// Slot returns the slot of the proposal.
func (v *VersionedProposal) Slot() (phase0.Slot, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Slot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.Slot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Block == nil {
			return 0, errors.New("no deneb block")
		}

		return v.Deneb.Block.Slot, nil
	default:
		return 0, errors.New("unsupported version")
	}
}

// ProposerIndex returns the proposer index of the proposal.
func (v *VersionedProposal) ProposerIndex() (phase0.ValidatorIndex, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix block")
		}

		return v.Bellatrix.ProposerIndex, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella block")
		}

		return v.Capella.ProposerIndex, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Block == nil {
			return 0, errors.New("no deneb block")
		}

		return v.Deneb.Block.ProposerIndex, nil
	default:
		return 0, errors.New("unknown version")
	}
}

// RandaoReveal returns the RANDAO reveal of the proposal.
func (v *VersionedProposal) RandaoReveal() (phase0.BLSSignature, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Body == nil {
			return phase0.BLSSignature{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Body.RANDAOReveal, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return phase0.BLSSignature{}, errors.New("no capella block")
		}

		return v.Capella.Body.RANDAOReveal, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Block == nil || v.Deneb.Block.Body == nil {
			return phase0.BLSSignature{}, errors.New("no deneb block")
		}

		return v.Deneb.Block.Body.RANDAOReveal, nil
	default:
		return phase0.BLSSignature{}, errors.New("unsupported version")
	}
}

// Graffiti returns the graffiti of the proposal.
func (v *VersionedProposal) Graffiti() ([32]byte, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Body == nil {
			return [32]byte{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Body.Graffiti, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return [32]byte{}, errors.New("no capella block")
		}

		return v.Capella.Body.Graffiti, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Block == nil || v.Deneb.Block.Body == nil {
			return [32]byte{}, errors.New("no deneb block")
		}

		return v.Deneb.Block.Body.Graffiti, nil
	default:
		return [32]byte{}, errors.New("unsupported version")
	}
}

// Attestations returns the attestations of the proposal.
func (v *VersionedProposal) Attestations() ([]*phase0.Attestation, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil || v.Bellatrix.Body == nil {
			return nil, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Body.Attestations, nil
	case spec.DataVersionCapella:
		if v.Capella == nil || v.Capella.Body == nil {
			return nil, errors.New("no capella block")
		}

		return v.Capella.Body.Attestations, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Block == nil || v.Deneb.Block.Body == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Block.Body.Attestations, nil
	default:
		return nil, errors.New("unsupported version")
	}
}

// Root returns the root of the proposal.
func (v *VersionedProposal) Root() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Block == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}

		return v.Deneb.Block.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// BodyRoot returns the body root of the proposal.
func (v *VersionedProposal) BodyRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.Body.HashTreeRoot()
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.Body.HashTreeRoot()
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Block == nil || v.Deneb.Block.Body == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}

		return v.Deneb.Block.Body.HashTreeRoot()
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// ParentRoot returns the parent root of the proposal.
func (v *VersionedProposal) ParentRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.ParentRoot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.ParentRoot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Block == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}

		return v.Deneb.Block.ParentRoot, nil
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// StateRoot returns the state root of the proposal.
func (v *VersionedProposal) StateRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return phase0.Root{}, errors.New("no bellatrix block")
		}

		return v.Bellatrix.StateRoot, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return phase0.Root{}, errors.New("no capella block")
		}

		return v.Capella.StateRoot, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil || v.Deneb.Block == nil {
			return phase0.Root{}, errors.New("no deneb block")
		}

		return v.Deneb.Block.StateRoot, nil
	default:
		return phase0.Root{}, errors.New("unsupported version")
	}
}

// Transactions returns the transactions of the proposal.
func (v *VersionedProposal) Transactions() ([]bellatrix.Transaction, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return nil, errors.New("no bellatrix block")
		}
		if v.Bellatrix.Body == nil {
			return nil, errors.New("no bellatrix block body")
		}
		if v.Bellatrix.Body.ExecutionPayload == nil {
			return nil, errors.New("no bellatrix block body execution payload header")
		}

		return v.Bellatrix.Body.ExecutionPayload.Transactions, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return nil, errors.New("no capella block")
		}
		if v.Capella.Body == nil {
			return nil, errors.New("no capella block body")
		}
		if v.Capella.Body.ExecutionPayload == nil {
			return nil, errors.New("no capella block body execution payload header")
		}

		return v.Capella.Body.ExecutionPayload.Transactions, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil ||
			v.Deneb.Block.Body.ExecutionPayload == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.Block.Body.ExecutionPayload.Transactions, nil
	default:
		return nil, errors.New("unsupported version")
	}
}

// FeeRecipient returns the fee recipient of the proposal.
func (v *VersionedProposal) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no bellatrix block")
		}
		if v.Bellatrix.Body == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no bellatrix block body")
		}
		if v.Bellatrix.Body.ExecutionPayload == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no bellatrix block body execution payload header")
		}

		return v.Bellatrix.Body.ExecutionPayload.FeeRecipient, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no capella block")
		}
		if v.Capella.Body == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no capella block body")
		}
		if v.Capella.Body.ExecutionPayload == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no capella block body execution payload header")
		}

		return v.Capella.Body.ExecutionPayload.FeeRecipient, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil ||
			v.Deneb.Block.Body.ExecutionPayload == nil {
			return bellatrix.ExecutionAddress{}, errors.New("no deneb block")
		}

		return v.Deneb.Block.Body.ExecutionPayload.FeeRecipient, nil
	default:
		return bellatrix.ExecutionAddress{}, errors.New("unsupported version")
	}
}

// Timestamp returns the timestamp of the proposal.
func (v *VersionedProposal) Timestamp() (uint64, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil {
			return 0, errors.New("no bellatrix block")
		}
		if v.Bellatrix.Body == nil {
			return 0, errors.New("no bellatrix block body")
		}
		if v.Bellatrix.Body.ExecutionPayload == nil {
			return 0, errors.New("no bellatrix block body execution payload header")
		}

		return v.Bellatrix.Body.ExecutionPayload.Timestamp, nil
	case spec.DataVersionCapella:
		if v.Capella == nil {
			return 0, errors.New("no capella block")
		}
		if v.Capella.Body == nil {
			return 0, errors.New("no capella block body")
		}
		if v.Capella.Body.ExecutionPayload == nil {
			return 0, errors.New("no capella block body execution payload header")
		}

		return v.Capella.Body.ExecutionPayload.Timestamp, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil ||
			v.Deneb.Block.Body.ExecutionPayload == nil {
			return 0, errors.New("no deneb block")
		}

		return v.Deneb.Block.Body.ExecutionPayload.Timestamp, nil
	default:
		return 0, errors.New("unsupported version")
	}
}

// BlobSidecars returns the blob sidecars of the proposal.
func (v *VersionedProposal) BlobSidecars() ([]*deneb.BlobSidecar, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		return make([]*deneb.BlobSidecar, 0), nil
	case spec.DataVersionAltair:
		return make([]*deneb.BlobSidecar, 0), nil
	case spec.DataVersionBellatrix:
		return make([]*deneb.BlobSidecar, 0), nil
	case spec.DataVersionCapella:
		return make([]*deneb.BlobSidecar, 0), nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return nil, errors.New("no deneb block")
		}

		return v.Deneb.BlobSidecars, nil
	default:
		return nil, errors.New("unknown version")
	}
}

// String returns a string version of the structure.
func (v *VersionedProposal) String() string {
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
		return "unknown version"
	}
}
