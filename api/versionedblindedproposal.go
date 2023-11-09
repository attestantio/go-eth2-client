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
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// VersionedBlindedProposal contains a versioned blinded proposal.
type VersionedBlindedProposal struct {
	Version   spec.DataVersion
	Bellatrix *apiv1bellatrix.BlindedBeaconBlock
	Capella   *apiv1capella.BlindedBeaconBlock
	Deneb     *apiv1deneb.BlindedBeaconBlock
}

// IsEmpty returns true if there is no proposal.
func (v *VersionedBlindedProposal) IsEmpty() bool {
	return v.Bellatrix == nil && v.Capella == nil && v.Deneb == nil
}

// Slot returns the slot of the blinded proposal.
func (v *VersionedBlindedProposal) Slot() (phase0.Slot, error) {
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
	default:
		return 0, ErrUnsupportedVersion
	}
}

// ProposerIndex returns the proposer index of the blinded proposal.
func (v *VersionedBlindedProposal) ProposerIndex() (phase0.ValidatorIndex, error) {
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
	default:
		return 0, ErrUnsupportedVersion
	}
}

// RandaoReveal returns the RANDAO reveal of the blinded proposal.
func (v *VersionedBlindedProposal) RandaoReveal() (phase0.BLSSignature, error) {
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
	default:
		return phase0.BLSSignature{}, ErrUnsupportedVersion
	}
}

// Graffiti returns the graffiti of the blinded proposal.
func (v *VersionedBlindedProposal) Graffiti() ([32]byte, error) {
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
	default:
		return [32]byte{}, ErrUnsupportedVersion
	}
}

// Attestations returns the attestations of the blinded proposal.
func (v *VersionedBlindedProposal) Attestations() ([]*phase0.Attestation, error) {
	switch v.Version {
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
			v.Deneb.Body == nil {
			return nil, ErrDataMissing
		}

		return v.Deneb.Body.Attestations, nil
	default:
		return nil, ErrUnsupportedVersion
	}
}

// Root returns the root of the blinded proposal.
func (v *VersionedBlindedProposal) Root() (phase0.Root, error) {
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
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// BodyRoot returns the body root of the blinded proposal.
func (v *VersionedBlindedProposal) BodyRoot() (phase0.Root, error) {
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
		if v.Deneb == nil ||
			v.Deneb.Body == nil {
			return phase0.Root{}, ErrDataMissing
		}

		return v.Deneb.Body.HashTreeRoot()
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// ParentRoot returns the parent root of the blinded proposal.
func (v *VersionedBlindedProposal) ParentRoot() (phase0.Root, error) {
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
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// StateRoot returns the state root of the blinded proposal.
func (v *VersionedBlindedProposal) StateRoot() (phase0.Root, error) {
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
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// TransactionsRoot returns the transactions root of the blinded proposal.
func (v *VersionedBlindedProposal) TransactionsRoot() (phase0.Root, error) {
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
	default:
		return phase0.Root{}, ErrUnsupportedVersion
	}
}

// FeeRecipient returns the fee recipient of the blinded proposal.
func (v *VersionedBlindedProposal) FeeRecipient() (bellatrix.ExecutionAddress, error) {
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
	default:
		return bellatrix.ExecutionAddress{}, ErrUnsupportedVersion
	}
}

// Timestamp returns the timestamp of the blinded proposal.
func (v *VersionedBlindedProposal) Timestamp() (uint64, error) {
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
	default:
		return 0, ErrUnsupportedVersion
	}
}

// String returns a string version of the structure.
func (v *VersionedBlindedProposal) String() string {
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
