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
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, ErrDataMissing
		}

		return v.Phase0.Slot, nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return 0, ErrDataMissing
		}

		return v.Altair.Slot, nil
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
		if v.Deneb == nil ||
			v.Deneb.Block == nil {
			return 0, ErrDataMissing
		}

		return v.Deneb.Block.Slot, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// ProposerIndex returns the proposer index of the proposal.
func (v *VersionedProposal) ProposerIndex() (phase0.ValidatorIndex, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return 0, ErrDataMissing
		}

		return v.Phase0.ProposerIndex, nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return 0, ErrDataMissing
		}

		return v.Altair.ProposerIndex, nil
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
		if v.Deneb == nil ||
			v.Deneb.Block == nil {
			return 0, ErrDataMissing
		}

		return v.Deneb.Block.ProposerIndex, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// RandaoReveal returns the RANDAO reveal of the proposal.
func (v *VersionedProposal) RandaoReveal() (phase0.BLSSignature, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil ||
			v.Phase0.Body == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}

		return v.Phase0.Body.RANDAOReveal, nil
	case spec.DataVersionAltair:
		if v.Altair == nil ||
			v.Altair.Body == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}

		return v.Altair.Body.RANDAOReveal, nil
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
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil {
			return phase0.BLSSignature{}, ErrDataMissing
		}

		return v.Deneb.Block.Body.RANDAOReveal, nil
	default:
		return phase0.BLSSignature{}, ErrUnsupportedVersion
	}
}

// Graffiti returns the graffiti of the proposal.
func (v *VersionedProposal) Graffiti() ([32]byte, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil ||
			v.Phase0.Body == nil {
			return [32]byte{}, ErrDataMissing
		}

		return v.Phase0.Body.Graffiti, nil
	case spec.DataVersionAltair:
		if v.Altair == nil ||
			v.Altair.Body == nil {
			return [32]byte{}, ErrDataMissing
		}

		return v.Altair.Body.Graffiti, nil
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
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil {
			return [32]byte{}, ErrDataMissing
		}

		return v.Deneb.Block.Body.Graffiti, nil
	default:
		return [32]byte{}, ErrUnsupportedVersion
	}
}

// Attestations returns the attestations of the proposal.
func (v *VersionedProposal) Attestations() ([]*phase0.Attestation, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil ||
			v.Phase0.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Phase0.Body.Attestations, nil
	case spec.DataVersionAltair:
		if v.Altair == nil ||
			v.Altair.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Altair.Body.Attestations, nil
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Bellatrix.Body.Attestations, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Capella.Body.Attestations, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Deneb.Block.Body.Attestations, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// Root returns the root of the proposal.
func (v *VersionedProposal) Root() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Phase0.HashTreeRoot()
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Altair.HashTreeRoot()
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
		if v.Deneb == nil ||
			v.Deneb.Block == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.Block.HashTreeRoot()
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// BodyRoot returns the body root of the proposal.
func (v *VersionedProposal) BodyRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Phase0.Body.HashTreeRoot()
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Altair.Body.HashTreeRoot()
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
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.Block.Body.HashTreeRoot()
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// ParentRoot returns the parent root of the proposal.
func (v *VersionedProposal) ParentRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Phase0.ParentRoot, nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Altair.ParentRoot, nil
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
		if v.Deneb == nil ||
			v.Deneb.Block == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.Block.ParentRoot, nil
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// StateRoot returns the state root of the proposal.
func (v *VersionedProposal) StateRoot() (phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Phase0.StateRoot, nil
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Altair.StateRoot, nil
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
		if v.Deneb == nil ||
			v.Deneb.Block == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.Block.StateRoot, nil
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// Transactions returns the transactions of the proposal.
func (v *VersionedProposal) Transactions() ([]bellatrix.Transaction, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil ||
			v.Bellatrix.Body.ExecutionPayload == nil {
			return nil, ErrDataMissing
		}

		return v.Bellatrix.Body.ExecutionPayload.Transactions, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayload == nil {
			return nil, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayload.Transactions, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil ||
			v.Deneb.Block.Body.ExecutionPayload == nil {
			return nil, ErrDataMissing
		}

		return v.Deneb.Block.Body.ExecutionPayload.Transactions, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// FeeRecipient returns the fee recipient of the proposal.
func (v *VersionedProposal) FeeRecipient() (bellatrix.ExecutionAddress, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil ||
			v.Bellatrix.Body.ExecutionPayload == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Bellatrix.Body.ExecutionPayload.FeeRecipient, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayload == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayload.FeeRecipient, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil ||
			v.Deneb.Block.Body.ExecutionPayload == nil {
			return bellatrix.ExecutionAddress{}, ErrDataMissing
		}

		return v.Deneb.Block.Body.ExecutionPayload.FeeRecipient, nil
	default:
		return bellatrix.ExecutionAddress{}, ErrUnsupportedVersion
	}
}

// Timestamp returns the timestamp of the proposal.
func (v *VersionedProposal) Timestamp() (uint64, error) {
	switch v.Version {
	case spec.DataVersionBellatrix:
		if v.Bellatrix == nil ||
			v.Bellatrix.Body == nil ||
			v.Bellatrix.Body.ExecutionPayload == nil {
			return 0, ErrDataMissing
		}

		return v.Bellatrix.Body.ExecutionPayload.Timestamp, nil
	case spec.DataVersionCapella:
		if v.Capella == nil ||
			v.Capella.Body == nil ||
			v.Capella.Body.ExecutionPayload == nil {
			return 0, ErrDataMissing
		}

		return v.Capella.Body.ExecutionPayload.Timestamp, nil
	case spec.DataVersionDeneb:
		if v.Deneb == nil ||
			v.Deneb.Block == nil ||
			v.Deneb.Block.Body == nil ||
			v.Deneb.Block.Body.ExecutionPayload == nil {
			return 0, ErrDataMissing
		}

		return v.Deneb.Block.Body.ExecutionPayload.Timestamp, nil
	default:
		return 0, ErrUnsupportedVersion
	}
}

// Blobs returns the blobs of the proposal.
func (v *VersionedProposal) Blobs() ([]deneb.Blob, error) {
	switch v.Version {
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return nil, ErrDataMissing
		}

		return v.Deneb.Blobs, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// KZGProofs returns the KZG proofs of the proposal.
func (v *VersionedProposal) KZGProofs() ([]deneb.KZGProof, error) {
	switch v.Version {
	case spec.DataVersionDeneb:
		if v.Deneb == nil {
			return nil, ErrDataMissing
		}

		return v.Deneb.KZGProofs, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// String returns a string version of the structure.
func (v *VersionedProposal) String() string {
	switch v.Version {
	case spec.DataVersionPhase0:
		if v.Phase0 == nil {
			return ""
		}

		return v.Phase0.String()
	case spec.DataVersionAltair:
		if v.Altair == nil {
			return ""
		}

		return v.Altair.String()
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
